package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

var _ echo.JSONSerializer = (*EchoSonicJSONSerializer)(nil)
var _ runtime.Marshaler = (*SonicJSONMarshaler)(nil)

// Response 如果要 Code 需要重定义 HTTPStatusFromCode
type Response struct {
	IsSuccess bool   `json:"isSuccess"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
}

func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	httpStatusCode := runtime.HTTPStatusFromCode(status.Code(err))
	st, _ := status.FromError(err)
	customResponse := &Response{
		Msg: st.Message(),
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)
	sonic.ConfigDefault.NewEncoder(writer).Encode(customResponse)
}

func ResponseRewriter(ctx context.Context, response proto.Message) (any, error) {
	return Response{
		IsSuccess: true,
		Data:      response,
	}, nil
}

type EchoSonicJSONSerializer struct{}

func (s EchoSonicJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := sonic.ConfigDefault.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (s EchoSonicJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := sonic.ConfigDefault.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}

type SonicJSONMarshaler struct{}

func NewRespMarshaler() runtime.Marshaler {
	return &SonicJSONMarshaler{}
}

func (g SonicJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return sonic.ConfigDefault.Marshal(v)
}

func (g SonicJSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return sonic.ConfigDefault.Unmarshal(data, v)
}

func (g SonicJSONMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return sonic.ConfigDefault.NewDecoder(r)
}

func (g SonicJSONMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return sonic.ConfigDefault.NewEncoder(w)
}

func (g SonicJSONMarshaler) ContentType(v interface{}) string {
	return "application/json"
}
