# review of monitor

## PART1. 本章小结

- 介绍监控prometheus的架构,组成和重要概念
- 引入docker-compose管理多个容器
- 订单微服务代码开发并接入监控

## PART2. 经验之谈

- docker-compose不仅能够管理容器,还能制作镜像(`docker-compose build`)
- 特别注意`docker-compose down`命令会清除数据,一定慎用
- 本章中采用配置文件记录监控目标,也可以采用配置中心的方式

## PART3. 课后思考

- 监控中采用配置文件的方式有哪些优劣?
- 监控横向扩展会使用哪种方式记录数据?
- docker-compose还能做什么?