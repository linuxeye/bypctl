package compose

import (
	"bypctl/pkg/cmd"
	"bypctl/pkg/global"
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
	// fmt.Println("filePaths---->", filePaths)
	// projectOptions, err := cli.NewProjectOptions(
	// 	filePaths,
	// 	cli.WithDiscardEnvFile,
	// 	cli.WithEnvFiles(global.Conf.System.EnvFile),
	// 	cli.WithDotEnv,
	// )
	// if err != nil {
	// 	return err
	// }
	// ctx := context.Background()
	// p, err := projectOptions.LoadProject(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// dockerCli, err := command.NewDockerCli()
	// if err != nil {
	// 	return err
	// }
	// clientOptions := flags.ClientOptions{
	// 	Hosts:     []string{client.DefaultDockerHost},
	// 	LogLevel:  "debug",
	// 	TLS:       false,
	// 	TLSVerify: false,
	// }
	//
	// if err := dockerCli.Initialize(&clientOptions); err != nil {
	// 	return err
	// }
	// newComposeService := pkgCompose.NewComposeService(dockerCli)
	// upOptions := api.UpOptions{
	// 	Create: api.CreateOptions{},
	// 	Start: api.StartOptions{
	// 		Project: p,
	// 		Wait:    true,
	// 	},
	// }
	// if err := newComposeService.Up(ctx, p, upOptions); err != nil {
	// 	return err
	// }

	// yaml, err := p.MarshalYAML()
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println("yaml---->\n", gconv.String(yaml))
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

func sliceToStrF(list []string) string {
	var newList []string
	for _, item := range list {
		newList = append(newList, "-f "+item)
	}
	return strings.Join(newList, " ")
}
