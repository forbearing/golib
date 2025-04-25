package middleware

import (
	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

var (
	cb           *gobreaker.CircuitBreaker
	RouteManager *routeParamsManager
)

func Init() (err error) {
	// Init circuit breaker
	cbCfg := config.App.Server.CircuitBreaker
	if cbCfg.MaxRequests == 0 {
		return errors.New("circuit breaker max_requests cannot be 0")
	}
	if cbCfg.MinRequests == 0 {
		return errors.New("circuit breaker min_requests cannot be 0")
	}
	if cbCfg.FailureRate <= 0 || cbCfg.FailureRate > 1 {
		return errors.New("circuit breaker failure_rate must be between 0 and 1")
	}

	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        cbCfg.Name,
		MaxRequests: cbCfg.MaxRequests,
		Interval:    cbCfg.Interval,
		Timeout:     cbCfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests < cbCfg.MinRequests {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= cbCfg.FailureRate
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			zap.S().Infow("circuit breaker state changed",
				"name", name,
				"from", from.String(),
				"to", to.String(),
			)
		},
	})
	zap.S().Infow("circuit breaker initialized",
		"name", cbCfg.Name,
		"max_requests", cbCfg.MaxRequests,
		"min_requests", cbCfg.MinRequests,
		"failure_rate", cbCfg.FailureRate,
		"interval", cbCfg.Interval,
		"timeout", cbCfg.Timeout,
	)

	// Init route params manager
	RouteManager = NewRouteParamsManager()

	return nil
}
