package zservice

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// 扩展 time 包
type Time sql.NullTime

func NewTime(ti time.Time) Time {
	return Time{Time: ti, Valid: ti.UnixMilli() > 0}
}

func TimeNull() Time {
	return NewTime(time.Time{})
}

func TimeNow() Time {
	return NewTime(time.Now())
}

func TimeDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return NewTime(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

func TimeUnix(sec int64, nsec int64) Time {
	return NewTime(time.Unix(sec, nsec))
}

func TimeUnixMilli(msec int64) Time {
	return NewTime(time.UnixMilli(msec))
}

func TimeUnixMicro(usec int64) Time {
	return NewTime(time.UnixMicro(usec))
}

func (ex *Time) Scan(value any) error {
	return (*sql.NullTime)(ex).Scan(value)
}

// Value implements the driver Valuer interface.
func (ex Time) Value() (driver.Value, error) {
	if !ex.Valid {
		return nil, nil
	}
	return ex.Time, nil
}

func (ex Time) After(u time.Time) bool {
	return ex.Time.After(u)
}

func (ex Time) AfterNow() bool {
	return ex.Time.After(time.Now())
}

func (ex Time) Before(u time.Time) bool {
	return ex.Time.Before(u)
}

func (ex Time) BeforeNow() bool {
	return ex.Time.Before(time.Now())
}

func (ex Time) Compare(u time.Time) int {
	return ex.Time.Compare(u)
}

func (ex Time) Equal(u time.Time) bool {
	return ex.Time.Equal(u)
}

func (ex Time) IsZero() bool {
	return ex.UnixNano() == 0
}

func (ex Time) Date() (year int, month time.Month, day int) {
	return ex.Time.Date()
}

func (ex Time) Year() int {
	return ex.Time.Year()
}

func (ex Time) Month() time.Month {
	return ex.Time.Month()
}

func (ex Time) Day() int {
	return ex.Time.Day()
}

func (ex Time) Weekday() time.Weekday {
	return ex.Time.Weekday()
}

func (ex Time) ISOWeek() (year int, week int) {
	return ex.Time.ISOWeek()
}

func (ex Time) Clock() (h, m, s int) {
	return ex.Time.Clock()
}

func (ex Time) Hour() int {
	return ex.Time.Hour()
}

func (ex Time) Minute() int {
	return ex.Time.Minute()
}

func (ex Time) Second() int {
	return ex.Time.Second()
}

func (ex Time) Nanosecond() int {
	return ex.Time.Nanosecond()
}

func (ex Time) YearDay() int {
	return ex.Time.YearDay()
}

func (ex Time) Add(d time.Duration) Time {
	return NewTime(ex.Time.Add(d))
}

func (ex Time) Sub(u time.Time) time.Duration {
	return ex.Time.Sub(u)
}

func (ex Time) AddDate(years int, months int, days int) Time {
	return NewTime(ex.Time.AddDate(years, months, days))
}

func (ex Time) UTC() Time {
	return NewTime(ex.Time.UTC())
}

func (ex Time) Local() Time {
	return NewTime(ex.Time.Local())
}

func (ex Time) In(loc *time.Location) Time {
	return NewTime(ex.Time.In(loc))
}

func (ex Time) Location() *time.Location {
	return ex.Time.Location()
}

func (ex Time) Zone() (name string, offset int) {
	return ex.Time.Zone()
}

func (ex Time) ZoneBounds() (start, end time.Time) {
	return ex.Time.ZoneBounds()
}

func (ex Time) Unix() int64 {
	return MaxInt64(ex.Time.Unix(), 0)
}

// 如果时间为 0  则返回 0
func (ex Time) UnixMilli() int64 {
	return MaxInt64(ex.Time.UnixMilli(), 0)
}

// 如果时间为 0  则返回 0
func (ex Time) UnixMicro() int64 {
	return MaxInt64(ex.Time.UnixMicro(), 0)
}

// 如果时间为 0  则返回 0
func (ex Time) UnixNano() int64 {
	return MaxInt64(ex.Time.UnixNano(), 0)
}

func (ex Time) MarshalBinary() ([]byte, error) {
	return ex.Time.MarshalBinary()
}

func (ex Time) UnmarshalBinary(data []byte) error {
	return ex.Time.UnmarshalBinary(data)
}

func (ex Time) GobEncode() ([]byte, error) {
	return ex.Time.GobEncode()
}

func (ex Time) GobDecode(data []byte) error {
	return ex.Time.GobDecode(data)
}

func (ex Time) MarshalJSON() ([]byte, error) {
	return ex.Time.MarshalJSON()
}

func (ex Time) UnmarshalJSON(data []byte) error {
	return ex.Time.UnmarshalJSON(data)
}

func (ex Time) MarshalText() ([]byte, error) {
	return ex.Time.MarshalText()
}

func (ex Time) UnmarshalText(data []byte) error {
	return ex.Time.UnmarshalText(data)
}

func (ex Time) IsDST() bool {
	return ex.Time.IsDST()
}
func (ex Time) Truncate(d time.Duration) Time {
	return NewTime(ex.Time.Truncate(d))
}

func (ex Time) Round(d time.Duration) Time {
	return NewTime(ex.Time.Round(d))
}
