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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// detailsCmd represents the details command
var detailsCmd = &cobra.Command{
	Use:   "details [youtube-url]",
	Short: "Get the details of a YouTube video",
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

		videoResponse, err := service.Videos.List([]string{"contentDetails", "snippet", "statistics"}).Id(videoId).Do()
		handleError(err)

		if len(videoResponse.Items) == 0 {
			fmt.Println("No video found for the given ID")
			return
		}

		video := videoResponse.Items[0]
		fmt.Println("Title:", video.Snippet.Title)
		fmt.Println("Description:", video.Snippet.Description)
		fmt.Println("View Count:", video.Statistics.ViewCount)

		// Fetching top comments
		commentResponse, err := service.CommentThreads.List([]string{"snippet"}).VideoId(videoId).Order("relevance").MaxResults(5).Do()
		if err != nil {
			handleError(err)
			os.Exit(1)
		}

		fmt.Println("Top Comments:")
		for _, item := range commentResponse.Items {
			comment := item.Snippet.TopLevelComment.Snippet
			fmt.Printf("Author: %s\nComment: %s\nLikes: %d\n\n", comment.AuthorDisplayName, comment.TextDisplay, comment.LikeCount)
		}
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// detailsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// detailsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
