package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	wshandler "shop/handler"
	"shop/rpc"
	"sync"
)
var Address = []string{
	"192.168.0.198:9092",
}

func StartKafkaConsumer(kafkaAdress string) {
	topic := []string{"nginx_log"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//Address.append(kafkaAdress)kafkaAdress
	Address = append(Address, kafkaAdress)
	//广播式消费：消费者1
	go clusterConsumer(wg, Address, topic, "group-1")

	go clusterConsumerRpc(wg, Address, topic, "group-2")


	//广播式消费：消费者2
	//go clusterConsumer(wg, Address, topic, "group-2")

	wg.Wait()
}


// 支持brokers cluster的消费者
func clusterConsumer(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		fmt.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}else{
		fmt.Println("消费者成功建立")
	}
	defer func(){
		consumer.Close()
		fmt.Println("消费者关闭")
	}()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Printf("消费者组%s: 出错，Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				wshandler.WebsocketChan <- string(msg.Value)

				fmt.Println("kafka 消费者消费消息")
				fmt.Fprintf(
					os.Stdout,
					"消费组ID=%s，主题=%s，分区=%d，offset=%d，key=%s，value=%s\n",
					groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}


func clusterConsumerRpc(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		fmt.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}else{
		fmt.Println("消费者成功建立")
	}
	defer func(){
		consumer.Close()
		fmt.Println("消费者关闭")
	}()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Printf("消费者组%s: 出错，Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				rpc.Client(string(msg.Value))

				fmt.Println("kafka 消费者消费消息")
				fmt.Fprintf(
					os.Stdout,
					"消费组ID=%s，主题=%s，分区=%d，offset=%d，key=%s，value=%s\n",
					groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}
