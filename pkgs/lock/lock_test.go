package lock

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"testing"
	"time"
)

func TestDistributeLock(t *testing.T) {
	assert := assert.New(t)
	lock := NewDistributeLock(context.Background(), "test", 99999999)
	name := "test_lock:99999999"
	assert.Equal(lock.name, name)
	// test lock
	assert.True(lock.Lock(time.Second * 5), "lock should success")
	assert.False(lock.Lock(time.Second * 5), "Should failed because previous lock")
	lock.UnLock()
	// test lock expire
	assert.True(lock.Lock(time.Millisecond), "lock should success")
	time.Sleep(time.Millisecond * 5)
	assert.True(lock.Lock(time.Second), "lock should expire")
	utils.GetRedis().Del(context.Background(), name)
}
