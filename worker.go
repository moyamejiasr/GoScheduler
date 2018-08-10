package scheduler

//Struct that contains a reference to
// the worker quit channel
type worker struct {
	quit chan bool
}

//A worker that executes a function received
// via channel from the queue list.
func newWorker(exec chan *Job) *worker {
	quit := make(chan bool)
	go func() {
		for {
			select {
			case f := <- exec:
				f.Func(f.Params...)
			case <- quit:
				return
			}
		}
	}()
	return &worker{quit}
}