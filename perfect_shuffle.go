package main

import (
	"math/rand"
	"time"
)

// 问题背景：
//     在开发一个小玩具的时候，发现了这个问题，仔细思考了一下还挺复杂的。
//     由于问题的随机性，导致算法网站上不会存在类似的问题，网上没有可以参考的解。
//     既然没有可参考的，那么就自己手撕吧，好在还是得到了一个自己比较满意的答案。
// 问题描述：
//     打乱数据量大于1的数组，打乱后的每个元素都不在打乱前的原位。
// 思路：
//     该问题有一个难点，即如果有三个元素，ABC，假设算法途中AB互换了位置，那么C就无论如何都无法得到正确的解了。
//     思路1：很简单的思路，一开始很容易往这种思路靠。遍历，给每个数字分配一个新的位置，解决不了难点。
//     思路2：把列表分成随机的两半，分别乱序之后再互换位置。不错的思路，但是会导致奇数的列表无法执行，简单的思路可以想到，随机挑选3个，他们串一下位置，其他的进行偶数算法。
//        例如：12345，随机挑选135进行操作，变成32541，再对24进行操作。但是这会导致奇数算法100%出现3数环，而n数环的概率则直接降低为0。破坏了随机性。
//        *环：即原数组中，有五个数字ABCDE，他们互相在另一个数字的原位置，ABCDE->BCDEA。
//        *严谨点（但不说人话）就是，每个数字的新位置所在的原数字的新位置，再去寻找这个位置下的原数字的新位置时，最终总是能找到自己。
//        *随机算法如果完全随机，那么出现5数环的概率应该不为0。且出现3数环的概率也不应为100%。
//
// 那么为了解决我们的问题，在思路2上，进行了优化。
// 首先偶数情况已经解决了这个问题，我们直接看奇数情况。
// 奇数的情况下，随机挑选一个非首位数据与首位（尾巴也一样）数据进行交换，剩下的执行偶数算法。
// 时间复杂度O(n)，空间复杂度O(n)
func main() {
	input := []int{0, 1, 2, 3, 4}
	res := PerfectShuffle(input)
	for index, value := range res {
		println("index:", index, "value:", value)
	}
}

// PerfectShuffle 完美洗牌，洗牌后所有数字不会在原位
// input 需要洗牌的数组
// res 洗牌后的情况
func PerfectShuffle(input []int) (res []int) {
	length := len(input)
	if length < 2 {
		return input
	}
	rand.Seed(time.Now().UnixNano())
	if isEven(length) {
		return evenPerfectShuffle(input)
	}
	luckyIndex := rand.Intn(length - 1)
	luckyNumber := input[luckyIndex]
	// 不这样写的话，会覆盖input中的数据。
	in := make([]int, length-1)
	copy(in, input[:length-1])
	in[luckyIndex] = input[length-1]
	res = evenPerfectShuffle(in)
	res = append(res, luckyNumber)
	return
}

// 对偶数数组进行完美洗牌
func evenPerfectShuffle(input []int) (res []int) {
	length := len(input)
	if !isEven(length) {
		panic("the input array is not even count!")
	}
	halfData1, halfData2 := randomHalfArray(input)
	halfData1.openChannel()
	halfData2.openChannel()
	res = make([]int, length)
	for i := 0; i < length; i++ {
		if _, ok := halfData1.Data[i]; ok {
			res[i], _ = halfData2.getElement()
			continue
		}
		if _, ok := halfData2.Data[i]; ok {
			res[i], _ = halfData1.getElement()
			continue
		}
		panic("index is not exists")
	}
	return
}

// 随机将数组分成两半。
// input : 需要被分割的数组。
// res1、res2 : 分割后的数组
func randomHalfArray(input []int) (res1, res2 ChannelMap) {
	length := len(input)
	if length == 0 {
		return
	}
	if !isEven(length) {
		panic("the input array is not even count!")
	}
	halfLength := length / 2
	res1map := make(map[int]int, halfLength)
	res2map := make(map[int]int, halfLength)
	randomHalfMap := getRandomHalfIndexList(length, halfLength)
	for index, value := range input {
		if _, ok := randomHalfMap[index]; ok {
			res1map[index] = value
		} else {
			res2map[index] = value
		}
	}
	res1.Data = res1map
	res2.Data = res2map
	return
}

// 给一个值，随机一个map，map中的key是大于等于0，小于入参的随机数。
// limit：随机数的上限（不含此值）
// length：随机数的数量
// res：map，长度为length，其中都是整数，最小值为0，最大值为limit - 1。
func getRandomHalfIndexList(limit, length int) (res map[int]interface{}) {
	var arr []int
	for i := 0; i < limit; i++ {
		arr = append(arr, i)
	}
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	arr = arr[:length]
	res = make(map[int]interface{}, length)
	for _, value := range arr {
		res[value] = nil
	}
	return
}

type ChannelMap struct {
	Data    map[int]int // key：该数字在原数组的位置。value：该数字的值。
	Channel chan int    // 用于从Data中取一个值，并且后续再取不会取到这个值
}

// 这里就是想试试golang的channel，所以才用这种方法来实现。
func (c *ChannelMap) openChannel() {
	c.Channel = make(chan int)
	go func(channelMap *ChannelMap) {
		for _, value := range channelMap.Data {
			channelMap.Channel <- value
		}
	}(c)
}

// 取一个值，每次取的值不会重复。
// 必须先调用openChannel()
func (c *ChannelMap) getElement() (res int, ok bool) {
	if c.Channel == nil {
		panic("channel is not init")
	}
	res, ok = <-c.Channel
	return
}

// 判断一个数字是否是偶数
func isEven(number int) bool {
	return number%2 == 0
}
