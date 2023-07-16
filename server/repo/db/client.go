// Package db 主要是使用mysql作为存储，提供增删改查接口
package db

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuanzhi-ai/luban/server/log"
)

const (
	maxConns     = 1000
	maxIdleConns = 1000
)

var dbClient *sql.DB

func init() {
	path := fmt.Sprintf("%v:%v@tcp(%v:%v)/db_action_4", os.Getenv("WANXIANG_DB_USER"),
		os.Getenv("WANXIANG_DB_PSWD"), os.Getenv("WANXIANG_DB_IP"), os.Getenv("WANXIANG_DB_PORT"))
	dbClient, err := sql.Open("mysql", path)
	if err != nil {
		log.Errorf("sql.Open err: %+v", err)
		panic("init sql err")
	}
	//设置数据库最大连接数
	dbClient.SetConnMaxLifetime(maxConns)
	//设置上数据库最大闲置连接数
	dbClient.SetMaxIdleConns(maxIdleConns)
	//验证连接
	if err := dbClient.Ping(); err != nil {
		log.Errorf("open database fail, path: %+v", path)
		panic("conn db err")
	}
}

// Query 请求存储获取结果
func Query(str string, rowStruct interface{}, args ...interface{}) ([]interface{}, error) {
	typ := reflect.TypeOf(rowStruct).Elem()
	var fieldNames []string
	for i := 0; i < typ.NumField(); i++ {
		if name, ok := typ.Field(i).Tag.Lookup("sql"); ok {
			fieldNames = append(fieldNames, name)
		}
	}
	if len(fieldNames) == 0 {
		return nil, fmt.Errorf("field empty")
	}
	str = strings.Replace(str, " * ", " "+strings.Join(fieldNames, ",")+" ", 1)

	rows, err := dbClient.Query(str, args...)
	if err != nil {
		return nil, err
	}

	var datas []interface{}
	for rows.Next() {
		row := reflect.New(typ)
		val := row.Elem()

		fieldValues := parseFieldVal(typ, val)
		err := rows.Scan(fieldValues...)
		if err != nil {
			return nil, err
		}
		datas = append(datas, row.Interface())
	}
	return datas, nil
}

// Exec 执行增删改
func Exec(str string, args ...interface{}) (rowsAffected int64, err error) {
	res, err := dbClient.Exec(str, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func parseFieldVal(elemType reflect.Type,
	elemVal reflect.Value) (fieldValues []interface{}) {
	fieldValues = make([]interface{}, 0)

	for i := 0; i < elemType.NumField(); i++ {
		field := elemVal.Field(i)
		var fieldVal interface{}
		if field.Kind() != reflect.Ptr && field.CanAddr() {
			fieldVal = field.Addr().Interface()
		} else {
			fieldVal = field.Interface()
		}
		fieldValues = append(fieldValues, fieldVal)
	}

	return
}
