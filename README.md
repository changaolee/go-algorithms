## go-algorithms

项目包含 LeetCode、AcWing 题解及常用算法模板。

### 代码生成工具

> 请在 cmd 目录下执行。

```shell
# 生成力扣指定一个题目的代码
go test -run GenLCProblemCode -args {ID}

# 生成力扣指定一场周赛/双周赛的代码
go test -run GenLCContestCode -args {Type} {ID}
```