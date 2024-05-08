package i18n

import (
	"bypctl/locale"
	"bypctl/pkg/global"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Translate(msgId, msg string) string {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(locale.LocaleFS, "active.zh-CN.toml")
	// fmt.Println("global.Conf.System.Lang Translate--->", msgId, global.Conf.System.Lang)
	localizer := i18n.NewLocalizer(bundle, global.Conf.System.Lang)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgId,
			Other: msg,
		},
	})
}

func Tf(msgId, msg string, data any) string {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(locale.LocaleFS, "active.zh-CN.toml")
	// fmt.Println("global.Conf.System.Lang Tf--->", global.Conf.System.Lang)
	localizer := i18n.NewLocalizer(bundle, global.Conf.System.Lang)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgId,
			Other: msg,
		},
		TemplateData: data,
	})
}
