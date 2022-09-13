package cmd

import (
	"PhishRod/internal/resource"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var webInjectDirectories []string
var resourceMatchConfiguration string
var outputFile string

func init() {

	packageCmd.Flags().StringSliceVarP(
		&webInjectDirectories,
		"directory",
		"d",
		[]string{},
		"Directory or directories (comma-separated list) of web-injects",
	)
	packageCmd.MarkFlagRequired("directory")

	packageCmd.Flags().StringVarP(
		&resourceMatchConfiguration,
		"matchfile",
		"f",
		"conf.xml",
		"Resource-matching configuration file path",
	)
	packageCmd.MarkFlagRequired("matchfile")

	packageCmd.Flags().StringVarP(
		&outputFile,
		"outputfile",
		"o",
		"",
		"Output file for encrypted resource archive",
	)
	packageCmd.MarkFlagRequired("outputfile")
}

func buildPackage(cmd *cobra.Command, args []string) {
	builder := resource.Builder{}
	// add matchfile and web-inject directories to builder, then builder.Build()
	_, err := os.Stat(resourceMatchConfiguration)
	if err != nil {
		fmt.Printf("error while opening resource match configuration file (%s): %v", resourceMatchConfiguration, err)
		os.Exit(1)
	}
	err = builder.AddMatchFile(resourceMatchConfiguration)
	if err != nil {
		fmt.Printf("error adding resource match configuration file (%s): %v", resourceMatchConfiguration, err)
		os.Exit(1)
	}

	for _, dir := range webInjectDirectories {
		if _, err = os.Stat(dir); err != nil {
			fmt.Printf("error while checking web-inject directory (%s): %v", dir, err)
			os.Exit(1)
		}
		err = builder.AddDirectory(dir)
		if err != nil {
			fmt.Printf("error while adding web-inject directory (%s): %v", dir, err)
			os.Exit(1)
		}
	}

	content, err := builder.Build()
	if err != nil {
		fmt.Printf("failed to build encrypted resource achive (%s): %v", outputFile, err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(outputFile, content, os.FileMode(644))
	if err != nil {
		fmt.Printf("failed to write output file (%s), %v", outputFile, err)
		os.Exit(1)
	}
}
