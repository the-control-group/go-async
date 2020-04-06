package async

// stdLib
import (
	"sync"
)

type ErrorHandler func(err error)

type ResponseHandler func(res interface{})

type FinallyHandler func()

type Promise struct {
	sync.WaitGroup
	handlers sync.WaitGroup
	routine Routine
	res Response
	cancel chan bool
	done chan bool
}

func PromiseAll(promises map[string]*Promise) (*Promise) {
	var multi = &MultiResponse{
		results: make(map[string]chan interface{}, len(promises)),
		errors: make(map[string]chan error, len(promises)),
	}
	var p = &Promise{
		routine: func(res chan<- interface{}, err chan<- error, done chan<- bool, cancel <-chan bool) {
			var wg sync.WaitGroup
			wg.Add(len(promises))
			for name, promise := range promises {
				var myName = name
				multi.results[myName] = make(chan interface{}, 1)
				multi.errors[myName] = make(chan error, 1)
				var myRes = multi.results[myName]
				var myErr = multi.errors[myName]
				promise.Then(func (res interface{}) {
					myRes <- res
					wg.Done()
				}).Catch(func(err error) {
					myErr <- err
					wg.Done()
				}).Finally(func(){
					close(myRes)
					close(myErr)
				})
			}
			wg.Wait()
			close(done)
		},
	}
	p.res = multi
	p.Add(1)
	go func(){
		var res chan interface{}
		var err chan error
		res, err, p.done, p.cancel = newRoutineChannels()
		p.routine(res, err, p.done, p.cancel)
		<-p.done
		close(res)
		close(err)
		p.Done()
		p.Wait()
	}()
	go func(){
		select {
			case <-p.cancel:
				for _, promise := range promises {
					promise.Cancel()
				}
			case <-p.done:
		}
	}()
	p.Then(func(res interface{}){
	})
	return p
}

func NewPromise(routine Routine) (p *Promise) {
	p = &Promise{
		routine: routine,
	}
	p.Add(1)
	go func() {
		var res chan interface{}
		var err chan error
		res, err, p.done, p.cancel = newRoutineChannels()
		p.routine(res, err, p.done, p.cancel)
		<-p.done
		close(res)
		close(err)
		p.res = &SingleResponse{
			res: <-res,
			err: <-err,
		}
		p.Done()
	}()
	return
}

func (p *Promise) Cancel() {
	close(p.cancel)
}

func (p *Promise) Then(fn ResponseHandler) *Promise {
	p.handlers.Add(1)
	go func() {
		p.Wait()
		if !p.res.HasError() {
			fn(p.res.Result())
		}
		p.handlers.Done()
	}()
	return p
}

func (p *Promise) Catch(fn ErrorHandler) *Promise {
	p.handlers.Add(1)
	go func() {
		p.Wait()
		if p.res.HasError() {
			fn(p.res.Error())
		}
		p.handlers.Done()
	}()
	return p
}

func (p *Promise) Finally(fn FinallyHandler) *Promise {
	go func() {
		p.Wait()
		p.handlers.Wait()
		fn()
	}()
	return p
}

func (p *Promise) WaitAll() *Promise {
	p.Wait()
	p.handlers.Wait()
	return p
}