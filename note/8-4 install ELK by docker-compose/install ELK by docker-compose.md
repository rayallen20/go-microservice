# install ELK by docker-compose

## PART1. 编排elasticsearch容器

### 1.1 编写ES的配置文件`elasticsearch.yml`

目录结构如下:

```
tree ./
./
├── docker-stack.yml	// 容器编排文件
└── elasticsearch
    └── config
        └── elasticsearch.yml		// ES的配置文件

2 directories, 2 files
```

`elasticsearch/config/elasticsearch.yml`:

```yaml
---
# 集群名称
cluster.name: "immoc-cluster"
# 主机IP
network.host: 0.0.0.0

# xpack是ES用于创建账号密码的
xpack.license.self_generated.type: trial
# 开启账号密码验证
xpack.security.enabled: true
xpack.monitoring.collection.enabled: true
```

### 1.2 编写容器编排文件`docker-stack.yml`

`docker-stack.yml`:

```yaml
version: '3.3'
services:
  # elasticsearch
  elasticsearch:
    image: cap1573/elasticsearch:7.9.3
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - ./elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
    environment:
      # ES允许使用的最大内存
      ES_JAVA_OPTS: "-Xmx256m -Xms256m"
      # ES的密码
      ELASTIC_PASSWORD: imoocpwd
      discovery.type: single-node
      network.publish_host: _eth0_
```

## PART2. 编排logstash容器

### 2.1 编写logstash的配置文件`logstash.yml`

目录结构如下:

```
tree ./
./
├── docker-stack.yml	// 容器编排文件
├── elasticsearch
│   └── config
│       └── elasticsearch.yml		// ES的配置文件
└── logstash
    └── config
        └── logstash.yml	// logstash的配置文件

4 directories, 3 files
```

`logstash/config/logstash `:

```yaml
---
http.host: "0.0.0.0"
xpack.monitoring.elasticsearch.hosts: ["http://elasticsearch:9200"]

xpack.monitoring.enabled: true
xpack.monitoring.elasticsearch.username: elastic
xpack.monitoring.elasticsearch.password: imoocwpd
```

`logstash/config/logstash.yml`:

```yaml
---
# 主机IP
http.host: "0.0.0.0"
# 指定ES的主机地址和端口
# 此处的elasticsearch 是service的名称
xpack.monitoring.elasticsearch.hosts: ["http://elasticsearch:9200"]

# 开启账号密码验证
xpack.monitoring.enabled: true
# ES的用户名
xpack.monitoring.elasticsearch.username: elastic
# ES的密码
xpack.monitoring.elasticsearch.password: imoocpwd
```

注:`---`表示yaml文件的开始

### 2.2 编写logstash的工作流水线`logstash.conf`

`logstash.conf`文件描述了logstash工作的3个阶段:

- input
- filter
- output

```
tree ./
./
├── docker-stack.yml	// 容器编排文件
├── elasticsearch
│   └── config
│       └── elasticsearch.yml		// ES的配置文件
└── logstash
    ├── config
    │   └── logstash.yml	// logstash的配置文件
    └── pipeline
        └── logstash.conf	// logstash的工作流水线

5 directories, 4 files
```

`logstash/pipeline/logstash.conf`:

```
input {
    beats {
        port => 5044
    }
    tcp {
        port => 5000
    }
}

output {
    elasticsearch {
        hosts => "elasticsearch:9200"
        user => "elastic"
        password => "imoocpwd"
        index => "%{[@metadata][-imooc]}-%{[@metadata][version]}-%{+YYYY.MM.dd}"
    }
}
```

此处没有定义filter阶段,表示logstash将所有接收到的消息都存入ES

### 2.3 编写容器编排文件`docker-stack.yml`

`docker-stack.yml`:

```yaml
version: '3.3'
services:
  # elasticsearch
  elasticsearch:
    image: cap1573/elasticsearch:7.9.3
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - ./elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
    environment:
      # ES允许使用的最大内存
      ES_JAVA_OPTS: "-Xmx256m -Xms256m"
      # ES的密码
      ELASTIC_PASSWORD: imoocpwd
      discovery.type: single-node
      network.publish_host: _eth0_

  # logstash
  logstash:
    image: cap1573/logstash:7.9.3
    ports:
      - "5044:5044"
      - "5000:5000"
      - "9600:9600"
    volumes:
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml
      - ./logstash/pipeline/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    environment:
      # LS允许使用的最大内存
      LS_JAVA_OPTS: "-Xmx256m -Xms256m"
```

## PART3. 编排kibana容器

### 3.1 编写kibana的配置文件`kibana.yml`

目录结构如下:

```
tree ./
./
├── docker-stack.yml	// 容器编排文件
├── elasticsearch
│   └── config
│       └── elasticsearch.yml		// ES的配置文件
├── kibana
│   └── config
│       └── kibana.yml		// kibana的配置文件
└── logstash
    ├── config
    │   └── logstash.yml	// logstash的配置文件
    └── pipeline
        └── logstash.conf	// logstash的工作流水线

7 directories, 5 files
```

`kibana/config/kibana.yml`:

```
---
# 主机名称
server.name: kibana
# kibana主机IP
server.host: 0.0.0.0
# ES主机地址
elasticsearch.hosts: ["http://elasticsearch:9200"]
# 是否开启UI界面
monitoring.ui.container.elasticsearch.enabled: true

# ES的用户名
elasticsearch.username: elastic
# ES的密码
elasticsearch.password: imoocpwd
```

### 3.3 编写容器编排文件`docker-stack.yml`

`docker-stack.yml`:

```yaml
version: '3.3'
services:
  # elasticsearch
  elasticsearch:
    image: cap1573/elasticsearch:7.9.3
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - ./elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
    environment:
      # ES允许使用的最大内存
      ES_JAVA_OPTS: "-Xmx256m -Xms256m"
      # ES的密码
      ELASTIC_PASSWORD: imoocpwd
      discovery.type: single-node
      network.publish_host: _eth0_

  # logstash
  logstash:
    image: cap1573/logstash:7.9.3
    ports:
      - "5044:5044"
      - "5000:5000"
      - "9600:9600"
    volumes:
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml
      - ./logstash/pipeline/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    environment:
      # LS允许使用的最大内存
      LS_JAVA_OPTS: "-Xmx256m -Xms256m"

  # kibana
  kibana:
    image: cap1573/kibana:7.9.3
    ports:
      - "5601:5601"
    volumes:
      - ./kibana/config/kibana.yml:/usr/share/kibana/config/kibana.yml
```