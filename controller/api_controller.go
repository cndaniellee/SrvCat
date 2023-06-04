package controller

import (
	"SrvCat/config"
	"SrvCat/response"
	"SrvCat/storage"
	"SrvCat/totp"
	"SrvCat/util"
	"encoding/base64"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/skip2/go-qrcode"
	"gopkg.in/go-playground/validator.v9"
	"image/color"
	"net"
	"time"
)

type ApiController struct {
	Ctx      iris.Context
	Validate *validator.Validate
}

func NewApiController() *ApiController {
	return &ApiController{
		Validate: validator.New(),
	}
}

func (c *ApiController) GetInit() {
	init, err := storage.Sqlite.GetInit()
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	resp := &response.InitResp{Init: init}
	if resp.Init {
		resp.Name = config.Config.Machine.Name
		resp.Secret = config.Config.Machine.Secret
		if resp.Name != "" && resp.Secret != "" {
			qr, err := qrcode.New(totp.GenerateURI(resp.Name, resp.Secret), qrcode.High)
			if err != nil {
				golog.Errorf("[Api]: %v", err)
				_, _ = c.Ctx.JSON(response.ServerError)
				return
			}
			qr.BackgroundColor = color.Transparent
			png, err := qr.PNG(200)
			if err != nil {
				golog.Errorf("[Api]: %v", err)
				_, _ = c.Ctx.JSON(response.ServerError)
				return
			}
			resp.Qrcode = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		}
	}
	_, _ = c.Ctx.JSON(response.DataResponse{Response: response.Success, Data: resp})
}

func (c *ApiController) PostInit() {
	err := storage.Sqlite.UpdateInit()
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	_, _ = c.Ctx.JSON(response.Success)
}

type PostValidateArg struct {
	Code int `json:"code" validate:"required"`
}

func (c *ApiController) PostValidate() {
	// 校验参数
	var arg PostValidateArg
	if err := util.ValidateJson(&arg, c.Ctx); err != nil {
		golog.Errorf("[Api]: %v", err)
		return
	}
	init, err := storage.Sqlite.GetInit()
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	if init {
		_, _ = c.Ctx.JSON(response.InitFirst)
		return
	}
	// 获取远端IP
	ip, _, err := net.SplitHostPort(c.Ctx.Request().RemoteAddr)
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	// 检查有效期内的验证
	verified, err := storage.Sqlite.GetVerified(ip, time.Now().Add(-time.Duration(config.Config.Settings.VerifyPeriod)*time.Minute).UnixMilli())
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	if verified {
		golog.Errorf("[Api]: Completed verify exists, ip: %v", ip)
		_, _ = c.Ctx.JSON(response.AlreadyVerify)
		return
	}
	// 检查超限
	if config.Config.Settings.LimitTimes != 0 {
		count, err := storage.Sqlite.GetAccessCount(time.Now().Add(-time.Duration(config.Config.Settings.LimitPeriod) * time.Minute).UnixMilli())
		if err != nil {
			golog.Errorf("[Api]: %v", err)
			_, _ = c.Ctx.JSON(response.ServerError)
			return
		}
		if count >= config.Config.Settings.LimitTimes {
			golog.Errorf("[Api]: Validate reach limit, ip: %v", ip)
			_, _ = c.Ctx.JSON(response.ReachLimit)
			return
		}
	}
	// 验证
	result, err := totp.ValidateCode(arg.Code)
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	if !result {
		err = storage.Sqlite.AddAccess(ip, storage.AccessCategoryCheck)
		if err != nil {
			golog.Errorf("[Api]: %v", err)
			_, _ = c.Ctx.JSON(response.ServerError)
			return
		}
		golog.Infof("[Api]: Check failed from ip: %v", ip)
		_, _ = c.Ctx.JSON(response.InvalidCode)
		return
	}
	// 追加IP地址到允许列表
	err = storage.Sqlite.AddAccess(ip, storage.AccessCategoryVerified)
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	_, _ = c.Ctx.JSON(response.DataResponse{Response: response.Success, Data: config.Config.Forwards})
}

func (c *ApiController) GetCheck() {
	init, err := storage.Sqlite.GetInit()
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	if init {
		_, _ = c.Ctx.JSON(response.InitFirst)
		return
	}
	// 获取远端IP
	ip, _, err := net.SplitHostPort(c.Ctx.Request().RemoteAddr)
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	resp := response.CheckResp{Time: time.Now().UnixMilli()}
	verified, err := storage.Sqlite.GetVerified(ip, time.Now().Add(-time.Duration(config.Config.Settings.VerifyPeriod)*time.Minute).UnixMilli())
	if err != nil {
		golog.Errorf("[Api]: %v", err)
		_, _ = c.Ctx.JSON(response.ServerError)
		return
	}
	if verified {
		resp.Exist = true
		resp.Forwards = config.Config.Forwards
	}
	_, _ = c.Ctx.JSON(response.DataResponse{Response: response.Success, Data: resp})
}
