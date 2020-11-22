# cloudgo
A Go web server with static file handling, json response and template rendering.

## Considerations for Choosing A Web Framework
此次使用了 [gin](https://github.com/gin-gonic/gin) 作为 Web 框架，主要考虑有：
1. 是否“好写”：对于路由器，`gin` 提供了类似于 Node.js Express 路由风格的写法，对于有过 Node.js Web 开发经验的人来说极易上手。
2. 是否有足够的性能：作为“重量型”的框架，`gin` 显然没有极致性能，但也提供了我们认为的较好的性能。这点无论是官方的 Benchmark 还是在我们后续自己使用 Apache Benchmark 进行的测试中都可以看出来。
3. 流行的程度：从 Stars 数上看，`gin` 无疑碾压了其他的框架（[mingrammer/go-web-framework-stars](https://github.com/mingrammer/go-web-framework-stars)）。使用框架的人数多，社区的活跃程度高，意味着遇到问题时可以在 Stack Overflow 等网站上找到相应解决方案的几率也就更大，减少开发时间成本。
4. 社区生态：对于“重量型”框架，是否有良好的“插件”支持以简化各种常见 Web 开发任务的实现也是重要考虑因素之一。`gin` 有一个中间件 collection repository，列举了社区提供的一些 `gin` 中间件：[gin-gonic/contrib](https://github.com/gin-gonic/contrib)
5. 文档是否可用：粗略来看 `gin` 的文档应该没有比较坑的地方。

## Testing

测试环境：

```
GitHub Codespaces Standard Instance (Linux), with 4 cores, 8 GB RAM and 32 GB storage

Debian GNU/Linux 9.13 (stretch)

go version go1.15.3 linux/amd64
```

### With `curl`

```
codespace ➜ ~/workspace/cloudgo (main ✗) $ curl http://localhost:8080/ -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.52.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Date: Sun, 22 Nov 2020 17:20:18 GMT
< Content-Length: 484
< 
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Home</title>
</head>
<body>
    <p id="current-time"></p>
    <form action="/" method="post">
        <input type="text" name="username" placeholder="用户名">
        <input type="text" name="password" placeholder="密码">
        <input type="submit">
    </form>

    <script src="/assets/script.js"></script>
</body>
* Curl_http_done: called premature == 0
* Connection #0 to host localhost left intact
```

```
codespace ➜ ~/workspace/cloudgo (main) $ curl http://localhost:8080/now -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /now HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.52.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 22 Nov 2020 17:20:44 GMT
< Content-Length: 30
< 
* Curl_http_done: called premature == 0
* Connection #0 to host localhost left intact
{"time":"2020/11/22 17:20:44"}
```

```
curl -d "username=test&password=123" http://localhost:8080/ -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> POST / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.52.1
> Accept: */*
> Content-Length: 26
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 26 out of 26 bytes
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Date: Sun, 22 Nov 2020 17:23:16 GMT
< Content-Length: 224
< 
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Result</title>
</head>
<body>
    Your username is test
</body>
* Curl_http_done: called premature == 0
* Connection #0 to host localhost left intact
```

### With `ab`

主要使用 `ab` 的两个参数，`-n` 指定总请求数，`-c` 指定并发请求数。

#### 100 并发，1000 总请求

```
codespace ➜ ~/workspace/cloudgo (main ✗) $ ab -n 1000 -c 100 http://localhost:8080/now
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /now
Document Length:        30 bytes

Concurrency Level:      100
Time taken for tests:   0.063 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      153000 bytes
HTML transferred:       30000 bytes
Requests per second:    15908.62 [#/sec] (mean)
Time per request:       6.286 [ms] (mean)
Time per request:       0.063 [ms] (mean, across all concurrent requests)
Transfer rate:          2376.97 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.8      2       4
Processing:     1    4   1.9      4      12
Waiting:        0    3   1.7      3      12
Total:          2    6   1.8      6      13

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      7
  75%      7
  80%      7
  90%      8
  95%     10
  98%     11
  99%     12
 100%     13 (longest request)
 ```

在测试结果中，我们关心：

1. 是否有失败的请求。显然对于这么一个简单的 HTTP handler，这项完全不应该出现。
2. 平均响应时间。从用户发出 HTTP 请求，到 HTTP 响应结束所需的总时间。这也是对于用户而言感知最强的一个数据。像本报告中平均是 `6.119 ms`。
3. 95 ~ 99 百分位最大请求时间。为什么要有一个百分位呢？因为总体的平均数有时会带来可观的误差，而当我们面向用户提供服务时，我们需要在意的是绝大部分请求（例如 95%, 99%）的响应时间，剩下的请求中有可能出现了极端数字，比如几百 ms 相对于绝大部分的几 ms 影响不大，如果采用平均值则反而可能会放大了这种影响。See https://www.elastic.co/cn/blog/averages-can-dangerous-use-percentile
4. 平均响应时间（考虑并发因素）。这里对应的是报告里的 `0.061 ms`，这个数字比不考虑并发因素的要小很多，为什么呢？这个数字反映的是服务器 accept 一个 HTTP 请求的速度（有流水线的味道），我们需要关注的是这个数字是否会随并发数增大而增大，如果变化不大则表明性能是可接受的。

#### 1000 并发，100000 总请求
```
codespace ➜ ~/workspace/cloudgo (main ✗) $ ab -n 10000 -c 1000 http://localhost:8080/now
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /now
Document Length:        30 bytes

Concurrency Level:      1000
Time taken for tests:   0.651 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1530000 bytes
HTML transferred:       300000 bytes
Requests per second:    15370.31 [#/sec] (mean)
Time per request:       65.061 [ms] (mean)
Time per request:       0.065 [ms] (mean, across all concurrent requests)
Transfer rate:          2296.54 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   22   7.3     24      34
Processing:    11   41  17.4     38     145
Waiting:        9   33  17.2     30     137
Total:         34   63  12.1     61     149

Percentage of the requests served within a certain time (ms)
  50%     61
  66%     64
  75%     67
  80%     69
  90%     79
  95%     87
  98%     93
  99%    102
 100%    149 (longest request)
 ```

 对比上一个测试结果，我们可以发现 `Time per request` 增大了很多（这是必然的，并发量增大后排队等待处理的请求便多了起来），但是服务器没有崩掉，99 百分数延迟在 100ms 左右，也是可接受的范围。