package config

import "github.com/micro/go-config/source"

type Options struct {
	Apps        map[string]interface{}
	AppName     string
	EnableWatch bool
	Sources     []source.Source
}

type Option func(o *Options)

func WithSource(src source.Source) Option {
	return func(o *Options) {
		o.Sources = append(o.Sources, src)
	}
}

func WithApp(appName string) Option {
	return func(o *Options) {
		o.AppName = appName
	}
}
