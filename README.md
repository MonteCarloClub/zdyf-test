# 国家重点研发计划测试客户端

## 0. 拉取与构建

```bash
git clone https://github.com/MonteCarloClub/zdyf-test.git
cd zdyf-test
make all # 或 make zdyf-test
```

## 1. `/abe/dabe/user2`API登录性能

此客户端实现了2种测试方法：在客户端并发连接，在服务端模拟并发连接。

### 1.1. `/abe/dabe/user2_dry_run`：在客户端并发连接

在客户端发送10000000个并发连接，每约1000000个连接完成后打印日志。

```bash
./bin/zdyf-test user2DryRun -l 1000000
```

### 1.2. `/abe/dabe/user2_batch_dry_run`：在服务端模拟并发连接

在客户端发送模拟请求。

```bash
./bin/zdyf-test user2BatchDryRun
```
