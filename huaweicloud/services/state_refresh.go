package services

import (
	"time"
)

// CalculatePollInterval is a method to calculate the query interval time based on the maximum waiting time.
// Timeout:   5m 10m 15m 20m 25m 30m 35m 40m 45m 50m 55m 60m
// Interval: 10s 10s 12s 13s 14s 16s 17s 18s 20s 21s 22s 24s
func CalculatePollInterval(t time.Duration) time.Duration {
	var (
		k    = 220 * time.Second // slope
		base = 450 * time.Second // the minimum timeout is 7.5m
	)

	return ((t-base)/k + 10) * time.Second
}
