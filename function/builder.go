package function

import (
	"strings"

	"github.com/berquerant/dql/arithmetic"
	"github.com/berquerant/dql/cast"
	"github.com/berquerant/dql/compare"
	"github.com/berquerant/gogrep"
)

type (
	Factory func() Function

	FactoryBuilder interface {
		Factory(name string) (Factory, bool)
	}

	factoryBuidler struct {
		caster        cast.Caster
		artCalculator arithmetic.Calculator
		comparer      compare.Comparer
		grepper       gogrep.Grepper
	}
)

func NewFactoryBuilder(
	caster cast.Caster,
	artCalculator arithmetic.Calculator,
	comparer compare.Comparer,
	grepper gogrep.Grepper,
) FactoryBuilder {
	return &factoryBuidler{
		caster:        caster,
		artCalculator: artCalculator,
		comparer:      comparer,
		grepper:       grepper,
	}
}

func (s *factoryBuidler) Factory(name string) (Factory, bool) {
	switch strings.ToLower(name) {
	case "now":
		return NewNow, true
	case "cast":
		return func() Function { return NewCast(s.caster) }, true
	case "int2bin":
		return NewInt2Bin, true
	case "bin2int":
		return NewBin2Int, true
	case "ext":
		return NewExt, true
	case "dir":
		return NewDir, true
	case "base":
		return NewBase, true
	case "len":
		return NewLen, true
	case "floor":
		return NewFloor, true
	case "ceil":
		return NewCeil, true
	case "depth":
		return NewDepth, true
	case "grep":
		return func() Function { return NewGrep(s.grepper) }, true
	case "pow":
		return func() Function { return NewPow(s.artCalculator) }, true
	case "count":
		return func() Function { return NewCount() }, true
	case "min":
		return func() Function { return NewMin(s.comparer) }, true
	case "max":
		return func() Function { return NewMax(s.comparer) }, true
	case "product":
		return func() Function { return NewProduct(s.artCalculator) }, true
	case "sum":
		return func() Function { return NewSum(s.artCalculator) }, true
	case "avg":
		return func() Function { return NewAvg(s.artCalculator, NewSum(s.artCalculator)) }, true
	default:
		return nil, false
	}
}
