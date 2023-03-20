package goCommsNetDialer

import (
	"context"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"time"
)

func NewSingleNetDialApp(
	name string,
	options ...common.INetManagerSettingsApply) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			Name: name,
			Callback: func() (messages.IApp, goConn.ICancellationContext, error) {
				dialSettings := &DialAppSettings{
					NetManagerSettings: common.NewNetManagerSettings(1),
					canDial:            nil,
				}
				namedLogger := params.ZapLogger.Named(name)
				ctx, cancelFunc := context.WithCancel(params.ParentContext)
				cancellationContext, err := goConn.NewCancellationContext(name, cancelFunc, ctx, namedLogger, nil)
				if err != nil {
					return nil, nil, err
				}
				for _, option := range options {
					if option == nil {
						continue
					}
					if dialAppSettingsApply, ok := option.(INetDialAppSettingsApply); ok {
						err := dialAppSettingsApply.apply(dialSettings)
						if err != nil {
							return nil, cancellationContext, err
						}
					} else {
						err := option.ApplyNetManagerSettings(&dialSettings.NetManagerSettings)
						if err != nil {
							return nil, cancellationContext, err
						}
					}
				}

				callbackForConnectionInstance, err := dialSettings.Build()
				if err != nil {
					return nil, nil, err
				}

				connectionOptions := common.ConnectionApp(
					time.Hour,
					time.Hour,
					name,
					name,
					params,
					cancellationContext,
					namedLogger,
					callbackForConnectionInstance,
					fx.Options(dialSettings.MoreOptions...),
					fx.Provide(fx.Annotated{Target: newSingleNetDialManager}),
					fx.Provide(
						fx.Annotated{
							Target: func() *DialAppSettings {
								return dialSettings
							},
						},
					),
					common.InvokeCancelContext(),
					fx.Invoke(
						func(
							params struct {
								fx.In
								NetManager          *netSingleDialManager
								CancelFunction      context.CancelFunc
								CancellationContext goConn.ICancellationContext
								Lifecycle           fx.Lifecycle
							},
						) {
							hook := fx.Hook{
								OnStart: func(ctx context.Context) error {
									return params.NetManager.Start(ctx)
								},
								OnStop: func(ctx context.Context) error {
									params.CancelFunction()
									return params.NetManager.Stop(ctx)
								},
							}
							params.Lifecycle.Append(hook)
						},
					),
				)
				fxApp := fx.New(connectionOptions)
				return fxApp, cancellationContext, fxApp.Err()
			},
		}
	}
}
