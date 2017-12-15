package response

const (
	statusOk  = 200
	statusErr = 400
)

type IResponse interface {
	SendOk(dataok interface{})
	SendError(code int, msg string)
}

/*******************************************************************************
//接口执行OK
*******************************************************************************/
type ResponseOk struct {
	Status int         `json:"code"`
	Data   interface{} `json:"data"`
}

/*******************************************************************************
//接口执行失败
*******************************************************************************/
type ResponseError struct {
	Status  int    `json:"code"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type Response struct {
	RsOk  *ResponseOk    `json:"responseOk"`
	RsErr *ResponseError `json:"responseError"`
}

func (r *Response) NewResponse(code int, msg string, data interface{}) {
	r.RsErr = new(ResponseError)
	r.RsOk = new(ResponseOk)
	r.RsErr.ErrCode = (int)(code)
	r.RsErr.ErrMsg = msg
	r.RsErr.Status = statusErr
	r.RsOk.Status = statusOk
	r.RsOk.Data = data
}

func (r *Response) SendOk(dataok interface{}) {
	r.RsOk = new(ResponseOk)
	r.RsOk.Status = statusOk
	r.RsOk.Data = dataok
}

func (r *Response) SendError(code int, msg string) {
	r.RsErr = new(ResponseError)
	r.RsErr.ErrCode = code
	r.RsErr.ErrMsg = msg
	r.RsErr.Status = statusErr
}

func (r *Response) GetOK() *ResponseOk {
	return r.RsOk
}
func (r *Response) GetError() *ResponseError {
	return r.RsErr
}
