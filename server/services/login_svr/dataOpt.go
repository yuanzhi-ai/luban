package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/repo/db"
)

const (
	UserInfoDB    = "db_user_info"
	UserInfoTable = "tb_user_info"
)

var cl *db.Client

func init() {
	var err error
	cl, err = db.NewClient(os.Getenv("WANXIANG_DB_USER"), os.Getenv("WANXIANG_DB_PSWD"), os.Getenv("WANXIANG_DB_IP"), os.Getenv("WANXIANG_DB_PORT"), UserInfoDB)
	if err != nil {
		panic(fmt.Sprintf("new client err: %+v", err))
	}
}

// RegisterInfo 用户注册信息
type UserInfo struct {
	Uid         sql.NullString `sql:"f_uid"`           // 用户的uid 手机号md5
	Phone       sql.NullString `sql:"f_phone"`         // 用户加密后的手机号
	S2          sql.NullString `sql:"f_s2"`            // md5(phone+md5(pswd)) 用于登录校验
	Nickname    sql.NullString `sql:"f_nickname"`      // 用户昵称
	RegisterTs  sql.NullInt64  `sql:"f_register_ts"`   // 用户注册时间 (秒)
	LastLoginTs sql.NullInt64  `sql:"f_last_login_ts"` //用户上次登录时间 ( 秒)
}

// generatorUserRegister 生成用户注册信息数据
func generatorUserRegister(phone string, pswd string) (*UserInfo, error) {
	enctpryPhone, err := comm.AesEncrptyPhone(phone)
	if err != nil {
		return nil, err
	}
	ts := time.Now().Unix()
	userInfo := &UserInfo{
		Uid:         sql.NullString{String: comm.Md5Encode(phone), Valid: true},
		Phone:       sql.NullString{String: enctpryPhone, Valid: true},
		S2:          sql.NullString{String: comm.CalculateS2(phone, pswd), Valid: true},
		Nickname:    sql.NullString{String: phone[len(phone)-4:], Valid: true}, // 昵称默认取手机号后4位
		RegisterTs:  sql.NullInt64{Int64: ts, Valid: true},
		LastLoginTs: sql.NullInt64{Int64: ts, Valid: true},
	}
	return userInfo, nil

}

// 插入用户注册信息
func insterUserRegristerInfo(ctx context.Context, info *UserInfo) error {
	sql := fmt.Sprintf("insert into %s.%s (f_uid,f_phone,f_s2,f_nickname,f_register_ts,f_last_login_ts) value(?,?,?,?,?,?)", UserInfoDB, UserInfoTable)
	rowAffect, err := cl.Exec(sql, info.Uid.String, info.Phone.String, info.S2.String, info.Nickname.String, info.RegisterTs.Int64, info.LastLoginTs.Int64)
	if err != nil || rowAffect != 1 {
		return fmt.Errorf("invalid insert user info err:%v, rowAffect:%v, user info:%+v", err, rowAffect, info)
	}
	return nil
}

type S2 struct {
	S2 sql.NullString `sql:"f_s2"` // md5(phone+md5(pswd)) 用于登录校验
}

// 获取用户的DB_A1信息
// @param uid 用户的uid
func getUserS2(ctx context.Context, phone string) (string, error) {
	uid := comm.Md5Encode(phone)
	sql := fmt.Sprintf("select f_s2 from %v.%v where f_uid= ?", UserInfoDB, UserInfoTable)
	rows, err := cl.Query(sql, (*S2)(nil), uid)
	if err != nil || len(rows) != 1 {
		return "", fmt.Errorf("query phone s2 fail. rows:%+v err:%v ", rows, err)
	}
	res, ok := rows[0].(*S2)
	if !ok {
		return "", fmt.Errorf("query user s2 fail. row:%+v", rows)
	}
	return res.S2.String, nil
}

// updateLoginTs 更新登录时间戳
func updateLoginTs(ctx context.Context, phone string) error {
	uid := comm.Md5Encode(phone)
	loginTs := time.Now().Unix()
	sql := fmt.Sprintf("update %s.%s set f_last_login_ts=? where f_uid=?", UserInfoDB, UserInfoTable)
	rowAffect, err := cl.Exec(sql, loginTs, uid)
	if err != nil || rowAffect != 1 {
		return fmt.Errorf("insert login ts fail. uid:%v err:%v", uid, err)
	}
	return nil
}

// 更新密码
func updatePswd(ctx context.Context, phone string, newPswd string) error {
	uid := comm.Md5Encode(phone)
	dbS2 := comm.CalculateS2(phone, newPswd)
	sql := fmt.Sprintf("update %s.%s set f_s2=? where f_uid=?", UserInfoDB, UserInfoTable)
	rowAffect, err := cl.Exec(sql, dbS2, uid)
	if err != nil || rowAffect != 1 {
		return fmt.Errorf("update user pswd fail. uid:%v err:%v", uid, err)
	}
	return nil
}
