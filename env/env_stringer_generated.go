// Code generated by "stringer -type Type -output env_stringer_generated.go"; DO NOT EDIT.

package env

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeUnknown-0]
	_ = x[TypeData-1]
	_ = x[TypeExpr-2]
	_ = x[TypeDataList-3]
}

const _Type_name = "TypeUnknownTypeDataTypeExprTypeDataList"

var _Type_index = [...]uint8{0, 11, 19, 27, 39}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
