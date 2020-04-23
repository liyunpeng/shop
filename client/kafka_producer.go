package client

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"shop/custchan"
	"shop/util"
)

type KafkaProducer struct {
	producerClientI     interface{}
}


var KafkaProducerObj *KafkaProducer

func NewKafkaProducer(kafkaAddr string, isAsync bool) (kafkaProducer *KafkaProducer, err error) {
	kafkaProducer = &KafkaProducer{

	}

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	//等待kafka broker返回成功的响应,只有上面的RequireAcks设置是sarama.WaitForAll这里才有用.
	config.Producer.Return.Successes = true
	//等待kafka broker返回失败的响应,
	config.Producer.Return.Errors = true

	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1

	if isAsync {
		kafkaProducer.producerClientI, err = sarama.NewAsyncProducer([]string{kafkaAddr}, config)
		fmt.Println("创建kafaka 异步生产者")
	} else {
		kafkaProducer.producerClientI, err = sarama.NewSyncProducer([]string{kafkaAddr}, config)
		fmt.Println("创建kafaka 同步生产者")
	}

	if err != nil {
		logs.Error("init KafkaProducerObj producerClient err: %v", err)
		return
	}

	if isAsync {
		go func(p sarama.AsyncProducer) {
			for {
				select {
				case suc := <-p.Successes():
					fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				case fail := <-p.Errors():
					fmt.Println("err: ", fail.Err)
				}
			}
		}(kafkaProducer.producerClientI.(sarama.AsyncProducer))
	}

	return
}

func (k *KafkaProducer) sendMsgToKfk(isAsync bool) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()

	for v := range custchan.KafkaProducerMsgChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.Topic
		msg.Value = sarama.StringEncoder(v.Line)
		//也可将字符串转化为字节数组
		//msg.Value = sarama.ByteEncoder(value)

		fmt.Println("kafka生产者向kafka broker发送消息，消息字符串=",
			msg.Value, ", 消息主题=", msg.Topic)

		var err error
		if isAsync {
			k.producerClientI.(sarama.AsyncProducer).Input() <- msg
		}else{
			_, _, err = k.producerClientI.(sarama.SyncProducer).SendMessage(msg)
		}

		fmt.Println("kafka生产者发送消息完成")

		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}

	}

	fmt.Println("生产者退出")
}

func StartKafkaProducer(kafkaAddr string, threadNum int, isAync bool) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	var err error
	KafkaProducerObj, err = NewKafkaProducer(kafkaAddr, isAync)
	fmt.Println("kafka broker 地址=", kafkaAddr)
	if err != nil {
		panic("连接kafka broker错误 ")
	} else {
		fmt.Println("成功连接kafka broker")
	}
	for i := 0; i < threadNum; i++ {
		fmt.Println("启动Kafka发送消息的协程")
		util.WaitGroup.Add(1)
		go KafkaProducerObj.sendMsgToKfk( isAync)
	}
}
