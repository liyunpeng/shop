package gopoll

type Job interface {
	// 一个数据接口，所有的数据都要实现该接口，才能被传递进来
	//实现Job接口的一个数据实例，需要实现一个Do()方法，对数据的处理就在这个Do()方法中。
	Do()
}
