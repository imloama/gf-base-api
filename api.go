package baseapi


import (
	"github.com/gogf/gf/net/ghttp"
)

const CODE_OK = 200
const CODE_ERROR = 500

// 数据返回通用JSON数据结构
type APIResult struct {
	Code    int         `json:"code" example:"200"`    // 错误码，200表示成功，其它表示失败
	Message string      `json:"msg"`     // 提示信息
	ErrCode string      `json:"errCode"` // 业务错误码
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
}

type TableAPIResult struct {
	Code    int         `json:"code" example:"200"`    // 错误码，200表示成功，其它表示失败
	Message string      `json:"msg"`     // 提示信息
	ErrCode string      `json:"errCode"` // 业务错误码
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
	Current int         `json:"current"` //	当前页码
	Page    int         `json:"page"`    //多少页
	Total   int         `json:"total"`   //	总页数
}

type PaginationData struct {
	List    interface{} `json:"list"`
	Cursor  int         `json:"cursor"`
}

// 标准返回结果数据结构封装。
func Json(r *ghttp.Request, code int, message string, errCode string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(APIResult{
		Code:    code,
		Message: message,
		ErrCode: errCode,
		Data:    responseData,
	})
}

// 返回JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, code int, msg string, errCode string, data ...interface{}) {
	Json(r, code, msg, errCode, data...)
	r.Exit()
}

func JsonPagination(r *ghttp.Request, code int, message string, errCode string, cursor int, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(APIResult{
		Code:    code,
		Message: message,
		ErrCode: errCode,
		Data:    PaginationData{List: responseData, Cursor: cursor},
	})
}

func JsonPaginationExit(r *ghttp.Request, code int, msg string, errCode string, cursor int, data ...interface{}) {
	JsonPagination(r, code, msg, errCode, cursor, data...)
	r.Exit()
}

/**
 * 正常返回API
 */
func OK(r *ghttp.Request, data interface{}) {
	r.Response.WriteJson(APIResult{
		Code:    CODE_OK,
		Message: "", // cannot use nil as type string in field value
		Data:    data,
	})
}

/**
 * 错误返回结果
 */
func Fail(r *ghttp.Request, message string, errorCode string) {
	r.Response.WriteJson(APIResult{
		Code:    CODE_ERROR,
		Message: message,
		ErrCode: errorCode,
		Data:    nil,
	})
}

/**
 * 表格接口返回正常的结果
 */
func OkTable(r *ghttp.Request, current int, page int, total int, list interface{}) {
	r.Response.WriteJson(TableAPIResult{
		Code:    CODE_OK,
		Message: "ok",
		Current: current,
		Page:    page,
		Total:   total,
		Data:    list,
	})
}
