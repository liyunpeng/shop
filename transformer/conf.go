package transformer

/*
加载系统配置
*/

type Conf struct {
	App      App
	Mysql    Mysql
	Mongodb  Mongodb
	Sqlite   Sqlite
	Redis    Redis
	TestData TestData
	Kafka    Kafka
	Etcd    EtcdConf
	Grpc    GrpcConf
	Consul   ConsulConf
}

type App struct {
	Name          string
	URl           string
	Port          string
	LoggerLevel   string
	DirverType    string
	//CreateSysData int
}

type Mysql struct {
	DirverName string
	Connect    string
	CasbinConnect    string

	Name       string
	TName      string
	MaxIdle		int
	MaxOpen		int
}

type Mongodb struct {
	Connect string
}

type Sqlite struct {
	DirverName string
	Connect    string
	TConnect   string
}

type Redis struct {
	Addr     string
	Password string
	DB       string
}

type TestData struct {
	UserName        string
	Name            string
	Pwd             string
	DataBaseDriver  string
	DataBaseConnect string
}

type Kafka struct {
	Addr     string
}

type EtcdConf struct {
	Addr     string
}

type GrpcConf struct {
	Addr     string
}

type ConsulConf struct {
	Addr     string
	Port     string
}
