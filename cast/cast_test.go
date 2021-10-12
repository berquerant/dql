package cast_test

import (
	"testing"

	"github.com/berquerant/dql/cast"
	"github.com/berquerant/dql/data"
	"github.com/stretchr/testify/assert"
)

func TestCaster(t *testing.T) {
	for _, tc := range []*struct {
		title string
		input data.Data
		to    cast.Type
		want  data.Data
		isErr bool
	}{
		{
			title: "int to int",
			input: data.FromInt(1),
			to:    cast.TypeInt,
			want:  data.FromInt(1),
		},
		{
			title: "float to int",
			input: data.FromFloat(1.2),
			to:    cast.TypeInt,
			want:  data.FromInt(1),
		},
		{
			title: "string to int",
			input: data.FromString("10"),
			to:    cast.TypeInt,
			want:  data.FromInt(10),
		},
		{
			title: "string to int error",
			input: data.FromString("one"),
			to:    cast.TypeInt,
			isErr: true,
		},
		{
			title: "bool true to int",
			input: data.FromBool(true),
			to:    cast.TypeInt,
			want:  data.FromInt(1),
		},
		{
			title: "bool false to int",
			input: data.FromBool(false),
			to:    cast.TypeInt,
			want:  data.FromInt(0),
		},
		{
			title: "int to float",
			input: data.FromInt(1),
			to:    cast.TypeFloat,
			want:  data.FromFloat(1),
		},
		{
			title: "float to float",
			input: data.FromFloat(1),
			to:    cast.TypeFloat,
			want:  data.FromFloat(1),
		},
		{
			title: "string to float",
			input: data.FromString("1.2"),
			to:    cast.TypeFloat,
			want:  data.FromFloat(1.2),
		},
		{
			title: "string to float error",
			input: data.FromString("one point two"),
			to:    cast.TypeFloat,
			isErr: true,
		},
		{
			title: "bool true to float",
			input: data.FromBool(true),
			to:    cast.TypeFloat,
			want:  data.FromFloat(1),
		},
		{
			title: "bool false to float",
			input: data.FromBool(false),
			to:    cast.TypeFloat,
			want:  data.FromFloat(0),
		},
		{
			title: "int to string",
			input: data.FromInt(1),
			to:    cast.TypeString,
			want:  data.FromString("1"),
		},
		{
			title: "float to string",
			input: data.FromFloat(1.2),
			to:    cast.TypeString,
			want:  data.FromString("1.2"),
		},
		{
			title: "string to string",
			input: data.FromString("s"),
			to:    cast.TypeString,
			want:  data.FromString("s"),
		},
		{
			title: "bool true to string",
			input: data.FromBool(true),
			to:    cast.TypeString,
			want:  data.FromString("true"),
		},
		{
			title: "bool false to string",
			input: data.FromBool(false),
			to:    cast.TypeString,
			want:  data.FromString("false"),
		},
		{
			title: "int to bool true",
			input: data.FromInt(1),
			to:    cast.TypeBool,
			want:  data.FromBool(true),
		},
		{
			title: "int to bool false",
			input: data.FromInt(0),
			to:    cast.TypeBool,
			want:  data.FromBool(false),
		},
		{
			title: "float to bool true",
			input: data.FromFloat(1.2),
			to:    cast.TypeBool,
			want:  data.FromBool(true),
		},
		{
			title: "floa to bool false",
			input: data.FromFloat(0),
			to:    cast.TypeBool,
			want:  data.FromBool(false),
		},
		{
			title: "string to bool true",
			input: data.FromString("s"),
			to:    cast.TypeBool,
			want:  data.FromBool(true),
		},
		{
			title: "string to bool false",
			input: data.FromString(""),
			to:    cast.TypeBool,
			want:  data.FromBool(false),
		},
		{
			title: "bool true to bool true",
			input: data.FromBool(true),
			to:    cast.TypeBool,
			want:  data.FromBool(true),
		},
		{
			title: "bool false to bool false",
			input: data.FromBool(false),
			to:    cast.TypeBool,
			want:  data.FromBool(false),
		},
		{
			title: "string to timestamp",
			input: data.FromString("2021-09-25T10:00:00Z"),
			to:    cast.TypeTimestamp,
			want:  data.FromInt(1632564000),
		},
		{
			title: "string to timestamp error",
			input: data.FromString(""),
			to:    cast.TypeTimestamp,
			isErr: true,
		},
		{
			title: "int to time",
			input: data.FromInt(1632564000),
			to:    cast.TypeTime,
			want:  data.FromString("2021-09-25T19:00:00+09:00"),
		},
		{
			title: "float to time",
			input: data.FromFloat(1632564000.2),
			to:    cast.TypeTime,
			want:  data.FromString("2021-09-25T19:00:00+09:00"),
		},
		{
			title: "int to duration",
			input: data.FromInt(3600),
			to:    cast.TypeDuration,
			want:  data.FromString("1h0m0s"),
		},
		{
			title: "float to duration",
			input: data.FromFloat(3600.2),
			to:    cast.TypeDuration,
			want:  data.FromString("1h0m0s"),
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, err := cast.New().Cast(tc.input, tc.to)
			if tc.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want.Value(), got.Value())
		})
	}
}
