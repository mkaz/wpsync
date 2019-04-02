package main

import (
	"fmt"
	"strings"
	"time"
)

// WPTime sets up a custom time format to support
// marshalling and unmarshalling json dates coming
// from WordPress

type WPTime struct {
	time.Time
}

const wptLayout = "2006-01-02T15:04:05"

func (wpt *WPTime) UnmarshalJSON(buf []byte) (err error) {
	wpt.Time, err = time.Parse(wptLayout, strings.Trim(string(buf), `"`))
	return
}

func (wpt *WPTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", wpt.Time.Format(wptLayout))), nil
}
