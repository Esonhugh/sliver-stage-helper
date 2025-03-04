package starListen

import "github.com/spf13/cobra"

func init() {
	StarListenCmd.Flags().StringVar(&Opt.BeaconProfile, "profile", "default", "beacon profile")
}

var Opt = struct {
	BeaconProfile    string
	RawBytesFromFile string
}{}

var StarListenCmd = &cobra.Command{
	Use:   "startListen",
	Short: "startListen",
	Long:  "startListen",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}
