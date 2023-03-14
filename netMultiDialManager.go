package goCommsNetDialer

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/url"
	"time"
)

type NetMultiDialManager struct {
	netDialManager
}

func (self NetMultiDialManager) Dial(releaseFunc func()) (messages.IApp, common.ICancellationContext, string, error) {
	var dm iDialManager = &net.Dialer{
		Timeout: time.Second * 30,
	}
	if self.UseProxy {
		var err error
		dm, err = self.sock5(dm)
		if err != nil {
			return nil, nil, "", err
		}
	}
	return self.dialll(dm, releaseFunc)
}

func NewMultiNetDialManager(
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	ConnectionManager goConnectionManager.IService,
	CancelCtx context.Context,
	CancellationContext common.ICancellationContext,
	ZapLogger *zap.Logger,
	UniqueSessionNumber interfaces.IUniqueReferenceService,
	ConnectionName string,
	ConnectionInstancePrefix string,
	AdditionalFxOptionsForConnectionInstance func() fx.Option,
	GoFunctionCounter GoFunctionCounter.IService,
) (NetMultiDialManager, error) {

	if ConnectionManager.State() != IFxService.Started {
		return NetMultiDialManager{}, IFxService.NewServiceStateError(
			ConnectionManager.ServiceName(),
			"Service in incorrect state",
			IFxService.Started,
			ConnectionManager.State())
	}

	netDialManagerInstance, err := newNetDialManager(
		ConnectionName,
		ConnectionInstancePrefix,
		UseProxy,
		ProxyUrl,
		ConnectionUrl,
		CancelCtx,
		CancellationContext,
		ConnectionManager,
		ZapLogger,
		UniqueSessionNumber,
		AdditionalFxOptionsForConnectionInstance,
		GoFunctionCounter,
	)
	if err != nil {
		return NetMultiDialManager{}, err
	}

	return NetMultiDialManager{
		netDialManager: netDialManagerInstance,
	}, nil
}
