package limit

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

//单机限流器的使用

func TestLimit(t *testing.T) {
	//每秒生成5个 令牌桶最多存在10个
	limit := rate.NewLimiter(5, 10)
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	//获取10个令牌，当超过令牌桶的阈值或者获取不到则会进行阻塞，主要用于控制速度，并且存在阻塞ddl，超过ddl则丢弃
	err := limit.WaitN(ctx, 10)
	_ = err

	//获取1个令牌
	err = limit.Wait(ctx)
	_ = err

	//非阻塞式，返回是否允许，允许则会扣除，主要用来丢弃流量
	can := limit.AllowN(time.Now(), 10)
	_ = can

	//主要用于限速
	r := limit.ReserveN(time.Now(), 1)
	if !r.OK() {
		// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		return
	}
	time.Sleep(r.Delay())

	//会归还获取到的令牌
	r.Cancel()
	//dosomething()

}
