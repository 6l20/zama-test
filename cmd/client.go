/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/6l20/zama-test/client"
	"github.com/6l20/zama-test/client/config"
	"github.com/6l20/zama-test/client/stores"
	"github.com/6l20/zama-test/client/usecases"
	"github.com/6l20/zama-test/common/log/zap"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
		fileDirPath := "test/data/"
		clientConfig, err := config.EnvClientConfig()
		if err != nil {
			fmt.Println("Error getting client config:", err)
			return
		}
		logConfig, err := zap.EnvConfig()
		if err != nil {
			fmt.Println("Error getting log config:", err)
			return
		}
		logger, err := zap.NewLogger(logConfig)
		if err != nil {
			fmt.Println("Error getting logger:", err)
			return
		}

		store := stores.NewLocalStore(fileDirPath, clientConfig.RootFile)

		client, err := client.NewClient(logger, *clientConfig, store)
		if err != nil {
			logger.Error("Error creating client:", err)
		}

		useCases := usecases.NewClientUseCases(logger, client)

		useCases.GenerateMerkleTree()

		useCases.UploadFilesFromDir(fileDirPath)

		proof, err := useCases.GetMerkleProofForFile(1)
		if err != nil {
			logger.Error("Error getting proof:", err)
		}
		logger.Info("Proof:", "proof", proof)

		useCases.VerifyMerkleProof("")
	},	
		
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
