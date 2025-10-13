package main

/*
1. 每秒随机生成一批（100个以内）键值对<key,value>，其中key必须是小写字母 a~z的范围，value必须是 0.1~5.1的浮点数
2. 通过chan将键值对传递给下游
3. 下游从channel消费，对<key, value>做全局统计，要求10个并发（即同时处理10个<key, value>），要求最终输出每个key的value最大值、最小值、平均值、和值以及分位值，分位值是加分项
4. 程序收到SIGHUP、SIGINT、SIGQUIT信号，停止生成器、统计器，并打印出统计结果
*/
func main() {

}
