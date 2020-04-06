package async

import (
	"testing"
)

func TestNewPromise (t *testing.T) {
	NewPromise(func(res chan<- interface{}, err chan<- error, done chan<- bool, cancel <-chan bool){
		res <- 5
		close(done)
	}).Then(func(res interface{}) {
		if res != 5 {
			t.Error("Result was not 5")
		}
		t.Log("done")
	}).Catch(func(err error) {
		t.Error(err)
	}).WaitAll()
}

func TestPromiseAll (t *testing.T) {
	PromiseAll(map[string]*Promise{
		"five": NewPromise(func(res chan<- interface{}, err chan<- error, done chan<- bool, cancel <-chan bool){
			res <- 5
			close(done)
		}),
		"six": NewPromise(func(res chan<- interface{}, err chan<- error, done chan<- bool, cancel <-chan bool){
			res <- 6
			close(done)
		}),
	}).Then(func(res interface{}) {
		results, ok := res.(map[string]interface{})
		if !ok {
			t.Error("Could not cast result")
		}
		if results["five"] != 5 {
			t.Error("five was not 5")
		}
		if results["six"] != 6 {
			t.Error("six was not 6")
		}
		t.Log("done")
	}).Catch(func(err error) {
		t.Error(err)
	}).WaitAll()
}