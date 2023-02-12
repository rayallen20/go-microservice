# working principle of logstash

## PART1. logstash工作原理

3个阶段:inputs -> filters -> outputs

- input:数据输入阶段.会把数据输入到logstash
- filter:数据清洗阶段.数据中间处理,对数据进行操作
- outputs:数据输出阶段.outputs是logstash处理管道的最末端组件

## PART2. logstash-inout阶段常见的输入

- file:从文件系统的文件中读取,类似于`tail -f`
- syslog:在514端口上监听系统日志消息,并根据RFC3164标准进行解析
- beats:从Filebeat中读取

## PART3. logstash-input常用样例

```
input {
	beats {
		port => 5044
	}
	tcp {
		port => 5000
	}
}
```

## PART4. logstash-filter数据中间件处理插件grok(可解析任意文本数据)

- GROK基础语法如:`%{SYNTAX:SEMANTIC}`
- SYNTAX:代表匹配值的类型
- SEMANTIC:代表存储该值的一个变量名称

例:`%{ERROR|DEBUG|INFO|WARN:log_level}`

## PART5. logstash-output数据输出

- 输出到kafka和ES,也可以输出到redis

例:输出至ES

```
output {
	elasticsearch {
		hosts => "elasticsearch:9200"
		user => "elastic"
		password => "changename"
		index => "%{[@metadata][-imooc]}-%{[@metadata][version]}-%{+YYYY.MM.dd}"
	}
}
```