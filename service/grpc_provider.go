package service

import (
	"bufio"
	"context"
	"os"
	"shop/logger"
	pb "shop/rpc/proto"
	"shop/util"
)

// 定义服务端实现约定的接口
type UserInfoService struct {
	file *os.File
}

func (u *UserInfoService) bufioWriteFile(s string) {
	var err error
	content := []byte(s)
	newWriter := bufio.NewWriterSize(u.file, 1024)
	if _, err = newWriter.Write(content); err != nil {
		logger.Info.Println(err)
	}
	if err = newWriter.Flush(); err != nil {
		logger.Info.Println(err)
	}
	logger.Info.Println("write file successful")
}

type WorkerPoolJob struct {
	Paraeter string
}

func (w *WorkerPoolJob) Do() {
	util.PrintFuncName()

	logger.Info.Println("Do parameter=", w.Paraeter)
}

// 实现服务端需要首先的接口
func (s *UserInfoService) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	name := req.Name
	// 在数据库查用户信息
	if name == "zhangsan" {
		resp = &pb.UserResponse{
			Id:    1,
			Name:  name,
			Age:   22,
			Hobby: []string{"Sing", "run", "basketball"},
		}
	} else {
		resp = &pb.UserResponse{
			Id:    1,
			Name:  name,
			Age:   22,
			Hobby: []string{"1", "2", "3"},
		}
	}
	s.bufioWriteFile(name)
	err = nil

	w := & WorkerPoolJob{
		Paraeter: name,
	}
	GrpcWorkerPool.JobQueue  <- w
	return
}
func (u *UserInfoService) init() {
	var err error
	u.file, err = os.OpenFile("./aaa.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.Info.Println(err)
	}
}
func (u *UserInfoService) destroy() {
	u.file.Close()
}
