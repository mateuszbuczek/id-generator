package main

import (
	"bufio"
	"fmt"
	"id-generator/src/generator"
	"os"
	"time"
)

func main() {
	num := 1_000_000
	result := make([]string, 0, num)

	startTime := time.Now().UnixMicro()

	data := make(chan int64)
	worker := generator.Worker{WorkerId: 1, GeneratedIds: &data, ThreadCount: 16}
	worker.Generate(int64(num))

	for i := 0; i < num; i++ {
		select {
		case generatedID := <-*worker.GeneratedIds:
			result = append(result, fmt.Sprintf("id: %d, epoch: %d, workerId: %d, threadId: %d, counter: %d",
				generatedID, generatedID>>24, (generatedID&0xf0_00_00)>>20, (generatedID&0x0f_00_00)>>16, generatedID&0xff_ff))
		case <-time.After(1 * time.Second):
			os.Exit(1)
		}
	}

	endTime := time.Now().UnixMicro()

	fmt.Printf("generating %d ids took %d ms", num, (endTime-startTime)/1000.0)

	// write to a file
	if true {
		file, _ := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		datawriter := bufio.NewWriter(file)
		for _, data := range result {
			_, _ = datawriter.WriteString(data + "\n")
		}

		_ = datawriter.Flush()
		_ = file.Close()
	}

	os.Exit(0)
}
