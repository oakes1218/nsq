package pclog

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type InFormatter struct {
}

func (f *InFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["time"] = entry.Time.Format(time.RFC3339Nano)
	data["message"] = entry.Message
	v, ok := data["severity"].(Level)
	if !ok {
		return nil, errors.New("err level type error")
	}
	l := LevelResolve(v)
	data["level"] = l

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error, %v", err)
	}
	return append(serialized, '\n'), nil
}
