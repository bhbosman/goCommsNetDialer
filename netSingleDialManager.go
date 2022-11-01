package goCommsNetDialer

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"net"
	"net/url"
	"time"
)

type netSingleDialManager struct {
	netDialManager
	MaxConnection int
}

func (self *netSingleDialManager) Start(_ context.Context) error {
	var dm iDialManager = &net.Dialer{
		Timeout: time.Second * 30,
	}
	if self.UseProxy {
		var err error
		dm, err = self.sock5(dm)
		if err != nil {
			return err
		}
	}
	return self.GoFunctionCounter.GoRun(
		"netSingleDialManager.start.Dial",
		func() {
			sem := semaphore.NewWeighted(int64(self.MaxConnection))
			releaseFunc := func() {
				sem.Release(1)
			}
		loop:
			for self.CancelCtx.Err() == nil {
				if sem.Acquire(self.CancelCtx, 1) != nil {
					break loop
				}
				instanceApp, cancellationContext, connctionId, err := self.dialll(dm, releaseFunc)
				if err != nil {
					self.ZapLogger.Error("Error on dial", zap.Error(err))
					if opErr, ok := err.(*net.OpError); ok {
						if dnsError, ok := opErr.Err.(*net.DNSError); ok {
							self.ZapLogger.Error("Further information", zap.Error(dnsError))
						}
					}
					time.Sleep(5 * time.Second)
					continue
				}

				err = instanceApp.Start(context.Background())
				if err != nil {
					self.ZapLogger.Error("ddddd", zap.Error(err))
				}
				// This is the cancellation function that must be called that will call the fxApp.Stop()
				// this can be triggered from two places
				// 1. from the connection itself, by calling the CancelFunc on the connecion stack
				// 2. from the service manager that can shut the whole dialing connection down
				// For all two of this instances, we have to register the same method and make sure it is only executed once

				cc := []goCommsDefinitions.ICancellationContext{
					self.CancellationContext,
					cancellationContext,
				}
				cancelFunction := func(connectionId string, CancellationContext ...goCommsDefinitions.ICancellationContext) func() {
					b := false
					return func() {
						if !b {
							b = true
							stopErr := instanceApp.Stop(context.Background())
							if stopErr != nil {
								self.ZapLogger.Error(
									"Stopping error. not really a problem. informational",
									zap.Error(stopErr))
							}
							for _, instance := range CancellationContext {
								_ = instance.Remove(connectionId)
							}
						}
					}
				}(connctionId, cc...)

				for _, c := range cc {
					added, err := c.Add(connctionId, cancelFunction)
					if !added {
						self.ZapLogger.Error("could not be added")
					}
					if err != nil {
						self.ZapLogger.Error("ddddd", zap.Error(err))
					}
				}
			}

			self.ZapLogger.Info("Exit loop")
		},
	)
}

func (self *netSingleDialManager) Stop(_ context.Context) error {
	return nil
}

func newSingleNetDialManager(
	params struct {
		fx.In
		UseProxy                                 bool     `name:"UseProxy"`
		ConnectionUrl                            *url.URL `name:"ConnectionUrl"`
		ProxyUrl                                 *url.URL `name:"ProxyUrl"`
		ConnectionManager                        goConnectionManager.IService
		CancelCtx                                context.Context
		Options                                  *DialAppSettings
		ZapLogger                                *zap.Logger
		UniqueSessionNumber                      interfaces.IUniqueReferenceService
		ConnectionName                           string `name:"ConnectionName"`
		ConnectionInstancePrefix                 string `name:"ConnectionInstancePrefix"`
		AdditionalFxOptionsForConnectionInstance func() fx.Option
		GoFunctionCounter                        GoFunctionCounter.IService
		CancellationContext                      goCommsDefinitions.ICancellationContext
	}) (*netSingleDialManager, error) {

	if params.ConnectionManager.State() != IFxService.Started {
		return nil, IFxService.NewServiceStateError(
			params.ConnectionManager.ServiceName(),
			"Service in incorrect state",
			IFxService.Started,
			params.ConnectionManager.State())
	}

	netDialManagerInstance, err := newNetDialManager(
		params.ConnectionName,
		params.ConnectionInstancePrefix,
		params.UseProxy,
		params.ProxyUrl,
		params.ConnectionUrl,
		params.CancelCtx,
		params.CancellationContext,
		params.ConnectionManager,
		params.ZapLogger,
		params.UniqueSessionNumber,
		params.AdditionalFxOptionsForConnectionInstance,
		params.GoFunctionCounter,
	)
	if err != nil {
		return nil, err
	}

	return &netSingleDialManager{
		netDialManager: netDialManagerInstance,
		MaxConnection:  params.Options.MaxConnections,
	}, nil
}
