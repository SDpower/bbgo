package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/c9s/bbgo/pkg/cmd/cmdutil"
)

//  Config struct we can move into config after it work fine.
type storgeEngine struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
	Host   string `json:"host,omitempty" yaml:"host,omitempty"`
	Token  string `json:"token,omitempty" yaml:"token,omitempty"`
	Bucket string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
}
type streamSubscription struct {
	On     string   `json:"on,omitempty" yaml:"on,omitempty"`
	Symbol []string `json:"symbol,omitempty" yaml:"symbol,omitempty"`
}
type StreamStorge struct {
	Engine       storgeEngine
	Subscription []streamSubscription
}
type StreamStorgeConfig struct {
	StreamStorge StreamStorge `json:"streamStorge,omitempty" yaml:"streamStorge,omitempty"`
}

var streamConfig *StreamStorgeConfig

// go run ./cmd/bbgo stream-subscribe
var streamSubscribeCmd = &cobra.Command{
	Use:   "stream-subscribe",
	Short: "Stream Subscribe to storge",
	Long:  "This command will subscribe tick data stor in streamstorge from exchange",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		// load config file nicely
		if len(configFile) > 0 {
			// if config file exists, use the config loaded from the config file.
			// otherwise, use a empty config object
			if _, err := os.Stat(configFile); err == nil {
				// load successfully
				content, err := ioutil.ReadFile(configFile)
				if err != nil {
					return err
				}

				if err := yaml.Unmarshal(content, &streamConfig); err != nil {
					return err
				}

			} else if os.IsNotExist(err) {
				return errors.New("config file doesn't exist, we should use the empty config")
			} else {
				// other error
				return errors.Wrapf(err, "config file load error: %s", configFile)
			}
		}

		//config check sure like tihs
		emptyConfig := &StreamStorgeConfig{}
		if fmt.Sprintf("%v", emptyConfig) == fmt.Sprintf("%v", streamConfig) {
			return errors.New("--config option is required")
		}
		log.Debugf("configFile: %+v", configFile)
		log.Debugf("streamConfig: %+v", streamConfig)
		cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
		return nil
	},
}

func init() {
	//override default "config" Flag values
	streamSubscribeCmd.Flags().String("config", "config/stream-storge.yaml", "stream-storge config file")
	RootCmd.AddCommand(streamSubscribeCmd)
}
