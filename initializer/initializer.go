package initializer

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	U "github.com/daifiyum/cat-box/common"
)

type AppConfig struct {
	Workspace bool
	Port      int
}

func Initialize() error {
	cfg := parseFlags()

	if err := setupWorkspace(cfg.Workspace); err != nil {
		return fmt.Errorf("工作区初始化失败: %w", err)
	}

	if cfg.Port != 0 {
		U.Port = fmt.Sprint(cfg.Port)
	}

	if err := initLogging(); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	return nil
}

func parseFlags() *AppConfig {
	workspace := flag.Bool("workspace", false, "启用工作区模式，将工作目录切换至可执行文件所在目录")
	port := flag.Int("port", 0, "设置端口号")
	flag.Parse()

	return &AppConfig{
		Workspace: *workspace,
		Port:      *port,
	}
}

func setupWorkspace(enable bool) error {
	if enable {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		dir := filepath.Dir(exePath)
		if err := os.Chdir(dir); err != nil {
			return err
		}
	}
	return nil
}

func initLogging() error {
	f, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}
