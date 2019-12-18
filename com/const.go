package com

import (
	"github.com/hanjingo/util"
)

const START_TIME_STR string = "2019-01-01 00:00:00"

var START_TIME = util.TimeStampToTime(START_TIME_STR)

//未知编码id
const UNKNOWN_CODECID uint8 = 0
