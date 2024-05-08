package util

import (
	"bufio"
	"bypctl/pkg/i18n"
	"fmt"
	"os"
	"strings"
)

func ReaderTf(msgId, msg string, data any) string {
	// fmt.Println("global.Conf.System.Lang ReaderTf--->", global.Conf.System.Lang)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(i18n.Tf(msgId, msg, data))
	name, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}
