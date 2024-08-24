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
	"os"

	"github.com/spf13/cobra"
)

// commentsCmd represents the comments command
var commentsCmd = &cobra.Command{
	Use:   "comments [youtube-url]",
	Short: "Get comments from a YouTube video",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := createContext()

		url := args[0]

		videoId, err := getVideoId(url)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}

		service := getService(ctx)

		var comments []string
		call := service.CommentThreads.List([]string{"snippet", "replies"}).VideoId(videoId).TextFormat("plainText").MaxResults(100)
		response, err := call.Do()
		handleError(err)

		for _, item := range response.Items {
			topLevelComment := item.Snippet.TopLevelComment.Snippet.TextDisplay
			comments = append(comments, topLevelComment)

			if item.Replies != nil {
				for _, reply := range item.Replies.Comments {
					replyText := reply.Snippet.TextDisplay
					comments = append(comments, "    - "+replyText)
				}
			}
		}

		jsonOutput, err := json.MarshalIndent(comments, "", " ")
		handleError(err)

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	rootCmd.AddCommand(commentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
