package starListen

import (
	ctx "context"
	"os"

	c "github.com/Esonhugh/sliver-linux-tcp-stager-helper/cmd/sliverStager/cmd"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	StartListenCmd.Flags().StringVar(&Opt.BeaconProfile, "profile", "default", "beacon profile")
	StartListenCmd.Flags().StringVar(&Opt.RawBytesFromFile, "rawBytesFromFile", "", "raw bytes from file")
}

var log = logrus.WithField("cmd", "startListen")

var Opt = struct {
	BeaconProfile    string
	RawBytesFromFile string
}{}

var StartListenCmd = &cobra.Command{
	Use:   "startListen",
	Short: "startListen",
	Long:  "startListen",
	Run: func(cmd *cobra.Command, args []string) {

		var data []byte
		if Opt.RawBytesFromFile != "" {
			var err error
			data, err = os.ReadFile(Opt.RawBytesFromFile)
			if err != nil {
				log.Fatalf("read file %s failed: %v", Opt.RawBytesFromFile, err)
			}
		}

		if Opt.BeaconProfile != "" {
			var in *clientpb.ImplantProfile
			in.Name = Opt.BeaconProfile
			res, err := c.Client.SaveImplantProfile(ctx.Background(), in)
			if err != nil {
				log.Fatalf("profile %s saved failed, err: %v", Opt.BeaconProfile, err)
			}
			// ToDo:
			_ = res.GetName()
		}

		in := &clientpb.StagerListenerReq{
			Host: "",
			Port: 0,
			Data: data,
		}
		c.Client.StartTCPStagerListener(ctx.Background(), in)
	},
}
