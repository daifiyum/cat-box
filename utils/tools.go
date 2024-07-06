package utils

import (
	"github.com/sagernet/sing/common/wininet"
	"golang.org/x/sys/windows"
)

// 清除系统代理
func DisableProxy() error {
	err := wininet.ClearSystemProxy()
	if err != nil {
		return err
	}
	return nil
}

// 检查权限
func IsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	isMember, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return isMember
}
