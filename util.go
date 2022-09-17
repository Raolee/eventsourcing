package eventsourcing

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func JsonString(v any) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

// panic 을 잡아주는 핸들러, err arg 에 에러를 대입한다.
func handleError(err *error) {
	if r := recover(); r != nil {
		e := ConvertRecoverToError(r)
		*err = e
	}
}

// ConvertRecoverToError | recover 된 value 를 어떻게든 에러로 변환해서 리턴한다. Stack 트레이스 추가는 덤.
func ConvertRecoverToError(r interface{}) error {
	switch x := r.(type) {
	case string:
		return errors.WithStack(errors.New(x))
	case error:
		return errors.WithStack(x)
	default:
		return errors.WithStack(errors.New(fmt.Sprint(x)))
	}
}
