package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func main() {
	// 顺序的执行时间几乎总是比乱序的执行时间短的
	// 其原因是，cpu的执行速度很快，经常会有空闲
	// 当遇到分支时，cpu会选择其中一个分支先进行计算，如果选错了，则退回重新计算。如果选对了，则节省了大量时间。
	// 选择有两种方式，动态选择和静态选择。
	// 静态的方式就是随机猜测。
	// 动态的会根据之前的选择结果，来判断本次到底走哪个分支。
	// 具体来说就是有四个档位，ABCD，AB的情况走顺序分支，CD的情况走跳转分支
	// A：前两次都是顺序分支
	// B：A的情况下，遇到了一次跳转分支，或者C的情况下，遇到了一次顺序分支。
	// C：D的情况下，遇到了一次顺序分支，或者B的情况下，遇到了一次跳转分支。
	// D：前两次都是跳转分支
	// 可能有点绕，可以想象成ABCD是一个数组，每遇到一次跳转分支就向D移动1格，每遇到一次顺序分支就向A一格。
	// 那么testTime执行的时候，顺序比乱序执行时间短就很好解释了，因为顺序的情况下，CPU几乎总是会"猜对"下一次要走的分支。
	testTime(3000000)
}

func testTime(numLimit int) {
	var list []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numLimit; i++ {
		list = append(list, rand.Intn(numLimit))
	}

	start := time.Now().UnixNano()
	process(list)
	println(fmt.Sprintf("乱序，执行时间%d", time.Now().UnixNano() - start)) // 乱序，执行时间47125500

	sort.Ints(list)
	start = time.Now().UnixNano()
	process(list)
	println(fmt.Sprintf("顺序，执行时间%d", time.Now().UnixNano() - start)) // 顺序，执行时间7008000
}

func process(list []int)  {
	halfLength := len(list) / 2
	for _, v := range list {
		if v > halfLength {
			continue
		}
		math.Abs(float64(v)) // 这一行代码需要消耗cpu，不能过于简单。且使用函数时，不能使用带有if的，不然就打乱了cpu对于分支的预测。
	}
}
