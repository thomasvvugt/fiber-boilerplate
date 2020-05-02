package configuration

type HelmetConfiguration struct {
	Enabled bool
	XSSProtection string
	ContentTypeNosniff string
	XFrameOptions string
	HSTSMaxAge int
	HSTSExcludeSubdomains bool
	HSTSPreloadEnabled bool
	ContentSecurityPolicy string
	CSPReportOnly bool
	ReferrerPolicy string
}
