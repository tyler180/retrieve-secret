package retrievesecrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// SecretType enum-like constants:
const (
	SecretTypePlaintext = "plaintext"
	SecretTypeJSON      = "json"
)

// RetrieveSecret fetches a secret from AWS Secrets Manager and returns a
// single key-value pair in a map[string]string.
//
//	secretName: the name or ARN of the secret you want to retrieve
//	secretType: either "plaintext" or "json"
//	keyName:    if secretType is "json", specify the JSON key to extract from the secret
//
// Return:
//
//	map[string]string: a single key-value pair
//	error:             if something goes wrong
func RetrieveSecret(ctx context.Context, secretName, secretType, keyName string) (map[string]string, error) {
	// Load the default AWS config from environment, shared config, etc.
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sm.NewFromConfig(cfg)

	// Fetch the secret from Secrets Manager
	out, err := client.GetSecretValue(ctx, &sm.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		var notFoundErr *smTypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return nil, fmt.Errorf("secret not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve secret: %w", err)
	}

	// SecretString could be plaintext or JSON
	secretString := aws.ToString(out.SecretString)
	if secretString == "" {
		return nil, fmt.Errorf("secret string is empty or not set")
	}

	switch secretType {
	case SecretTypePlaintext:
		// Return the entire secret under a key of your choice.
		// For instance, use "plaintext" if you don’t have a meaningful keyName.
		kvMap := map[string]string{
			keyName: secretString,
		}
		return kvMap, nil

	case SecretTypeJSON:
		// Parse JSON and extract the specific key.
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(secretString), &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON secret: %w", err)
		}

		// If keyName is empty, return all key-value pairs
		if keyName == "" {
			result := make(map[string]string)
			for k, v := range data {
				strVal, ok := v.(string)
				if ok {
					result[k] = strVal
				} else {
					continue
				}
			}
			return result, nil
		}

		val, ok := data[keyName]
		if !ok {
			return nil, fmt.Errorf("key %q not found in JSON secret", keyName)
		}

		// Ensure the value is a string
		strVal, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("key %q is not a string in JSON secret", keyName)
		}

		// Return the single key-value pair
		kvMap := map[string]string{
			keyName: strVal,
		}
		return kvMap, nil

	default:
		return nil, fmt.Errorf("unknown secretType: %q (use %q or %q)",
			secretType, SecretTypePlaintext, SecretTypeJSON)
	}
}
