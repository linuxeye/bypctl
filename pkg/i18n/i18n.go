package i18n

import (
	"bypctl/pkg/global"
	"context"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

func Translate(content string) string {
	ctx := gi18n.WithLanguage(context.TODO(), global.Conf.System.Lang)
	return gi18n.Translate(ctx, content)
}

func Tf(content string, values ...any) string {
	ctx := gi18n.WithLanguage(context.TODO(), global.Conf.System.Lang)
	return gi18n.Tf(ctx, content, values...)
}
