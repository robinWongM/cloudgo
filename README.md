# cloudgo
A Go web server with a `/hello/:id` route.

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
codespace ➜ ~/workspace/cloudgo (main) $ curl http://localhost:8080/hello/your -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /hello/your HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.52.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Thu, 19 Nov 2020 14:06:08 GMT
< Content-Length: 21
< 
* Curl_http_done: called premature == 0
* Connection #0 to host localhost left intact
{"Test":"Hello your"}
```

### With `ab`

主要使用 `ab` 的两个参数，`-n` 指定总请求数，`-c` 指定并发请求数。

#### 100 并发，1000 总请求

```
codespace ➜ ~/workspace/cloudgo (main) $ ab -n 1000 -c 100 http://localhost:8080/hello/your
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

Document Path:          /hello/your
Document Length:        21 bytes

Concurrency Level:      100
Time taken for tests:   0.061 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      144000 bytes
HTML transferred:       21000 bytes
Requests per second:    16342.81 [#/sec] (mean)
Time per request:       6.119 [ms] (mean)
Time per request:       0.061 [ms] (mean, across all concurrent requests)
Transfer rate:          2298.21 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.8      1       4
Processing:     1    4   2.5      4      12
Waiting:        0    3   2.3      3      12
Total:          1    6   2.3      6      14
WARNING: The median and mean for the initial connection time are not within a normal deviation
        These results are probably not that reliable.

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      6
  75%      7
  80%      7
  90%      9
  95%     11
  98%     12
  99%     13
 100%     14 (longest request)
 ```

在测试结果中，我们关心：

1. 是否有失败的请求。显然对于这么一个简单的 HTTP handler，这项完全不应该出现。
2. 平均响应时间。从用户发出 HTTP 请求，到 HTTP 响应结束所需的总时间。这也是对于用户而言感知最强的一个数据。像本报告中平均是 `6.119 ms`。
3. 95 ~ 99 百分位最大请求时间。为什么要有一个百分位呢？因为总体的平均数有时会带来可观的误差，而当我们面向用户提供服务时，我们需要在意的是绝大部分请求（例如 95%, 99%）的响应时间，剩下的请求中有可能出现了极端数字，比如几百 ms 相对于绝大部分的几 ms 影响不大，如果采用平均值则反而可能会放大了这种影响。See https://www.elastic.co/cn/blog/averages-can-dangerous-use-percentile
4. 平均响应时间（考虑并发因素）。这里对应的是报告里的 `0.061 ms`，这个数字比不考虑并发因素的要小很多，为什么呢？这个数字反映的是服务器 accept 一个 HTTP 请求的速度（有流水线的味道），我们需要关注的是这个数字是否会随并发数增大而增大，如果变化不大则表明性能是可接受的。

 #### 1000 并发，100000 总请求
 ```
 codespace ➜ ~/workspace/cloudgo (main ✗) $ ab -n 100000 -c 1000 http://localhost:8080/hello/your
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /hello/your
Document Length:        21 bytes

Concurrency Level:      1000
Time taken for tests:   5.369 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      14400000 bytes
HTML transferred:       2100000 bytes
Requests per second:    18626.08 [#/sec] (mean)
Time per request:       53.688 [ms] (mean)
Time per request:       0.054 [ms] (mean, across all concurrent requests)
Transfer rate:          2619.29 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   20   5.0     21      34
Processing:     7   33  10.1     32     150
Waiting:        0   26  10.2     25     144
Total:          7   53   9.1     53     174

Percentage of the requests served within a certain time (ms)
  50%     53
  66%     55
  75%     57
  80%     58
  90%     62
  95%     68
  98%     77
  99%     86
 100%    174 (longest request)
 ```

 对比上一个测试结果，我们可以发现 `Time per request` 增大了很多（这是必然的，并发量增大后排队等待处理的请求便多了起来），但是服务器没有崩掉，99 百分数延迟在 86ms，也是可接受的范围。