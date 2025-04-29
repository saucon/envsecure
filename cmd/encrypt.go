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
	SilenceErrors: false,
	SilenceUsage:  false,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

		keyFile, err := cmd.Flags().GetString("keyfile")
		if err != nil {
			return err
		}
		algo, err := cmd.Flags().GetString("algo")
		if err != nil {
			return err
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
				return errors.New("error read key file, " + err.Error())
			}
			keySecret = string(fileBytes)
		default:
			return err
		}

		if keySecret == "" {
			return errors.New("error key or keyfile is required")
		}

		err = encrypt.EncryptEnv(secureAlgo, filePath, keySecret)
		if err != nil {
			return err
		}

		fp := strings.Split(filePath, "/")
		filename := fp[len(fp)-1]

		cmd.Println()
		cmd.Println("Config encrypted and saved successfully. File: secure." + filename)
		return nil
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
