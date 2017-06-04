package component

import (
	"github.com/vaporz/turbo"
	"reflect"
	"net/http"
	"github.com/vaporz/turbo-example/yourservice/gen/proto"
	"strconv"
	"errors"
	i "github.com/vaporz/turbo-example/yourservice/interceptor"
	"google.golang.org/grpc"
)

func GrpcClient(conn *grpc.ClientConn) interface{} {
	return proto.NewYourServiceClient(conn)
}

func InitComponents() {
	turbo.Intercept([]string{"GET"}, "/hello", i.LogInterceptor{})
	turbo.Intercept([]string{"GET"}, "/eat_apple/{num:[0-9]+}", i.LogInterceptor{})
	turbo.Intercept([]string{"GET"}, "/a/a", i.LogInterceptor{Msg: "url interceptor"})
	turbo.Intercept([]string{}, "/a/", i.LogInterceptor{Msg: "path interceptor"})
	turbo.SetPreprocessor("/eat_apple/{num:[0-9]+}", preEatApple)
	//turbo.SetHijacker("/eat_apple/{num:[0-9]+}", hijackEatApple)
	turbo.SetPostprocessor("/eat_apple/{num:[0-9]+}", postEatApple)

	//turbo.RegisterMessageFieldConvertor(new(proto.CommonValues), convertCommonValues)

	turbo.WithErrorHandler(errorHandler)
}

func errorHandler(resp http.ResponseWriter, req *http.Request, err error) {
	resp.Write([]byte("from errorHandler:" + err.Error()))
}

func convertCommonValues(req *http.Request) reflect.Value {
	result := &proto.CommonValues{}
	result.SomeId = 1111111
	return reflect.ValueOf(result)
}

func hijackEatApple(resp http.ResponseWriter, req *http.Request) {
	client := turbo.GrpcService().(proto.YourServiceClient)
	r := new(proto.EatAppleRequest)
	r.Num = 999
	res, err := client.EatApple(req.Context(), r)
	if err == nil {
		resp.Write([]byte(res.String()))
	} else {
		resp.Write([]byte(err.Error()))
	}
}

func preEatApple(resp http.ResponseWriter, req *http.Request) error {
	num, err := strconv.Atoi(req.Form["num"][0])
	if err != nil {
		resp.Write([]byte("'num' is not numberic"))
		return errors.New("invalid num")
	}
	if num > 5 {
		resp.Write([]byte("Too many apples!"))
		return errors.New("Too many apples")
	}
	return nil
}

func postEatApple(resp http.ResponseWriter, req *http.Request, serviceResp interface{}, err error) {
	sr := serviceResp.(*proto.EatAppleResponse)
	resp.Write([]byte("this is from postprocesser, message=" + sr.Message))
}