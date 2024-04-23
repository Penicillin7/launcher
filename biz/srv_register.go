package biz

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func NewClient() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	err = registerService(client)
	if err != nil {
		log.Fatal(err)
	}

	// 捕获退出信号，用于取消服务注册
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	go http.ListenAndServe(":8080", nil)

	//s := StartGrpcServer()

	// 等待退出信号
	<-stop

	// 取消服务注册
	err = deregisterService(client)
	if err != nil {
		log.Fatal(err)
	}
}

// registerService 注册服务到 Consul
func registerService(client *api.Client) error {
	// 构建服务注册参数
	service := &api.AgentServiceRegistration{
		ID:      "my-service",
		Name:    "my-service",
		Port:    8080,
		Address: "127.0.0.1",
		Check: &api.AgentServiceCheck{
			TCP:      "127.0.0.1:8080", // 检查服务是否可用的地址
			Interval: "10s",            // 检查间隔
		},
	}

	// 发送服务注册请求
	err := client.Agent().ServiceRegister(service)
	if err != nil {
		return err
	}

	fmt.Println("Service registered successfully")
	return nil
}

// deregisterService 从 Consul 注销服务
func deregisterService(client *api.Client) error {
	// 构建服务注销参数
	err := client.Agent().ServiceDeregister("my-service")
	if err != nil {
		return err
	}

	fmt.Println("Service deregistered successfully")
	return nil
}
