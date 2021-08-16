测试元素

> 完成
>
> 测试

> 比如这样



当我们知道上面这些原理后，就可以利用它来对一些场景进行性能优化：

> map[string]int -> map[[12]byte]int

因为 string 底层有指针，所以当 string 作为 map 的 key 时，GC 阶段会扫描整个 map；而数组 `[12]byte` 是一个值类型，不会被 GC 扫描。

我们用两种方法来验证优化效果。

```go
package main

import (
 "fmt"
 "io"
 "net/http"
 _ "net/http/pprof"
)

// var m = map[[12]byte]int{}
var m = map[string]int{}

func init()  {
 for i := 0; i < 1000000; i++ {
  // var arr [12]byte
  // copy(arr[:], fmt.Sprint(i))
  // m[arr] = i

  m[fmt.Sprint(i)] = i
 }
}

func sayHello(wr http.ResponseWriter, r *http.Request) {
 io.WriteString(wr ,"hello")
}

func main() {
 http.HandleFunc("/", sayHello)
 err := http.ListenAndServe(":8000", nil)
 if err != nil {
  fmt.Println(err)
 }
}
```



<h1 align="center">主动</h1>

# 主动触发 GC

这里的测试代码来自文章**《尽量不要在大 map 中保存指针》**[1]：

```go
func MapWithPointer() {
    const N = 10000000
    m := make(map[string]string)
    for i := 0; i < N; i++ {
        n := strconv.Itoa(i)
        m[n] = n
    }
    now := time.Now()
    runtime.GC()     
    fmt.Printf("With a map of strings, GC took: %s\n", time.Since(now))

    // 引用一下防止被 GC 回收掉
    _ = m["0"]
}

func MapWithoutPointer() {
    const N = 10000000
    m := make(map[int]int)
    for i := 0; i < N; i++ {
        str := strconv.Itoa(i)
        // hash string to int
        n, _ := strconv.Atoi(str)
        m[n] = n
    }
    now := time.Now()
    runtime.GC()
    fmt.Printf("With a map of int, GC took: %s\n", time.Since(now))

    _ = m[0]
}

func TestMapWithPointer(t *testing.T) {
    MapWithPointer()
}

func TestMapWithoutPointer(t *testing.T) {
    MapWithoutPointer()
}
```

直接用了 2 个不同类型的 map：前者 key 和 value 都是 string 类型，后者 key 和 value 都是 int 类型。整个 map 大小为 1kw。

测试结果：

```
=== RUN   TestMapWithPointer
With a map of strings, GC took: 150.078ms
--- PASS: TestMapWithPointer (4.22s)
=== RUN   TestMapWithoutPointer
With a map of int, GC took: 4.9581ms
--- PASS: TestMapWithoutPointer (2.33s)
PASS
```

于是验证了 string 相对于 int 这种值类型对 GC 的消耗更大。正如这篇文章的标题所说：

> Go语言使用 map 时尽量不要在 big map 中保存指针。

# 用 pprof 看对象数

第二种方式就是直接开个 pprof 来看 heap profile。这次我们将 string 类型的 key 优化成数组类型：

```bash
package main

import (
 "fmt"
 "io"
 "net/http"
 _ "net/http/pprof"
)

// var m = map[[12]byte]int{}
var m = map[string]int{}

func init()  {
 for i := 0; i < 1000000; i++ {
  // var arr [12]byte
  // copy(arr[:], fmt.Sprint(i))
  // m[arr] = i

  m[fmt.Sprint(i)] = i
 }
}

func sayHello(wr http.ResponseWriter, r *http.Request) {
 io.WriteString(wr ,"hello")
}

func main() {
 http.HandleFunc("/", sayHello)
 err := http.ListenAndServe(":8000", nil)
 if err != nil {
  fmt.Println(err)
 }
}
```

注意，去掉代码里的注释即可将 key 从 string 优化成数组类型。

直接在 init 里构建 map，然后开 pprof 看 profile：

![图片](https://mmbiz.qpic.cn/mmbiz_png/ASQrEXvmx60aJtP1sGia9p9wmE10XxMbZJ9pv6iatUqCT2NgYiaPtfa6GgAlHLZquG9j5odNIF4XJRqVT3CGaZicyg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

key 为 string

![图片](https://mmbiz.qpic.cn/mmbiz_png/ASQrEXvmx60aJtP1sGia9p9wmE10XxMbZYYDhEic1qjHV48pjrhUYn3yAn4w8X5kbfnMDLds73H61LJeeiaaGRiaKg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

key 为数组

对象数从 33w 下降到 1.5w，效果非常明显。

> map 的 key 和 value 要不要在 GC 里扫描，和类型是有关的。数组类型是个值类型，string 底层也是指针。
>
> 不过要注意，key/value 大于 128B 的时候，会退化成指针类型。
>
> 那么问题来了，什么是指针类型呢？**所有显式 *T 以及内部有 pointer 的对像都是指针类型。
>
> ——来自董神的 map [优化文章](https://mp.weixin.qq.com/s?__biz=Mzg5MTYyNzM3OQ==&mid=2247484159&idx=1&sn=94a524d651aabdb9a42b2a5e1d5011ef&scene=21#wechat_redirect)

关于超过 128 字节的情况，源码里也有说明：

```bash
 // Maximum key or elem size to keep inline (instead of mallocing per element).
 maxKeySize  = 128
 maxElemSize = 128
```

# 总结

当 map 的 key/value 是非指针类型时，GC 不会对所有的 bucket 进行扫描。如果线上服务使用了一个超大的 map ，会因此提升性能。

为了不让 overflow 的 bucket 被 GC 错误地回收掉，在 hmap 里用 extra.overflow 指针指向它，从而在三色标记里将其标记为黑色。

如果你用了 key 是 string 类型的 map，并且恰好这些 string 是定长的，那么就可以用 key 为数组类型的 map 来优化它。

通过主动调用 GC 以及开 pprof 都可观察优化效果。

好了，这就是今天全部的内容了~ 我是小X，我们下期再见~

------

欢迎关注曹大的 TechPaper 以及码农桃花源~

### 参考资料

[1]《尽量不要在大 map 中保存指针》: *https://www.jianshu.com/p/5903323a7110*
