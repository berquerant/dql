package env

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/data"
)

type Type int

//go:generate stringer -type Type -output env_stringer_generated.go

const (
	TypeUnknown Type = iota
	// TypeData means that Data contains actual data.
	TypeData
	// TypeExpr means that Data contains AST.
	TypeExpr
	// TypeDataList means that Data contains actual multiple data.
	TypeDataList
)

func (s Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

type (
	Data interface {
		Type() Type
		Data() data.Data
		Expr() ast.Expr
		DataList() []data.Data
		Value() interface{}
	}

	dataImpl struct {
		typ   Type
		value interface{}
	}
)

func (s *dataImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  s.typ,
		"value": s.value,
	})
}

func FromData(v data.Data) Data {
	return &dataImpl{
		typ:   TypeData,
		value: v,
	}
}

func FromDataList(v []data.Data) Data {
	return &dataImpl{
		typ:   TypeDataList,
		value: v,
	}
}

func FromExpr(v ast.Expr) Data {
	return &dataImpl{
		typ:   TypeExpr,
		value: v,
	}
}

func (s *dataImpl) Type() Type         { return s.typ }
func (s *dataImpl) Value() interface{} { return s.value }
func (s *dataImpl) Data() data.Data {
	if x, ok := s.value.(data.Data); ok {
		return x
	}
	return nil
}
func (s *dataImpl) Expr() ast.Expr {
	if x, ok := s.value.(ast.Expr); ok {
		return x
	}
	return nil
}
func (s *dataImpl) DataList() []data.Data {
	if x, ok := s.value.([]data.Data); ok {
		return x
	}
	return nil
}

type (
	Map interface {
		Get(key string) (Data, bool)
		Set(key string, value Data)
		Ref(key, ref string)
		Clone() Map
	}

	mapImpl struct {
		table map[string]Data
		refs  map[string]string
		mux   sync.RWMutex
	}
)

func New() Map {
	return &mapImpl{
		table: map[string]Data{},
		refs:  map[string]string{},
	}
}

func (s *mapImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"table": s.table,
		"refs":  s.refs,
	})
}

func (s *mapImpl) Clone() Map {
	s.mux.RLock()
	defer s.mux.RUnlock()
	table := make(map[string]Data, len(s.table))
	for k, v := range s.table {
		table[k] = v // shallow copy
	}
	refs := make(map[string]string, len(s.refs))
	for k, v := range s.refs {
		refs[k] = v
	}
	return &mapImpl{
		table: table,
		refs:  refs,
	}
}

func (s *mapImpl) Get(key string) (Data, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if ref, ok := s.refs[key]; ok {
		return s.get(ref)
	}
	return s.get(key)
}

func (s *mapImpl) get(key string) (Data, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	v, ok := s.table[key]
	return v, ok
}

func (s *mapImpl) Set(key string, value Data) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if ref, ok := s.refs[key]; ok {
		s.table[ref] = value
		return
	}
	s.table[key] = value
}

func (s *mapImpl) Ref(key, ref string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.refs[key] = ref
}
