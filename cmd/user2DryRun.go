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
	"sync"
	"sync/atomic"
	"time"

	"github.com/MonteCarloClub/log"
	"github.com/spf13/cobra"
)

const (
	User2DryRunUrlFormat = `http://10.176.40.46:8080/abe/dabe/user2_dry_run?filename=%v&password=%v`
	CountOfRetry         = 3
)

var (
	filename          string
	password          string
	count             int
	waitQueueCapacity int64
	logStep           int
)

// user2DryRunCmd represents the user2DryRun command
var user2DryRunCmd = &cobra.Command{
	Use:   "user2DryRun",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf(User2DryRunUrlFormat, url.QueryEscape(filename), url.QueryEscape(password))
		log.Info("request server...", "url", url, "count", count, "log step", logStep)
		var wg sync.WaitGroup
		// "1" - print end log
		wg.Add(count + 1)
		startTs := time.Now().UnixNano()

		// this goroutine prints logs
		var retryCount int64
		statusChan := make(chan bool, count)
		go func() {
			var latestLogCount int
			for {
				reqCount := len(statusChan)
				if reqCount >= count {
					endTs := time.Now().UnixNano()
					timeInS := (endTs - startTs) / 1e9
					log.Info("all request(s) finished", "count of request(s)", reqCount, "time/s", timeInS)
					log.Info("test results", "retry count", retryCount)
					wg.Done()
					return
				}
				if reqCount > 0 && reqCount < count && reqCount >= logStep*(latestLogCount/logStep+1) {
					latestLogCount = reqCount
					log.Info(fmt.Sprintf("%v request(s) finished", reqCount))
				}
			}
		}()

		// this goroutine tries to connect
		var waitQueueSize int64
		go func() {
			for {
				if waitQueueSize >= waitQueueCapacity {
					continue
				}
				resp, _ := http.Get(url)
				if resp != nil && resp.StatusCode == http.StatusOK {
					statusChan <- true
				} else {
					atomic.AddInt64(&waitQueueSize, 1)
					// retry infinitely
					for j := 0; j < CountOfRetry; j++ {
						atomic.AddInt64(&retryCount, 1)
						retryResp, _ := http.Get(url)
						if retryResp != nil && retryResp.StatusCode == http.StatusOK {
							statusChan <- true
							atomic.AddInt64(&waitQueueSize, -1)
							return
						}
					}
					statusChan <- false
					atomic.AddInt64(&waitQueueSize, -1)
				}
				wg.Done()
			}
		}()
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(user2DryRunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// user2DryRunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// user2DryRunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	user2DryRunCmd.Flags().StringVarP(&filename, "filename", "f", "filename", "Filename, default filename")
	user2DryRunCmd.Flags().StringVarP(&password, "password", "p", "password", "Password, default password")
	user2DryRunCmd.Flags().IntVarP(&count, "count", "c", 10000000, "Count of connection(s) to zdyf3, default 10000000")
	user2DryRunCmd.Flags().Int64VarP(&waitQueueCapacity, "wait-queue-capacity", "w", 1, "Maximum of connection(s) waiting, default 1")
	user2DryRunCmd.Flags().IntVarP(&logStep, "log-step", "l", 1, "Print log at every step, default 1 - print log at every request")
}
