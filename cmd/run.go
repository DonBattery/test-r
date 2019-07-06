package cmd

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/donbattery/test-r/model"
	"github.com/spf13/cobra"
)

var timeout int64

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run runs the test ofc",
	Run: func(cmd *cobra.Command, args []string) {
		runTest()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().Int64VarP(&timeout, "timeout", "T", 10, "Timeout of the test in seconds")
}

func runTest() {
	pool := model.NewPool()
	for i := 0; i < 10; i++ {
		pool.AddJob(model.NewTask())
	}
	pool.InitJobs()
	pool.WaitAll()
	log.Info("All task finished gg bb")

	rep := model.Report{
		Apps: []model.App{
			model.App{
				Name:     "App 1",
				Duration: time.Second * 111,
				Builds: []*model.Build{
					&model.Build{
						Name:     "Build 1",
						Duration: time.Second * 15,
						Result:   1,
					},
					&model.Build{
						Name:     "Build 2",
						Duration: time.Second * 13,
						Result:   0,
					},
				},
			},
		},
	}

	fmt.Printf("%#v", rep)
	// model.JUnitReportXML(&rep, false, os.Stdout)
}
