package ioc

import (
	"fmt"
	"math/rand"
	"testing"
)

type GreeterService interface {
	Greet() string
}

type DefaultGreeter struct {
	name string
}

func (greeter *DefaultGreeter) Greet() string {
	return fmt.Sprintf("%s: Hello, World!", greeter.name)
}

type AnotherGreeter struct{}

func (greeter *AnotherGreeter) Greet() string {
	return "AnotherGreeter: Hello, World!"
}

func TestSingletons(t *testing.T) {
	t.Run("resolves a singleton instance", func(t *testing.T) {
		reset()
		instance := DefaultGreeter{
			name: "DefaultGreeter",
		}

		RegisterSingletonInstance[GreeterService](&instance)
		RegisterSingletonInstance[DefaultGreeter](DefaultGreeter{
			name: "foo",
		})
		Build()

		greeter := Resolve[GreeterService]()

		if greeter == nil {
			t.Error("Greeter should not be nil")
		}

		if greeter != &instance {
			t.Error("Greeter should be the same instance")
		}

		greeting := greeter.Greet()
		if greeting != "DefaultGreeter: Hello, World!" {
			t.Errorf("Greeter should greet 'Hello, World!', but got '%s'", greeting)
		}
	})
	t.Run("resolves all singletons", func(t *testing.T) {
		reset()
		names := []string{"a", "b", "c"}
		for _, name := range names {
			RegisterSingletonInstance[GreeterService](&DefaultGreeter{name})
		}
		Build()

		greeters := ResolveAll[GreeterService]()

		if len(greeters) != 3 {
			t.Errorf("Expected 3 greeters, but got %d", len(greeters))
		}

		for i, name := range names {
			greeting := greeters[i].Greet()
			if greeting != fmt.Sprintf("%s: Hello, World!", name) {
				t.Errorf("Greeter should greet 'Hello, World!', but got '%s'", greeting)
			}
		}
	})
}

func TestTransients(t *testing.T) {
	t.Run("resolves a transient instance", func(t *testing.T) {
		reset()
		RegisterFactory[GreeterService](func() GreeterService {
			return &DefaultGreeter{
				name: randomString(5),
			}
		}, Transient)
		Build()

		greeter := Resolve[GreeterService]()
		greeter2 := Resolve[GreeterService]()

		if greeter == nil || greeter2 == nil {
			t.Error("Greeters are nil")
		}

		greeting := greeter.Greet()
		greeting2 := greeter2.Greet()
		if greeting == greeting2 {
			t.Errorf("Greetings should not be the same instance")
		}
	})

	t.Run("resolves all transient instances", func(t *testing.T) {
		reset()
		RegisterFactory[GreeterService](func() GreeterService {
			return &DefaultGreeter{
				name: randomString(5),
			}
		}, Transient)
		RegisterFactory[GreeterService](func() GreeterService {
			return &AnotherGreeter{}
		}, Transient)
		Build()

		greeters := ResolveAll[GreeterService]()

		if len(greeters) != 2 {
			t.Errorf("Expected 2 greeters, but got %d", len(greeters))
		}

		// check that greeters[0] is DefaultGreeter
		if _, ok := greeters[0].(*DefaultGreeter); !ok {
			t.Errorf("greeters[0] should be DefaultGreeter")
		}

		// check that greeters[1] is AnotherGreeter
		if _, ok := greeters[1].(*AnotherGreeter); !ok {
			t.Errorf("greeters[1] should be AnotherGreeter")
		}
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
