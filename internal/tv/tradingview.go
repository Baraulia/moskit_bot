package tv

type Interval string
type LineType string

const DefaultInterval = Interval1D

const (
	IntervalAll   Interval = "All"
	Interval1Min  Interval = "1m"
	Interval5Min  Interval = "5m"
	Interval15Min Interval = "15m"
	Interval30Min Interval = "30m"
	Interval60Min Interval = "1h"
	Interval1H    Interval = "1h"
	Interval2H    Interval = "2h"
	Interval4H    Interval = "4h"
	Interval1D    Interval = "1d"
	Interval1W    Interval = "1w"
	Interval1M    Interval = "1M"
)

const (
	BearOrderBlock   LineType = "BearOB"
	BullOrderBlock   LineType = "BullOB"
	BearBreakerBlock LineType = "BearBB"
	BullBreakerBlock LineType = "BullBB"
	SupportLevel     LineType = "SupLev"
	ResistanceLevel  LineType = "ResLev"
)

func (l LineType) PrepareColumn(timeframe Interval) string {
	return string(l) + "|" + string(timeframe)
}

func (i Interval) ForColumn() string {
	switch i {
	case Interval1Min:
		return "|1"
	case Interval5Min:
		return "|5"
	case Interval15Min:
		return "|15"
	case Interval30Min:
		return "|30"
	case Interval1H:
		return "|60"
	case Interval2H:
		return "|120"
	case Interval4H:
		return "|240"
	case Interval1W:
		return "|1W"
	case Interval1M:
		return "|1M"
	}
	return ""
}

func ForResponse(i string) string {
	switch i {
	case "1":
		return " 1min"
	case "5":
		return " 5min"
	case "15":
		return " 15min"
	case "30":
		return " 30min"
	case "60":
		return " 1H"
	case "120":
		return " 2H"
	case "240":
		return " 4H"
	case "1W":
		return " 1W"
	case "1M":
		return " 1Month"
	}
	return ""
}
