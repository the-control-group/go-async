package async

// import "log"

type Routine func(res chan<- interface{}, err chan<- error, done chan<- bool, cancel <-chan bool)

func SendReturns (res interface{}, err error) (func(resCh chan<- interface{}, errCh chan<- error)) {
	return func(resCh chan<- interface{}, errCh chan<- error) {
		resCh <- res
		errCh <- err
	}
}

func newRoutineChannels() (res chan interface{}, err chan error, done chan bool, cancel chan bool) {
	res = make(chan interface{}, 1)
	err = make(chan error, 1)
	done = make(chan bool, 1)
	cancel = make(chan bool, 1)
	return res, err, done, cancel
}