package bytesize

import "testing"

func TestPow(t *testing.T) {
	type inputType struct {
		n      int64
		pow    int64
		result int64
	}
	testData := []inputType{
		{n: 3,    pow: 1,     result: 3},
		{n: 10,   pow: 2,     result: 100},
		{n: 10,   pow: 3,     result: 1000},
		{n: 6,    pow: 9,     result: 10077696},
		{n: -2,   pow: 5,     result: -32},
	}
	for _, td := range testData {
		if result := pow(td.n, td.pow); result != td.result {
			t.Errorf("pow(%d, %d) = %d, want %d", td.n, td.pow, result, td.result)
		}
	}

}
func TestDigitToInt64(t *testing.T) {
	type inputType struct {
		in byte
		out int64
	}
	testData := []inputType{
		{in: '0', out: 0},
		{in: '1', out: 1},
		{in: '2', out: 2},
		{in: '3', out: 3},
		{in: '4', out: 4},
		{in: '5', out: 5},
		{in: '6', out: 6},
		{in: '7', out: 7},
		{in: '8', out: 8},
		{in: '9', out: 9},
	}
	for _, td := range testData {
		if result := digitToInt64(td.in); result != td.out {
			t.Errorf("digitToInt64(%b) = %d, want %d", td.in, result, td.out)
		}
	}
}

func TestNumToInt64(t *testing.T) {
	type inputType struct {
		in []byte
		out int64
	}
	testData := []inputType{
		{in: []byte("190"), out: 190},
		{in: []byte("1"), out: 1},
		{in: []byte("0"), out: 0},
		{in: []byte("88999222"), out: 88999222},
		{in: []byte("1234567890"), out: 1234567890},
		{in: []byte("9038201098732466790"), out: 9038201098732466790},
		{in: []byte(" "), out: -1},
		{in: []byte(" 8372"), out: -1},
		{in: []byte("892Ki"), out: -1},
		{in: []byte("as89283"), out: -1},
		{in: []byte("1G"), out: -1},
	}
	for _, td := range testData {
		if result := numToInt64(td.in); result != td.out {
			t.Errorf("numToInt64(%s) = %d, want %d", td.in, result, td.out)
		}
	}
}

func TestByteSize(t *testing.T) {
	tests := map[string]int64{
		"67K":     67000,
		"10M":     10000000,
		"2G":     2000000000,

		"89Ki":    91136,
		"8291Mi":  8693743616,
		"5Gi":     5368709120,

		"893":     893,
		"3":       3,
		"0":       0,
		"0Gi":     0,

		"-1G":     -1,
		"-2G":     -1,
		"-5Ki":    -1,
		"akjf":    -1,
		"8Kii":    -1,
		"1A":      -1,
		" 3":      -1,
	}

	for input, want := range tests {
		result := Convert(input) 
		if result != want {
			t.Errorf(`ByteSize(%s) = %d, want %d`, input, result, want)
		}
	}
}

