package vutil

import (
	"strings"
	"time"
)

func TimeParseLocal(in string) (out time.Time, err error) {
	birStrLen := strings.Count(in, "") - 1
	switch birStrLen {
	case 10:
		//2006-01-02
		return timeParse10(in)
	case 19:
		//2006-01-02 15:04:05
		return timeParse19(in)
	default:
	}
	return
}

// yyyy-mm-dd | mm-dd-yyyy 格式
func timeParse10(in string) (out time.Time, err error) {
	chars := []rune(in)
	fc := ""
	for k, char := range chars {
		if char == '-' || char == '/' || char == ':' || char == '.' {
			fc = string(char)
		}
		switch {
		case (k == 0) && ((char <= ('Z') && char >= ('A')) || (char <= ('z') && char >= ('a'))):
			// 这个是  Jun-MM-YYYY 格式
		case (k == 2) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return timeParseMMDDYYYY(in, string(char))
		case (k == 4) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return timeParseYYYYMMDD(in, string(char))
		}
	}
	return timeParseYYYYMMDD(in, fc)
}

// yyyy-mm-dd hh:mm:ss | mm-dd-yyyy hh:mm:ss 格式
func timeParse19(in string) (out time.Time, err error) {
	chars := []rune(in)
	fc := ""
	for k, char := range chars {
		if char == '-' || char == '/' || char == ':' || char == '.' {
			fc = string(char)
		}
		switch {
		case (k == 0) && ((char <= ('Z') && char >= ('A')) || (char <= ('z') && char >= ('a'))):
			// 这个是  Jun-MM-YYYY 格式
		case (k == 2) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return timeParseMMDDYYYYHHMMSS(in, string(char))
		case (k == 4) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return timeParseYYYYMMDDHHMMSS(in, string(char))
		}
	}
	return timeParseYYYYMMDDHHMMSS(in, fc)
}

func timeParseYYYYMMDD(in string, sub string) (out time.Time, err error) {
	layout := "2006" + sub + "01" + sub + "02"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func timeParseMMDDYYYY(in string, sub string) (out time.Time, err error) {
	layout := "01" + sub + "02" + sub + "2006"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func timeParseYYYYMMDDHHMMSS(in string, sub string) (out time.Time, err error) {
	layout := "2006" + sub + "01" + sub + "02 15:04:05"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func timeParseMMDDYYYYHHMMSS(in string, sub string) (out time.Time, err error) {
	layout := "01" + sub + "02" + sub + "2006 15:04:05"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}
