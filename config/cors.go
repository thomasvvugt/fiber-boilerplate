package configuration

type CORSConfiguration struct {
	Enabled bool
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
	AllowCredentials bool
	ExposeHeaders string
	MaxAge int
}
