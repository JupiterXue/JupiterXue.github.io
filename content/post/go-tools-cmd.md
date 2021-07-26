---
title: "Go Tools Cmd"
date: 2021-07-26T11:00:46+08:00
draft: true
---

## Go 小工具

### 执行命令行
```go
package main

import (
	"flag"
	"fmt"
	"runtime"
    "os/exec"
    "strings"
)

func main() {
	// flag 包使用方法：flag.Type("flagName",defaultValue,"help message") *Type
	var name = flag.String("name","ls","info: 命令")
	var args = flag.String("args","-h","info: 多个参数")
    flag.Parse()
	fmt.Println(*name)
	fmt.Println(*args)
	cmd := *name + " " + *args
	fmt.Println("Command: %s", cmd)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if runtime.GOOS == "windows" && err != nil {
		fmt.Println("Command error: %s", err)
		out, err = exec.Command("cmd", "/C", cmd).CombinedOutput()
	}
	checkOutput := strings.Contains(string(out), "不是内部或外部命令，也不是可运行的程序")
	if err != nil || checkOutput {
		fmt.Println("Execute Command Error: %s", err)
	} else {
		fmt.Println("Successfully Execute Command Output: %s", out)
	}
}
```
实现在 Windows/Linux 下执行命令行的操作
例如：

```bash
ls -a -l
total 1024
dr-xr-x---. 10 root root    4096 Jul 26 11:12 .
dr-xr-xr-x. 19 root root    4096 Jul 26 01:30 ..
-rw-------.  1 root root    3135 Mar  8 12:44 xxx.cfg
-rw-------   1 root root    8789 Jul 26 11:10 .bash_history
-rw-r--r--.  1 root root      18 Dec 29  2013 .bash_logout
...
```