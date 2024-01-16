package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

type Request[TResponse interface{}] interface {
	producesResponse(response TResponse)
}

type requestContext[TRequest Request[TResponse], TResponse interface{}] struct {
	Request *TRequest
	Ctx     context.Context
	Data    map[string]any
}

type RequestHandler[TRequest Request[TResponse], TResponse interface{}] interface {
	Handle(context *requestContext[TRequest, TResponse]) (TResponse, error)
}

var requestHandlers = make(map[reflect.Type]reflect.Value)

func RegisterRequestHandler[TRequest Request[TResponse], TResponse interface{}](handler RequestHandler[TRequest, TResponse]) {
	var zero [0]TRequest
	requestType := reflect.TypeOf(zero).Elem()
	// fmt.Println(fmt.Sprintf("Registering handler for request type %v %v", requestType.Name(), requestType.Kind()))

	requestHandlers[requestType] = reflect.ValueOf(handler.Handle)
}

func Send[TRequest Request[TResponse], TResponse interface{}](
	ctx context.Context,
	request *TRequest,
) (TResponse, error) {
	requestContext := &requestContext[TRequest, TResponse]{
		Request: request,
		Ctx:     ctx,
	}

	requestType := reflect.TypeOf(request).Elem()
	handler, ok := requestHandlers[requestType]
	if !ok {
		var res TResponse
		return res, fmt.Errorf("No handler registered for request type %v", requestType)
	}

	result := handler.Call([]reflect.Value{reflect.ValueOf(requestContext)})

	response := result[0].Interface()

	if !result[1].IsNil() {
		err := result[1].Interface()
		return response.(TResponse), err.(error)
	}

	return response.(TResponse), nil
}
