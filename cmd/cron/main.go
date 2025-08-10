package main

import (
	"context"
	"github.com/xxl-job/xxl-job-executor-go"
	"log"
	"time"
)

func main() {
	// 从环境变量获取调度中心地址和执行器端口
	// 这种方式非常适合容器化部署
	adminAddr := "http://127.0.0.1:8089/xxl-job-admin"
	executorPort := "9998"

	// 初始化执行器
	exec := xxl.NewExecutor(
		xxl.ServerAddr(adminAddr),
		xxl.ExecutorIp("192.168.124.19"), // 为空时会自动获取IP，容器化部署时很有用
		xxl.ExecutorPort(executorPort),
		xxl.RegistryKey("golang-executor-sample"), // 注册到调度中心的AppName，必须与调度中心配置一致
		xxl.AccessToken("default_token"),
	)
	exec.Init()

	// 注册任务处理器
	// JobHandler 对应调度中心 "JobHandler" 字段的值
	exec.RegTask("demoJobHandler.golang", DemoJobHandler)
	exec.RegTask("shardingJobHandler.golang", ShardingJobHandler)

	log.Printf("Starting Golang XXL-Job Executor on port %s...", executorPort)
	log.Printf("Registering to admin at: %s", adminAddr)

	// 运行执行器
	err := exec.Run()
	if err != nil {
		log.Fatalf("Failed to run executor: %v", err)
	}
}

// Demo任务1：简单打印日志
func DemoJobHandler(c context.Context, param *xxl.RunReq) string {
	log.Println("------ [Golang] XXL-JOB DemoJobHandler executed ------")
	log.Printf("Param: %s", param.ExecutorParams)
	log.Printf("LogID: %d", param.LogID)
	log.Println("------ [Golang] Job finished ------")
	time.Sleep(10 * time.Second)
	return "Golang executor task finished successfully."
}

// Demo任务2：演示分片
func ShardingJobHandler(c context.Context, param *xxl.RunReq) string {
	log.Println("------ [Golang] XXL-JOB ShardingJobHandler executed ------")
	log.Printf("Sharding total: %d", param.BroadcastTotal)
	log.Printf("Sharding index: %d", param.BroadcastIndex)
	log.Println("------ [Golang] Sharding job finished ------")
	return "Golang sharding task finished."
}
