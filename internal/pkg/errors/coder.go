package errors

import (
	"fmt"
	"github.com/json-iterator/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path/filepath"
	"runtime"
	"strings"
)

// 业务代码通用的错误
type ErrorCoder interface {
	Error() string
	Code() uint32
	Msg() string
	Where() string // 第一次生成这个错的地方, 第一次: 当newCoder和wrap一个非errorCoder的时候
}

// Grpc的错误
type GRPCStatuser interface {
	GRPCStatus() *status.Status
	Error() string
}

// 不要将此结构体作为NewError返回结果: https://golang.org/doc/faq#nil_error
type ErrorCode struct {
	code  uint32
	msg   string
	where string
}

type errorCodeUnmarshal struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

const (
	Separator = ":: "
)

func Unmarshal(bs []byte) (err ErrorCoder) {
	var eu errorCodeUnmarshal
	e := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bs, &eu)
	if e != nil {
		err = NewCodere(500, e, fmt.Sprintf("Unmarshal Fail, row: %s", bs))
		return
	}

	err = NewCoderWhere(eu.Code, 2, eu.Msg)
	return
}

// 错误，附带code
func (e *ErrorCode) Error() string {
	if e == nil {
		return ""
	}

	if e.code == 0 {
		return e.msg
	}
	msg := e.msg
	if msg == "" {
		msg = "empty error"
	}
	return fmt.Sprintf("code = %d ; msg = %s", e.code, msg)
}

// 不带code的错误消息
func (e *ErrorCode) Msg() string {
	if e == nil {
		return ""
	}
	return e.msg
}

func (e *ErrorCode) Code() uint32 {
	if e == nil {
		return 0
	}
	return e.code
}
func (e *ErrorCode) Where() string {
	if e == nil {
		return ""
	}
	return e.where
}

type CallDepth int

type ExtendMsg string

// NewCoder 返回一个竹子的高级错误, 所有竹子项目都应该使用这个错误.
// 入参是多个interface, 将通过类型判断是含义. 灵感来自: [Upspin 中的错误处理 —— 来自 Rob Pike](https://studygolang.com/articles/12045)
//
// 可以传递的类型有
// - int/int32/uint/uint32: code码
// - string: `消息`, 一般是"动作", 如调用短信微服务发送短信可以这样表述: "Sms.SendSms", 不建议传递多个string, 如果传递了多个则会用", "隔开.
// - error: `错误` 又分了3种 [支持一个]
//   - ErrorCoder: 此方法也可以当Wrap使用, 所以可以传递一个高级的错误, 这种情况下 错误码将会被继承
//   - GRPCStatuser: 同ErrorCode, 处理Grpc的高级错误.
//   - error: 普通错误
//
// ErrorCode.Error() 返回的格式为: code = 400, msg = xxx.
// 其中msg由`消息`和`错误`组成, 消息在前, 错误在后, 由":: "符号拼接.
func NewCoder(args ...interface{}) ErrorCoder {
	var code uint32 // code 码
	var msgs []string
	var deepError []string // 在wrap时有用, 代表最核心(底层)的错误, 例如err.Error(), 会放在最后.
	var extend string      // 扩展信息, 不wrap, 而是放在error之后
	var where string
	var dep = CallDepth(1)

	for _, v := range args {
		switch t := v.(type) {
		case CallDepth:
			// for callDepth(where)
			dep = t
		case int32, int, int64, uint32, uint, uint64:
			// for code
			code = toUint32(t)
		case string:
			// for msg, 如果传递了多条, 最后会用`, `隔开
			msgs = append(msgs, t)
		case ExtendMsg:
			extend = string(t)
		case []string:
			// for msg
			msgs = append(msgs, t...)
		case ErrorCoder:
			deepError = append(deepError, t.Msg())
			tc := t.Code()
			if tc != 0 {
				code = tc
			}
			where = t.Where()
		case GRPCStatuser:
			s := t.GRPCStatus()
			deepError = append(deepError, s.Message())

			if s.Code() == codes.Unknown {
			} else if s.Code() < 20 {
				// 只要是grpc自带的错误就说明是系统错误
				code = 500
				deepError = append(deepError, fmt.Sprintf("%d:%s", s.Code(), s.Message()))
			} else {
				code = uint32(s.Code())
			}
		case error:
			deepError = append(deepError, t.Error())
		default:
			// 其他类型就用%+v格式化成字符串
			msgs = append(msgs, fmt.Sprintf("%+v", t))
		}
	}

	deepErrorStr := strings.Join(deepError, ",")

	// 特殊处理xorm返回的错误, 忽略这个错误
	if strings.Contains(deepErrorStr, "No content found to be updated") {
		return nil
	}

	// msg + 上一层的error + 后缀

	var msg string
	if len(msgs) != 0 {
		msg = strings.Join(msgs, ", ")
	}

	if deepErrorStr != "" {
		if msg != "" {
			msg += Separator + deepErrorStr
		} else {
			msg = deepErrorStr
		}
	}
	if extend != "" {
		if msg != "" {
			msg += ", " + extend
		} else {
			msg = extend
		}
	}

	if where == "" {
		where = caller(dep, false)
	}

	return &ErrorCode{code: code, msg: msg, where: where}
}

func toUint32(i interface{}) uint32 {
	switch t := i.(type) {
	case int32:
		return uint32(t)
	case int:
		return uint32(t)
	case uint32:
		return t
	case uint:
		return uint32(t)
	}
	return 0
}

func NewCoderWhere(code uint32, callDepth int, msg string, extMsg ...string) ErrorCoder {
	return NewCoder(CallDepth(callDepth+1), code, extMsg, msg)
}

func NewCodere(code uint32, err error, extMsg ...string) ErrorCoder {
	return NewCoder(CallDepth(2), code, err, extMsg)
}

// Wrap 为error添加一个说明, 当这个err不确定是否应该报500或者是由其他服务返回时使用
// 如果err是ErrorCoder或者GRPCStatuser, code将继承, 否则code为0
func Wrap(err error, extMsg ...string) ErrorCoder {
	return NewCoder(err, CallDepth(2), extMsg)
}

func Extend(err error, ext string) ErrorCoder {
	return NewCoder(err, CallDepth(2), ExtendMsg(ext))
}

// 拼装多个error为一个
// 当所有的error都是nil的时候, 也返回nil
// where: 是concat的发生地
// code: 选取第一个code
// msg: 用 | 分割
func Concat(es ...error) ErrorCoder {
	var gmsg string
	var gcode uint32
	var gwhere string
	for _, err := range es {
		var msg string
		var code uint32

		switch v := err.(type) {
		case ErrorCoder:
			msg = v.Msg()
			code = v.Code()
		case GRPCStatuser:
			s := v.GRPCStatus()
			if s.Code() == codes.Unknown {
				code = 0
			} else if s.Code() < 20 {
				// 只要是grpc自带的错误就说明是系统错误
				code = 500
			} else {
				code = uint32(s.Code())
			}
			msg = s.Message()
		default:
			msg = v.Error()
			code = 0
		}

		// 选取第一个code
		if gcode == 0 && code != 0 {
			gcode = code
		}
		gmsg += msg + " | "
	}
	if len(gmsg) != 0 {
		gmsg = gmsg[:len(gmsg)-3]
	}

	gwhere = caller(1, false)
	return &ErrorCode{code: gcode, msg: gmsg, where: gwhere}
}

// Wrap 为error添加一个说明, 当这个err不确定是否应该报500或者是由其他服务返回时使用
// 如果err是ErrorCoder或者GRPCStatuser, code将继承, 否则code为0
// 默认callDepth为1, 可自定义callDepth.
func WrapWhere(err error, callDepth int, extMsg ...string) ErrorCoder {
	return NewCoder(err, CallDepth(callDepth+1), extMsg)
}

func caller(calldepth CallDepth, short bool) string {
	_, file, line, ok := runtime.Caller(int(calldepth) + 1)
	if !ok {
		file = "???"
		line = 0
	} else if short {
		file = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d", file, line)
}

func New(msg string) ErrorCoder {
	where := caller(1, false)
	return &ErrorCode{code: 0, msg: msg, where: where}
}

// GetCode 返回err的code, 如果不是ErrorCode则返回0
func GetCode(err error) uint32 {
	switch t := err.(type) {
	case ErrorCoder:
		return t.Code()
	case GRPCStatuser:
		s := t.GRPCStatus()
		if s.Code() == codes.Unknown {
			return 500
		} else if s.Code() < 20 {
			// 只要是grpc自带的错误就说明是系统错误
			return 500
		} else {
			return uint32(s.Code())
		}
	}

	return 0
}
