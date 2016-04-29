package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"huijujiayuan.com/err_code"
)

const (
	SQL_NUM     = 150 //SQL批处理条数
	MAX_CLIENT  = 400 //最大链接个数
	INIT_CLIENT = 10  //初始化链接个数
)

var (
	PonyDB *xorm.Engine //业务数据库
)

func init() {
	mysqlUrl := "huiju:2015-_h_uiJ_uHj@tcp(rds2hhr4tsha6c4p5059.mysql.rds.aliyuncs.com:3306)/"
	PonyDB, _ = xorm.NewEngine("mysql", mysqlUrl+"hjuser?charset=utf8")
	PonyDB.ShowSQL = true
	PonyDB.SetMaxIdleConns(INIT_CLIENT)
	PonyDB.SetMaxOpenConns(MAX_CLIENT)

}
func NoData(b bool) error {
	if b {
		return nil
	}
	return err_code.NotFound
}

//检查是否是空数据
func CheckNil(err error) bool {
	ret := false
	switch t := err.(type) {

	case *err_code.Error:
		if t != nil && t.Code == err_code.NotFound.Code {
			ret = true
		}

	}
	return ret
}

//错误消息定义
func NoDataMsg(b bool, msg string) error {
	if b {
		return nil
	}
	return err_code.RestErr(err_code.NotFound, msg)
}

/***********************************/
