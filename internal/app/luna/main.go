package luna

import (
	"strconv"
	"strings"
)

func IsLuna(value string) (bool, error) {
	strArr := strings.Split(value, "")
	var arr []int

	for _, v := range strArr {
		res, err := strconv.Atoi(v)
		if err != nil {
			return false, err
		}
		arr = append(arr, res)
	}

	i := len(arr) % 2
	for ; i < len(arr); i += 2 {
		arr[i] *= 2
		if arr[i] > 9 {
			arr[i] -= 9
		}
	}

	summ := 0
	for i = 0; i < len(arr); i++ {
		summ += arr[i]
	}

	if summ%10 == 0 {
		return true, nil
	}

	return false, nil
}
