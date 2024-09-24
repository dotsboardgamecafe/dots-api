package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

//GetTimeLocationWIB get WIB location
func GetTimeLocationWIB() *time.Location {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	return wib
}

// ToUTCfromGMT7 ...
func ToUTCfromGMT7(strTime string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	date, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, location)
	if err != nil {
		fmt.Printf("\nerror when parse strTime [%s] -> err: %v\n", strTime, err)
		return time.Now(), err
	}

	return date.In(time.UTC), nil
}

// FromUTCLocationToGMT7 ...
func FromUTCLocationToGMT7(date time.Time) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	return date.In(location), nil
}

// FromGMT7LocationUTCMin7 ...
func FromGMT7LocationUTCMin7(date time.Time) (time.Time, error) {
	date = date.Add(time.Hour * -7)
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	date = date.In(location)

	return date.In(time.UTC), nil
}

// TimeElapsed ...
func TimeElapsed(param time.Time) string {
	var text string
	var parts []string

	now := time.Now()

	currentYear, currentMonth, currentDay := now.Date()
	currentHour, currentMinute, currentSecond := now.Clock()

	paramYear, paramMonth, paramDay := param.Date()
	paramHour, paramMinute, paramSecond := param.Clock()

	year := math.Abs(float64(int(currentYear - paramYear)))
	month := math.Abs(float64(int(currentMonth - paramMonth)))
	day := math.Abs(float64(int(currentDay - paramDay)))
	hour := math.Abs(float64(int(currentHour - paramHour)))
	minute := math.Abs(float64(int(currentMinute - paramMinute)))
	second := math.Abs(float64(int(currentSecond - paramSecond)))
	week := math.Floor(day / 7)

	s := func(x float64) string {
		if int(x) == 1 {
			return ""
		}
		return "s"
	}

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+" Year"+s(year))
	}
	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+" Month"+s(month))
	}
	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+" Week"+s(week))
	}
	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+" Day"+s(day))
	}
	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+" Hour"+s(hour))
	}
	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+" Minute"+s(minute))
	}
	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+" Second"+s(second))
	}
	if len(parts) == 0 {
		return "Now"
	}
	if now.After(param) {
		text = " Ago"
	} else {
		text = " After"
	}

	return parts[0] + text
}

// This function counts the
// number of leap years
// since the starting of time
// to the current year that
// is passed
func leapYears(date time.Time) (leaps int) {

	// returns year, month,
	// date of a time object
	y, m, _ := date.Date()

	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}

// The function calculates the
// difference between two dates and times
// and returns the days, hours, minutes,
// seconds between two dates

func GetTimeDifference(a, b time.Time) (days, hours, minutes, seconds int) {

	// month-wise days
	monthDays := [12]int{31, 28, 31, 30, 31,
		30, 31, 31, 30, 31, 30, 31}

	// extracting years, months,
	// days of two dates
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()

	// extracting hours, minutes,
	// seconds of two times
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()

	// totalDays since the
	// beginning = year*365 + number_of_days
	totalDays1 := y1*365 + d1

	// adding days of the months
	// before the current month
	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}

	// counting leap years since
	// beginning to the year "a"
	// and adding that many extra
	// days to the totaldays
	totalDays1 += leapYears(a)

	// Similar procedure for second date
	totalDays2 := y2*365 + d2

	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}

	totalDays2 += leapYears(b)

	// Number of days between two days
	days = totalDays2 - totalDays1

	// calculating hour, minutes,
	// seconds differences
	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1

	// if seconds difference goes below 0,
	// add 60 and decrement number of minutes
	if seconds < 0 {
		seconds += 60
		minutes--
	}

	// performing similar operations
	// on minutes and hours
	if minutes < 0 {
		minutes += 60
		hours--
	}

	// performing similar operations
	// on hours and days
	if hours < 0 {
		hours += 24
		days--
	}

	return days, hours, minutes, seconds
}
