package job

import (
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func getMockDispatcher(amount int) Dispatcher {
	D := Dispatcher{}
	D.SetPersistStorage(db)
	if amount != 0 {
		D.StartDispatcher(amount)
	}
	return D
}

var w Worker
var ww []Worker

func TestStartDispatcher(t *testing.T) {
	d := getMockDispatcher(4)
	assert.Equal(t, d.Workers.Size(), 4)
	assert.Equal(t, len(d.Waiting), 0)
	if work, ok := d.Workers.Get(0); ok {
		work := work.(Worker)
		work.Stop()
	}
	testjob := &Job{}
	now := time.Now().Add(1 * time.Second)
	testjob.Schedule = fmt.Sprintf("%s/%s/%s", "R", now.Format(time.RFC3339), INTERVAL)
	testjob.InternalType = OSCOMMAND
	testjob.Command = pq.StringArray{"python", "-V"}
	_ = testjob.NewJob(newjob, &d)
	_ = testjob.NewJob(newjob, &d)
	_ = testjob.NewJob(newjob, &d)
	_ = testjob.NewJob(newjob, &d)
	_ = testjob.NewJob(newjob, &d)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 3, len(d.WorkerQueue))
}

func add(d *Dispatcher) {

}
