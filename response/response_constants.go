package response

/**
定义返回值常量，便于前端处理
*/

var Success = Response{0, "请求成功"}

/**
用户端报错
*/

var InitFirst = Response{1000, "请先初始化"}

var InvalidCode = Response{1001, "无效校验码"}

/**
通用报错
*/

var InvalidParameter = Response{3000, "无效参数"}

var MissingParameter = Response{3001, "参数缺失"}

var AlreadyVerify = Response{4000, "已完成过验证"}

var ReachLimit = Response{4001, "验证次数超限"}

var ServerError = Response{5000, "服务器错误"}
