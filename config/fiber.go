package configuration

import "time"

type FiberConfiguration struct {
	Prefork bool
	ServerHeader string
	StrictRouting bool
	CaseSensitive bool
	Immutable bool
	BodyLimit int
	Concurrency int
	DisableKeepalive bool
	DisableDefaultDate bool
	DisableDefaultContentType bool
	DisableStartupMessage bool
	ETag bool
	TemplateEngineName string
	TemplateEngine func(raw string, bind interface{}) (string, error)
	TemplateFolder string
	TemplateExtension string
	ReadTimeoutSeconds int
	WriteTimeoutSeconds int
	IdleTimeoutSeconds int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
}
