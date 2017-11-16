package util

import (
	"log"
)

// TokenOrTicketRefreshBufferPeriod is the duration before token/ticket
// actually expires that our cache is removed and server request is
// able to be triggered.
var TokenOrTicketRefreshBufferPeriod = 1500

// {Min|Max}imumCacheLife - if set to a positive number, this number will
// override the `expires' field returned by WeChat server and act as the
// actual TTL of entries in our database.
var (
	MinimumCacheLife = 0
	MaximumCacheLife = 0
)

const absoluteLowerLimit = 15

// CalculateTTL tailors the expire value from WeChat server to our configuration.
func CalculateTTL(expires int) int {
	log.Println("TTL from WeChat server:", expires)
	expires = expires - TokenOrTicketRefreshBufferPeriod
	if !(MinimumCacheLife > 0 && MaximumCacheLife > 0 && MinimumCacheLife > MaximumCacheLife) {
		if MinimumCacheLife > 0 && expires < MinimumCacheLife {
			expires = MinimumCacheLife
		}
		if MaximumCacheLife > 0 && expires > MaximumCacheLife {
			expires = MinimumCacheLife
		}
	}
	if expires < absoluteLowerLimit {
		expires = absoluteLowerLimit
	}
	log.Println("TTL after tailoring:", expires)
	return expires
}
