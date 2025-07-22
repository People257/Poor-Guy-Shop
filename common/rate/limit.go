package rate

import "time"

type Limit struct {
	windowSize time.Duration
	limit      int
}

func NewLimit(windowSize time.Duration, limit int) Limit {
	return Limit{windowSize: windowSize, limit: limit}
}

func PerSecond(limit int) Limit {
	return NewLimit(time.Second, limit)
}

func PerMinute(limit int) Limit {
	return NewLimit(time.Minute, limit)
}

func PerHour(limit int) Limit {
	return NewLimit(time.Hour, limit)
}

func PerDay(limit int) Limit {
	return NewLimit(24*time.Hour, limit)
}

func (l *Limit) WindowSize() time.Duration {
	return l.windowSize
}

func (l *Limit) Limit() int {
	return l.limit
}
