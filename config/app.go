package configuration

type AppConfiguration struct {
	ListenAddress interface{}
	SuppressWWW bool
	ForceHTTPS bool
}
