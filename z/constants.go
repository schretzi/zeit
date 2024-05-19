package z

const (
	TextStatusRunning string = ":status:running"
	TextEntry         string = ":entry:"
	TextTask          string = ":task:"
	TextProject       string = ":project:"
	TextImports       string = ":imports:sha1"
)

const ExampleDateIso string = "2006-01-02 15:04 -0700"

const ErrorString string = "%s %+v\n"

const (
	FlagNoColors string = "no-colors"
	FlagDebug    string = "debug"
)

const (
	TFAbsTwelveHour     int = 0
	TFAbsTwentyfourHour int = 1
	TFRelHourMinute     int = 2
	TFRelHourFraction   int = 3
)

const (
	FinishWithMetadata int = 0
	FinishOnlyTime     int = 1
)
