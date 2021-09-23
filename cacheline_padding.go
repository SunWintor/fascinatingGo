package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Obj struct {
	num  uint64
	num2 uint64
}

type Obj2 struct {
	num  uint64
	temp [7]uint64 // 注释掉之后，执行速度一致
	num2 uint64
}

var (
	runTimes = int64(7000000)
    wg sync.WaitGroup
	obj = Obj{}
	obj2 = Obj2{}
)

// cpu分为L1 L2 L3三级缓存，其中L1和L2每个核各有一个，而L3则所有核共享。
// 当L1或者L2更新时，将会触发同步机制MESI（缓存一致性协议） *注意，缓存一致性协议不仅这一个。
// 同步是有开销的
// 缓存中每行数据长度为64byte，如果修改这64byte中的任意数据，都会触发整条数据的缓存同步。
// Obj中，uint64的大小是8byte，远小于64byte，导致更新num的时候，会对num2也进行同步
// 而Obj2中，一共是有8byte的数据+7*8byte的缓冲，保证了num和num2一定在不同的数据行中
// 更新num不会触发num2的同步，所以Obj2中，temp插在num和num2中间，看似无影响，实际提高了执行速度。
func main() {
	// runTime1 = 194097000
	// runTime2 = 29778000
	testCacheLine()
}

func testCacheLine() {
	testRun(func(){atomic.AddUint64(&obj.num, 1)}, func(){atomic.AddUint64(&obj.num2, 1)}, "runTime1")
	testRun(func(){atomic.AddUint64(&obj2.num, 1)}, func(){atomic.AddUint64(&obj2.num2, 1)}, "runTime2")
}

func testRun(addFunc, addFunc2 func(), funcName string) {
	runTime := time.Now().UnixNano()
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := int64(0); i < runTimes; i++ {
			addFunc()
		}
	}()
	go func() {
		defer wg.Done()
		for i := int64(0); i < runTimes; i++ {
			addFunc2()
		}
	}()
	wg.Wait()
	println(fmt.Sprintf("%s = %+v", funcName, time.Now().UnixNano()-runTime))
}