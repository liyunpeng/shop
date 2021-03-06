#### 1. 解决因zookeeper容器启动错误导致的kafka容器启动不了
* 问题描述：
docker logs 看到这样的log:  
```
Will not attempt to authenticate using SASL (unknown error) ,
```
运行kafka的时候,提示连接zookeeper错误

* 原因分析：
后来发现是启动zookeeper命令错误, 这个错误的命令是：  
```
docker run -d --name zookeeper -p 2181:2181 -t wurstmeister/zookeeper
```
很多网上是这么写的启动zookeeper，虽然启动了，但是其他程序无法访问

* 解决方法：修改zookeeper启动命令
原有zookeper容器结束掉， 启动新的zookeeper容器， 正确的命令为：  
``` 
docker run -itd --name zookeeper -p 2181:2181 wurstmeister/zookeeper
```
1). -it 表示允许与容器实例进行交互，如果不加-it, 容器里的服务连接会拒绝  
2). -p 端口映射，把容器的端口和宿主机的端口进行映射

#### 2. 启动kafka容器
```
docker run -d --name kafka -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=192.168.0.223:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://192.168.0.223:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -t wurstmeister/kafka
```

#### 3. 进入kafka容器，创建主题，生产者向主题发送消息，消费者从主题消费消息，查看主题下消息。
```
$ docker exec -it ef0857453b57 /bin/sh
/ # ls
bin    etc    kafka  lib64  mnt    proc   run    srv    tmp    var
dev    home   lib    media  opt    root   sbin   sys    usr
```

kafka脚本自动在path里， 按tab可以看到这些命令
```shell script
/ # kafka-
kafka-acls.sh                        kafka-delete-records.sh              kafka-replica-verification.sh
kafka-broker-api-versions.sh         kafka-dump-log.sh                    kafka-run-class.sh
kafka-configs.sh                     kafka-leader-election.sh             kafka-server-start.sh
kafka-console-consumer.sh            kafka-log-dirs.sh                    kafka-server-stop.sh
kafka-console-producer.sh            kafka-mirror-maker.sh                kafka-streams-application-reset.sh
kafka-consumer-groups.sh             kafka-preferred-replica-election.sh  kafka-topics.sh
kafka-consumer-perf-test.sh          kafka-producer-perf-test.sh          kafka-verifiable-consumer.sh
kafka-delegation-tokens.sh           kafka-reassign-partitions.sh         kafka-verifiable-producer.sh

#创建主题， zookeeper地址， 要指定分区数，副本数，主题的名字
/ # kafka-topics.sh --create --zookeeper 192.168.0.223:2181 --replication-factor 1 --partitions 1 --topic MyTopic
Created topic MyTopic.
MyTopic

#生产者向主题发消息
/ # kafka-console-producer.sh --broker-list 192.168.0.223:9092 --topic MyTopic
>1111111111
>ffffff

#打开另一终端，进入容器， 消费者接受消息：
bogon:shop admin1$ docker exec -it ef0857453b57 /bin/sh
/ # kafka-console-consumer.sh --bootstrap-server PLAINTEXT://192.168.0.223:9092 --topic MyTopic --from-beginning
1111111111
ffffff

#查看主题下的所有消息：
/ # kafka-console-consumer.sh --bootstrap-server 192.168.0.223:9092 --topic nginx_log  --from-beginning
100
300

#查看kafka的日志目录:
/kafka # ls kafka-logs-ef0857453b57/
MyTopic-0                         __consumer_offsets-25             __consumer_offsets-42
__consumer_offsets-0              __consumer_offsets-26             __consumer_offsets-43
__consumer_offsets-1              __consumer_offsets-27             __consumer_offsets-44
__consumer_offsets-10             __consumer_offsets-28             __consumer_offsets-45
__consumer_offsets-11             __consumer_offsets-29             __consumer_offsets-46
__consumer_offsets-12             __consumer_offsets-3              __consumer_offsets-47
__consumer_offsets-13             __consumer_offsets-30             __consumer_offsets-48
__consumer_offsets-14             __consumer_offsets-31             __consumer_offsets-49
__consumer_offsets-15             __consumer_offsets-32             __consumer_offsets-5
__consumer_offsets-16             __consumer_offsets-33             __consumer_offsets-6
__consumer_offsets-17             __consumer_offsets-34             __consumer_offsets-7
__consumer_offsets-18             __consumer_offsets-35             __consumer_offsets-8
__consumer_offsets-19             __consumer_offsets-36             __consumer_offsets-9
__consumer_offsets-2              __consumer_offsets-37             cleaner-offset-checkpoint
__consumer_offsets-20             __consumer_offsets-38             log-start-offset-checkpoint
__consumer_offsets-21             __consumer_offsets-39             meta.properties
__consumer_offsets-22             __consumer_offsets-4              nginx_log-0
__consumer_offsets-23             __consumer_offsets-40             recovery-point-offset-checkpoint
__consumer_offsets-24             __consumer_offsets-41             replication-offset-checkpoint

/kafka/kafka-logs-ef0857453b57/nginx_log-0 # ls -l
total 12
-rw-r--r--    1 root     root      10485760 Mar 14 23:28 00000000000000000000.index
-rw-r--r--    1 root     root          1846 Mar 15 01:53 00000000000000000000.log
-rw-r--r--    1 root     root      10485756 Mar 14 23:28 00000000000000000000.timeindex
-rw-r--r--    1 root     root            10 Mar 14 23:28 00000000000000000003.snapshot
-rw-r--r--    1 root     root             8 Mar 14 23:28 leader-epoch-checkpoint
```

