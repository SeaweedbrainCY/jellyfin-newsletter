package clock

import "time"

// clock package is used to inject time dependency accross the app.
// It is mainly use to mock time.Now()

type Interface interface {
	Now() time.Time
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}
