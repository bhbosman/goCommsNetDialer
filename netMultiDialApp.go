//go:build exclude

package goCommsNetDialer

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"net/url"
	"time"
)

func NewMultiNetDialApp(
	name string,
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	connectionInstancePrefix string,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	options ...common.INetManagerSettingsApply) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			ServiceId:         serviceIdentifier,
			ServiceDependency: serviceDependentOn,
			Name:              name,
			Callback: func() (messages.IApp, context.CancelFunc, error) {
				cancelFunc := func() {}
				dialSettings := &dialAppSettings{
					NetManagerSettings: common.NewNetManagerSettings(1),
					userContext:        nil,
					canDial:            nil,
				}

				for _, option := range options {
					if option == nil {
						continue
					}
					if dialAppSettingsApply, ok := option.(iNetDialAppSettingsApply); ok {
						err := dialAppSettingsApply.apply(dialSettings)
						if err != nil {
							return nil, cancelFunc, err
						}
					} else {
						err := option.ApplyNetManagerSettings(&dialSettings.NetManagerSettings)
						if err != nil {
							return nil, cancelFunc, err
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
					connectionInstancePrefix,
					params,
					callbackForConnectionInstance,
					fx.Options(dialSettings.MoreOptions...),
					goCommsDefinitions.ProvideUrl("ConnectionUrl", ConnectionUrl),
					goCommsDefinitions.ProvideUrl("ProxyUrl", ProxyUrl),
					goCommsDefinitions.ProvideBool("UseProxy", UseProxy),
					fx.Provide(fx.Annotated{Target: newMultiNetDialManager}),
					fx.Provide(fx.Annotated{
						Target: func() *dialAppSettings {
							return dialSettings
						},
					}),
					fx.Invoke(
						func(
							params struct {
								fx.In
								NetManager     *NetMultiDialManager
								CancelFunction context.CancelFunc
								Lifecycle      fx.Lifecycle
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
				return fxApp, cancelFunc, fxApp.Err()
			},
		}
	}
}
