package timer

import "time"

// Clock is a wrapper global time for testing purpose
type Clock struct {
	faked bool
	time  time.Time
}

func (c *Clock) Set(time time.Time) {
	c.faked = true
	c.time = time
}

func (c *Clock) Now() time.Time {
	if c.faked {
		return c.time
	}

	return time.Now()
}
