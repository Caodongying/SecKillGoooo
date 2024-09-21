package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// create Redis connection pool
func init() {
	RedisPool = &redis.Pool{
		MaxIdle: 10, // maximum connection at initialization
		MaxActive: 0, // maximum connection numbers, no limitation
		IdleTimeout: 300 * time.Second, // timeout, terminate the connection
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:9108")
		},
	}
}

// sec kill
var lock sync.Mutex
var listName = "seckill" // assume that seckill only serves ten users
const maximumSell = 10

func main() {
	createQueue()
}

// mimic sec kill
func createQueue() {
	var wg sync.WaitGroup
	redisDB := RedisPool.Get()
	defer redisDB.Close()

	userIDs := []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
	userCount := len(userIDs)
	ch := make(chan int, userCount)

	f := func(id int, done func()) {
		defer func() {
			done()
		}()
		lock.Lock()
		defer lock.Unlock()

		listNameLen, _ := redis.Int64(redisDB.Do("LLen", listName)) // ignore the error for now
		if listNameLen < maximumSell {
			redisDB.Do("RPUSH", listName, fmt.Sprintf("%d@%v", id, time.Now()))
			fmt.Println(id, "Successfully Sec Killed!")
		} else {
			fmt.Println("Event finished :()")
		}

		ch <- id
	}

	wg.Add(userCount)

	for _, v := range userIDs {
		go f(v, wg.Done)
	}

	for i := 0; i < userCount; i++ {
		<- ch
	}

	close(ch)
	wg.Wait()
}