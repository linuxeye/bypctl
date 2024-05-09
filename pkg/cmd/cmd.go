package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/scottames/cmder"
	"os/exec"
	"strings"
	"time"
)

func Exec(cmdStr string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecCronjobWithTimeOut(cmdStr string, workdir string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = workdir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}

	var errMsg string
	if len(stderr.String()) != 0 {
		errMsg = fmt.Sprintf("stderr:\n %s", stderr.String())
	}
	if len(stdout.String()) != 0 {
		if len(errMsg) != 0 {
			errMsg = fmt.Sprintf("%s \n\n; stdout:\n %s", errMsg, stdout.String())
		} else {
			errMsg = fmt.Sprintf("stdout:\n %s", stdout.String())
		}
	}
	return errMsg, err
}

func Execf(cmdStr string, a ...any) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdStr, a...))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		var errMsg string
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecWithCheck(name string, a ...string) (string, error) {
	cmd := exec.Command(name, a...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		var errMsg string
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecScript(scriptPath, workDir string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = workDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}
	if err != nil {
		var errMsg string
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func CheckIllegal(args ...string) bool {
	if args == nil {
		return false
	}
	for _, arg := range args {
		if strings.Contains(arg, "&") || strings.Contains(arg, "|") || strings.Contains(arg, ";") ||
			strings.Contains(arg, "$") || strings.Contains(arg, "'") || strings.Contains(arg, "`") ||
			strings.Contains(arg, "(") || strings.Contains(arg, ")") || strings.Contains(arg, "\"") {
			return true
		}
	}
	return false
}

func HasNoPasswordSudo() bool {
	command := exec.Command("sudo", "-n", "ls")
	err := command.Run()
	return err == nil
}

func SudoHandleCmd() string {
	cmd := exec.Command("sudo", "-n", "ls")
	if err := cmd.Run(); err == nil {
		return "sudo "
	}
	return ""
}

func Which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func ExecCmdOutput(cmdStr string, a ...any) error {
	cmd := cmder.New("bash", "-c", fmt.Sprintf(cmdStr, a...))
	cmd.Silent()
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
