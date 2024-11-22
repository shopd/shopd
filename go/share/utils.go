package share

import "time"

const VersionFormatSeconds = "2006-01-02-15-04-05"

func NowVersion() string {
	return time.Now().UTC().Format(VersionFormatSeconds)
}
