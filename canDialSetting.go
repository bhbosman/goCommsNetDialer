package goCommsNetDialer

import (
	"github.com/bhbosman/gocomms/common"
)

type canDialSetting struct {
	canDial []ICanDial
}

func (self canDialSetting) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func CanDial(canDial ...ICanDial) *canDialSetting {
	return &canDialSetting{canDial: canDial}
}

func (self canDialSetting) apply(settings *DialAppSettings) error {
	err := self.ApplyNetManagerSettings(&settings.NetManagerSettings)
	if err != nil {
		return err
	}
	for _, cd := range self.canDial {
		settings.canDial = append(settings.canDial, cd)
	}
	return nil
}
