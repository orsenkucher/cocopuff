package env

import "context"

type ctxKey int

const (
	tagsCtx ctxKey = iota
	serviceCtx
	versionCtx
	releaseCtx
	deploymentCtx
)

// With environment
func With(ctx context.Context, service, deployment, version string, release bool) context.Context {
	ctx = context.WithValue(ctx, serviceCtx, service)
	ctx = context.WithValue(ctx, deploymentCtx, deployment)
	ctx = context.WithValue(ctx, versionCtx, version)
	ctx = context.WithValue(ctx, releaseCtx, release)
	ctx = context.WithValue(ctx, tagsCtx, []string{deployment, version})
	return ctx
}
