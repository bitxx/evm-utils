package dateutil

import "time"

// GetDistanceOfTwoDate
// @param timeStart
// @param timeEnd
// @return int64
func GetDistanceOfTwoDate(timeStart, timeEnd time.Time) int64 {
	before := timeStart.Unix()
	after := timeEnd.Unix()
	return (after - before) / 86400
}

func BeginTime(param time.Time) time.Time {
	timeStr := param.Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return time.Unix(t.Unix(), 0)
}

func EndTimeNum(param time.Time) time.Time {
	timeStr := param.Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return time.Unix(t.Unix()+86399, 999)
}

func ParseTimestampToTime(timestamp int64, location string) (time.Time, error) {
	if location == "" {
		location = "Local"
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timestamp, 0).In(loc), nil
}

func ParseStrToTimestamp(timeStr, location string, flag int) (int64, error) {
	if location == "" {
		location = "Local"
	}
	var t int64
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}
	switch flag {
	case 1:
		t1, err := time.ParseInLocation("2006.01.02 15:04:05", timeStr, loc)
		t = t1.Unix()
		return t, err
	case 2:
		t1, err := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
		t = t1.Unix()
		return t, err
	case 3:
		t1, err := time.ParseInLocation("2006-01-02", timeStr, loc)
		t = t1.Unix()
		return t, err
	case 4:
		t1, err := time.ParseInLocation("2006.01.02", timeStr, loc)
		t = t1.Unix()
		return t, err
	default:
		t1, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
		t = t1.Unix()
		return t, err
	}
}

func ParseStrToTime(timeStr, location string, flag int) (time.Time, error) {
	if location == "" {
		location = "Local"
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	switch flag {
	case 1:
		return time.ParseInLocation("2006.01.02 15:04:05", timeStr, loc)
	case 2:
		return time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
	case 3:
		return time.ParseInLocation("2006-01-02", timeStr, loc)
	case 4:
		return time.ParseInLocation("2006.01.02", timeStr, loc)
	default:
		return time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	}
}

// ConvertToStrByPrt
// @Description:
// @param dateTime
// @param flag
// @return string
func ConvertToStrByPrt(dateTime *time.Time, flag int) string {
	if dateTime == nil {
		return ""
	}
	switch flag {
	case 1:
		return dateTime.Format("2006-01-02")
	case 2:
		return dateTime.Format("2006-01-02 15:04")
	}
	return dateTime.Format("2006-01-02 15:04:05")
}

func ConvertToStr(dateTime time.Time, flag int) string {
	switch flag {
	case 1:
		return dateTime.Format("2006-01-02")
	case 2:
		return dateTime.Format("2006-01-02 15:04")
	case 3:
		return dateTime.Format("2006_01_02_15_04_05")
	case 4:
		return dateTime.Format("2006-01-02T15:04:05.000Z")
	}
	return dateTime.Format("2006-01-02 15:04:05")
}

//
//  ConvertToStrByPrt
//  @Description: 传入的地址是指针，避免外部频繁判断是否为空
//  @param dateTime
//  @param flag
//  @return string
//
/*func ConvertToStrByPrt(dateTime *time.Time, flag int) string {
	if dateTime==nil{
		return ""
	}
	switch flag {
	case 1:
		return dateTime.Format("2006-01-02")
	case 2:
		return dateTime.Format("2006-01-02 15:04")
	}
	return dateTime.Format("2006-01-02 15:04:05")
}*/
