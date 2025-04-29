package encrypt

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Secure interface {
	Encrypt(plainText string, secretKey string) (string, error)
	Decrypt(cipherText string, secretKey string) (string, error)
}

func EncryptEnv(secure Secure, filepath string, secretKey string) error {
	var configMap map[string]any

	fp := strings.Split(filepath, "/")
	filename := fp[len(fp)-1]

	// 1. Read YAML file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading file, %s", err.Error()))
	}

	// 2. Unmarshal into a map
	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return errors.New(fmt.Sprintf("error unmarshalling data, %s", err.Error()))
	}
	//printMap(configMap, "")

	err = encryptInterface(secure, configMap, secretKey)
	if err != nil {
		return err
	}

	// Create the output file
	file, err := os.Create("secure." + filename)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating file, %s", err.Error()))
	}
	defer file.Close()

	// Create a YAML encoder
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	// Encode the modified node
	if err = encoder.Encode(configMap); err != nil {
		return errors.New(fmt.Sprintf("error encoding YAML, %s", err.Error()))
	}

	return nil
}

// encryptInterface encrypts values inside maps or slices recursively
func encryptInterface(secure Secure, i interface{}, secretKey string) error {
	switch v := i.(type) {
	case map[string]interface{}:
		for k, val := range v {
			switch child := val.(type) {
			case map[string]interface{}, []interface{}:
				// Recurse deeper
				err := encryptInterface(secure, child, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting interface, %s", err.Error()))
				}
			case string:
				enc, err := secure.Encrypt(child, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting key %s, %s", k, err.Error()))
				}
				v[k] = enc
			case nil: // do nothing for nil
			default:
				s := fmt.Sprintf("%v", child)
				enc, err := secure.Encrypt(s, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting key %s, %s", k, err.Error()))
				}
				v[k] = enc
			}
		}
	case []interface{}:
		for idx, item := range v {
			switch child := item.(type) {
			case map[string]interface{}, []interface{}:
				err := encryptInterface(secure, child, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting interface, %s", err.Error()))
				}
			case string:
				enc, err := secure.Encrypt(child, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting key, %s", err.Error()))
				}
				v[idx] = enc
			case nil: // do nothing for nil
			default:
				s := fmt.Sprintf("%v", child)
				enc, err := secure.Encrypt(s, secretKey)
				if err != nil {
					return errors.New(fmt.Sprintf("error encrypting key, %s", err.Error()))
				}
				v[idx] = enc
			}
		}
	default:
		// do nothing for other types
	}
	return nil
}
