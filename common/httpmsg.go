package common

type HttpMsg struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewHttpMsg(code int, flag bool, msg string, data interface{}) *HttpMsg {
	return &HttpMsg{
		Code:    code,
		Success: flag,
		Message: msg,
		Data:    data,
	}
}

func NewHttpMsgData(data interface{}) *HttpMsg {
	return &HttpMsg{
		Code:    SUCCESS,
		Success: true,
		Message: "",
		Data:    data,
	}
}
