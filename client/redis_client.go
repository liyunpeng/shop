package client

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"os/signal"
	"shop/config"
	"shop/logger"
	"syscall"
	"time"
)

var (
	RedisPool *redis.Pool
	Conn      redis.Conn
)

func StartRedisClient() {
	redisHost := config.TransformConfiguration.Redis.Addr
	logger.Info.Println("redis连接池初始化")
	RedisPool = newPool(redisHost)
	Conn = RedisPool.Get()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			logger.Info.Println("TCP REDIS 建立连接")
			if err != nil {
				logger.Info.Println("TCP REDIS 连接异常，err=", err)
				return nil, err
			} else {
				logger.Info.Println("TCP REDIS 连接成功")
				return c, err
			}
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			logger.Info.Println("tcp redis 连接测试 err=", err)
			return err
		},
	}
}

func RedisPoolclose() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		RedisPool.Close()
		os.Exit(0)
	}()
}

func RedisSet(key string, value string) error {
	var data []byte
	data, err := redis.Bytes(Conn.Do("SET", key, value))
	if err != nil {
		return fmt.Errorf("error get key %s: %v", key, err)
	}
	logger.Info.Println("redis set data :", string(data))
	return nil
}
type RedisUser struct {
	Id	string
	Name	string
 	Address		string
	Order	string
	ShopCar	string
}

func RedisUserHSet( userid string, k string , v string) { //userid string,  microServiceName string, ad0dress string, priority string ){
	hashName := "user:" + userid
	value, err := redis.String(Conn.Do("hmset", hashName,
		k, v,
	))

	if err != nil {
		logger.Info.Println("Hash failed:", err)
	} else {
		fmt.Printf("Hash result: \n %v \n", value)
	}
}
func RedisUserHMSet( user *RedisUser) {
	hashName := "user:" + user.Id
	value, err := redis.String(Conn.Do("hmset", hashName,
		"Name", user.Name,
		"Address", user.Address,
		"Order", user.Order,
		"Shopcar", user.ShopCar,
	))

	if err != nil {
		logger.Info.Println("Hash failed:", err)
	} else {
		fmt.Printf("Hash result: \n %v \n", value)
	}
}

func RedisGet(key string) ([]byte, error) {
	var data []byte
	data, err := redis.Bytes(Conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}

	//logger.Info.Println("redis get key=", string(data))
	return data, err
}

func RedisSetString(k string, v string) {
	_, err := Conn.Do("SET", k, v)
	if err != nil {
		logger.Info.Println("redis set failed:", err)
	}
}
func RedisDelString(k string, v string) {
	_, err := Conn.Do("DEL", k, v)
	if err != nil {
		logger.Info.Println("redis set failed:", err)
	}
}

func StringTest(k string, v string) {
	var err error
	_, err = Conn.Do("SET", "mykey", "superWang", "EX", "5")
	if err != nil {
		logger.Info.Println("redis set failed:", err)
	}

	var username string
	username, err = redis.String(Conn.Do("GET", "mykey"))
	if err != nil {
		logger.Info.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

	//time.Sleep(8 * time.Second)

	username, err = redis.String(Conn.Do("GET", "mykey"))
	if err != nil {
		logger.Info.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

	/*
		运行结果：
		redis set data : OK
		method= redis
		Get mykey: superWang
		redis get failed: redigo: nil returned
	*/
}

func Lua() {

	value, err := redis.String(Conn.Do("eval", "return 'hello world'", 0))

	if err != nil {
		logger.Info.Println("lua failed:", err)
	} else {
		fmt.Printf("lua result: %v \n", value)
	}
	/*
		运行结果：
		lua result: hello world
	*/
}

func Lpush() {

	listName := "redlist1"
	Conn.Do("lpush", listName, "qqq")
	Conn.Do("lpush", listName, "www")

	values, _ := redis.Values(Conn.Do("lrange", listName, "0", "100"))

	for _, v := range values {
		logger.Info.Println("v:", string(v.([]byte)))
	}
	/*
		运行结果：
		v: www
		v: qqq
	*/

	var v1 string
	redis.Scan(values, &v1)
	logger.Info.Println("v1=", v1)
	redis.Scan(values, &v1)
	logger.Info.Println("v1=", v1)
}

func Hashsetget() {

	//runtime.PrintStrack()

	hashName := "user:1000"
	value, err := redis.String(Conn.Do("hmset", hashName,
		"username", "antirez",
		"birthyear", "1977",
		"verified", "1"))

	if err != nil {
		logger.Info.Println("Hash failed:", err)
	} else {
		fmt.Printf("Hash result: \n %v \n", value)
	}

	values, err := redis.Values(Conn.Do("hgetall", hashName))
	if err != nil {
		logger.Info.Println("Hash hgetall failed:", err)
	} else {
		fmt.Printf("Hash hgetall result: \n %v \n", value)
		for _, v := range values {
			logger.Info.Println("hgetall:", string(v.([]byte)))
		}
		/*
			运行结果：
			Hash hgetall result:
			 OK
			hgetall: username
			hgetall: antirez
			hgetall: birthyear
			hgetall: 1977
			hgetall: verified
			hgetall: 1
		*/
	}
}

func Subscribe() {
	c := RedisPool.Get()
	psc := redis.PubSubConn{c}
	psc.Subscribe("redChatRoom")

	defer func() {
		c.Close()
		psc.Unsubscribe("redChatRoom")
	}()
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: messages: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			//fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			continue
		case error:
			logger.Info.Println(v)
			return
		}
	}
}

/**
 *redis发布信息
 *
 */
func Pubscribe(s string) {
	c := RedisPool.Get()
	defer c.Close()

	_, err := c.Do("PUBLISH", "redChatRoom", s)
	if err != nil {
		logger.Info.Println("pub err: ", err)
		return
	}
}

func Info() {

	value, err := redis.String(Conn.Do("info", "server"))

	if err != nil {
		logger.Info.Println("Info failed:", err)
	} else {
		fmt.Printf("Info result: \n %v \n", value)
	}
	/*
		运行结果：
		Info result:
		 # Server
		redis_version:3.0.500
		redis_git_sha1:00000000
		redis_git_dirty:0
		redis_build_id:e18d443f4404273c
		redis_mode:standalone
		os:Windows
		arch_bits:64
		multiplexing_api:WinSock_IOCP
		process_id:8056
		run_id:6145cb7e81c495ab8ea78f306bd11df9956b3ccd
		tcp_port:6379
		uptime_in_seconds:1064
		uptime_in_days:0
		hz:10
		lru_clock:12751945
		config_file:
	*/
}

func TestAll() {
	//StringTest()
	Lua()
	Info()
	Lpush()
	Hashsetget()

	Pubscribe("message publich")
	Subscribe()
}

