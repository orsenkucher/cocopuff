package env

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const (
	sugarCtx ctxKey = iota
	tagsCtx
	serviceCtx
	versionCtx
	releaseCtx
	deploymentCtx
)

// With environment
func With(ctx context.Context, sugar *zap.SugaredLogger, service, deployment, version string, release bool) context.Context {
	sugar = sugar.With(zap.String("package", "env"))
	ctx = context.WithValue(ctx, sugarCtx, sugar)
	ctx = context.WithValue(ctx, serviceCtx, service)
	ctx = context.WithValue(ctx, versionCtx, version)
	ctx = context.WithValue(ctx, releaseCtx, release)
	ctx = context.WithValue(ctx, deploymentCtx, deployment)
	ctx = context.WithValue(ctx, tagsCtx, []string{deployment, version})
	return ctx
}

func TagsFor(ctx context.Context) []string {
	if tags, ok := ctx.Value(tagsCtx).([]string); ok {
		return tags
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanic("fail to retrieve tags", zap.String("function", "TagsFor"))
	}

	return nil
}

func ServiceFor(ctx context.Context) string {
	if service, ok := ctx.Value(tagsCtx).(string); ok {
		return service
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanic("fail to retrieve service", zap.String("function", "ServiceFor"))
	}

	return ""
}

func DeploymentFor(ctx context.Context) string {
	if deployment, ok := ctx.Value(deploymentCtx).(string); ok {
		return deployment
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanic("fail to retrieve deployment", zap.String("function", "DeploymentFor"))
	}

	return ""
}

func VersionFor(ctx context.Context) string {
	if version, ok := ctx.Value(versionCtx).(string); ok {
		return version
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanic("fail to retrieve version", zap.String("function", "VersionFor"))
	}

	return ""
}

func ReleaseFor(ctx context.Context) bool {
	if release, ok := ctx.Value(versionCtx).(bool); ok {
		return release
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanic("fail to retrieve release", zap.String("function", "ReleaseFor"))
	}

	return true
}
