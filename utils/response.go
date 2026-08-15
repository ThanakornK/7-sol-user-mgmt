package utils

type ResponseStruct struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func NewResponseStruct(status bool, msg string, data interface{}) *ResponseStruct {
	return &ResponseStruct{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
}

func NewSuccessResponseStruct(msg string, data interface{}) *ResponseStruct {
	return NewResponseStruct(true, msg, data)
}

func NewErrorResponseStruct(msg string, error interface{}) *ResponseStruct {
	return &ResponseStruct{
		Status: false,
		Msg:    msg,
		Error:  error,
	}
}
