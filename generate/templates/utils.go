package templates

import (
	"fmt"
	"strings"
)

func BulkReplace(content string, pairs map[string]string) string {

	for key, value := range pairs {
		replaceKey := fmt.Sprintf("{{%s}}", key)
		content = strings.ReplaceAll(content, replaceKey, value)
	}

	return content
}
