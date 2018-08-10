package scheduler

import (
	"time"
)

func main() {
	//Init timer and set workers size
	// and channel buffer size
	t := New(2, 2) //2 worker with a channel of buffer size 2
	//Start work every x time
	t.Every(100 * time.Millisecond)

	//Append a job (Notice interval type)
	rj := t.AppendJob(&Job{
		Time:  time.Time{},
		Lapse: 1 * time.Second,
		Func: func(i ...interface{}) {
			println("hello ", i[0].(string))
		},
		Params: []interface{}{"world"},
	})
	//Append a job (Notice the one time execution
	// since there is no lapse defined)
	t.AppendJob(&Job{
		Time:  time.Now().Add(1 * time.Second),
		Func: func(i ...interface{}) {
			println("xd")
		},
	})
	time.Sleep(5 * time.Second)
	//Now we remove the first job
	//The second one has been removed automatically
	t.RemoveJob(rj)
	time.Sleep(5 * time.Second)
	//Stop
	t.Stop()
}
