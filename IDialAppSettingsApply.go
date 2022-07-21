package goCommsNetDialer

import "github.com/bhbosman/gocomms/common"

type INetDialAppSettingsApply interface {
	common.INetManagerSettingsApply
	apply(settings *DialAppSettings) error
}
