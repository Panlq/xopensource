# 思考

## 1. 如何理解如下注释

// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext.

首个 return err 的 func，将会 cancel error group，是会将所有其他 goroutine 退出嘛？

验证脚本 errgroup_test.go

解答：
errgroup, 有两种初始化方式，一种零值 Group，一种 WithContextGroup

- Group: 启动多个 goroutine，g.Wait 会一直等到所有 g 执行完，然后抛出第一个报错 g 的 err 内容，出错并不会停止其他 g
- WithContextGrou: 带有 cancel contextgroup，同样 g.Wait 会等待所有 g 执行完，然后抛出第一个报错 g 的 err 内容，出错并不会停止其他 g，文档中说到的`The first call to return a non-nil error cancels the group` 其意思是，会 cancel group 持有的 ctx，而对于启动 g 所有者来说，要自己处理 ctx 的取消信号，并非杀手其他 g。
