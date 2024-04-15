/// @coding utf-8
/// @author errorcpp@qq.com
///   nohup command > /dev/null 2>&1 &  # 无任何重定向

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rifflock/lfshook"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	// Add any additional fields you need
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 没有logger_name这个东西
	func_name := ""
	file_name := ""
	line := 0
	if entry.HasCaller() {
		func_name = entry.Caller.Func.Name()
		file_name = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	return []byte(
		fmt.Sprintf(
			"[%s][%s][tid=%d][%s:%s:%d] [%s]\n",
			entry.Time.Format("2006-01-02 15:04:05"),
			entry.Level.String(),
			// Add thread ID if available
			// entry.Data["tid"],
			// Add filename, function name, and line number
			file_name, // Replace with entry.Caller.File
			func_name, // Replace with entry.Caller.Function
			line,      // Replace with entry.Caller.Line
			entry.Message,
		),
	), nil
}

func SetupLog() {
	// 创建日志记录器
	logger := logrus.StandardLogger() // logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	// // 设置日志级别
	// logger.SetLevel(logrus.DebugLevel)
	//logger.SetFormatter(&CustomFormatter{})
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000 +0800",
		// 函数名和行号
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.Function, f.Line)
		},
	})
	logger.SetReportCaller(true)

	logger.SetOutput(os.Stdout)
	// // 添加控制台输出hook，如果默认就已经有控制台输出了，这里反而重复添加
	// consoleWriter := os.Stdout
	// consoleHook := &writer.Hook{ // 控制台输出
	// 	Writer:    consoleWriter,
	// 	LogLevels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	// }
	// logger.AddHook(consoleHook)

	// 获取进程名
	process_name := filepath.Base(os.Args[0])
	process_name = strings.TrimSuffix(process_name, filepath.Ext(process_name))

	// 设置文件日志
	log_file_name := process_name + ".log"
	log_path := "./logs/"
	file_logger := lumberjack.Logger{
		Filename:   filepath.Join(log_path, log_file_name),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   false,
	}
	file_log_hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  &file_logger,
			logrus.WarnLevel:  &file_logger,
			logrus.ErrorLevel: &file_logger,
			logrus.FatalLevel: &file_logger,
			logrus.PanicLevel: &file_logger,
		},
		&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000 +0800",
			// 函数名和行号
			// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 	return "", fmt.Sprintf("%s:%d", f.Function, f.Line)
			// },
		}, // 使用文本格式
	)
	// 添加文件日志钩子
	logger.AddHook(file_log_hook)
}

func BytesToMB(cnt uint64) float64 {
	return float64(cnt) / 1024 / 1024
}

func BytesToGB(cnt uint64) float64 {
	return float64(cnt) / 1024 / 1024 / 1024
}

// Used + Free != Total
// Total：代表系统内存的总量
// Available：代表系统当前可用的内存
// Used：代表已被系统或应用程序使用的内存
// Free：代表系统当前未被使用的内存
// Buffers：代表用于存放缓冲数据的内存量
// Cached：代表被文件系统缓存使用的内存量
// Active：代表活跃的内存
// Inactive：代表不活跃的内存
// 如果你将 Total - (Available + Buffers + Cached)，你将得到已被使用的内存量，这样最终就能加起来等于 Total

func main() {
	SetupLog()
	for {
		vmem, _ := mem.VirtualMemory()
		// 虚拟内存(在windows指标里边就是物理内存)
		used_percent := vmem.UsedPercent
		//fmt.Printf("virtual memory usage_percent/total: %.2f%%/%.2f\n", used_percent, BytesToGB(vmem.Total))
		logrus.Debugf("virtual memory usage_percent/total: %.2f%%/%.2f\n", used_percent, BytesToGB(vmem.Total))
		// 可用虚拟内存
		logrus.Debugf("virtual memory used/usable: %.2f/%.2f\n", BytesToGB(vmem.Used), BytesToGB(vmem.Free))
		// 交换分区
		swapmem, _ := mem.SwapMemory()
		swap_used_percent := swapmem.UsedPercent
		logrus.Debugf("swap memory usage: %.2f%%\n", swap_used_percent)
		// 空闲内存低于指定值执行shell
		free_mb := BytesToMB(vmem.Free)
		// 剩余可用内存
		available_mb := BytesToMB(vmem.Available)
		if available_mb < 500 {
			cmd := "ps aux | grep \".vscode-server\" | awk '{print $2}' | xargs kill -9"
			logrus.Infof("free memory is low: free=%.2f,available=%.2f, run_cmd=%s", free_mb, available_mb, cmd)
			output, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				logrus.Info("run cmd faild: %s", output)
			} else {
				logrus.Info("run result: %s", output)
			}
		}
		time.Sleep(9 * time.Second)
	}
}
