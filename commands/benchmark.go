// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"github.com/spf13/cobra"
	"os"
	"runtime/pprof"
)

var cpuProfilefile string
var memProfilefile string
var benchmarkTimes int

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark hugo by building a site a number of times.",
	Long: `Hugo can build a site many times over and analyze the running process
creating a benchmark.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InitializeConfig(); err != nil {
			return err
		}

		return bench(cmd, args)
	},
}

func init() {
	benchmarkCmd.Flags().StringVar(&cpuProfilefile, "cpuprofile", "", "path/filename for the CPU profile file")
	benchmarkCmd.Flags().StringVar(&memProfilefile, "memprofile", "", "path/filename for the memory profile file")

	benchmarkCmd.Flags().IntVarP(&benchmarkTimes, "count", "n", 13, "number of times to build the site")
}

func bench(cmd *cobra.Command, args []string) error {

	if memProfilefile != "" {
		f, err := os.Create(memProfilefile)

		if err != nil {
			return err
		}
		for i := 0; i < benchmarkTimes; i++ {
			_ = buildSite()
		}
		pprof.WriteHeapProfile(f)
		f.Close()

	} else {
		if cpuProfilefile == "" {
			cpuProfilefile = "/tmp/hugo-cpuprofile"
		}
		f, err := os.Create(cpuProfilefile)

		if err != nil {
			return err
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		for i := 0; i < benchmarkTimes; i++ {
			_ = buildSite()
		}
	}

	return nil

}
