// Package comm 公共库，
// 这里存放错误码
package comm

const (
	/*通用错误码*/
	SuccessCode int32 = 0
	InputErr    int32 = 1007 // 入参数错误
	/*chat svr 错误码
	2000 ~ 2999
	*/
	RemoteGetAllNpcInfoErr int32 = 2000  // 远程rpc获取npc信息错误
	UnknownNPC             int32 = 20001 // 未知npc
	SessionExist           int32 = 20002 // 会话已存在
	SessionNotExist        int32 = 20003 // 会话不存在
)
