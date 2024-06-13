package compose

import (
	"bypctl/pkg/cmd"
	"bypctl/pkg/global"
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"os/exec"
	"strings"
)

func Pull(filePaths, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s pull %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Ps(filePaths []string, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s ps %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Up(filePaths, args []string, detachMode bool) error {
	var detach string
	if detachMode {
		detach = "-d"
	}
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s up %s %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "), detach)
	return err
}

func Logs(filePaths, args []string, followMode bool) error {
	var follow string
	if followMode {
		follow = "-f"
	}
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s logs %s %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "), follow)
	return err
}

func Down(filePaths, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s down %s --remove-orphans", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Start(filePaths, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s start %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Stop(filePaths, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s stop %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Restart(filePaths, args []string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s restart %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), strings.Join(args, " "))
	return err
}

func Operate(filePath, operation string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s -f %s %s", global.Conf.System.EnvFile, filePath, operation)
	return err
}

func Exec(filePaths []string, app, command string) error {
	err := cmd.ExecCmdOutput("docker-compose --env-file %s %s exec %s %s", global.Conf.System.EnvFile, sliceToStrF(filePaths), app, command)
	return err
}

func GetServices(filePaths []string) ([]string, error) {
	command := exec.Command("bash", "-c", fmt.Sprintf("docker-compose --env-file %s %s ps --services", global.Conf.System.EnvFile, sliceToStrF(filePaths)))
	if err := command.Run(); err != nil {
		return nil, err
	}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	return strings.Split(gconv.String(stdout), "\n"), nil
}

func sliceToStrF(list []string) string {
	var newList []string
	for _, item := range list {
		newList = append(newList, "-f "+item)
	}
	return strings.Join(newList, " ")
}
