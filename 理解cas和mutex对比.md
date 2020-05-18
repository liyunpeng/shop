### 一. CAS概述
Go中的sync/atomic基于CAS做到lock free。  
* 在race condition特别严重的时候，也就是两个goroutine一直在抢着修改同一个对象的时候，CAS的表现和加锁mutex的效率差不多，时高时低。
* 在没有race condition的时候，CAS的表现比mutex耗时要低一个数量级，原因CAS在compare失败的时候回重试，所以在race严重的时候retry多，耗时偶尔炒过mutex。

但是在race没有或者小的时候，效率就体现出来了，因为没有retry。
而mutex还是需要加锁和解锁。耗时是无锁时间的两倍时间。

总结起来，多数情况下能用CAS做到lock free是比较好的，因为没有那么大量的race condition.
所以
* 没有大量race 竞争的情况， 即并发度比较少的情况，用cas.
* 有大量race 竞争的情况，即并发度比较高的情况， 用mutex锁.


#### 二. 被锁住的资源是一个整型变量的并发处理
如果需要被锁住的资源是一个整型变量， 有两种并发处理：
* 用原子操作：
```
atomic.AddInt32(x, -1)
atomic.AddInt32(x, 1)
```
原子操作，是cas方式，是不上锁的。
* 锁的方式：
```
lock.Lock()
*x-- （或*x++）
lock.Unlock()
```
理论上是原子操作因为少了上锁解锁和routine少了睡眠唤醒的动作， 比锁要快很多。
实验证明：不管是有竞争还是无竞争， 原子操作耗时是锁耗时的一半时间：
```
D:\goworkspace\shop\test>go test -v cas_test.go
=== RUN   TestCasMutex
No lock:  2155405
19.0011ms
Mutex lock with condition race:  0
1.3000743s
Atomic CAS with condition race:  0
395.0226ms
Mutex lock without condition race:  10000000
117.0067ms
Atomic CAS without condition race:  10000000
59.0034ms
--- PASS: TestCasMutex (1.89s)
PASS
ok      command-line-arguments  2.525s
```
