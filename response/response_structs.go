package response

type DataResponse struct {
	Response
	Data interface{} `json:"data"`
}

type List struct {
	List     interface{} `json:"list"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}

type InitResp struct {
	Init   bool   `json:"init"`
	Name   string `json:"name,omitempty"`
	Secret string `json:"secret,omitempty"`
	Qrcode string `json:"qrcode,omitempty"`
}

type CheckResp struct {
	Time     int64    `json:"time"`
	Exist    bool     `json:"exist"`
	Forwards []string `json:"forwards,omitempty"`
}
