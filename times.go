package toolbox

import (
	"errors"
	"time"
)

type timesUtils struct {
}

// 时间类
func (us *utils) NewTimesUtils() *timesUtils {
	return &timesUtils{}
}

// 当前时间
func (ti *timesUtils) NowTime() time.Time {
	nowTimeFormat := time.Now().Format("2006-01-02 15:04:05")
	nowTime, _ := time.ParseInLocation("2006-01-02 15:04:05", nowTimeFormat, time.Local)
	return nowTime
}

// 转换时间格式
func (ti *timesUtils) FormatTimeString(t string) (time.Time, error) {
	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	res, _ := time.ParseInLocation(Layout, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation("01/02 03:04:05PM '06 -0700", t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}

	res, _ = time.ParseInLocation("Mon, 2 Jan 2006 15:04:05 -0700 (CST)", t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC1123, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC1123Z, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}

	res, _ = time.ParseInLocation(time.ANSIC, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.UnixDate, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RubyDate, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC822, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC822Z, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC850, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC3339, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.RFC3339Nano, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.Kitchen, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.Stamp, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.StampMilli, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.StampMicro, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	res, _ = time.ParseInLocation(time.StampNano, t, loc)
	if int(res.Unix()) > 0 {
		return res, nil
	}
	return time.Time{}, errors.New("fail")
}
