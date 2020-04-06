package async

//stdLib
import (
	"fmt"
	"sync"
)

type Response interface {
	HasError() bool
	Error() error
	HasResult() bool
	Result() interface{}
}

type SingleResponse struct {
	res interface{}
	err error
}

func (r *SingleResponse) HasError() bool {
	return r.err != nil
}

func (r *SingleResponse) Error() error {
	return r.err
}

func (r *SingleResponse) HasResult() bool {
	return r.res != nil
}

func (r *SingleResponse) Result() interface{} {
	return r.res
}

type MultiResponse struct {
	sync.Mutex
	results map[string]chan interface{}
	errors map[string]chan error
	resultsRead map[string]interface{}
	errorsRead map[string]error
}

func (r *MultiResponse) HasError() bool {
	var errors = r.Errors()
	for _, err := range errors {
		if err != nil {
			return true
		}
	}
	return false
}

func (r *MultiResponse) Error() (error) {
	return fmt.Errorf("%v", r.Errors())
}

func (r *MultiResponse) Errors() (map[string]error) {
	r.Lock()
	defer r.Unlock()
	if r.errorsRead != nil {
		return r.errorsRead
	}
	r.errorsRead = make(map[string]error, len(r.errors))
	for name, ch := range r.errors {
		r.errorsRead[name], _ = <-ch
	}
	return r.errorsRead
}

func (r *MultiResponse) HasResult() bool {
	var results = r.Results()
	for _, res := range results {
		if res != nil {
			return true
		}
	}
	return false
}

func (r *MultiResponse) Results() (map[string]interface{}) {
	r.Lock()
	defer r.Unlock()
	if r.resultsRead != nil {
		return r.resultsRead
	}
	r.resultsRead = make(map[string]interface{}, len(r.results))
	for name, ch := range r.results {
		for res := range ch {
			r.resultsRead[name] = res
		}
	}
	return r.resultsRead
}

func (r *MultiResponse) Result() (interface{}) {
	return r.Results()
}