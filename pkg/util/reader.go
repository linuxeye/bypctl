package util

import (
	"bufio"
	"bypctl/pkg/i18n"
	"fmt"
	"os"
	"strings"
)

func ReaderTf(format string, values ...any) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(i18n.Tf(format, values...))
	name, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}
