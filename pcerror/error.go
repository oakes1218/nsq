package pcerror

import (
	"encoding/json"
	"time"

	"google.golang.org/grpc"
)

// Service Service
type Service string

var Extrainfo map[string]interface{}

const (
	// TEST
	TEST Service = "TEST"

	// Unknown Unknown
	Unknown Service = "UNKNOWN"
)

// Code
const (
	TestCode = 12501
)

// var reError = regexp.MustCompile(`\[Error (\d+)\]\t([A-Z]{1,}):([\w\s]{1,}),Extra:(\S+.)\t(?:Origin:(\[\w+.*\].\w+.*|\w+.*))?`)
var defaultError = Error{
	Msg:       "inerror parse error",
	Code:      129999999,
	ExtraInfo: make(map[string]interface{}),
	Service:   Unknown,
	OriginErr: "",
}

// Errors Errors
type Errors struct {
	s Service
}

// Wrap Wrap
func (e *Errors) Wrap(code int, msg string, extrainfo map[string]interface{}, err error) (er Error) {
	return Wrap(code, msg, e.s, extrainfo, err)
}

// Parse Parse
func (e *Errors) Parse(s interface{}) (er Error) {
	return Parse(s)
}

// New New
func New(s Service) (er *Errors) {
	return &Errors{
		s: s,
	}
}

// Error Error
type Error struct {
	Msg       string                 `json:"msg"`
	Code      int                    `json:"code"`
	ExtraInfo map[string]interface{} `json:"extrainfo"`
	Time      int64                  `json:"time"`
	Service   Service                `json:"service"`
	OriginErr string                 `json:"origin_err"`
}

// SwitchSysCode SwitchSysCode
func SwitchSysCode(s Service) (code int) {
	switch s {
	case TEST:
		code = TestCode
	default:
		code = 12999
	}
	return
}

func (e Error) Error() (s string) {
	var b []byte
	b, _ = json.Marshal(e)
	return string(b)
}

// Wrap Wrap
func Wrap(code int, msg string, service Service, extrainfo map[string]interface{}, err error) (er Error) {
	c := SwitchSysCode(service)
	code = (c * 10000) + code
	var errs string
	if err != nil {
		errs = err.Error()
	} else {
		errs = ""
	}
	er = Error{
		Code:      code,
		Msg:       msg,
		ExtraInfo: extrainfo,
		Service:   service,
		Time:      time.Now().Unix(),
		OriginErr: errs,
	}
	return
}

// Parse Parse
func Parse(s interface{}) (e Error) {
	var ss string
	switch v := s.(type) {
	case string:
		ss = v
	case error:
		ss = grpc.ErrorDesc(v)
	default:
		ss = "can't parse"
	}
	e = Error{}
	err := json.Unmarshal([]byte(ss), &e)
	if err != nil {
		e = defaultError
	}
	return

}
