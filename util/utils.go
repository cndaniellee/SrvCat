package util

import (
	"SrvCat/response"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"gopkg.in/go-playground/validator.v9"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var validate = validator.New()

func ValidateForm(arg interface{}, ctx iris.Context) error {
	if err := ctx.ReadForm(arg); err != nil {
		_, _ = ctx.JSON(response.InvalidParameter)
		return err
	}
	err := validate.Struct(arg)
	if err != nil {
		_, _ = ctx.JSON(response.MissingParameter)
		return err
	}
	return nil
}

func ValidateJson(arg interface{}, ctx iris.Context) error {
	if err := ctx.ReadJSON(arg); err != nil {
		_, _ = ctx.JSON(response.InvalidParameter)
		return err
	}
	err := validate.Struct(arg)
	if err != nil {
		_, _ = ctx.JSON(response.MissingParameter)
		return err
	}
	return nil
}

// 打印重大错误
func FailOnException(msg string, err error) {
	if err != nil {
		golog.Fatalf("%s: \n%s", msg, err)
	}
}

func UnmarshalJson(text string) map[string]string {
	result := map[string]string{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		golog.Errorf("Unmarshal json err: %v", err)
	}
	return result
}

// 生成MD5
func StrToMd5(source string) string {
	h := md5.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}

// 生成随机字符串
func GetRandomString(length int, model string) string {
	char := ""
	if strings.Contains(model, "a") {
		char += "abcdefghijklmnopqrstuvwxyz"
	}
	if strings.Contains(model, "A") {
		char += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if strings.Contains(model, "0") {
		char += "0123456789"
	}
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < length; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).Anonymous && t.Field(i).Tag.Get("map") != "transient" {
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return data
}

func Struct2SnakeMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).Anonymous && t.Field(i).Tag.Get("map") != "transient" {
			data[ToSnakeString(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}
	return data
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func ToSnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func ToCamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 复制结构体中的相同字段，source传结构体，target传结构体指针
func CopyProperty(source interface{}, target interface{}) {
	if source != nil {
		var copyObjV = reflect.ValueOf(target)
		//反射获取属性的类型和值
		var objT = reflect.TypeOf(source)
		var objV = reflect.ValueOf(source)
		//对应设置的数值
		for i := 0; i < objT.NumField(); i++ {
			if objT.Field(i).Anonymous || objT.Field(i).Tag.Get("map") == "transient" {
				continue
			}
			property := copyObjV.Elem().FieldByName(objT.Field(i).Name)
			if property.IsValid() {
				property.Set(objV.Field(i))
			}
		}
	}
}

func DesensitizeEmail(source string) string {
	reg, err := regexp.Compile("(\\w+)\\w{3}@(\\w+)")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(source, "$1***@$2")
}

func DesensitizePhone(source string) string {
	reg, err := regexp.Compile("(\\d{3})\\d*(\\d{4})")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(source, "$1****$2")
}

func DesensitizeCert(source string) string {
	reg, err := regexp.Compile("(\\d{6})\\d*([\\d(Xx)]{4})")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(source, "$1****$2")
}

var contentType = map[string]string{
	"jpg":  "image/jpeg",
	"gif":  "image/gif",
	"png":  "image/png",
	"jpeg": "image/jpeg",
	"mp4":  "video/mp4",
	"mpg":  "video/mpeg",
	"mpeg": "video/mpeg",
	"mov":  "video/quicktime",
	"avi":  "video/x-msvideo",
}

func Ext2ContentType(extension string) (string, error) {
	if result, ok := contentType[extension]; ok {
		return result, nil
	}
	return "", errors.New("invalid extension")
}

func UintIndexOf(target uint, items []uint) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
