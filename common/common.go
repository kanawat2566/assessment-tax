package common

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func MsgWithNumber(msg string, num float64) string {
	p := message.NewPrinter(language.English)
	n := number.Decimal(num, number.MaxFractionDigits(1))
	return p.Sprintf(msg+" %v", n)
}

func MsgWithInt(msg string, num int) string {
	return MsgWithNumber(msg, float64(num))
}
