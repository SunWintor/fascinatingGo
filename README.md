# fascinatingGo
收集一些平时发现的很炫酷的东西，以及一些我自己有意思的想法。

如果想看更加有意思的，可以解读的东西更多的，可以看我另一个项目：https://github.com/SunWintor/SWT

---

* perfect_shuffle.go 完美洗牌，排序后每个元素都不在原位置。原本以为是个简单的问题，越写越发现其实并不容易。
* bit_round_float32.go 关于浮点数四舍五入的实现（做完之后和golang官方库的比了一下，速度几乎一样，看了一下源码，似乎是用类似的思路实现的）
* cacheline_padding.go 关于CPU对于缓存一致性对于性能影响的猜想和验证。
* disorder.go 关于CPU指令重排序。
