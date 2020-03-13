package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var kafkaSender = &KafkaSender{}

type KafkaService interface {
	RunServer()
}

type kafkaService struct {
	Sendors []*KafkaSender
}

func NewKafkaService(kafkaAddr string, threadNum int) *kafkaService {
	k := &kafkaService{
		Sendors: make([]*KafkaSender, 5, 10),
	}
	kafkaSender, _ = NewKafkaSend(kafkaAddr, threadNum)
	k.Sendors[0] = kafkaSender

	// TODO  comsumer
	return k
}

type Message struct {
	line  string
	topic string
}

type KafkaSender struct {
	producerClient sarama.SyncProducer
	lineChan       chan *Message
}

// NewKafkaSend is
func NewKafkaSend(kafkaAddr string, threadNum int) (kafkaSender *KafkaSender, err error) {
	kafkaSender = &KafkaSender{
		lineChan: make(chan *Message, 10000),
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // wait kafkaSender ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafkaSender producerClient err: %v", err)
		return
	}
	kafkaSender.producerClient = client

	for i := 0; i < threadNum; i++ {
		fmt.Println("启动执行Kafka发送消息的协程")
		waitGroup.Add(1)
		go kafkaSender.sendMsgToKfk()
	}
	return
}

func (k *KafkaSender) sendMsgToKfk() {
	defer waitGroup.Done()

	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)

		_, _, err := k.producerClient.SendMessage(msg)

		fmt.Println("kafka生产者向kafka broker发送消息，消息字符串=",
			msg.Value, ", 消息主题=", msg.Topic)

		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
	}
}

func (k *KafkaSender) addMessage(line string, topic string) (err error) {
	k.lineChan <- &Message{line: line, topic: topic}
	return
}
