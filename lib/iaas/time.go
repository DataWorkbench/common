package iaas

import (
	"time"
)

const (
	iaasTimeLayout = "2006-01-02T15:04:05Z"
)

func TimeStringToTimestampSecond(s string) int64 {
	t, err := time.Parse(iaasTimeLayout, s)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}

func TimestampSecondToTimeString(i int64) string {
	return time.Unix(i, 0).Format(iaasTimeLayout)
}
