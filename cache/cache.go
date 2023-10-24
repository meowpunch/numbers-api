package cache

import "time"

type Cache interface {
	Get(key string) ([]int, bool)
	Set(key string, value []int, duration time.Duration)
}
