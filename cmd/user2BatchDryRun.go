/*
Copyright © 2023 Zhang Zhanpeng <zhangregister@outlook.com>

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
	"net/http"
	"net/url"
	"strconv"

	"github.com/MonteCarloClub/log"
	"github.com/spf13/cobra"
)

const (
	User2BatchDryRunUrlFormat = "http://10.176.40.46:8080/abe/dabe/user2_batch_dry_run?batch_size=%v"
)

var (
	batchSize int
)

// user2BatchDryRunCmd represents the user2BatchDryRun command
var user2BatchDryRunCmd = &cobra.Command{
	Use:   "user2BatchDryRun",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf(User2BatchDryRunUrlFormat, url.QueryEscape(strconv.Itoa(batchSize)))
		log.Info("request server...", "url", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error("fail to GET server", "url", url, "err", err)
			return
		}
		if resp.StatusCode == http.StatusOK {
			log.Info("/user2_batch_dry_run: ok")
		} else {
			log.Warn("/user2_batch_dry_run: not ok", "status code", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(user2BatchDryRunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// user2BatchDryRunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// user2BatchDryRunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	user2BatchDryRunCmd.Flags().IntVarP(&batchSize, "batch-size", "s", 10000000, "Count of connection(s) simulated by zdyf3, default 10000000")
}
