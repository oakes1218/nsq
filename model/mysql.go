package model

import (
	"fmt"
	"nsqco/pclog"
	"time"

	"github.com/jinzhu/gorm"
)

type Options struct {
	UserName     string
	Pass         string
	Host         string
	Port         string
	DB           string
	Debug        bool
	MaxIdeConns  int
	MaxOpenConns int
	MaxLifetime  int
}

type GormLogger struct{}

var UserM *gorm.DB

func InitUser() {
	UserM = ConnDB(Options{
		UserName:     "root",
		Pass:         "password",
		Host:         "127.0.0.1",
		Port:         "3306",
		DB:           "user",
		Debug:        true,
		MaxIdeConns:  10,
		MaxOpenConns: 100,
		MaxLifetime:  30,
	})
}

const UserTable = "user"

type User struct {
	ID        int64     `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"id,omitempty"`
	Name      string    `gorm:"unique_index:name" json:"name,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

func ConnDB(opt Options) *gorm.DB {
	addr := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true",
		opt.UserName,
		opt.Pass,
		opt.Host,
		opt.Port,
		opt.DB,
	)

	conn, err := gorm.Open("mysql", addr)
	if err != nil {
		panic(err)
	}

	conn.DB().SetMaxIdleConns(opt.MaxIdeConns)
	conn.DB().SetMaxOpenConns(opt.MaxOpenConns)
	conn.DB().SetConnMaxLifetime(time.Duration(opt.MaxLifetime) * time.Second)

	conn.SetLogger(&GormLogger{})
	conn.LogMode(opt.Debug)

	return conn
}

// Print handles log events from Gorm for the custom logger.
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		pclog.Pclog.WithFields(
			pclog.Fields{
				"module":      "gorm",
				"type":        "sql",
				"sql_rows":    v[5],
				"sql_src_ref": v[1],
				"sql_values":  v[4],
			},
		).Info(v[3])
	case "log":
		pclog.Pclog.WithFields(pclog.Fields{"module": "gorm", "type": "log"}).Info(v[2])
	}
}

func CreateUser(u *User) error {
	res := UserM.Table(UserTable).Create(&u)
	return res.Error
}
