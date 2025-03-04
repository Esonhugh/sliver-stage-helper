package cmd

import (
	"fmt"
	"os"

	"github.com/Esonhugh/sliver-linux-tcp-stager-helper/pkg/sliverClient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Opts = struct {
		Verbose        int
		ConfigFilePath string
		StagerType     string
		ListenerURL    string
		Advanced       string
		Format         string
	}{}
	Client *sliverClient.Client
)

func init() {
	RootCmd.PersistentFlags().IntVarP(&Opts.Verbose, "verbose", "v", 0, "verbose level (-v debug | -vv trace)")
	RootCmd.PersistentFlags().StringVarP(&Opts.ConfigFilePath, "config", "c", "~/.sliver", "config file path")
	RootCmd.PersistentFlags().StringVarP(&Opts.StagerType, "stager-type", "t", "linux-x64-tcp", "stager type: linux-x64-tcp")
	RootCmd.PersistentFlags().StringVarP(&Opts.ListenerURL, "url", "u", "tcp://127.1:4444", "listener url, like tcp://127.0.0.1:4444")
	RootCmd.PersistentFlags().StringVarP(&Opts.Format, "format", "f", "raw", "format")
	RootCmd.PersistentFlags().StringVarP(&Opts.Advanced, "advanced", "a", "", "advanced options")
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
