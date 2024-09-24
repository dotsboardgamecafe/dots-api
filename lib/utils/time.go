package utils

import "time"

func AddTime(originTime time.Time, timeDuration time.Duration) time.Time {
	return originTime.Add(timeDuration)
}

func ReduceTime(originTime time.Time, timeDuration time.Duration) time.Time {
	return originTime.Add(-timeDuration)
}
