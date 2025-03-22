package list

import (
	"context"
	"encoding/json"

	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient/protobuf/commonpb"
	log "github.com/sirupsen/logrus"

	c "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	"github.com/spf13/cobra"
)

func init() {

	ListCmd.AddCommand(ListProfileCmd, ListJobsCmd)
	c.RootCmd.AddCommand(ListCmd)
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  "list",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func toJson(v interface{}) string {
	s, _ := json.MarshalIndent(v, "\t", "  ")
	return string(s)
}

var ListProfileCmd = &cobra.Command{
	Use:     "profiles",
	Aliases: []string{"profile", "p"},
	Short:   "list profiles",
	Long:    "list profiles",
	Run: func(cmd *cobra.Command, args []string) {
		pfs := c.Client.ListImplantProfiles()
		if len(pfs) == 0 {
			log.Infof("Empty sliver implant profiles")
			return
		}
		for _, pf := range pfs {
			log.Infof("Implant Profile Name: %s\nDetails:\n\t%s\n", pf.Name, toJson(pf.Config))
		}
	},
}

var ListJobsCmd = &cobra.Command{
	Use:     "jobs",
	Aliases: []string{"job", "j"},
	Short:   "list jobs",
	Long:    "list jobs",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, err := c.Client.GetJobs(context.Background(), &commonpb.Empty{})
		if err != nil {
			log.Fatalf("Error getting jobs: %v", err)
			return
		}
		rjobs := jobs.GetActive()
		if len(rjobs) == 0 {
			log.Warnf("No jobs found")
		}
		for _, job := range rjobs {
			log.Infof("Job ID: %v\nDetails:\n\t%s\n", job.ID, toJson(job))
		}
	},
}
