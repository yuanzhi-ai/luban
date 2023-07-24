package auth

const (
	RPCGetVerCode              = "GetVerCode"
	fGetVerCode                = 1
	RPCVerifyCode              = "VerifyCode"
	fVerifyCode                = 1
	PRCGetMachineVerify        = "GetMachineVerify"
	fGetMachineVerify          = 1
	RPCSendMachineVerifyResult = "SendMachineVerifyResult"
	fSendMachineVerifyResult   = 1
	RPCSendTextVerCode         = "SendTextVerCode"
	fSendTextVerCode           = 60
	RPCUserPswdLogin           = "UserPswdLogin"
	fUserPswdLogin             = 1
	RPCUserPhoneLogin          = "UserPhoneLogin"
	fUserPhoneLogin            = 1
	RPCUserRegister            = "UserRegister"
	fUserRegister              = 1
	RPCResetPswd               = "ResetPswd"
	fResetPswd                 = 1
	defaultSize                = 100000
	JwtLoginPurpose            = "login"
	HomePageJwtExpTs           = 2 * 60 * 60 // 主站2小时登录态有效期
	TouristJwtExpTs            = 5 * 60      // 游客态5分钟登录态有效期
)
