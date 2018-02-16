package job

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	FUNCTIONCALL  = 1
	HTTPREQUEST   = 2
	OSCOMMAND     = 3
	ErrCmdIsEmpty = errors.New("No Command was provided")
)

//WorkRequest possible
type WorkRequest struct {
	wJob *Job
	id   uuid.UUID
	when *Timer
}

var callable Funcs

func NewWorker(id int, workerQueue chan chan WorkRequest, verbose bool) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		message:     make(chan string),
		verbose:     verbose}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	message     chan string
	verbose     bool
}

// Start  the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			//time.Sleep(time.Second*60)
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				Type := work.wJob.InternalType
				stats := &JobStats{Session: work.id.String()}
				work.wJob.lock.Lock()
				work.wJob.LastRunAt = time.Now()
				work.wJob.lock.Unlock()
				work.wJob.lock.RLock()
				switch Type {
				case FUNCTIONCALL:
					s := make([]interface{}, len(work.wJob.Arguments))
					for i, j := range work.wJob.Arguments {
						s[i] = j
					}
					stats.RanAt = time.Now()
					vals, err := callable.Call(work.wJob.FunctionName, s...)
					if err != nil {
						log.Errorln(err)
						stats.Success = false
						stats.ErrorMessage = err.Error()
						break
					}
					stats.Success = true
					stats.ExecutionDuration = time.Now().Sub(stats.RanAt)
					log.Debugf("Function:%s with parameters (%v) got values %v", work.wJob.FunctionName, s, vals)
				case OSCOMMAND:
					args := work.wJob.Command
					if len(args) == 0 {
						stats.Success = false
						stats.ErrorMessage = ErrCmdIsEmpty.Error()
						log.Debugf("Exiting Early")
						break
					}
					stats.RanAt = time.Now()
					log.Debugf("%v\n", args)
					cmd := exec.Command(args[0], args[1:]...)
					if work.wJob.Pipeoutput || w.verbose {
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
					}
					err := cmd.Run()
					if err != nil {
						log.Errorln(err)
						stats.Success = false
						stats.ErrorMessage = err.Error()
						break
					}
					stats.Success = true
					stats.ExecutionDuration = time.Now().Sub(stats.RanAt)
				default:
					//log.Println("")
					log.Debugf("Cannot do requested work unknown type: %s", work.wJob.Type)
				}
				work.wJob.lock.RUnlock()
				work.wJob.AddStats(stats)
				StatsQueue <- stats
				CheckQueue <- work.wJob
			case <-w.QuitChan:
				// We have been asked to stop.
				log.Infof("worker %d stopping", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

//Dispatcher keeps track of the workers
type Dispatcher struct {
	Workers     *arraylist.List
	WorkerQueue chan chan WorkRequest
	timer       *time.Timer
	quit        chan bool
	RegFuncs    Funcs
	db          *gorm.DB
	Waiting     *rbt.Tree
	verbose     bool
}

var WorkQueue = make(chan WorkRequest, 100)
var Message = make(chan string, 100)
var CheckQueue = make(chan *Job, 10)
var StatsQueue = make(chan *JobStats, 10)

func (d *Dispatcher) RegisterFunctions(f Funcs) {
	callable = f
}

//StartDispatcher will
func (d *Dispatcher) StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	d.WorkerQueue = make(chan chan WorkRequest, nworkers)
	d.quit = make(chan bool, 1)
	WorkerQueue := d.WorkerQueue
	d.Workers = arraylist.New()
	d.Waiting = rbt.NewWithStringComparator()
	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Infof("Starting worker %d", i+1)
		worker := NewWorker(i+1, d.WorkerQueue, d.verbose)
		worker.Start()
		d.Workers.Add(worker)
	}
	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			case work := <-CheckQueue:
				go func(j *Job, d *Dispatcher) {
					j.CheckSchedule(d)
					d.db.Save(j)
				}(work, d)
			case stats := <-StatsQueue:
				//log.Info("CALLING")
				_ = stats
				//go func(stats *JobStats, d *Dispatcher) {
				//if err := d.db.Create(stats).Error; err != nil {
				//	log.Debugf("Error saving stats")
				//}
				//}(stats, d)
			}
		}
	}()
}

func (d *Dispatcher) Block(t time.Duration) {
	time.AfterFunc(t, func() {
		d.quit <- true
	})
	//tickChan := time.NewTicker(500 * time.Millisecond).C
	for {
		select {
		case <-d.quit:
			return
		case me := <-Message:
			fmt.Printf(me)
			//case now := <-tickChan:
			//fmt.Printf("\r\rUpdate jobs waiting %d, avaiable workers %d, %v", len(d.Waiting), len(d.WorkerQueue), now)
		}
	}
}
func (d *Dispatcher) AddJob(w WorkRequest) {
	log.Infof("Job %s added to queue, will run shortly", w.wJob.Name)
	WorkQueue <- w
}

func (d *Dispatcher) RemoveWaiting(t *WorkRequest) {
	_, found := d.Waiting.Get(t.id.String())
	if found {
		t.when.Stop()
		d.AddJob(*t)
		d.Waiting.Remove(t.id.String())
	}
}
func (d *Dispatcher) AddWorkRequest(w WorkRequest, t time.Duration) {
	w.wJob.lock.Lock()
	defer w.wJob.lock.Unlock()
	f := func() {
		d.RemoveWaiting(&w)
	}
	w.when = NewAfterFunc(t, f)
	w.wJob.jobTimer = w.when.timer
	d.Waiting.Put(w.id.String(), w)
}

func (d *Dispatcher) RemoveWorkRequest(j *Job) bool {
	it := d.Waiting.Iterator()
	for it.Next() {
		value := it.Value()
		val := value.(WorkRequest)
		if val.wJob.ID == j.ID {
			val.when.Stop()
			val.wJob.jobTimer.Stop()
			d.Waiting.Remove(val.id.String())
			return true
		}
	}
	return false
}

func (d *Dispatcher) AddFutureJob(w *Job, t time.Duration) {
	d.db.Save(w)
	work := WorkRequest{wJob: w, id: uuid.NewV1()}
	d.AddWorkRequest(work, t)
}

func (d *Dispatcher) SetPersistStorage(db *gorm.DB) {
	d.db = db
}

func (d *Dispatcher) SetVerbose(v bool) {
	d.verbose = v
}

func (d *Dispatcher) AddPendingJobs() {
	jobs := []Job{}
	d.db.Preload("Stats").Find(&jobs)
	for i := 0; i < len(jobs); i++ {
		t := &jobs[i]
		t.CheckSchedule(d)
		d.db.Save(t)
	}
}
