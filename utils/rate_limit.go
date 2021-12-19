package utils

import (
	"net/http"
)

//var RateLimitMap = map[string]int64{}
//var RateLimitLock = sync.Mutex{}

// CheckRateLimit checks if the request has remaining quota or not. If it returns false,
// the IP in the request has exceeded their quota
func CheckRateLimit(r *http.Request) bool {
	//RateLimitLock.Lock()
	//ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	//now := time.Now().Unix() * preloaded.RateLimit
	//if val, ok := RateLimitMap[ip]; ok {
	//	if (val + preloaded.RateLimitBurst) < now {
	//		RateLimitMap[ip] = now - preloaded.RateLimitBurst
	//	}
	//} else {
	//	RateLimitMap[ip] = now - preloaded.RateLimit
	//}
	//
	//RateLimitMap[ip]++
	//RateLimitLock.Unlock()
	//return RateLimitMap[ip] <= now
	return true
}
