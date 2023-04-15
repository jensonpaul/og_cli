package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opsgenie",
	Short: "A client for querying data from Opsgenie",
	Long: `An Opsgenie client.
	
This cli uses viper to configure it. In order to provide an API key for Opsgenie you can either
expose an env var with the name 'OPS_APIKEY' or you can create a config file in your working directory 
or your home folder. The home folder config is stored in a hidden folder called '.opsgenie' and the config
is called 'config.yaml'.

The key for the API key to set in the config is called 'apikey'.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("apiKey")
		if apiKey == "" {
			return fmt.Errorf("no API key set")
		}
		opsgenieConfig = &client.Config{ApiKey: apiKey}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string
var opsgenieConfig *client.Config

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opsgenie/config.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "Your Opsgenie API key")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetEnvPrefix("OPS")
	viper.BindEnv("apiKey")
	viper.BindPFlag("apiKey", rootCmd.Flags().Lookup("api-key"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		confDir := filepath.Join(home, ".opsgenie")

		if err := os.Mkdir(confDir, os.ModeDir|os.ModePerm); err != nil && !os.IsExist(err) {
			log.Printf("error creating config directory: %v", err)
			os.Exit(1)
		}

		viper.AddConfigPath(confDir)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				log.Printf("error writing config file: %v", err)
				os.Exit(1)
			}
			return
			// if err := viper.ReadInConfig(); err != nil {
			// 	log.Println("Can't read config:", err)
			// 	os.Exit(1)
			// }
		}
		log.Println("Can't read config:", err)
		os.Exit(1)
	}
}
