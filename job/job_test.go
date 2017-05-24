package job

import (
	"encoding/json"
	"fmt"
	"os/user"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	REPEAT   = "R2"
	INTERVAL = "PT10S"
)
var newjob = []byte(
	`{"Name":"Youtube Channel"}`,
)

var newjobResults = []byte(
	`"Name":"Youtube Channel"`,
)
var Dispatch *Dispatcher
var db *gorm.DB
var err error

func init() {
	Dispatch = &Dispatcher{}
	db, err = gorm.Open("postgres", "postgres://williamfl:williamfl@192.168.1.25/wjob")
	Dispatch.SetPersistStorage(db)
	Dispatch.StartDispatcher(2)
}

func getMockJob(testjob *Job) error {
	owner, er := user.Current()
	if er != nil {
		fmt.Println("Cannot get the current user setting to blank")
	}
	testjob.Owner = owner.Name
	testjob.Schedule = fmt.Sprintf("%s/%s/%s", REPEAT, time.Now().Add(5*time.Second).Format(time.RFC3339), INTERVAL)
	err := testjob.NewJob(newjob, Dispatch)
	return err
}
func TestNewJob(t *testing.T) {
	testjob := &Job{}
	err := getMockJob(testjob)
	assert.NoError(t, err)
	assert.Equal(t, db.Create(testjob).GetErrors(), []error{}, "Should not contains any errors")
	assert.NotEmpty(t, testjob.ID, "Job came back with a bad ID")
	tt, err := json.Marshal(testjob)
	if err != nil {
		fmt.Println("Error")
	}
	assert.Contains(t, string(tt), string(newjobResults), "Jobs is suppose to contains these features")
	db.Delete(&testjob)
}

func TestParseSchedule(t *testing.T) {
	testjob := &Job{}
	testjob.Schedule = fmt.Sprintf("%s/%s/%s", REPEAT, time.Now().Add(-5*time.Second).Format(time.RFC3339), INTERVAL)
	err := testjob.NewJob(newjob, Dispatch)
	assert.Error(t, err)
}

func TestGetWaitingDuration(t *testing.T) {
	testjob := &Job{}
	now := time.Now().Add(5 * time.Second)
	testjob.Schedule = fmt.Sprintf("%s/%s/%s", "R0", now.Format(time.RFC3339), INTERVAL)
	err := testjob.NewJob(newjob, Dispatch)
	later := testjob.GetWaitDuration()
	assert.WithinDuration(t, now, now.Add(later), 5*time.Second, "Time should be within 5 seconds")
	assert.NoError(t, err)
	testjob.ScheduleTime = time.Now().Add(time.Second * -5)
	later = testjob.GetWaitDuration()
	assert.Equal(t, 0*time.Second, later)
}
