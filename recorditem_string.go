// Code generated by "stringer -type=RecordItemType -output recorditem_string.go"; DO NOT EDIT

package lara

import "fmt"

const _RecordItemType_name = "LaborMaterial"

var _RecordItemType_index = [...]uint8{0, 5, 13}

func (i RecordItemType) String() string {
	if i < 0 || i >= RecordItemType(len(_RecordItemType_index)-1) {
		return fmt.Sprintf("RecordItemType(%d)", i)
	}
	return _RecordItemType_name[_RecordItemType_index[i]:_RecordItemType_index[i+1]]
}
