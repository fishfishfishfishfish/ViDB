/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	_ "net/http/pprof" // 自动注册 pprof HTTP 路由
	"os"
	"os/signal"
	"syscall"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "letus-vidb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initPprof()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	enablePprof bool
	pprofPort   string
)

func init() {
	// 注册持久化标志（所有子命令继承）
	rootCmd.PersistentFlags().BoolVar(&enablePprof, "pprof", false, "启用 pprof 性能分析")
	rootCmd.PersistentFlags().StringVar(&pprofPort, "pprof-port", "6060", "pprof 服务端口（默认 6060）")
}

// initPprof 启动 pprof 服务（修复函数名拼写错误）
func initPprof() {
	if enablePprof {
		go func() {
			addr := ":" + pprofPort
			// 启动 HTTP 服务（pprof 路由已通过 import 自动注册）
			log.Printf("启动 pprof 服务: http://localhost%s/debug/pprof", addr)
			// ListenAndServe 本身会阻塞，无需额外 delay()
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Printf("pprof 服务启动失败: %v", err)
			}
		}()
	}
}

// initMetrics 启动 Prometheus  metrics 服务
func initMetrics() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting metrics server on :8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
}

func delay() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 等待退出信号
}
