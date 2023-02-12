package validation

import (
	"errors"
	"strconv"
)

func OrderValid(order string) (bool, error) {

	if len(order) < 1 {
		return false, errors.New("400")
	}

	num, err := strconv.Atoi(order)
	if err != nil {
		return false, errors.New("400")
	}

	res := LuhnValid(num)

	return res, nil
}

// Valid check number is valid or not based on Luhn algorithm
func LuhnValid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
