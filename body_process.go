package main

import (
	"regexp"
	"strings"
)

func process(body []byte) []byte {
	// Remove <script> and <style> blocks
	reScript := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	reStyle := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
	body = reScript.ReplaceAll(body, []byte(""))
	body = reStyle.ReplaceAll(body, []byte(""))

	// Remove HTML comments
	reComments := regexp.MustCompile(`<!--.*?-->`)
	body = reComments.ReplaceAll(body, []byte(""))

	// Optionally, remove all HTML tags
	reTags := regexp.MustCompile(`<[^>]+>`)
	text := reTags.ReplaceAll(body, []byte(""))
	// fmt.Println("Processed Body:", string(text))
	// fmt.Println()
	// fmt.Println("Processed Body Length:", len(text))
	// fmt.Println()

	// Truncate to max length
	maxLen := 10000
	if len(text) > maxLen {
		text = text[:maxLen]
	}
	return text
}

func strip(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
	return content

}
