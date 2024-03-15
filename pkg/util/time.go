package util

import (
	"bypctl/pkg/cmd"
	"fmt"
	"runtime"
	"time"
)

// FormatDate .
func FormatDate(tm time.Time) string {
	return tm.Format("2006-01-02")
}

// FormatTime .
func FormatTime(tm time.Time) string {
	return tm.Format("2006-01-02 15:04:05")
}

// TimestampToTime .
func TimestampToTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// FormatStr2Time .
func FormatStr2Time(strTime string) (time.Time, error) {
	formatTime, err := time.Parse("2006-01-02 15:04:05", strTime)
	if err != nil {
		return time.Now(), err
	}
	return formatTime, nil
}

func UpdateSystemTime(dateTime time.Time) error {
	goos := runtime.GOOS
	switch goos {
	case "linux":
		stdout, err := cmd.Execf(`%s date -s "%s"`, cmd.SudoHandleCmd(), dateTime.String())
		if err != nil {
			return fmt.Errorf("update system time failed, stdout: %s, err: %v", stdout, err)
		}
	case "darwin":
		stdout, err := cmd.Execf(`%s date "%s"`, cmd.SudoHandleCmd(), dateTime.Format("010215042006.05"))
		if err != nil {
			return fmt.Errorf("update system time failed, stdout: %s, err: %v", stdout, err)
		}
	default:
		return fmt.Errorf("the current system architecture %v does not support synchronization", goos)
	}
	return nil
}

func LoadSystemTimeZone() string {
	loc := time.Now().Location()
	if _, err := time.LoadLocation(loc.String()); err != nil {
		return "Asia/Shanghai"
	}
	return loc.String()
}

func UpdateSystemTimeZone(timezone string) error {
	goos := runtime.GOOS
	switch goos {
	case "linux":
		stdout, err := cmd.Execf(`%s timedatectl set-timezone "%s"`, cmd.SudoHandleCmd(), timezone)
		if err != nil {
			return fmt.Errorf("update system time zone failed, stdout: %s, err: %v", stdout, err)
		}
	case "darwin":
		stdout, err := cmd.Execf(`%s systemsetup -settimezone "%s"`, cmd.SudoHandleCmd(), timezone)
		if err != nil {
			return fmt.Errorf("update system time zone failed, stdout: %s, err: %v", stdout, err)
		}
	default:
		return fmt.Errorf("the current system architecture %v does not support synchronization", goos)
	}
	return nil
}
