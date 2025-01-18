package toolbox

import (
	"encoding/json"
	"strconv"
)

type valueUtils struct {
}

func (u *utils) NewValueUtils() *valueUtils {
	return &valueUtils{}
}

type Value struct {
	i *int64
	s *string
}

func (*valueUtils) NewInt(i int64) *Value {
	return &Value{i: &i}
}

func (*valueUtils) NewString(s string) *Value {
	return &Value{s: &s}
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

func (b *Value) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
		b.s = &s
	} else {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		b.i = &v
	}
	return nil
}

func (b Value) MarshalJSON() ([]byte, error) {
	if b.i != nil {
		v := strconv.FormatInt(*b.i, 10)
		return []byte("\"" + v + "\""), nil
	} else if b.s != nil {
		v := "\"" + *b.s + "\""
		return []byte(v), nil
	}
	return []byte(`""`), nil
}

func (b *Value) IsInt64() bool {
	return b.i != nil
}
func (b *Value) IsString() bool {
	return b.s != nil
}

func (b *Value) ToString() string {
	if b.IsInt64() {
		return strconv.FormatInt(*b.i, 10)
	} else if b.IsString() {
		return *b.s
	}
	return ""
}

func (b *Value) ToInt64() (int64, error) {
	if b.s != nil {
		intI, err := strconv.Atoi(*b.s)
		if err != nil {
			return 0, err
		}
		return int64(intI), nil
	} else {
		return *b.i, nil
	}
}

var _ json.Unmarshaler = &Value{}
var _ json.Marshaler = &Value{}
