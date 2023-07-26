// Package comm 公共库，
// 这里存放错误码
package comm

const (
	/*通用状态码*/
	SuccessCode  int32 = 666  // 成功
	InputErr     int32 = 1007 // 入参数错误
	NetWorkErr   int32 = 1008 // 网路错误
	RedisErr     int32 = 1009 // redis 失败
	SmsErr       int32 = 1010 // 发送短信消息失败
	JWTErr       int32 = 1011 // jwt验证失败
	FrequencyErr int32 = 1012 // 接口限频
	/*verify svr 错误码
	2000 ~ 2999
	*/
	GeneratorCaptchaErr int32 = 2000  // 生成验证码失败
	VerifyCaptchaErr    int32 = 20001 // 检验验证码失败
	CaptchaUsed         int32 = 20002 // 验证码已经用过了

	/*login svr 错误码
	3000 ~ 3999
	*/
	GetVerCodeErr             int32 = 3000 // 获取人机图形验证码失败
	VerifyCodeErr             int32 = 3001 // 验证人机图形验证码失败
	CapchaCodeWrong           int32 = 3002 // 验证码错误
	QueryPhoneCodeErr         int32 = 3003 // 查询手机验证码错误
	CodeExpiration            int32 = 3004 // 验证码过期
	PhoneCodeWrong            int32 = 3005 // 验证码错误
	GeneratorUserInfoErr      int32 = 3006 // 生成用户注册信息错误
	InsertUserRegisterInfoErr int32 = 3007 // 写入用户注册信息错误
	QueryS2Err                int32 = 3008 //查询s2失败
	DercyA1Err                int32 = 3009 // 解码a1失败
	HexDBS2Err                int32 = 3010 //16进制解码db_s2失败
	DecodeA1Err               int32 = 3011 // 解码a1失败
	LoginErr                  int32 = 3012 //登录失败
	PswdLegalErr              int32 = 3013 // 密码不合法
	UpdatePswdErr             int32 = 3014 // 数据库更新密码失败
	SendSmsVerCodeErr         int32 = 3015 //发送短信验证码失败
	RegisterErr               int32 = 3016 //注册失败
	ResetPswdErr              int32 = 3017 // 重置密码失败

)
