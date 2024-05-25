# ExpFuzzWordlist
ExpFuzz字典

## 使用方法：

```shell
go build main.go
```

```shell
# 对指定的url进行扫描
./main -u <url>
```

```shell
# 导出所有的poc路径
./main -t <filename>.txt
```

导出的poc可以使用ffuf、dirsearch等工具进行扫描。

## 本项目参考以下项目：

https://github.com/wy876/POC
