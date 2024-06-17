//go:build wireinject
// +build wireinject

package fetcher

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/mongo"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/google/wire"
)

type mongoConfig resolver.ConfigMap
