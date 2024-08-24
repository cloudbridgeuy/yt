package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type key string

const apiKeyCtxKey = key("apiKey")

func createContext() context.Context {
	apiKey := viper.GetString("youtube_api_key")

	if apiKey == "" {
		handleError(fmt.Errorf("Please provide a YouTube API key."))
	}

	ctx := context.Background()
	return context.WithValue(ctx, apiKeyCtxKey, apiKey)
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("%v", err.Error())
		os.Exit(1)
	}
}

func getVideoId(url string) (string, error) {
	pattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if len(match) == 0 {
		return "", fmt.Errorf("No video ID found in URL")
	}
	return match[1], nil
}

func getService(ctx context.Context) *youtube.Service {
	apiKey, ok := ctx.Value(apiKeyCtxKey).(string)
	if !ok {
		handleError(fmt.Errorf("API key not found in context"))
	}
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	handleError(err)
	return service
}

func parseDuration(durationStr string) (int, error) {
	matches := regexp.MustCompile(`(?i)PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`).FindStringSubmatch(durationStr)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration string: %s", durationStr)
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])

	return hours*60 + minutes + seconds/60, nil
}
