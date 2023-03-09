package response

/**
约定返回的基本格式
*/

type Response struct {
	Code int    `json:"code"`
	Note string `json:"note"`
}
