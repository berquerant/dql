package chrono

import "time"

func Now() Time { return &timeImpl{t: time.Now()} }

func TimeFromString(v string) (Time, error) {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, err
	}
	return &timeImpl{
		t: t,
	}, nil
}

func TimeFromTimestamp(v int) Time { return &timeImpl{t: time.Unix(int64(v), 0)} }

type (
	// Time is wrapper of time.Time.
	Time interface {
		Unix() int
		String() string
	}

	timeImpl struct {
		t time.Time
	}
)

func (s *timeImpl) Unix() int      { return int(s.t.Unix()) }
func (s *timeImpl) String() string { return s.t.Format(time.RFC3339) }

func DurationFromString(v string) (Duration, error) {
	d, err := time.ParseDuration(v)
	if err != nil {
		return nil, err
	}
	return &duration{
		d: d,
	}, nil
}

func DurationFromInt(second int) Duration { return &duration{d: time.Duration(second) * time.Second} }

type (
	// Duration is wrapper of time.Duration.
	Duration interface {
		String() string
	}

	duration struct {
		d time.Duration
	}
)

func (s *duration) String() string { return s.d.String() }
