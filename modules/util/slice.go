package util

import (
	"fmt"
	"reflect"
	"strconv"
)

// ToIntSlice covert slice of interface with value of string, number to int slice
func ToIntSlice(slice []interface{}) ([]int, error) {
	intSlice := make([]int, len(slice))
	var err error
	for i, data := range slice {
		switch v := data.(type) {
		case string:
			intSlice[i], err = strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
		case float32, float64, int, int32, int64:
			intSlice[i], err = strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("value should be string or number, found %v", reflect.TypeOf(data))
		}
	}
	return intSlice, nil
}
