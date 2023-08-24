package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		cpu := LoadCpu()
		time.Sleep(200 * time.Millisecond)

		for _, r := range LoadTemps() {
			fmt.Println(r.name)
			for _, t := range r.temps {
				fmt.Printf("%s: %.00f°C\n", t.label, float32(t.input)/1000)
			}
		}
		mem := LoadMemInfo()
		fmt.Printf("RAM: %.2f GB/%.2f GB	SWAP: %.2f GB/%.2f GB\n", mem.MemTotalGB()-mem.MemAvailableGB(), mem.MemTotalGB(), mem.SwapTotalGB()-mem.SwapFreeGB(), mem.SwapTotalGB())

		loadavg := GetLoadavg()
		fmt.Printf("%.2f %.2f %.2f\n", loadavg.Loadavg1, loadavg.Loadavg5, loadavg.Loadavg15)

		cpufreq := LoadCpufreq()
		current := LoadCpu()

		for i := 0; i < len(cpufreq)/4; i++ {
			fmt.Printf("%d: %d%% %.0f MHz	%d: %d%% %.0f MHz	%d: %d%% %.0f MHz	%d: %d%% %.0f MHz\n", cpufreq[i].Processor, current.cores[i].Usage(cpu.cores[i]), cpufreq[i].Freq, cpufreq[i+3].Processor, current.cores[i+3].Usage(cpu.cores[i+3]), cpufreq[i+3].Freq, cpufreq[i+6].Processor, current.cores[i+6].Usage(cpu.cores[i+6]), cpufreq[i+6].Freq, cpufreq[i+9].Processor, current.cores[i+9].Usage(cpu.cores[i+9]), cpufreq[i+9].Freq)
		}
		fmt.Println()

		time.Sleep(1 * time.Second)
	}
}
