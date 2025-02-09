package toolbox

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

type valueUtils struct{}

func (u *utils) NewValueUtils() *valueUtils {
	return &valueUtils{}
}

type Value struct {
	i *int64
	s *string
	f *float64
}

func (v *valueUtils) NewInt(i int64) *Value {
	return &Value{i: &i}
}

func (v *valueUtils) NewString(s string) *Value {
	return &Value{s: &s}
}

func (v *valueUtils) NewFloat(f float64) *Value {
	return &Value{f: &f}
}

func (b *Value) Int64() int64 {
	if b == nil {
		return 0
	}
	if b.i != nil {
		return *b.i
	}
	return 0
}

func (b *Value) String() string {
	if b == nil {
		return ""
	}
	if b.s != nil {
		return *b.s
	}
	return ""
}

func (b *Value) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		b.s = &s
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		b.i = &i
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		b.f = &f
		return nil
	}

	return errors.New("unsupported value type")
}

func (b Value) MarshalJSON() ([]byte, error) {
	if b.i != nil {
		s := strconv.FormatInt(*b.i, 10)
		return json.Marshal(s)
	} else if b.s != nil {
		return json.Marshal(*b.s)
	} else if b.f != nil {
		s := strconv.FormatFloat(*b.f, 'f', -1, 64)
		return json.Marshal(s)
	}
	return json.Marshal("")
}

func (b *Value) IsInt64() bool {
	return b.i != nil
}

func (b *Value) IsString() bool {
	return b.s != nil
}

func (b *Value) IsFloat64() bool {
	return b.f != nil
}

func (b *Value) ToString() string {
	if b == nil {
		return ""
	}
	if b.s != nil {
		return *b.s
	} else if b.i != nil {
		return strconv.FormatInt(*b.i, 10)
	} else if b.f != nil {
		return strconv.FormatFloat(*b.f, 'f', -1, 64)
	}
	return ""
}

func (b *Value) ToInt64() (int64, error) {
	if b == nil {
		return 0, errors.New("nil Value")
	}
	if b.i != nil {
		return *b.i, nil
	} else if b.s != nil {
		return strconv.ParseInt(*b.s, 10, 64)
	} else if b.f != nil {
		if math.Trunc(*b.f) != *b.f {
			return 0, errors.New("float value is not an integer")
		}
		return int64(*b.f), nil
	}
	return 0, errors.New("no value available")
}

func (b *Value) ToFloat64() (float64, error) {
	if b == nil {
		return 0, errors.New("nil Value")
	}
	if b.f != nil {
		return *b.f, nil
	} else if b.i != nil {
		return float64(*b.i), nil
	} else if b.s != nil {
		f, err := strconv.ParseFloat(*b.s, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}
	return 0, errors.New("no value available")
}

var _ json.Unmarshaler = &Value{}
var _ json.Marshaler = &Value{}
