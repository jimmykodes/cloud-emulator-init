package runner

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/jimmykodes/cloud-emulator-init/internal/atomicerr"
	"github.com/jimmykodes/cloud-emulator-init/internal/aws"
	"github.com/jimmykodes/cloud-emulator-init/internal/config"
	"github.com/jimmykodes/cloud-emulator-init/internal/gcp"
)

const (
	configFileFlag = "config-file"
)

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("runner", pflag.ExitOnError)
	fs.String(configFileFlag, "conf.yaml", "file containing resources to create")

	return fs
}

func RunE(cmd *cobra.Command, _ []string) error {
	configFile := viper.GetString(configFileFlag)
	if configFile == "" {
		return fmt.Errorf("missing config file location")
	}

	conf, err := config.New(configFile)
	if err != nil {
		return err
	}

	var (
		atomicErr atomicerr.Error
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		atomicErr.Append(aws.RunE(cmd.Context(), conf.AWS))
	}()

	go func() {
		defer wg.Done()
		atomicErr.Append(gcp.RunE(cmd.Context(), conf.GCP))
	}()

	wg.Wait()
	return atomicErr.Err()
}
