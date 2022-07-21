package goCommsNetDialer

import "github.com/bhbosman/gocomms/common"

type DialAppSettings struct {
	common.NetManagerSettings
	//userContext interface{}
	canDial []ICanDial
}
