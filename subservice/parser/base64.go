package parser

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64URLSafe(content string) ([]byte, error) {
	content = strings.ReplaceAll(content, " ", "-")
	content = strings.ReplaceAll(content, "/", "_")
	content = strings.ReplaceAll(content, "+", "-")
	content = strings.ReplaceAll(content, "=", "")
	result, err := base64.RawURLEncoding.DecodeString(content)
	return result, err
}
