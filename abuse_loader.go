package masker

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// AbuseWordLoader 提供從不同來源載入敏感詞清單的方法。
type AbuseWordLoader struct{}

// NewAbuseWordLoader 建立一個新的敏感詞載入器。
func NewAbuseWordLoader() *AbuseWordLoader {
	return &AbuseWordLoader{}
}

// LoadFromFile 從本機檔案載入敏感詞，每行一個詞，空行與 # 開頭的註解行會被略過。
func (l *AbuseWordLoader) LoadFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return l.loadFromReader(file)
}

// LoadFromReader 從 io.Reader 載入敏感詞，每行一個詞，空行與 # 開頭的註解行會被略過。
func (l *AbuseWordLoader) LoadFromReader(reader io.Reader) ([]string, error) {
	return l.loadFromReader(reader)
}

// LoadFromString 從字串載入敏感詞，每行一個詞，空行與 # 開頭的註解行會被略過。
func (l *AbuseWordLoader) LoadFromString(content string) ([]string, error) {
	reader := strings.NewReader(content)
	return l.loadFromReader(reader)
}

// LoadFromSlice 從字串 slice 載入敏感詞，會去除前後空白並略過空字串，適用於詞清單已在記憶體中的情況。
func (l *AbuseWordLoader) LoadFromSlice(words []string) []string {
	var result []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			result = append(result, word)
		}
	}
	return result
}

// loadFromReader 是處理 reader 的內部方法。
func (l *AbuseWordLoader) loadFromReader(reader io.Reader) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
