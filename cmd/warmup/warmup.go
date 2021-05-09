package main

import (
	"context"
	"github.com/weijunji/go-lottery/internal/lottery"
)

func main()  {
	id := uint64(928)
	lottery.GetLotteryDuration(context.Background(), id)
	lottery.GetRate(context.Background(), id)
	lottery.GetLotteryTimes(context.Background(), id)
}
