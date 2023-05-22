package pclog

import (
	"nsqco/pcerror"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var Pclog *Logger

type Fields map[string]interface{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func New() *Logger {
	log := logrus.New()
	log.Formatter = &InFormatter{}

	return &Logger{
		Log:   log,
		level: LevelDefault,
	}
}

type Logger struct {
	Log   *logrus.Logger
	level Level
}

// func (l *Logger) WithGrpcError(s interface{}) (e *Entry) {
// 	err := pcerror.Parse(s)
// 	return l.WithError(err)
// }
func (l *Logger) WithFields(fields map[string]interface{}) (e *Entry) {

	//	clean 保留字
	delete(fields, "code")
	delete(fields, "err_msg")
	delete(fields, "extrainfo")
	delete(fields, "time")
	delete(fields, "service")
	delete(fields, "origing_err")
	delete(fields, "error")

	return &Entry{
		Entry: l.Log.WithFields(logrus.Fields(ConvertInt64ToString(fields))),
		level: l.level,
	}
}

func (l *Logger) WithError(err error) (e *Entry) {
	switch v := err.(type) {
	case pcerror.Error:
		e = &Entry{
			Entry: l.Log.WithFields(logrus.Fields{
				"code":       v.Code,
				"err_msg":    v.Msg,
				"extrainfo":  ConvertInt64ToString(v.ExtraInfo),
				"time":       v.Time,
				"service":    v.Service,
				"origin_err": v.OriginErr,
			}),
			level: l.level,
		}

	default:
		e = &Entry{
			Entry: l.Log.WithError(err),
			level: l.level,
		}

	}
	return
}
func (l *Logger) SetLogLevel(lv Level) (lg *Logger) {
	l.level = lv
	return l
}

func (l *Logger) Println(args ...interface{}) {
	if l.level <= LevelDefault {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelDefault,
		}).Println(args...)
	}
}
func (l *Logger) Printf(format string, args ...interface{}) {
	if l.level <= LevelDefault {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelDefault,
		}).Printf(format, args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.level <= LevelDebug {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelDebug,
		}).Println(args...)
	}
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelDebug,
		}).Printf(format, args...)
	}
}
func (l *Logger) Info(args ...interface{}) {
	if l.level <= LevelInfo {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelInfo,
		}).Println(args...)
	}
}
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelInfo,
		}).Printf(format, args...)
	}
}
func (l *Logger) Notice(args ...interface{}) {
	if l.level <= LevelNotice {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelNotice,
		}).Println(args...)
	}
}
func (l *Logger) Noticef(format string, args ...interface{}) {
	if l.level <= LevelNotice {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelNotice,
		}).Printf(format, args...)
	}
}
func (l *Logger) Warn(args ...interface{}) {
	if l.level <= LevelWarn {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelWarn,
		}).Println(args...)
	}
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelWarn,
		}).Printf(format, args...)
	}
}
func (l *Logger) Error(args ...interface{}) {
	if l.level <= LevelError {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelError,
		}).Println(args...)
	}
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelError,
		}).Printf(format, args...)
	}
}
func (l *Logger) Critical(args ...interface{}) {
	if l.level <= LevelCritical {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelCritical,
		}).Println(args...)
	}
}
func (l *Logger) Criticalf(format string, args ...interface{}) {
	if l.level <= LevelCritical {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelCritical,
		}).Printf(format, args...)
	}
}
func (l *Logger) Alert(args ...interface{}) {
	if l.level <= LevelAlert {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelAlert,
		}).Println(args...)
	}
}
func (l *Logger) Alertf(format string, args ...interface{}) {
	if l.level <= LevelAlert {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelAlert,
		}).Printf(format, args...)
	}
}
func (l *Logger) Emergency(args ...interface{}) {
	if l.level <= LevelEmergency {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelEmergency,
		}).Println(args...)
	}
}
func (l *Logger) Emergencyf(format string, args ...interface{}) {
	if l.level <= LevelEmergency {
		l.Log.WithFields(logrus.Fields{
			"severity": LevelEmergency,
		}).Printf(format, args...)
	}
}
