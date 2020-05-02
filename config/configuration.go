package configuration

type Configuration struct {
	App AppConfiguration
	Database DatabaseConfiguration
	Fiber FiberConfiguration
	Public PublicConfiguration
	Logger LoggerConfiguration
	Recover RecoverConfiguration
	Compression CompressionConfiguration
	CORS CORSConfiguration
	Helmet HelmetConfiguration
}
