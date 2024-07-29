package lang

import (
	"embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go-rest-api/internal/config"
	"golang.org/x/text/language"
	"strings"
)

var (
	EnglishTranslate *i18n.Localizer
	ChineseTranslate *i18n.Localizer
	Bundles          *i18n.Bundle

	//go:embed locale/*.toml
	LocaleFS embed.FS
)

func Init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//fmt.Printf("bundle %v", bundle)
	loadMessageFiles(bundle)
	Bundles = bundle
	EnglishTranslate = i18n.NewLocalizer(Bundles, language.English.String())
	ChineseTranslate = i18n.NewLocalizer(Bundles, language.Chinese.String())
}

func loadMessageFiles(bundle *i18n.Bundle) {
	files := []string{
		"locale/en.toml",
		"locale/zh-CN.toml",
	}

	for _, file := range files {
		if _, err := bundle.LoadMessageFileFS(LocaleFS, file); err != nil {
			config.Logger.Error(fmt.Sprintf("Failed to load message file %s: [Error] %v", file, err))
		}
	}
}

func T(lang string, message string) string {
	translate := getLangTran(lang)
	if translate == nil {
		return message
	}
	msg, err := translate.Localize(&i18n.LocalizeConfig{
		MessageID: message,
	})
	if err != nil {
		return message
	}

	return msg
}

func GetLang(c *gin.Context) string {
	if c == nil {
		return ""
	}
	locale := c.Request.Header.Get("language")

	return locale
}

func getLangTran(lang string) *i18n.Localizer {
	if lang == "zh" || strings.Contains(lang, "zh") {
		return ChineseTranslate
	}
	return EnglishTranslate
}
