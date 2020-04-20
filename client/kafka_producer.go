package client

import (
	"fmt"
		"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"shop/custchan"
	"shop/util"
)



type KafkaProducer struct {
	producerClient       sarama.SyncProducer

}
var KafkaProducerObj *KafkaProducer

func NewKafkaProducer(kafkaAddr string) (kafkaProducer *KafkaProducer, err error) {
	kafkaProducer = &KafkaProducer{

	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // wait KafkaProducerObj ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init KafkaProducerObj producerClient err: %v", err)
		return
	}
	kafkaProducer.producerClient = client
	return
}

func (k *KafkaProducer) sendMsgToKfk() {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()

	for v := range custchan.KafkaProducerMsgChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.Topic
		msg.Value = sarama.StringEncoder(v.Line)

		fmt.Println("kafka生产者向kafka broker发送消息，消息字符串=",
			msg.Value, ", 消息主题=", msg.Topic)

		_, _, err := k.producerClient.SendMessage(msg)

		fmt.Println("kafka生产者发送消息完成")

		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
	}

	fmt.Println("生产者退出")
}



func StartKafkaProducer(kafkaAddr string, threadNum int) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	var err error
	KafkaProducerObj, err = NewKafkaProducer(kafkaAddr)
	fmt.Println("kafka broker 地址=", kafkaAddr)
	if  err != nil {
		panic("连接kafka broker错误 ")
	} else {
		fmt.Println("成功连接kafka broker")
	}
	for i := 0; i < threadNum; i++ {
		fmt.Println("启动Kafka发送消息的协程")
		util.WaitGroup.Add(1)
		go KafkaProducerObj.sendMsgToKfk()
	}
}

