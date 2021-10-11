package main

import (
	"fmt"
	"math"
)

// 背景：
//     群里闲聊的时候，看到有人用js自己手写四舍五入代码，其他人在评论效率太低。
//     我想到是不是可以判断小数部分的比特位来提高执行速度，于是便开始写了第一版代码。
//     然而，float的存储结构似乎和我想得不太一样，第一版代码出了点bug，所以就有了这个第二版代码。
// float32结构：
//     float32顾名思义，是由32位组成的浮点数，结构如下：
//     A BBBBBBBB CCCCCCCCCCCCCCCCCCCCCCC
//     A代表符号位，标识整个浮点数是正数还是负数
//     C代表该float32的值，因为浮点数一定是1.xxxxxx这种形式的，所以只保存了xxxxxx部分。
//     B代表指数，最终float32的值是C乘以2的B - 127次方 （减去127是保证可以不用判断符号，float64是1023）
// 代码思路：
//     其实整个float32的值，就是C，B来控制偏移量而已。可以理解B控制小数点在什么位置。
//     同时，二进制的0.1就是十进制的0.5，那么很简单，只需要判断B位移后的右侧那个比特位是不是1就能确定是否要进位了。
func main() {
	println(fmt.Sprintf("%f",roundFloat32(10.50001))) // 11.000000
	println(fmt.Sprintf("%f",roundFloat32(10.49999))) // 10.000000
	println(fmt.Sprintf("%f",roundFloat32(-10.50001))) // -11.000000
	println(fmt.Sprintf("%f",roundFloat32(-10.49999))) // -10.000000
	// 注意，这里会因为浮点数的精度问题，产生异常。
	// 因为二进制的浮点数是有限的，如果标识值的一些位在float32的C中，处于C长度外的位置，那么就会丢失精度。
	// 如果两个数字的C只有这一点是有区别的，那么在内存中存储的比特位也就完全一致，此时他们就是“相等”的。
	println(math.Float32bits(1000.49999) == math.Float32bits(1000.49998)) // true
}

const (
	_bit_11111111111111111111111          = 0x7fffff
	_bit_10000000000000000000000000000000 = 0x80000000
	_bit_01111111100000000000000000000000 = 0x7f800000
	_bit_127                              = 0x7f
)

func roundFloat32(number float32) (res float32) {
	// 获取整个float的bit位，没什么好说的
	bits := math.Float32bits(number)
	// 通过按位与，拿到所有C的值，也就是尾数部分。
	mantissaBits := bits & _bit_11111111111111111111111
	// 通过按位与，拿到所有B的值，也就是指数部分
	exponentBits := (bits & _bit_01111111100000000000000000000000) >> 23
	// 计算出B实际上代表的指数，注意，直接减会因为exponentBits是uint，导致负数变整数，必须转义一下。（不同go版本表现可能不一致）
	exponent := int32(exponentBits) - _bit_127
	// 指数小于-1，代表小数点一定在1的左侧。（因为C的值是.xxxxxx，默认有一个1在前面，如果往左偏移2位，那就变成了0.01xxxxxx，必然是小于0.5的，-0.5同理）
	if exponent < -1 {
		return 0
	}
	// 获取B实际上要偏移多少位
	offset := uint(22 - exponent)
	// 把小数点后面的值都清空（右移到溢出，再移回来，通过这种丢失信息的方式清空）
	res = math.Float32frombits(bits >> (offset + 1) << (offset + 1))
	carry := float32(0.0)
	// mantissaBits&(1<<offset) 找到小数点右侧的那个数字，如果是1，那么这个float32的小数位必然大于0.5。
	// = -1的情况下，1.xxxxxx会变成0.1xxxxxx，必然是1
	if exponent == -1 || mantissaBits&(1<<offset) > 0 {
		// 因为不同的符号进位的方向是不同的，所以这里判断一下符号再进位
		if bits&_bit_10000000000000000000000000000000 > 0 {
			carry = -1
		} else {
			carry = 1
		}
	}
	return res + carry
}
