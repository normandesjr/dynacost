package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/normandesjr/dynacost/pkg/configkeys"
	"github.com/normandesjr/dynacost/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dynacost <table1>...<tablen>",
	Short: "Interactive DynamoDB Cost",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootAction(args)
	},
}

func rootAction(args []string) error {
	tl := &scan.TableList{}
	for _, t := range args {
		tl.Add(t)
	}

	// TODO: start app
	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dynacost.yaml)")

	rootCmd.Flags().DurationP(configkeys.Interval, "i", 5*time.Minute, "minutes to wait between updates")
	rootCmd.Flags().StringP(configkeys.Region, "r", "sa-east-1", "default AWS region")
	rootCmd.Flags().StringP(configkeys.Profile, "p", "", "AWS profile")

	viper.BindPFlag(configkeys.Interval, rootCmd.Flags().Lookup(configkeys.Interval))
	viper.BindPFlag(configkeys.Region, rootCmd.Flags().Lookup(configkeys.Region))
	viper.BindPFlag(configkeys.Profile, rootCmd.Flags().Lookup(configkeys.Profile))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dynacost" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dynacost")
	}

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("DYNACOST")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
