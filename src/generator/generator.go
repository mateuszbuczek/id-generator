package generator

import (
	"fmt"
	"math"
	"time"
)

type Worker struct {
	GeneratedIds *chan int64
	WorkerId     int64
	ThreadCount  int
}

var START_EPOCH int64 = 1646922484

var ID_BIT_SIZE = 64
var EPOCH_BIT_SIZE int = 40
var WORKER_ID_BIT_SIZE int = 4
var THREAD_ID_BIT_SIZE int = 4
var COUNTER_BIT_SIZE int = 16

func (w *Worker) Generate(num int64) {
	var numPerThread = num/int64(w.ThreadCount) + 1

	if float64(numPerThread) > math.Pow(2, float64(COUNTER_BIT_SIZE)) {
		panic(fmt.Sprintf("invalid config: COUNTER_BIT_SIZE: %f, numOfIdsPerThread; %d", math.Pow(2, float64(COUNTER_BIT_SIZE)), numPerThread))
	}

	for threadId := 0; threadId < w.ThreadCount; threadId++ {
		go generate(int64(threadId), w.WorkerId, w.GeneratedIds, numPerThread)
	}
}

func generate(threadId int64, workerId int64, output *chan int64, num int64) {
	nowMilli := time.Now().UnixMilli()
	epoch := nowMilli - START_EPOCH

	for i := 0; i < int(num); i++ {
		id := epoch << (ID_BIT_SIZE - EPOCH_BIT_SIZE)
		id |= workerId << (ID_BIT_SIZE - EPOCH_BIT_SIZE - WORKER_ID_BIT_SIZE)
		id |= threadId << (ID_BIT_SIZE - EPOCH_BIT_SIZE - WORKER_ID_BIT_SIZE - THREAD_ID_BIT_SIZE)
		id |= int64(i)
		*output <- id
	}
}
