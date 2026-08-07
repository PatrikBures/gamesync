package fileval

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// Reads content from a file at p into target
//
// if the file does not exist it returns os.ErrNotExist
func Read[T int | int32 | int64 | string](p string, target *T) error {
	var valueString string
	if content, err := os.ReadFile(p); err != nil {
		return err
	} else {
		valueString = string(content)
		valueString = strings.TrimSpace(valueString)
	}

	errMsg := fmt.Errorf("loading cache '%s'", p)
	switch t := any(target).(type) {
	case *int:
		i, err := strconv.ParseInt(valueString, 10, bits.UintSize)
		if err != nil {
			return errors.Join(errMsg, err)
		}
		*t = int(i)
	case *int32:
		i, err := strconv.ParseInt(valueString, 10, 32)
		if err != nil {
			return errors.Join(errMsg, err)
		}
		*t = int32(i)
	case *int64:
		i, err := strconv.ParseInt(valueString, 10, 64)
		if err != nil {
			return errors.Join(errMsg, err)
		}
		*t = int64(i)
	case *string:
		*t = valueString
	default:
		panic(fmt.Sprintf("undefined type when loading cache: %T", target))
	}
	return nil
}

// Sets content of p with value
func Write[T int | int32 | int64 | string](p string, value T) (err error) {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	var valueString string

	switch t := any(value).(type) {
	case int, int32, int64:
		valueString = fmt.Sprint(value)
	case string:
		valueString = t
	}
	if _, err = f.WriteString(valueString); err != nil {
		return err
	}

	return nil
}

