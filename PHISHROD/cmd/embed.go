package cmd

import (
	"PhishRod/internal"
	"PhishRod/internal/resource"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var targetFilePath string
var resourceFilePath string
var resourceName string

func init() {
	embedCmd.Flags().StringVarP(
		&targetFilePath,
		"target",
		"i",
		"",
		"Target Portable Executable (PE) file where the resource will be embedded",
	)
	embedCmd.MarkFlagRequired("target")

	embedCmd.Flags().StringVarP(
		&resourceFilePath,
		"resource",
		"f",
		"",
		"Resource file that will be embedded in the target file",
	)
	embedCmd.MarkFlagRequired("resource")

	embedCmd.Flags().StringVarP(
		&resourceName,
		"name",
		"n",
		"500",
		"Name of resource to embed in target file",
	)
}

func embedResource(cmd *cobra.Command, args []string) {
	for _, file := range []string{targetFilePath, resourceFilePath} {
		exists, err := internal.FileExists(file)
		if err != nil {
			fmt.Printf("error while checking existince of file (%s): %v\n", file, err)
			os.Exit(1)
		}

		if !exists {
			fmt.Printf("[-] Required file does not exist: %s\n", file)
		}
	}

	err := resource.Embed(targetFilePath, resourceFilePath, resourceName)
	if err != nil {
		fmt.Printf("error while embedding resource file (%s) into target file (%s): %v\n",
			resourceFilePath,
			targetFilePath,
			err,
		)
		os.Exit(1)
	}

	fmt.Printf("[+] Embedded resource (%s) into file (%s)\n", resourceFilePath, targetFilePath)
}
