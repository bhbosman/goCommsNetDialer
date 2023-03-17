package goCommsNetDialer

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/netBase"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"net/url"
)

type netDialManager struct {
	netBase.ConnNetManager
}

func (self *netDialManager) sock5(dialManager iDialManager) (iDialManager, error) {
	return proxy.SOCKS5(self.ProxyUrl.Scheme, self.ProxyUrl.Host, nil, dialManager)
}

func (self *netDialManager) dialll(dm iDialManager, releaseFunc func()) (messages.IApp, common.ICancellationContext, string, error) {
	conn, err := dm.Dial("tcp4", self.ConnectionUrl.Host)
	if err != nil {
		if releaseFunc != nil {
			releaseFunc()
		}
		self.ZapLogger.Error(
			"Connection failed due to",
			zap.Error(err))
		return nil, nil, "", err
	}

	temp := conn
	conn, err = common.NewNetConnWithSemaphoreWrapper(
		conn,
		releaseFunc,
	)
	if err != nil {
		self.ZapLogger.Error(
			"error in creating connection with semaphore",
			zap.Error(err))
		_ = temp.Close()
		if releaseFunc != nil {
			releaseFunc()
		}
		return nil, nil, "", err
	}

	connectionInstance := netBase.NewConnectionInstance(
		self.ConnectionUrl,
		self.UniqueSessionNumber,
		self.ConnectionManager,
		self.CancelCtx,
		self.AdditionalFxOptionsForConnectionInstance,
		self.ZapLogger,
	)
	connectionId := self.UniqueSessionNumber.Next(self.ConnectionInstancePrefix)
	instanceApp, instanceAppCtx, cancellationContext, err := connectionInstance.NewConnectionInstance(
		connectionId,
		self.GoFunctionCounter,
		model.ClientConnection,
		conn,
	)
	if instanceAppCtx != nil {
		err = multierr.Append(err, instanceAppCtx.Err())
	}
	onErr := func() {
		if cancellationContext != nil {
			cancellationContext.Cancel("asdasdas")
		}
	}
	if err != nil {
		onErr()
		return nil, nil, "", err
	}
	return instanceApp, cancellationContext, connectionId, nil
}

func newNetDialManager(
	name string,
	connectionInstancePrefix string,
	useProxy bool,
	proxyUrl *url.URL,
	connectionUrl *url.URL,
	cancelCtx context.Context,
	CancellationContext common.ICancellationContext,
	connectionManager goConnectionManager.IService,
	ZapLogger *zap.Logger,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	GoFunctionCounter GoFunctionCounter.IService,
) (netDialManager, error) {
	newConnNetManager, err := netBase.NewConnNetManager(
		name,
		connectionInstancePrefix,
		useProxy,
		proxyUrl,
		connectionUrl,
		cancelCtx,
		CancellationContext,
		connectionManager,
		ZapLogger,
		uniqueSessionNumber,
		additionalFxOptionsForConnectionInstance,
		GoFunctionCounter,
	)
	if err != nil {
		return netDialManager{}, err
	}

	result := netDialManager{
		ConnNetManager: newConnNetManager,
	}
	return result, nil
}
