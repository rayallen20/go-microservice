# example of docker-compose

## PART1. 将之前用到的容器以docker-compose方式运行

docker-compose.yml如下:

```yaml
# 声明版本
version: "3"
# 定义服务 这些服务可以在同一个网络内进行访问
services:
  # 注册中心 与 配置中心
  consul:
    # 指定使用的镜像
    image: cap1573/consul:latest
    # 指定端口映射
    ports:
      # 容器端口:宿主机端口
      - "8500:8500"
  # 链路追踪
  jaeger:
    image: cap1573/jaeger:latest
    ports:
      - "6831:6831/udp"
      - "16686:16686"
  # 熔断器
  hystrix-dashboard:
    image: cap1573/hystrix-dashboard
    ports:
      - "9002:9002"
```

前台运行测试:

```
docker-compose -f docker-compose.yml up
```

```
docker-compose % docker-compose -f docker-compose.yml up
[+] Running 4/4
 ⠿ Network docker-compose_default                Created                                                                                                                                                                                                            0.1s
 ⠿ Container docker-compose-hystrix-dashboard-1  Created                                                                                                                                                                                                            0.1s
 ⠿ Container docker-compose-consul-1             Created                                                                                                                                                                                                            0.1s
 ⠿ Container docker-compose-jaeger-1             Created                                                                                                                                                                                                            0.1s
Attaching to docker-compose-consul-1, docker-compose-hystrix-dashboard-1, docker-compose-jaeger-1
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6373227,"caller":"healthcheck/handler.go:99","msg":"Health Check server started","http-port":14269,"status":"unavailable"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6414206,"caller":"memory/factory.go:55","msg":"Memory storage configuration","configuration":{"MaxTraces":0}}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6461158,"caller":"tchannel/builder.go:94","msg":"Enabling service discovery","service":"jaeger-collector"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6461997,"caller":"peerlistmgr/peer_list_mgr.go:111","msg":"Registering active peer","peer":"127.0.0.1:14267"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6473787,"caller":"standalone/main.go:187","msg":"Starting agent"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.6488128,"caller":"standalone/main.go:227","msg":"Starting jaeger-collector TChannel server","port":14267}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.648885,"caller":"standalone/main.go:237","msg":"Starting jaeger-collector HTTP server","http-port":14268}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.8864074,"caller":"standalone/main.go:298","msg":"Registering metrics handler with jaeger-query HTTP server","route":"/metrics"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.8867087,"caller":"standalone/main.go:304","msg":"Starting jaeger-query HTTP server","port":16686}
docker-compose-jaeger-1             | {"level":"info","ts":1675957231.886737,"caller":"healthcheck/handler.go:133","msg":"Health Check state change","status":"ready"}
docker-compose-consul-1             | ==> Starting Consul agent...
docker-compose-consul-1             |            Version: '1.9.11'
docker-compose-consul-1             |            Node ID: 'bda9797a-51b0-74b5-aeb4-5047dedf229c'
docker-compose-consul-1             |          Node name: 'ae6d8aef30fe'
docker-compose-consul-1             |         Datacenter: 'dc1' (Segment: '<all>')
docker-compose-consul-1             |             Server: true (Bootstrap: false)
docker-compose-consul-1             |        Client Addr: [0.0.0.0] (HTTP: 8500, HTTPS: -1, gRPC: 8502, DNS: 8600)
docker-compose-consul-1             |       Cluster Addr: 127.0.0.1 (LAN: 8301, WAN: 8302)
docker-compose-consul-1             |            Encrypt: Gossip: false, TLS-Outgoing: false, TLS-Incoming: false, Auto-Encrypt-TLS: false
docker-compose-consul-1             | 
docker-compose-consul-1             | ==> Log data will now stream in as it occurs:
docker-compose-consul-1             | 
docker-compose-consul-1             | 2023-02-09T15:40:31.990Z [INFO]  agent.server.raft: initial configuration: index=1 servers="[{Suffrage:Voter ID:bda9797a-51b0-74b5-aeb4-5047dedf229c Address:127.0.0.1:8300}]"
docker-compose-consul-1             | 2023-02-09T15:40:31.991Z [INFO]  agent.server.raft: entering follower state: follower="Node at 127.0.0.1:8300 [Follower]" leader=
docker-compose-consul-1             | 2023-02-09T15:40:31.995Z [INFO]  agent.server.serf.wan: serf: EventMemberJoin: ae6d8aef30fe.dc1 127.0.0.1
docker-compose-consul-1             | 2023-02-09T15:40:31.996Z [INFO]  agent.server.serf.lan: serf: EventMemberJoin: ae6d8aef30fe 127.0.0.1
docker-compose-consul-1             | 2023-02-09T15:40:31.996Z [INFO]  agent.router: Initializing LAN area manager
docker-compose-consul-1             | 2023-02-09T15:40:31.998Z [INFO]  agent.server: Handled event for server in area: event=member-join server=ae6d8aef30fe.dc1 area=wan
docker-compose-consul-1             | 2023-02-09T15:40:31.998Z [INFO]  agent.server: Adding LAN server: server="ae6d8aef30fe (Addr: tcp/127.0.0.1:8300) (DC: dc1)"
docker-compose-consul-1             | 2023-02-09T15:40:32.000Z [INFO]  agent: Started DNS server: address=0.0.0.0:8600 network=tcp
docker-compose-consul-1             | 2023-02-09T15:40:32.002Z [INFO]  agent: Started DNS server: address=0.0.0.0:8600 network=udp
docker-compose-consul-1             | 2023-02-09T15:40:32.003Z [INFO]  agent: Starting server: address=[::]:8500 network=tcp protocol=http
docker-compose-consul-1             | 2023-02-09T15:40:32.004Z [WARN]  agent: DEPRECATED Backwards compatibility with pre-1.9 metrics enabled. These metrics will be removed in a future version of Consul. Set `telemetry { disable_compat_1.9 = true }` to disable them.
docker-compose-consul-1             | 2023-02-09T15:40:32.005Z [INFO]  agent: Started gRPC server: address=[::]:8502 network=tcp
docker-compose-consul-1             | 2023-02-09T15:40:32.005Z [INFO]  agent: started state syncer
docker-compose-consul-1             | 2023-02-09T15:40:32.005Z [INFO]  agent: Consul agent running!
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [WARN]  agent.server.raft: heartbeat timeout reached, starting election: last-leader=
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [INFO]  agent.server.raft: entering candidate state: node="Node at 127.0.0.1:8300 [Candidate]" term=2
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [DEBUG] agent.server.raft: votes: needed=1
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [DEBUG] agent.server.raft: vote granted: from=bda9797a-51b0-74b5-aeb4-5047dedf229c term=2 tally=1
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [INFO]  agent.server.raft: election won: tally=1
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [INFO]  agent.server.raft: entering leader state: leader="Node at 127.0.0.1:8300 [Leader]"
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [INFO]  agent.server: cluster leadership acquired
docker-compose-consul-1             | 2023-02-09T15:40:32.047Z [DEBUG] agent.server: Cannot upgrade to new ACLs: leaderMode=0 mode=0 found=true leader=127.0.0.1:8300
docker-compose-consul-1             | 2023-02-09T15:40:32.050Z [INFO]  agent.server: New leader elected: payload=ae6d8aef30fe
docker-compose-consul-1             | 2023-02-09T15:40:32.051Z [INFO]  agent.leader: started routine: routine="federation state anti-entropy"
docker-compose-consul-1             | 2023-02-09T15:40:32.051Z [INFO]  agent.leader: started routine: routine="federation state pruning"
docker-compose-consul-1             | 2023-02-09T15:40:32.052Z [DEBUG] connect.ca.consul: consul CA provider configured: id=07:80:c8:de:f6:41:86:29:8f:9c:b8:17:d6:48:c2:d5:c5:5c:7f:0c:03:f7:cf:97:5a:a7:c1:68:aa:23:ae:81 is_primary=true
docker-compose-consul-1             | 2023-02-09T15:40:32.055Z [DEBUG] agent.server.autopilot: autopilot is now running
docker-compose-consul-1             | 2023-02-09T15:40:32.055Z [DEBUG] agent.server.autopilot: state update routine is now running
docker-compose-consul-1             | 2023-02-09T15:40:32.064Z [INFO]  agent.server.connect: initialized primary datacenter CA with provider: provider=consul
docker-compose-consul-1             | 2023-02-09T15:40:32.064Z [INFO]  agent.leader: started routine: routine="intermediate cert renew watch"
docker-compose-consul-1             | 2023-02-09T15:40:32.064Z [INFO]  agent.leader: started routine: routine="CA root pruning"
docker-compose-consul-1             | 2023-02-09T15:40:32.064Z [DEBUG] agent.server: successfully established leadership: duration=16.493426ms
docker-compose-consul-1             | 2023-02-09T15:40:32.064Z [INFO]  agent.server: member joined, marking health alive: member=ae6d8aef30fe
docker-compose-consul-1             | 2023-02-09T15:40:32.392Z [DEBUG] agent: Skipping remote check since it is managed automatically: check=serfHealth
docker-compose-consul-1             | 2023-02-09T15:40:32.393Z [INFO]  agent: Synced node info
docker-compose-consul-1             | 2023-02-09T15:40:32.393Z [DEBUG] agent: Node info in sync
docker-compose-consul-1             | 2023-02-09T15:40:32.485Z [INFO]  agent.server: federation state anti-entropy synced
docker-compose-jaeger-1             | {"level":"info","ts":1675957232.6467905,"caller":"peerlistmgr/peer_list_mgr.go:157","msg":"Not enough connected peers","connected":0,"required":1}
docker-compose-jaeger-1             | {"level":"info","ts":1675957232.6468854,"caller":"peerlistmgr/peer_list_mgr.go:166","msg":"Trying to connect to peer","host:port":"127.0.0.1:14267"}
docker-compose-jaeger-1             | {"level":"info","ts":1675957232.6486213,"caller":"peerlistmgr/peer_list_mgr.go:176","msg":"Connected to peer","host:port":"[::]:14267"}
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:33.027  INFO 1 --- [           main] c.l.HysterixDashboardApplication         : Starting HysterixDashboardApplication v0.0.1-SNAPSHOT on 56ea256254cb with PID 1 (/hysterix-dashboard.jar started by root in /)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:33.037  INFO 1 --- [           main] c.l.HysterixDashboardApplication         : No active profile set, falling back to default profiles: default
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:33.126  INFO 1 --- [           main] s.c.a.AnnotationConfigApplicationContext : Refreshing org.springframework.context.annotation.AnnotationConfigApplicationContext@478267af: startup date [Thu Feb 09 15:40:33 UTC 2023]; root of context hierarchy
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:33.437  INFO 1 --- [           main] trationDelegate$BeanPostProcessorChecker : Bean 'configurationPropertiesRebinderAutoConfiguration' of type [class org.springframework.cloud.autoconfigure.ConfigurationPropertiesRebinderAutoConfiguration$$EnhancerBySpringCGLIB$$5f6aca0f] is not eligible for getting processed by all BeanPostProcessors (for example: not eligible for auto-proxying)
docker-compose-consul-1             | 2023-02-09T15:40:33.797Z [DEBUG] agent: Skipping remote check since it is managed automatically: check=serfHealth
docker-compose-consul-1             | 2023-02-09T15:40:33.797Z [DEBUG] agent: Node info in sync
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:33.849  INFO 1 --- [           main] c.l.HysterixDashboardApplication         : Started HysterixDashboardApplication in 1.151 seconds (JVM running for 2.219)
docker-compose-hystrix-dashboard-1  | 
docker-compose-hystrix-dashboard-1  |   .   ____          _            __ _ _
docker-compose-hystrix-dashboard-1  |  /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
docker-compose-hystrix-dashboard-1  | ( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
docker-compose-hystrix-dashboard-1  |  \\/  ___)| |_)| | | | | || (_| |  ) ) ) )
docker-compose-hystrix-dashboard-1  |   '  |____| .__|_| |_|_| |_\__, | / / / /
docker-compose-hystrix-dashboard-1  |  =========|_|==============|___/=/_/_/_/
docker-compose-hystrix-dashboard-1  |  :: Spring Boot ::        (v1.3.3.RELEASE)
docker-compose-hystrix-dashboard-1  | 
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:34.074  INFO 1 --- [           main] c.l.HysterixDashboardApplication         : No active profile set, falling back to default profiles: default
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:34.092  INFO 1 --- [           main] ationConfigEmbeddedWebApplicationContext : Refreshing org.springframework.boot.context.embedded.AnnotationConfigEmbeddedWebApplicationContext@4de1bcbe: startup date [Thu Feb 09 15:40:34 UTC 2023]; parent: org.springframework.context.annotation.AnnotationConfigApplicationContext@478267af
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:34.612  INFO 1 --- [           main] o.s.b.f.s.DefaultListableBeanFactory     : Overriding bean definition for bean 'beanNameViewResolver' with a different definition: replacing [Root bean: class [null]; scope=; abstract=false; lazyInit=false; autowireMode=3; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=org.springframework.boot.autoconfigure.web.ErrorMvcAutoConfiguration$WhitelabelErrorViewConfiguration; factoryMethodName=beanNameViewResolver; initMethodName=null; destroyMethodName=(inferred); defined in class path resource [org/springframework/boot/autoconfigure/web/ErrorMvcAutoConfiguration$WhitelabelErrorViewConfiguration.class]] with [Root bean: class [null]; scope=; abstract=false; lazyInit=false; autowireMode=3; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=org.springframework.boot.autoconfigure.web.WebMvcAutoConfiguration$WebMvcAutoConfigurationAdapter; factoryMethodName=beanNameViewResolver; initMethodName=null; destroyMethodName=(inferred); defined in class path resource [org/springframework/boot/autoconfigure/web/WebMvcAutoConfiguration$WebMvcAutoConfigurationAdapter.class]]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:34.725  INFO 1 --- [           main] o.s.cloud.context.scope.GenericScope     : BeanFactory id=ed9cb8f1-65d7-38ff-bd08-faf757ada316
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:34.750  INFO 1 --- [           main] trationDelegate$BeanPostProcessorChecker : Bean 'org.springframework.cloud.autoconfigure.ConfigurationPropertiesRebinderAutoConfiguration' of type [class org.springframework.cloud.autoconfigure.ConfigurationPropertiesRebinderAutoConfiguration$$EnhancerBySpringCGLIB$$5f6aca0f] is not eligible for getting processed by all BeanPostProcessors (for example: not eligible for auto-proxying)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.018  INFO 1 --- [           main] s.b.c.e.t.TomcatEmbeddedServletContainer : Tomcat initialized with port(s): 9002 (http)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.029  INFO 1 --- [           main] o.apache.catalina.core.StandardService   : Starting service Tomcat
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.030  INFO 1 --- [           main] org.apache.catalina.core.StandardEngine  : Starting Servlet Engine: Apache Tomcat/8.0.32
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.126  INFO 1 --- [ost-startStop-1] o.a.c.c.C.[Tomcat].[localhost].[/]       : Initializing Spring embedded WebApplicationContext
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.127  INFO 1 --- [ost-startStop-1] o.s.web.context.ContextLoader            : Root WebApplicationContext: initialization completed in 1035 ms
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.922  INFO 1 --- [ost-startStop-1] o.s.b.c.e.ServletRegistrationBean        : Mapping servlet: 'proxyStreamServlet' to [/proxy.stream]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.925  INFO 1 --- [ost-startStop-1] o.s.b.c.e.ServletRegistrationBean        : Mapping servlet: 'dispatcherServlet' to [/]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.936  INFO 1 --- [ost-startStop-1] o.s.b.c.embedded.FilterRegistrationBean  : Mapping filter: 'characterEncodingFilter' to: [/*]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.937  INFO 1 --- [ost-startStop-1] o.s.b.c.embedded.FilterRegistrationBean  : Mapping filter: 'hiddenHttpMethodFilter' to: [/*]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.937  INFO 1 --- [ost-startStop-1] o.s.b.c.embedded.FilterRegistrationBean  : Mapping filter: 'httpPutFormContentFilter' to: [/*]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:35.938  INFO 1 --- [ost-startStop-1] o.s.b.c.embedded.FilterRegistrationBean  : Mapping filter: 'requestContextFilter' to: [/*]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.291  INFO 1 --- [           main] o.s.ui.freemarker.SpringTemplateLoader   : SpringTemplateLoader for FreeMarker: using resource loader [org.springframework.boot.context.embedded.AnnotationConfigEmbeddedWebApplicationContext@4de1bcbe: startup date [Thu Feb 09 15:40:34 UTC 2023]; parent: org.springframework.context.annotation.AnnotationConfigApplicationContext@478267af] and template loader path [classpath:/templates/]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.293  INFO 1 --- [           main] o.s.w.s.v.f.FreeMarkerConfigurer         : ClassTemplateLoader for Spring macros added to FreeMarker configuration
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.569  INFO 1 --- [           main] s.w.s.m.m.a.RequestMappingHandlerAdapter : Looking for @ControllerAdvice: org.springframework.boot.context.embedded.AnnotationConfigEmbeddedWebApplicationContext@4de1bcbe: startup date [Thu Feb 09 15:40:34 UTC 2023]; parent: org.springframework.context.annotation.AnnotationConfigApplicationContext@478267af
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.624  INFO 1 --- [           main] s.w.s.m.m.a.RequestMappingHandlerMapping : Mapped "{[/hystrix]}" onto public java.lang.String org.springframework.cloud.netflix.hystrix.dashboard.HystrixDashboardController.home(org.springframework.ui.Model,org.springframework.web.context.request.WebRequest)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.625  INFO 1 --- [           main] s.w.s.m.m.a.RequestMappingHandlerMapping : Mapped "{[/hystrix/{path}]}" onto public java.lang.String org.springframework.cloud.netflix.hystrix.dashboard.HystrixDashboardController.monitor(java.lang.String,org.springframework.ui.Model,org.springframework.web.context.request.WebRequest)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.627  INFO 1 --- [           main] s.w.s.m.m.a.RequestMappingHandlerMapping : Mapped "{[/error],produces=[text/html]}" onto public org.springframework.web.servlet.ModelAndView org.springframework.boot.autoconfigure.web.BasicErrorController.errorHtml(javax.servlet.http.HttpServletRequest,javax.servlet.http.HttpServletResponse)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.627  INFO 1 --- [           main] s.w.s.m.m.a.RequestMappingHandlerMapping : Mapped "{[/error]}" onto public org.springframework.http.ResponseEntity<java.util.Map<java.lang.String, java.lang.Object>> org.springframework.boot.autoconfigure.web.BasicErrorController.error(javax.servlet.http.HttpServletRequest)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.653  INFO 1 --- [           main] o.s.w.s.handler.SimpleUrlHandlerMapping  : Mapped URL path [/webjars/**] onto handler of type [class org.springframework.web.servlet.resource.ResourceHttpRequestHandler]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.653  INFO 1 --- [           main] o.s.w.s.handler.SimpleUrlHandlerMapping  : Mapped URL path [/**] onto handler of type [class org.springframework.web.servlet.resource.ResourceHttpRequestHandler]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.687  INFO 1 --- [           main] o.s.w.s.handler.SimpleUrlHandlerMapping  : Mapped URL path [/**/favicon.ico] onto handler of type [class org.springframework.web.servlet.resource.ResourceHttpRequestHandler]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.824  WARN 1 --- [           main] o.s.c.n.a.ArchaiusAutoConfiguration      : No spring.application.name found, defaulting to 'application'
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.828  WARN 1 --- [           main] c.n.c.sources.URLConfigurationSource     : No URLs will be polled as dynamic configuration sources.
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.828  INFO 1 --- [           main] c.n.c.sources.URLConfigurationSource     : To enable URLs as dynamic configuration sources, define System property archaius.configurationSource.additionalUrls or make config.properties available on classpath.
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.838  WARN 1 --- [           main] c.n.c.sources.URLConfigurationSource     : No URLs will be polled as dynamic configuration sources.
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.838  INFO 1 --- [           main] c.n.c.sources.URLConfigurationSource     : To enable URLs as dynamic configuration sources, define System property archaius.configurationSource.additionalUrls or make config.properties available on classpath.
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.883  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Registering beans for JMX exposure on startup
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.892  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Bean with name 'configurationPropertiesRebinder' has been autodetected for JMX exposure
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.893  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Bean with name 'refreshScope' has been autodetected for JMX exposure
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.894  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Bean with name 'environmentManager' has been autodetected for JMX exposure
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.897  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Located managed bean 'environmentManager': registering with JMX server as MBean [org.springframework.cloud.context.environment:name=environmentManager,type=EnvironmentManager]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.910  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Located managed bean 'refreshScope': registering with JMX server as MBean [org.springframework.cloud.context.scope.refresh:name=refreshScope,type=RefreshScope]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:36.918  INFO 1 --- [           main] o.s.j.e.a.AnnotationMBeanExporter        : Located managed bean 'configurationPropertiesRebinder': registering with JMX server as MBean [org.springframework.cloud.context.properties:name=configurationPropertiesRebinder,context=4de1bcbe,type=ConfigurationPropertiesRebinder]
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:37.073  INFO 1 --- [           main] s.b.c.e.t.TomcatEmbeddedServletContainer : Tomcat started on port(s): 9002 (http)
docker-compose-hystrix-dashboard-1  | 2023-02-09 15:40:37.075  INFO 1 --- [           main] c.l.HysterixDashboardApplication         : Started HysterixDashboardApplication in 4.478 seconds (JVM running for 5.445)
```

后台运行的启停:

```
docker-compose % docker-compose -f docker-compose.yml up -d
[+] Running 3/3
 ⠿ Container docker-compose-consul-1             Started                                                                                                                                                                                                            0.6s
 ⠿ Container docker-compose-jaeger-1             Started                                                                                                                                                                                                            0.6s
 ⠿ Container docker-compose-hystrix-dashboard-1  Started
```

```
docker-compose -f docker-compose.yml stop 
[+] Running 3/3
 ⠿ Container docker-compose-hystrix-dashboard-1  Stopped                                                                                                                                                                                                            0.3s
 ⠿ Container docker-compose-consul-1             Stopped                                                                                                                                                                                                            0.6s
 ⠿ Container docker-compose-jaeger-1             Stopped
```

```
docker-compose % docker-compose -f docker-compose.yml down
[+] Running 4/4
 ⠿ Container docker-compose-hystrix-dashboard-1  Removed                                                                                                                                                                                                            0.0s
 ⠿ Container docker-compose-consul-1             Removed                                                                                                                                                                                                            0.0s
 ⠿ Container docker-compose-jaeger-1             Removed                                                                                                                                                                                                            0.0s
 ⠿ Network docker-compose_default                Removed  
```