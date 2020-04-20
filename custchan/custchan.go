package custchan
type Message struct {
	Line  string
	Topic string
}

var ConfChan  = make(chan string, 10)

var KafkaProducerMsgChan  = make(chan *Message, 10000)

func AddKafkaProducerMsg(line string, topic string) (err error) {
	KafkaProducerMsgChan <- &Message{Line: line, Topic: topic}
	return
}
