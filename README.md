# selpg.go
<<<<<<< HEAD


=======

### 1. 设计思路
这个selpg的CLI实现的思路很简单：

（1）读取命令并解析对应的参数，初始化程序运行的变量和参数

（2）读入目标文件（通过重定向的方式或手动输入的方式）

（3）根据预设的参数对文件执行读写操作

（4）将对应输出输出到相应的文件中

### 2. 代码分析
##### （1）参数解析部分
```go
// 使用pflag进行参数的定义
flag.IntVarP(&s, "s", "s", 1, "start page")
flag.IntVarP(&e, "e", "e", 1, "end page(e >= s)")
flag.IntVarP(&lines, "l", "l", 72, "page length")
flag.BoolVarP(&changeLine, "f", "f", false, "change page by symbol '\\f'")
flag.Lookup("f").NoOptDefVal = "true"
flag.StringVarP(&pipe, "d", "d", "", "printer")
flag.Usage = func() {
    fmt.Fprintln(os.Stderr, "Format: ./selpg [-s n] [-e n] [-f | -l n] [-d p] [filename] [other options]")
    flag.PrintDefaults()
}

flag.Parse()

// 判断输入是文件重定向还是用户自主输入
if flag.NArg() == 1 {
    _, file_err := os.Stat(flag.Args()[0])
    if file_err != nil {
        return file_err
    }
    input_file = flag.Args()[0]
} else if flag.NArg() > 1 {
    flag.Usage()
    err = errors.New("ERROR: too many args")
}
```
##### （2）文件读写部分
利用os.Open和bufio读取文件内容部分
```go
f, err := os.Open(filename)
if filename == "" {
    f = os.Stdin
} else {
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
defer f.Close()
buf := bufio.NewReader(f)
```
处理以行数作为换页标志的情况
```go
for page = 1; page <= e; page += 1 {
    count = 1
    for {
        if count > lines {
            break
        }
        line, err := buf.ReadString('\n')
        count += 1
        if page < s {
            if err != nil {
                return err
            } else {
                continue
            }
        }
        // 输出
        if err == nil {
            out.Write([]byte(line))
        } else if err == io.EOF {
            out.Write([]byte(line))
        } else {
            return err
        }
    }
```
处理以换页符作为换页标志的情况
```go
page = 0
for page < e {
    line, err := buf.ReadString('\f')
    page += 1
    if page < s {
        if err != nil {
            return err
        } else {
            continue
        }
    }
    // 输出
    if err == nil {
        out.Write([]byte(line))
    } else if err == io.EOF {
        out.Write([]byte(line))
    } else {
        return err
    }
}
```
##### （3）main函数
```go
// 获取并处理参数
err = getArgs()
if err == nil {
	// 输出位置
    out = os.Stdout
    // 假如打印机选项不为空，则执行sendToPrinter()
    // 否则，直接读取文件并按照用户要求输出
    if pipe != "" {
        sendToPrinter()
    } else {
        err = readFile(input_file)
    }
}
```

### 代码测试

这部分我写在了CSDN博客上，详细内容可点击此[链接](https://blog.csdn.net/think_A_lot/article/details/82988219)查看
>>>>>>> 566c1579e5d9130f471dbda8a1ef8a8e4106ce05
