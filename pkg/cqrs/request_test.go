package cqrs

import (
	"context"
	"testing"
)

type NilHandlerRequest struct {
	Nil bool
}

func (r NilHandlerRequest) producesResponse(re struct{}) {}

type SayHelloRequest struct {
	Name string
}

func (r SayHelloRequest) producesResponse(re SayHelloResponse) {}

type SayHelloResponse struct {
	Message string
}

type SayHelloRequestHandler struct{}

func (h *SayHelloRequestHandler) Handle(ctx *requestContext[SayHelloRequest, SayHelloResponse]) (SayHelloResponse, error) {
	return SayHelloResponse{Message: "Hello " + ctx.Request.Name}, nil
}

type EmptyRequest struct{}

func (r EmptyRequest) producesResponse(res struct{}) {}

type EmptyRequestHandler struct{}

func (h *EmptyRequestHandler) Handle(ctx *requestContext[EmptyRequest, struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

// create test
func TestSend(t *testing.T) {
	RegisterRequestHandler(&SayHelloRequestHandler{})

	request := &SayHelloRequest{
		Name: "John",
	}

	response, err := Send(context.Background(), request)

	if err != nil {
		t.Errorf("Send() = %v, want %v", err, nil)
	}

	if response.Message != "Hello John" {
		t.Errorf("Send() = %v, want %v", response.Message, "Hello John")
	}
}

func TestSendWithNoHandler(t *testing.T) {
	request := &NilHandlerRequest{}

	_, err := Send(context.Background(), request)

	if err == nil || err.Error() != "No handler registered for request type cqrs.NilHandlerRequest" {
		t.Errorf("error is nil, want %v", "No handler registered for request type cqrs.NilHandlerRequest")
	}
}

func BenchmarkSend(b *testing.B) {
	RegisterRequestHandler(&EmptyRequestHandler{})
	request := &EmptyRequest{}

	b.Run("Send", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Send(context.Background(), request)
		}
	})
}
