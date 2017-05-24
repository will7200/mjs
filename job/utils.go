package job

import "time"

type Timer struct {
	timer *time.Timer
	end   time.Time
}

func NewTimer(t time.Duration) *Timer {
	return &Timer{time.NewTimer(t), time.Now().Add(t)}
}

func NewAfterFunc(t time.Duration, f func()) *Timer {
	return &Timer{time.AfterFunc(t, f), time.Now().Add(t)}
}

func (s *Timer) Reset(t time.Duration) {
	s.timer.Reset(t)
	s.end = time.Now().Add(t)
}

func (s *Timer) Stop() {
	s.timer.Stop()
}

func (s *Timer) HasPassed() bool {
	return time.Now().After(s.end)
}
