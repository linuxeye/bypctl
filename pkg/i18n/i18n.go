package i18n

import (
	"bypctl/pkg/global"
	"context"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

func Translate(format string) string {
	ctx := gi18n.WithLanguage(context.TODO(), global.Conf.System.Lang)
	return gi18n.Translate(ctx, format)
}

func Tf(format string, values ...any) string {
	ctx := gi18n.WithLanguage(context.TODO(), global.Conf.System.Lang)
	return gi18n.Tf(ctx, format, values...)
}
