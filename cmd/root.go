package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	cfgParallel int
)

var rootCmd = &cobra.Command{
	Use:   "bagel",
	Short: "Yogurt - Configuration Management",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /opt/bagel/bagel.yaml)")

	rootCmd.PersistentFlags().IntVarP(&cfgParallel, "parallel", "p", 10, "number of parallel threads")
	viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel"))

	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("bagel")
		viper.AddConfigPath("/opt/bagel")
		viper.AddConfigPath("$HOME/.bagel")
		viper.AddConfigPath(".")
	}

	viper.ReadInConfig()

	viper.SetDefault("site_dir", "/opt/bagel")
}

func Execute() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
