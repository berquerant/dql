// Code generated by "stringer -type Instr -output dig_stringer_generated.go"; DO NOT EDIT.

package dig

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InstrCancel-0]
	_ = x[InstrSkipDir-1]
	_ = x[InstrContinue-2]
}

const _Instr_name = "InstrCancelInstrSkipDirInstrContinue"

var _Instr_index = [...]uint8{0, 11, 23, 36}

func (i Instr) String() string {
	if i < 0 || i >= Instr(len(_Instr_index)-1) {
		return "Instr(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Instr_name[_Instr_index[i]:_Instr_index[i+1]]
}