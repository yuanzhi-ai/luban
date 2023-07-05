// Package comm 公共库，
// 这里存放错误码
package comm

const (
	/*通用错误码*/
	SuccessCode int32 = 0
	InputErr    int32 = 1007 // 入参数错误
	/*verify svr 错误码
	2000 ~ 2999
	*/
	GeneratorCaptchaErr int32 = 2000  // 生成验证码失败
	VerifyCaptchaErr    int32 = 20001 // 检验验证码失败
)
