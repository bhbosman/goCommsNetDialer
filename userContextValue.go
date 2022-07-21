package goCommsNetDialer

import "github.com/bhbosman/gocomms/common"

type userContextValue struct {
	userContext interface{}
}

func (self userContextValue) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self userContextValue) apply(settings *DialAppSettings) error {
	err := self.ApplyNetManagerSettings(&settings.NetManagerSettings)
	if err != nil {
		return err
	}
	//settings.userContext = self.userContext
	return nil
}

func UserContextValue(userContext interface{}) *userContextValue {
	return &userContextValue{userContext: userContext}
}
