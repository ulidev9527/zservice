package zservice

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// 扩展 time 包
type Ztime sql.NullTime

func ZtimeNew(ti time.Time) Ztime {
	return Ztime{Time: ti, Valid: ti.UnixMilli() > 0}
}

func ZtimeNull() Ztime {
	return ZtimeNew(time.Time{})
}

func ZtimeNow() Ztime {
	return ZtimeNew(time.Now())
}

func ZtimeDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Ztime {
	return ZtimeNew(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

func ZtimeUnix(sec int64, nsec int64) Ztime {
	return ZtimeNew(time.Unix(sec, nsec))
}

func ZtimeUnixMilli(msec int64) Ztime {
	return ZtimeNew(time.UnixMilli(msec))
}

func ZtimeUnixMicro(usec int64) Ztime {
	return ZtimeNew(time.UnixMicro(usec))
}

func (ex *Ztime) Scan(value any) error {
	return (*sql.NullTime)(ex).Scan(value)
}

// Value implements the driver Valuer interface.
func (ex Ztime) Value() (driver.Value, error) {
	if !ex.Valid {
		return nil, nil
	}
	return ex.Time, nil
}

func (ex Ztime) After(u time.Time) bool {
	return ex.Time.After(u)
}

func (ex Ztime) Before(u time.Time) bool {
	return ex.Time.Before(u)
}

func (ex Ztime) Compare(u time.Time) int {
	return ex.Time.Compare(u)
}

func (ex Ztime) Equal(u time.Time) bool {
	return ex.Time.Equal(u)
}

func (ex Ztime) IsZero() bool {
	return ex.UnixNano() == 0
}

func (ex Ztime) Date() (year int, month time.Month, day int) {
	return ex.Time.Date()
}

func (ex Ztime) Year() int {
	return ex.Time.Year()
}

func (ex Ztime) Month() time.Month {
	return ex.Time.Month()
}

func (ex Ztime) Day() int {
	return ex.Time.Day()
}

func (ex Ztime) Weekday() time.Weekday {
	return ex.Time.Weekday()
}

func (ex Ztime) ISOWeek() (year int, week int) {
	return ex.Time.ISOWeek()
}

func (ex Ztime) Clock() (h, m, s int) {
	return ex.Time.Clock()
}

func (ex Ztime) Hour() int {
	return ex.Time.Hour()
}

func (ex Ztime) Minute() int {
	return ex.Time.Minute()
}

func (ex Ztime) Second() int {
	return ex.Time.Second()
}

func (ex Ztime) Nanosecond() int {
	return ex.Time.Nanosecond()
}

func (ex Ztime) YearDay() int {
	return ex.Time.YearDay()
}

func (ex Ztime) Add(d time.Duration) Ztime {
	return ZtimeNew(ex.Time.Add(d))
}

func (ex Ztime) Sub(u time.Time) time.Duration {
	return ex.Time.Sub(u)
}

func (ex Ztime) AddDate(years int, months int, days int) Ztime {
	return ZtimeNew(ex.Time.AddDate(years, months, days))
}

func (ex Ztime) UTC() Ztime {
	return ZtimeNew(ex.Time.UTC())
}

func (ex Ztime) Local() Ztime {
	return ZtimeNew(ex.Time.Local())
}

func (ex Ztime) In(loc *time.Location) Ztime {
	return ZtimeNew(ex.Time.In(loc))
}

func (ex Ztime) Location() *time.Location {
	return ex.Time.Location()
}

func (ex Ztime) Zone() (name string, offset int) {
	return ex.Time.Zone()
}

func (ex Ztime) ZoneBounds() (start, end time.Time) {
	return ex.Time.ZoneBounds()
}

func (ex Ztime) Unix() int64 {
	return MaxInt64(ex.Time.Unix(), 0)
}

// 如果时间为 0  则返回 0
func (ex Ztime) UnixMilli() int64 {
	return MaxInt64(ex.Time.UnixMilli(), 0)
}

// 如果时间为 0  则返回 0
func (ex Ztime) UnixMicro() int64 {
	return MaxInt64(ex.Time.UnixMicro(), 0)
}

// 如果时间为 0  则返回 0
func (ex Ztime) UnixNano() int64 {
	return MaxInt64(ex.Time.UnixNano(), 0)
}

func (ex Ztime) MarshalBinary() ([]byte, error) {
	return ex.Time.MarshalBinary()
}

func (ex Ztime) UnmarshalBinary(data []byte) error {
	return ex.Time.UnmarshalBinary(data)
}

func (ex Ztime) GobEncode() ([]byte, error) {
	return ex.Time.GobEncode()
}

func (ex Ztime) GobDecode(data []byte) error {
	return ex.Time.GobDecode(data)
}

func (ex Ztime) MarshalJSON() ([]byte, error) {
	return ex.Time.MarshalJSON()
}

func (ex Ztime) UnmarshalJSON(data []byte) error {
	return ex.Time.UnmarshalJSON(data)
}

func (ex Ztime) MarshalText() ([]byte, error) {
	return ex.Time.MarshalText()
}

func (ex Ztime) UnmarshalText(data []byte) error {
	return ex.Time.UnmarshalText(data)
}

func (ex Ztime) IsDST() bool {
	return ex.Time.IsDST()
}
func (ex Ztime) Truncate(d time.Duration) Ztime {
	return ZtimeNew(ex.Time.Truncate(d))
}

func (ex Ztime) Round(d time.Duration) Ztime {
	return ZtimeNew(ex.Time.Round(d))
}
