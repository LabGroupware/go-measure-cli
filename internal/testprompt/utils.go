package testprompt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"golang.org/x/exp/rand"
)

// 数値を入力させるプロンプト
func PromptNumber(label string) (int, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateNumber,
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	value, _ := strconv.Atoi(result)
	return value, nil
}

// 文字列選択プロンプト
func PromptSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

// 時間間隔入力プロンプト
func PromptMillisecond(label string) (time.Duration, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateNumber,
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	ms, _ := strconv.Atoi(result)
	return time.Duration(ms) * time.Millisecond, nil
}

// 真偽値入力プロンプト
func PromptBool(label string) (bool, error) {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

// テキスト入力プロンプト
func PromptInput(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

// 数値バリデーション
func validateNumber(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("invalid input: %s", input)
	}
	return nil
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
