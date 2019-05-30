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
		return TimeParse10(in)
	case 19:
		//2006-01-02 15:04:05
		return TimeParse19(in)
	default:
	}
	return
}

// yyyy-mm-dd | mm-dd-yyyy 格式
func TimeParse10(in string) (out time.Time, err error) {
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
			return TimeParseMMDDYYYY(in, string(char))
		case (k == 4) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return TimeParseYYYYMMDD(in, string(char))
		}
	}
	return TimeParseYYYYMMDD(in, fc)
}

// yyyy-mm-dd hh:mm:ss | mm-dd-yyyy hh:mm:ss 格式
func TimeParse19(in string) (out time.Time, err error) {
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
			return TimeParseMMDDYYYYHHMMSS(in, string(char))
		case (k == 4) && (char == '-' || char == '/' || char == ':' || char == '.'):
			return TimeParseYYYYMMDDHHMMSS(in, string(char))
		}
	}
	return TimeParseYYYYMMDDHHMMSS(in, fc)
}

func TimeParseYYYYMMDD(in string, sub string) (out time.Time, err error) {
	layout := "2006" + sub + "01" + sub + "02"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func TimeParseMMDDYYYY(in string, sub string) (out time.Time, err error) {
	layout := "01" + sub + "02" + sub + "2006"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func TimeParseYYYYMMDDHHMMSS(in string, sub string) (out time.Time, err error) {
	layout := "2006" + sub + "01" + sub + "02 15:04:05"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func TimeParseMMDDYYYYHHMMSS(in string, sub string) (out time.Time, err error) {
	layout := "01" + sub + "02" + sub + "2006 15:04:05"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}
