package model

import "strconv"

type AccountLevel int64

func (al AccountLevel) String() string {
	return strconv.FormatInt(int64(al), 10)
}

const (
	AccountLevel1 AccountLevel = iota + 1
	AccountLevel2
	AccountLevel3
	AccountLevelVip
)
