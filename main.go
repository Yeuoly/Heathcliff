package main

import (
	"fmt"

	"github.com/Yeuoly/Heathcliff/runner"
	"github.com/Yeuoly/Heathcliff/server"
	"github.com/Yeuoly/Heathcliff/types"
)

func main() {
	//buffer := make(map[int]types.Buffer)
	server.AppendListener(func(id int, msg []byte, n int) int {
		exec_path, err := runner.RunCompile(&types.Task{
			Buf: &types.Buffer{
				Buf: msg[:n],
				N:   n,
			},
			Type: types.TASK_TYPE_ACM_C,
		})

		if err != nil {
			fmt.Println(err)
			return server.SIGNAL_OVER
		}

		result, err := runner.RunExec(exec_path)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result)

		return server.SIGNAL_OVER
	})

	server.Run()
}
