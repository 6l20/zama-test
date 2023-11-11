/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/6l20/zama-test/client"
	"github.com/6l20/zama-test/client/config"
	"github.com/6l20/zama-test/common/log/zap"
	"github.com/6l20/zama-test/common/merkle"
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

		files, err := os.ReadDir(fileDirPath)
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

		logger.Debug("Starting client")

		merkleManager := merkle.NewMerkleManager(fileDirPath+ "merkle-root.txt", logger)

		merkleTree, err := merkleManager.BuildMerkleTreeFromFS(fileDirPath)
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

		logger.Info("Merkle Tree", "root Hash", merkleTree.Root.Hash)

		proof := merkleManager.GenerateProof(1)
		logger.Info("Merkle Tree", "proof",  proof)


		verified := merkleManager.VerifyProof("57a7503b110edb69d272202911dcf347ef82f80eb71f307cc67af768baca92ca",*proof, merkleTree.Root.Hash)

		logger.Info("Merkle Tree", "verified", verified)

		merkleManager.StoreMerkleRoot()
	
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			err = client.UploadFile(fileDirPath + f.Name(), clientConfig.UploadUrl)
		if err != nil {
			logger.Error("Error uploading file:", err)
			return
		}
		logger.Info("File uploaded successfully", "file", f.Name())
		}
	},	
		
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
