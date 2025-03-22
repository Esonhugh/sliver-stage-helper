package cmd

import (
	"fmt"
	"os"

	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Opts = struct {
		Verbose        int
		ConfigFilePath string
	}{}
	Client *sliverClient.Client
)

func init() {
	RootCmd.PersistentFlags().CountVarP(&Opts.Verbose, "verbose", "v", "verbose level (-v debug | -vv trace)")
	RootCmd.PersistentFlags().StringVarP(&Opts.ConfigFilePath, "config", "c", os.Getenv("SLIVER_CLIENT_CONFIG"), "config file path")
}

var RootCmd = &cobra.Command{
	Use:   "sliverStager",
	Short: "Sliver Stager Helper",
	Long:  "Sliver Stager Helper",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(Opts.Verbose)
		cfg, err := sliverClient.ReadConfig(Opts.ConfigFilePath)
		if err != nil {
			log.Fatalf("can't load sliver config under %v, err: %v", Opts.ConfigFilePath, err)
		}
		Client, err = sliverClient.NewClient(cfg)
		if err != nil {
			log.Fatalf("can't create sliver client, err: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		Client.Close()
	},
}

var LevelMap = []log.Level{
	log.InfoLevel,  // 0
	log.DebugLevel, // 1
	log.TraceLevel, // 2
}

func SetLogLevel(level int) {
	if level > 2 {
		level = 2
	}
	log.SetLevel(LevelMap[level])
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
