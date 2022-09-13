package cmd

import (
	"PhishRod/internal/resource"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var inputArchivePath string
var outputDirectoryPath string

func init() {

	unpackCmd.Flags().StringVarP(
		&outputDirectoryPath,
		"output",
		"d",
		"",
		"Output directory",
	)
	unpackCmd.MarkFlagRequired("output")

	unpackCmd.Flags().StringVarP(
		&inputArchivePath,
		"archive",
		"i",
		"",
		"Input file as encrypted resource archive",
	)
	unpackCmd.MarkFlagRequired("archive")
}

func unpackArchive(cmd *cobra.Command, args []string) {
	archive := resource.Archive{Filepath: inputArchivePath}
	err := archive.Unpack(outputDirectoryPath)
	if err != nil {
		fmt.Printf("error unpacking archive file (%s) to directory (%s): %v\n", inputArchivePath, outputDirectoryPath, err)
		os.Exit(1)
	}
}
