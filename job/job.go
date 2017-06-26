package job

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/will7200/mjs/utils/iso8601"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	RFC3339WithoutTimezone = "2006-01-02T15:04:05"
)

type Job struct {
	Name string `json:"Name"` // command is required
	ID   string `json:"ID" gorm:"primary_key"`

	Owner string `json:"Owner"`

	Type            string `json:"Type"` //default is command but can be function
	InternalType    int
	Command         pq.StringArray `json:"Command,omitempty" gorm:"type:varchar(100)"`
	FunctionName    string         `json:"Function,omitempty"`
	Arguments       pq.StringArray `json:"Arguments,omitempty" gorm:"type:varchar(100)"`
	Description     string         `json:"Description,omitempty"`
	IsActive        bool           `json:"Active" sql:"type:boolean;default:true"`
	Schedule        string         `json:"Schedule"` //required
	ScheduleTime    time.Time
	Domain          string            `json:"Domain"`
	SubDomain       string            `json:"SubDomain"`
	Application     string            `json:"Application"`
	NextRunAt       time.Time         `json:"Next Run At"`
	Epsilon         string            `json:"epsilon,omitempty"`
	DelayDuration   *iso8601.Duration `gorm:"type:varchar(100)" json:"-"`
	EpsilonDuration *iso8601.Duration `gorm:"type:varchar(100)" json:"-"`
	LastRunAt       time.Time
	TimesToRepeat   int64
	jobTimer        *time.Timer
	Stats           []*JobStats `json:"stats,omitempty" gorm:"ForeignKey:ID;AssociationForeignKey:ID"`
	lock            sync.RWMutex
}

type JobStats struct {
	ID                string        `json:"job_id"`
	Session           string        `json:"Session" gorm:"primary_key"`
	RanAt             time.Time     `json:"ran_at"`
	NumberOfRetries   uint          `json:"number_of_retries"`
	Success           bool          `json:"success"`
	ExecutionDuration time.Duration `json:"execution_duration"`
	ErrorMessage      string        `json:"error,omitempty"`
}

func (j *Job) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV1().String())
	scope.DB().Model(j).UpdateColumn(Job{IsActive: true})
	return nil
}

func RawJob(v []byte, d *Dispatcher) (*Job, error) {
	j := &Job{}
	err := j.NewJob(v, d)
	return j, err
}
func (j *Job) NewJob(v []byte, d *Dispatcher) error {
	if err := json.Unmarshal(v, j); err != nil {
		log.Printf("Error occured when unmarshalling data: %s", err)
		return err
	}
	if err := j.ParseSchedule(); err != nil {
		return err
	}
	j.StartWaiting(d)
	return nil
}

func (j *Job) ParseSchedule() error {
	var err error
	splitTime := strings.Split(j.Schedule, "/")
	if len(splitTime) != 3 {
		return fmt.Errorf(
			"Schedule not formatted correctly. Should look like: R/2014-03-08T20:00:00Z/PT2H",
		)
	}

	// Handle Repeat Amount
	if splitTime[0] == "R" {
		// Repeat forever
		j.TimesToRepeat = -1
	} else {
		j.TimesToRepeat, err = strconv.ParseInt(strings.Split(splitTime[0], "R")[1], 10, 0)
		if err != nil {
			log.Errorf("Error converting TimesToRepeat to an int: %s", err)
			return err
		}
	}
	log.Debugf("TimesToRepeat: %d", j.TimesToRepeat)

	j.ScheduleTime, err = time.Parse(time.RFC3339, splitTime[1])
	if err != nil {
		j.ScheduleTime, err = time.Parse(RFC3339WithoutTimezone, splitTime[1])
		if err != nil {
			log.Errorf("Error converting scheduleTime to a time.Time: %s", err)
			return err
		}
	}
	if (time.Duration(j.ScheduleTime.UnixNano() - time.Now().UnixNano())) < 0 {
		return fmt.Errorf("Schedule time has passed on Job with id of %s", j.ID)
	}
	log.Debugf("Schedule Time: %s", j.ScheduleTime)

	if j.TimesToRepeat != 0 {
		j.DelayDuration, err = iso8601.FromString(splitTime[2])
		if err != nil {
			log.Errorf("Error converting delayDuration to a iso8601.Duration: %s", err)
			return err
		}
		log.Debugf("Delay Duration: %s", j.DelayDuration.ToDuration())
	}

	if j.Epsilon != "" {
		j.EpsilonDuration, err = iso8601.FromString(j.Epsilon)
		if err != nil {
			log.Errorf("Error converting j.Epsilon to iso8601.Duration: %s", err)
			return err
		}
	}

	return nil
}

func (j *Job) CheckSchedule(d *Dispatcher) {
	if !j.IsActive {
		return
	}
	if j.TimesToRepeat == -1 {
		j.StartWaiting(d)
		return
	}
	if len(j.Stats) < int(j.TimesToRepeat) {
		j.StartWaiting(d)
		return
	}
}
func (j *Job) StartWaiting(d *Dispatcher) {
	waitDuration := j.GetWaitDuration()

	log.Infof("Job Scheduled to run in: %s", waitDuration)
	j.lock.Lock()
	j.NextRunAt = time.Now().Add(waitDuration)
	j.lock.Unlock()
	d.AddFutureJob(j, waitDuration)
}

func (j *Job) GetWaitDuration() time.Duration {
	log.Debugf("%+v", j)
	waitDuration := time.Duration(j.ScheduleTime.UnixNano() - time.Now().UnixNano())

	if waitDuration < 0 {
		if j.TimesToRepeat == 0 {
			return 0
		}

		if j.LastRunAt.IsZero() {
			waitDuration = j.DelayDuration.ToDuration()
			t := j.ScheduleTime
			for {
				t = t.Add(waitDuration)
				if t.After(time.Now()) {
					break
				}
			}
			waitDuration = t.Sub(time.Now())
		} else {
			last := j.LastRunAt.Add(j.DelayDuration.ToDuration())
			waitDuration = last.Sub(time.Now())
		}
	}

	return waitDuration
}

func (j *Job) AddStats(s *JobStats) {
	j.lock.Lock()
	defer j.lock.Unlock()
	if s != nil {
		j.Stats = append(j.Stats, s)
	}
}
