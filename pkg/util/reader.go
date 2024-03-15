package util

import (
	"bufio"
	"bypctl/pkg/i18n"
	"fmt"
	"os"
	"strings"
)

func ReaderTf(content string, values ...any) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(i18n.Tf(content, values...))
	name, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}
