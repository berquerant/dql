package eval

import (
	"encoding/json"

	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/dig"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
)

/* source rows */

type (
	Row interface {
		Err() error
		Info() Info
	}

	Info interface {
		Name() string
		Size() int
		Mode() string
		ModTime() int
		IsDir() bool
		ToMap() map[string]data.Data
	}
)

type (
	row struct {
		info Info
	}
	errRow struct {
		err error
	}
)

func NewRow(v Info) Row       { return &row{info: v} }
func NewErrRow(err error) Row { return &errRow{err: err} }

func (s *row) Info() Info { return s.info }
func (*row) Err() error   { return nil }
func (s *row) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"info": s.info,
	})
}
func (*errRow) Info() Info   { return nil }
func (s *errRow) Err() error { return s.err }
func (s *errRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"err": s.err.Error(),
	})
}

type info struct {
	name    string
	size    int
	mode    string
	modTime int
	isDir   bool
}

func NewInfo(v dig.FileInfo) Info {
	return &info{
		name:    v.Name(),
		size:    int(v.Size()),
		mode:    v.Mode().String(),
		modTime: int(v.ModTime().Unix()),
		isDir:   v.IsDir(),
	}
}

func (s *info) Name() string { return s.name }
func (s *info) Size() int    { return s.size }
func (s *info) Mode() string { return s.mode }
func (s *info) ModTime() int { return s.modTime }
func (s *info) IsDir() bool  { return s.isDir }
func (s *info) ToMap() map[string]data.Data {
	return map[string]data.Data{
		"name":     data.FromString(s.name),
		"size":     data.FromInt(s.size),
		"mode":     data.FromString(s.mode),
		"mod_time": data.FromInt(s.modTime),
		"is_dir":   data.FromBool(s.isDir),
	}
}
func (s *info) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToMap())
}

func AppendRowToEnv(table env.Map, row Row) env.Map {
	x := table.Clone()
	for k, v := range row.Info().ToMap() {
		x.Set(k, env.FromData(v))
	}
	return x
}

/* grouped rows */

type (
	GRow interface {
		Type() RowType
		Raw() Row
		Grouped() GroupedRow
		Err() error
	}

	GroupedRow interface {
		Key() string
		Value() data.Data
		Rows() []Row
	}

	rawGRow struct {
		row Row
	}
	groupedGRow struct {
		grouped GroupedRow
	}
	errGRow struct {
		err error
	}
)

func NewRawGRow(row Row) GRow                { return &rawGRow{row: row} }
func NewGroupedGRow(grouped GroupedRow) GRow { return &groupedGRow{grouped: grouped} }
func NewErrGRow(err error) GRow              { return &errGRow{err: err} }

func (*rawGRow) Type() RowType             { return RawRowType }
func (s *rawGRow) Raw() Row                { return s.row }
func (*rawGRow) Grouped() GroupedRow       { return nil }
func (*rawGRow) Err() error                { return nil }
func (*groupedGRow) Type() RowType         { return GroupedRowType }
func (*groupedGRow) Raw() Row              { return nil }
func (s *groupedGRow) Grouped() GroupedRow { return s.grouped }
func (s *groupedGRow) Err() error          { return nil }
func (*errGRow) Type() RowType             { return ErrRowType }
func (*errGRow) Raw() Row                  { return nil }
func (*errGRow) Grouped() GroupedRow       { return nil }
func (s *errGRow) Err() error              { return s.err }

type groupedRow struct {
	key   string
	value data.Data
	rows  []Row
}

func NewGroupedRow(key string, value data.Data, rows []Row) GroupedRow {
	return &groupedRow{
		key:   key,
		value: value,
		rows:  rows,
	}
}

func (s *groupedRow) Key() string      { return s.key }
func (s *groupedRow) Value() data.Data { return s.value }
func (s *groupedRow) Rows() []Row      { return s.rows }

func AppendGroupedRowToEnv(table env.Map, row GroupedRow) env.Map {
	x := table.Clone()
	x.Set(row.Key(), env.FromData(row.Value()))

	grouped := map[string][]data.Data{}
	for i, r := range row.Rows() {
		d := r.Info().ToMap()
		for k, v := range d {
			if k == row.Key() {
				continue
			}
			if _, ok := grouped[k]; !ok {
				grouped[k] = make([]data.Data, len(row.Rows()))
			}
			grouped[k][i] = v
		}
	}
	for k, v := range grouped {
		x.Set(k, env.FromDataList(v))
	}
	return x
}

func AppendRowsToEnv(table env.Map, rows []Row) env.Map {
	x := table.Clone()
	grouped := map[string][]data.Data{}
	for i, r := range rows {
		for k, v := range r.Info().ToMap() {
			if _, ok := grouped[k]; !ok {
				grouped[k] = make([]data.Data, len(rows))
			}
			grouped[k][i] = v
		}
	}
	for k, v := range grouped {
		x.Set(k, env.FromDataList(v))
	}
	return x
}

func AppendGRowToEnv(table env.Map, row GRow) (env.Map, error) {
	switch row.Type() {
	case RawRowType:
		return AppendRowToEnv(table, row.Raw()), nil
	case GroupedRowType:
		return AppendGroupedRowToEnv(table, row.Grouped()), nil
	default:
		return nil, errors.Wrap(ErrUnknownRowType, row.Type().String())
	}
}

/* selected rows */

type (
	sRow struct {
		values []data.Data
	}
	errSRow struct {
		err error
	}
)

func NewSRow(values []data.Data) SRow { return &sRow{values: values} }
func NewErrSRow(err error) SRow       { return &errSRow{err: err} }

type SRow interface {
	Len() int
	Err() error
	Get(i int) data.Data
}

func (s *sRow) Len() int { return len(s.values) }
func (s *sRow) Get(i int) data.Data {
	if i < 0 || i >= len(s.values) {
		return nil
	}
	return s.values[i]
}
func (*sRow) Err() error { return nil }
func (s *sRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"values": s.values,
	})
}
func (*errSRow) Len() int            { return 0 }
func (*errSRow) Get(_ int) data.Data { return nil }
func (s *errSRow) Err() error        { return s.err }
func (s *errSRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"err": s.err.Error(),
	})
}
