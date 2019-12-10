package utils

import (
	"fmt"
	"time"
)

func IsBetween(begin, end time.Time) bool {
	now := time.Now().Unix()
	return now >= begin.Unix() && now <= end.Unix()
}

func IsExpire(end time.Time) bool {
	now := time.Now().Unix()
	return end.Unix() <= now
}

func GetDayStart(day time.Time) time.Time {
	start, _ := time.Parse(DatetimeFormat, day.Format(DateFormat)+" 00:00:00")
	return start
}

func GetDayEnd(day time.Time) time.Time {
	end, _ := time.Parse(DatetimeFormat, day.Format(DateFormat)+" 23:59:59")
	return end
}

func GetYesterday() time.Time {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	return yesterday
}

type JsonTime time.Time

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DatetimeFormat = "2006-01-02 15:04:05"
)

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+DatetimeFormat+`"`, string(data), time.Local)
	*t = JsonTime(now)
	return
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(DatetimeFormat))
	return []byte(stamp), nil
}

func (t JsonTime) String() string {
	return time.Time(t).Format(DatetimeFormat)
}
