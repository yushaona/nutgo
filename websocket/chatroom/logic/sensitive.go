package logic

import (
	"strings"

	"github.com/yushaona/nutgo/websocket/chatroom/global"
)

func FilterSensitive(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}
