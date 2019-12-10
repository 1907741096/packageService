package lib

import (
	"time"
	"strconv"
	"database/sql/driver"
	"fmt"
	"strings"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	Day        = time.Hour * 24
	Week       = Day * 7
	Hour       = time.Hour
)

type LocalTime struct {
	time.Time
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	//格式化秒
	seconds := t.Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

func (t *LocalTime) UnmarshalJSON(b []byte) (error) {
	strTime := strings.ReplaceAll(string(b), "\"", "")
	if strTime == EmptyString {
		tt, _ := time.Parse(TimeFormat, "1970-01-01 00:00:00")
		*t = LocalTime{Time: tt}
		return nil
	}

	tt, err := time.Parse(TimeFormat, strTime)
	if err == nil {
		*t = LocalTime{Time: tt}
		return nil
	}

	unix, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	// 修改t底层的值
	*t = LocalTime{Time: time.Unix(unix, 0)}
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t LocalTime) TimeAdd(add time.Duration) LocalTime {
	tt := t.Add(add)
	return LocalTime{Time: tt}
}

func Now() LocalTime { return LocalTime{Time: time.Now()} }

func BeforeDays(days int) LocalTime {
	if days < 0 {
		panic("days must be more than zero")
	}
	t := Now().AddDate(0, 0, -days)
	return LocalTime{Time: t}
}

func AfterDays(days int) LocalTime {
	if days < 0 {
		panic("days must be more than zero")
	}
	t := Now().AddDate(0, 0, days)
	return LocalTime{Time: t}
}

func (t *LocalTime) FormatLocal() string {
	return t.Format(TimeFormat)
}

func (t LocalTime) FullTime() LocalTime {
	timestamp := t.FormatFullTimeStamp()
	tt := time.Unix(timestamp, 0)
	return LocalTime{Time: tt}
}

func (t *LocalTime) FormatFullTimeStamp() int64 {
	return t.Unix() - int64(t.Second()) - int64(60*t.Minute())
}

func (t *LocalTime) Since() time.Duration {
	return time.Since(t.Time)
}

func FromUnix(timestamp int64) *LocalTime {
	t := time.Unix(timestamp, 0)
	return &LocalTime{Time: t}
}
