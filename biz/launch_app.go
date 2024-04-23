package biz

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var (
	appPath  = "/home/ubuntu/.local/share/Steam/steamcmd/cs2/game/bin/linuxsteamrt64/cs2"
	keywords = []string{
		"activated session on GC",
		"获取地图配置成功",
		"解析配置成功",
		"CSSharp: CGameSystem::Shutdown",
	}
)

func execScript(shFilePath, appPath, port, createRoomCmd string) error {
	cmd := exec.CommandContext(context.Background(), "/bin/bash", shFilePath, appPath, port)

	// 通过管道连接脚本的标准输入和输出
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating StdinPipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating StdoutPipe: %w", err)
	}

	// 启动脚本
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	// 创建一个 scanner 以读取脚本的标准输出
	scanner := bufio.NewScanner(stdout)

	// 持续读取脚本的标准输出
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Log:", line)
		// 检查日志中是否包含关键字
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				fmt.Println("Keyword found:", keyword)
				// 根据关键字发送新的指令或者打印日志
				switch keyword {
				case "activated session on GC":
					if err := sendCommandToScript(stdin, createRoomCmd); err != nil {
						fmt.Println("Error sending command to script:", err)
					}
				case "获取地图配置成功", "解析配置成功":
					fmt.Println("Map configuration:", line)
				case "CSSharp: CGameSystem::Shutdown":
					fmt.Println("Game shutdown:", line)
				}
			}
		}
	}

	// 检查扫描器是否出错
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading Stdout: %w", err)
	}

	// 等待脚本执行完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %w", err)
	}

	fmt.Println("Script execution finished")
	return nil

}

func sendCommandToScript(stdin io.WriteCloser, command string) error {
	if _, err := io.WriteString(stdin, command+"\n"); err != nil {
		return fmt.Errorf("error writing to Stdin: %w", err)
	}

	fmt.Println("Sent command to script: ", command)

	if err := stdin.Close(); err != nil {
		return fmt.Errorf("error closing Stdin: %w", err)
	}
	return nil
}

func LaunchCSGO() {
	createRoomCmd := "aHR0cHM6Ly9vcGVuLnluY3VlLmNvbS9hcHAvbWF0Y2hJbmZv aHR0cHM6Ly9vcGVuLnluY3VlLmNvbS9hcHAvcG9zdA== dGVzdA=="
	scriptPath := "../script/launch.sh"
	gamePort := "27015"
	if err := execScript(scriptPath, appPath, gamePort, createRoomCmd); err != nil {
		fmt.Println(err)
		return
	}
}
