package env

type CtxKey string

const (
	Tags       CtxKey = "tags"
	Service    CtxKey = "service"
	Version    CtxKey = "version"
	Release    CtxKey = "release"
	Deployment CtxKey = "deployment"
	Dataloader CtxKey = "dataloader"
)
