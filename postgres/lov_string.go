// Code generated by "stringer -type=listOfValuesType -output lov_string.go"; DO NOT EDIT

package postgres

import "fmt"

const _listOfValuesType_name = "titleunitgenderspeciesbreed"

var _listOfValuesType_index = [...]uint8{0, 5, 9, 15, 22, 27}

func (i listOfValuesType) String() string {
	if i < 0 || i >= listOfValuesType(len(_listOfValuesType_index)-1) {
		return fmt.Sprintf("listOfValuesType(%d)", i)
	}
	return _listOfValuesType_name[_listOfValuesType_index[i]:_listOfValuesType_index[i+1]]
}