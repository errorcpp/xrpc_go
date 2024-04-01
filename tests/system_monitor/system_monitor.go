/// @coding utf-8
/// @author errorcpp@qq.com

package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	for {
		vmem, _ := mem.VirtualMemory()
		// 使用率
		used_percent := vmem.UsedPercent
		fmt.Printf("virtual memory usage: %.2f%%\n", used_percent)
		time.Sleep(1 * time.Second)
	}
}
