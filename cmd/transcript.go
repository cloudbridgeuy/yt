/*
Copyright © 2024 Guzmán Monné guzman.monne@cloudbridge.com.uy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/spf13/cobra"
)

// transcriptCmd represents the transcript command
var transcriptCmd = &cobra.Command{
	Use:   "transcript [youtube-url]",
	Short: "Generate a transcript from a YouTube URL",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		videoId, err := getVideoId(url)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}

		var transcriptText string
		transcript, err := getTranscript(videoId)
		handleError(err)

		doc := soup.HTMLParse(transcript)
		textTags := doc.FindAll("text")
		var textBuilder strings.Builder
		for _, textTag := range textTags {
			textBuilder.WriteString(textTag.Text())
			textBuilder.WriteString(" ")
			transcriptText = textBuilder.String()
		}

		if transcriptText == "" {
			handleError(fmt.Errorf("Can't found transcription or it is empty"))
		}

		parsedString := html.UnescapeString(transcriptText)
		fmt.Println(parsedString)
	},
}

func init() {
	rootCmd.AddCommand(transcriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transcriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transcriptCmd.Flags().StringP("lang", "l", "en", "Language for the transcript")
}

func getTranscript(videoID string) (string, error) {
	url := "https://www.youtube.com/watch?v=" + videoID
	resp, err := soup.Get(url)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(resp)
	scriptTags := doc.FindAll("script")
	for _, scriptTag := range scriptTags {
		if strings.Contains(scriptTag.Text(), "captionTracks") {
			regex := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
			match := regex.FindStringSubmatch(scriptTag.Text())
			if len(match) > 1 {
				var captionTracks []struct {
					BaseURL string `json:"baseUrl"`
				}
				json.Unmarshal([]byte(match[1]), &captionTracks)
				if len(captionTracks) > 0 {
					transcriptURL := captionTracks[0].BaseURL
					transcriptResp, err := soup.Get(transcriptURL)
					if err != nil {
						return "", err
					}
					return transcriptResp, nil
				}
			}
		}
	}
	return "", fmt.Errorf("transcript not found")
}
