// Package db 主要是使用mysql作为存储，提供增删改查接口
package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuanzhi-ai/luban/server/log"
)

const (
	maxConns     = 1000
	maxIdleConns = 1000
)

// Client mysql存储对象
type Client struct {
	db *sql.DB
}

// NewClient 新建client
func NewClient(userName string, passwd string, ip string, port string, database string) (*Client, error) {
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, passwd, ip, port, database)
	db, err := sql.Open("mysql", path)
	if err != nil {
		log.Errorf("sql.Open err: %+v", err)
		return nil, err
	}
	//设置数据库最大连接数
	db.SetConnMaxLifetime(maxConns)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(maxIdleConns)
	//验证连接
	if err := db.Ping(); err != nil {
		log.Errorf("open database fail, path: %+v", path)
		return nil, err
	}
	return &Client{db: db}, nil
}

// Query 请求存储获取结果
func (d *Client) Query(str string, rowStruct interface{}, args ...interface{}) ([]interface{}, error) {
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

	rows, err := d.db.Query(str, args...)
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
func (d *Client) Exec(str string, args ...interface{}) (rowsAffected int64, err error) {
	res, err := d.db.Exec(str, args...)
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
