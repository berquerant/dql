package calc

import (
	"github.com/berquerant/dql/arithmetic"
	"github.com/berquerant/dql/bit"
	"github.com/berquerant/dql/cast"
	"github.com/berquerant/dql/compare"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/function"
	"github.com/berquerant/gogrep"
)

func NewNormal(env env.Map) Calculator {
	return NewWithCaller(env, function.NewCallerWithNames(
		function.NewFactoryBuilder(
			cast.New(),
			arithmetic.New(),
			compare.New(),
			gogrep.New(),
		), function.NormalFunctionNames()...,
	))
}

func NewAggregation(env env.Map) Calculator {
	functionNames := append(function.NormalFunctionNames(), function.AggregationFunctionNames()...)
	return NewWithCaller(env, function.NewCallerWithNames(
		function.NewFactoryBuilder(
			cast.New(),
			arithmetic.New(),
			compare.New(),
			gogrep.New(),
		), functionNames...,
	))
}

func NewWithCaller(env env.Map, caller function.Caller) Calculator {
	return New(
		compare.New(),
		arithmetic.New(),
		bit.New(),
		caller,
		env,
	)
}
