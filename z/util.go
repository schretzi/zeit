package z

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

var fractional bool

func fmtDuration(dur time.Duration) string {
	if fractional {
		return decimal.NewFromFloat(dur.Hours()).StringFixed(2)
	} else {
		dur = dur.Round(time.Second)
		hour := dur / time.Hour
		dur -= hour * time.Hour
		minute := dur / time.Minute
		return fmt.Sprintf("%d:%02d", hour, minute)
	}
}

func fmtHours(hours decimal.Decimal) string {
	if fractional {
		return hours.StringFixed(2)
	} else {
		return fmt.Sprintf(
			"%s:%02s",
			hours.Floor(), // hours
			hours.Sub(hours.Floor()).
				Mul(decimal.NewFromFloat(.6)).
				Mul(decimal.NewFromInt(100)).
				Floor())
	}
}
