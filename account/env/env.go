package env

type Key string

const (
	Tags       Key = "tags"
	Service    Key = "service"
	Version    Key = "version"
	Release    Key = "release"
	Deployment Key = "deployment"
)
