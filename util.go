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

func handleError(err *error) {
	if r := recover(); r != nil {
		e := ConvertRecoverToError(r)
		*err = e
	}
}

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
