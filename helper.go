package grpcsteps

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	xreflect "go.nhat.io/grpcmock/reflect"
	"go.nhat.io/grpcmock/service"
	"google.golang.org/grpc/codes"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func unmarshal(in interface{}, isSlice bool, data *string) (interface{}, error) {
	result := reflect.New(xreflect.UnwrapType(in))

	if isSlice {
		result = reflect.MakeSlice(reflect.SliceOf(result.Type()), 0, 0)
		result = reflect.New(result.Type())
	}

	if data == nil {
		if isSlice {
			return reflect.Zero(result.Elem().Type()).Interface(), nil
		}

		return reflect.Zero(result.Type()).Interface(), nil
	}

	if err := json.Unmarshal([]byte(*data), result.Interface()); err != nil {
		return nil, err
	}

	if isSlice {
		return result.Elem().Interface(), nil
	}

	return result.Interface(), nil
}

func toPayload(methodType service.Type, in interface{}, data *string) (interface{}, error) {
	isSlice := service.IsMethodClientStream(methodType) || service.IsMethodBidirectionalStream(methodType)

	return unmarshal(in, isSlice, data)
}

func toStatusCode(data string) (codes.Code, error) {
	data = fmt.Sprintf("%q", toUpperSnakeCase(data))

	var code codes.Code

	if err := code.UnmarshalJSON([]byte(data)); err != nil {
		return codes.Unknown, err
	}

	return code, nil
}

func toUpperSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToUpper(snake)
}
