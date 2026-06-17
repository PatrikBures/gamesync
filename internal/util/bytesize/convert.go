package bytesize

import (
	"errors"
)

const numDigits = 10
var validDigits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var (
	ErrInvalidDigit = errors.New("invalid digit")
)

func pow(n int64, p int64) int64 {
	var result int64
	result = 1
	for range p {
		result*=n
	}
	return result
}

func digitToInt64(byte byte) int64 {
	var i int64
	for ; i < numDigits; i++ {
		if validDigits[i] == byte {
			return i
		}
	}
	return -1
}

func numToInt64(bytes []byte) int64 {
	var num  int64
	var mult int64
	var idx  int64
	for mult = int64(len(bytes))-1; mult > -1; mult-- {
		b := bytes[idx]
		if n := digitToInt64(b); n > -1 {
			c := pow(10, mult) * n
			num += c
		} else {
			println("asdf", string(b))
			return -1
		}
		idx++
	}
	return num
}

func Convert(input string) (int64) {
	inputBytes := []byte(input)
	
	suffixLast := inputBytes[len(inputBytes)-1]

	if digitToInt64(suffixLast) > -1 {
		if num := numToInt64(inputBytes); num == -1 {
			return -1
		} else {
			return num
		}
	}

	suffix := inputBytes[len(inputBytes)-2]

	isMetric := false 
	if suffixLast != 'i'  {
		isMetric = true
		suffix = suffixLast
	}
	
	if isMetric {
		inputNum := inputBytes[:len(inputBytes)-1]
		num := numToInt64(inputNum)
		if num < 0 {
			return -1
		}
		return metric(suffix, num)
	}

	inputNum := inputBytes[:len(inputBytes)-2]
	num := numToInt64(inputNum)
	if num < 0 {
		return -1
	}
	return binary(suffix, num)
}

var suffixes = []byte{0, 'K', 'M', 'G'}
var suffixesLen = int64(len(suffixes))

func metric(suffix byte, num int64) int64 {
	var i int64
	for ; i < suffixesLen; i++ {
		s := suffixes[i]
		if s == suffix {
			return num * pow(1000, i)
		}
	}
	return -1
}

func binary(suffix byte, num int64) int64 {
	switch suffix {
	case 'K':
		return num << 10
	case 'M':
		return num << 20
	case 'G':
		return num << 30
	}
	return -1
}


