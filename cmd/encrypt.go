/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"github.com/saucon/envsecure/pkg/encrypt"
	"github.com/saucon/sauron/v2/pkg/secure"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt an env file",
	Long: `Encrypt an env file now only support yaml format. This command only support 
RSA and AEAD algorithms. For example:

envsecure encrypt -f sample/env.sample.yml --algo rsa --keyfile sample/public_key.pem
`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		keyFile, err := cmd.Flags().GetString("keyfile")
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		algo, err := cmd.Flags().GetString("algo")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		var keySecret string
		var secureAlgo secure.Secure
		switch algo {
		case "aead":
			secureAlgo = secure.NewSecureAEAD()
			keySecret = key
		case "rsa":
			secureAlgo = secure.NewSecureRSA()
			fileBytes, err := os.ReadFile(keyFile)
			if err != nil {
				cmd.PrintErr(errors.New("error read key file, " + err.Error()))
				return
			}
			keySecret = string(fileBytes)
		default:
			cmd.PrintErr(errors.New("invalid algorithm. Options: aead, rsa"))
			return
		}

		if keySecret == "" {
			cmd.PrintErr(errors.New("key or keyfile is required"))
			return
		}

		err = encrypt.EncryptEnv(secureAlgo, filePath, keySecret)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		fp := strings.Split(filePath, "/")
		filename := fp[len(fp)-1]

		cmd.Println("Config encrypted and saved successfully. File: secure." + filename)
		return
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	encryptCmd.Flags().StringP("file", "f", "", "File to encrypt. This file contains the config with plain text")
	encryptCmd.MarkFlagRequired("file")

	encryptCmd.Flags().StringP("key", "k", "", "Encryption key")

	encryptCmd.Flags().StringP("keyfile", "", "", "Encryption key file")

	encryptCmd.Flags().StringP("algo", "", "", "Algorithm to use for encryption. Options: aead, rsa")
	encryptCmd.MarkFlagRequired("algo")
}
