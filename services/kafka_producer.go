package services

import (
	"fmt"
		"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type Message struct {
	line  string
	topic string
}

type KafkaProducer struct {
	producerClient sarama.SyncProducer
	lineChan       chan *Message
}
var kafkaProducer *KafkaProducer

func NewKafkaProducer(kafkaAddr string) (kafkaProducer *KafkaProducer, err error) {
	kafkaProducer = &KafkaProducer{
		lineChan: make(chan *Message, 10000),
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // wait kafkaProducer ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafkaProducer producerClient err: %v", err)
		return
	}
	kafkaProducer.producerClient = client
	return
}

func (k *KafkaProducer) sendMsgToKfk() {
	defer waitGroup.Done()

	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)

		fmt.Println("kafka生产者向kafka broker发送消息，消息字符串=",
			msg.Value, ", 消息主题=", msg.Topic)

		_, _, err := k.producerClient.SendMessage(msg)

		fmt.Println("kafka生产者发送消息完成")

		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
	}
}

func (k *KafkaProducer) addMessage(line string, topic string) (err error) {
	k.lineChan <- &Message{line: line, topic: topic}
	return
}

func StartKafkaProducer(kafkaAddr string, threadNum int) {
	var err error
	kafkaProducer, err = NewKafkaProducer(kafkaAddr)
	fmt.Println("kafka broker 地址=", kafkaAddr)
	if ( err != nil) {
		panic("连接kafka broker错误 ")
	} else {
		fmt.Println("成功连接kafka broker")
	}
	for i := 0; i < threadNum; i++ {
		fmt.Println("启动Kafka发送消息的协程")
		waitGroup.Add(1)
		go kafkaProducer.sendMsgToKfk()
	}
}

