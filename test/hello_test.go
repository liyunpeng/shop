package test

import (
	"shop/client"
	"shop/service"
	"testing"
)

const helloPrefix = "Hello, "

/*
	被测试的函数
 */
func Hello(name string) string {
	return helloPrefix + name
}

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf(" got '%s' want '%s'", got, want)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("World1")
		want := "Hello, World1"
		assertCorrectMessage(t, got, want)
	})
}

func TestMicro(t *testing.T){

	client.MicroCallUser()

	service.StartMicroService()
}

