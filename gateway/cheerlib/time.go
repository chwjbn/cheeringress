package cheerlib

import "time"

func TimeGetNow() string {

	timeFmt := "2006-01-02 15:04:05"
	return time.Now().Format(timeFmt)
}

func TimeGetTime(t time.Time) string {
	timeFmt := "2006-01-02 15:04:05"
	return t.Format(timeFmt)
}

func TimeStrToTime(timeStr string) time.Time {

	timeFmt := "2006-01-02 15:04:05"

	timeData, timeErr := time.ParseInLocation(timeFmt, timeStr, time.Local)

	if timeErr != nil {
		timeData = time.Unix(0, 0)
	}

	return timeData
}

func TimeTimestamp() int64 {
	return time.Now().Unix()
}

func TimeUtcTimestamp() int64 {
	return time.Now().UTC().Unix()
}
