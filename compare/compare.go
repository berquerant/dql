package compare

import "regexp"

// Comparer provides comparison operations.
type Comparer interface {
	// Compare compares 2 values.
	Compare(left, right interface{}) Result
	// In returns true if target is in list.
	In(target, list interface{}) Result
	// Between returns true if target is in range (left lower to upper).
	Between(target, lower, upper interface{}) Result
	// Like returns true if target is like pattern.
	Like(target, pattern interface{}) Result
}

type Result int

//go:generate stringer -type Result -output compare_stringer_generated.go

const (
	// ResultUndefined means that they are not comparable.
	ResultUndefined Result = iota
	// ResultEqual means that left equals right.
	ResultEqual
	// ResultLessThan means that left is less than right.
	ResultLessThan
	// ResultGreaterThan means that left is greater than right.
	ResultGreaterThan
	// ResultIn means that target is in list or range.
	ResultIn
	// ResultNotIn means that target is not in list or range.
	ResultNotIn
	// ResultMatched means that target is like pattern.
	ResultMatched
	// ResultNotMatched means that target is not like pattern.
	ResultNotMatched
)

// New returns a new Comparer.
func New() Comparer { return &comparer{} }

type comparer struct{}

func (*comparer) Like(target, pattern interface{}) Result {
	t, ok := target.(string)
	if !ok {
		return ResultUndefined
	}
	p, ok := pattern.(string)
	if !ok {
		return ResultUndefined
	}
	matched, err := regexp.MatchString(p, t)
	if err != nil {
		return ResultUndefined
	}
	if matched {
		return ResultMatched
	}
	return ResultNotMatched
}

func (*comparer) Between(target, lower, upper interface{}) Result {
	switch target := target.(type) {
	case int:
		l, ok := lower.(int)
		if !ok {
			return ResultUndefined
		}
		u, ok := upper.(int)
		if !ok {
			return ResultUndefined
		}
		if target >= l && target <= u {
			return ResultIn
		}
		return ResultNotIn
	case float64:
		l, ok := lower.(float64)
		if !ok {
			return ResultUndefined
		}
		u, ok := upper.(float64)
		if !ok {
			return ResultUndefined
		}
		if target >= l && target <= u {
			return ResultIn
		}
		return ResultNotIn
	case string:
		l, ok := lower.(string)
		if !ok {
			return ResultUndefined
		}
		u, ok := upper.(string)
		if !ok {
			return ResultUndefined
		}
		if target >= l && target <= u {
			return ResultIn
		}
		return ResultNotIn
	}
	return ResultUndefined
}

func (*comparer) In(target, list interface{}) Result {
	switch target := target.(type) {
	case bool:
		l, ok := list.([]bool)
		if !ok {
			return ResultUndefined
		}
		for _, x := range l {
			if target == x {
				return ResultIn
			}
		}
		return ResultNotIn
	case int:
		l, ok := list.([]int)
		if !ok {
			return ResultUndefined
		}
		for _, x := range l {
			if target == x {
				return ResultIn
			}
		}
		return ResultNotIn
	case float64:
		l, ok := list.([]float64)
		if !ok {
			return ResultUndefined
		}
		for _, x := range l {
			if target == x {
				return ResultIn
			}
		}
		return ResultNotIn
	case string:
		l, ok := list.([]string)
		if !ok {
			return ResultUndefined
		}
		for _, x := range l {
			if target == x {
				return ResultIn
			}
		}
		return ResultNotIn
	}
	return ResultUndefined
}

func (s *comparer) Compare(left, right interface{}) Result {
	switch left := left.(type) {
	case bool:
		if r, ok := right.(bool); ok {
			return s.compareBool(left, r)
		}
	case int:
		if r, ok := right.(int); ok {
			return s.compareInt(left, r)
		}
	case float64:
		if r, ok := right.(float64); ok {
			return s.compareFloat(left, r)
		}
	case string:
		if r, ok := right.(string); ok {
			return s.compareString(left, r)
		}
	}
	return ResultUndefined
}

func (*comparer) compareBool(left, right bool) Result {
	if !left && right {
		return ResultLessThan
	}
	if left && !right {
		return ResultGreaterThan
	}
	return ResultEqual
}

func (*comparer) compareInt(left, right int) Result {
	if left < right {
		return ResultLessThan
	}
	if left > right {
		return ResultGreaterThan
	}
	return ResultEqual
}

func (*comparer) compareFloat(left, right float64) Result {
	if left < right {
		return ResultLessThan
	}
	if left > right {
		return ResultGreaterThan
	}
	return ResultEqual
}

func (*comparer) compareString(left, right string) Result {
	if left < right {
		return ResultLessThan
	}
	if left > right {
		return ResultGreaterThan
	}
	return ResultEqual
}
