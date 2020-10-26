package providers

import "go-fiber-v2-boilerplate/app/configuration"

var appConfig *configuration.Configuration

func SetConfiguration(config *configuration.Configuration)  {
	appConfig = config
}

func GetConfiguration() (config *configuration.Configuration) {
	return appConfig
}
