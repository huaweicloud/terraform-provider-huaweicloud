package huaweicloud

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func GetCredentials(c *Config) (*awsCredentials.Credentials, error) {
	// build a chain provider, lazy-evaluated by aws-sdk
	providers := []awsCredentials.Provider{
		&awsCredentials.StaticProvider{Value: awsCredentials.Value{
			AccessKeyID:     c.AccessKey,
			SecretAccessKey: c.SecretKey,
			SessionToken:    c.Token,
		}},
		&awsCredentials.EnvProvider{},
		&awsCredentials.SharedCredentialsProvider{
			Filename: "",
			Profile:  "",
		},
	}

	// Build isolated HTTP client to avoid issues with globally-shared settings
	client := cleanhttp.DefaultClient()

	// Keep the default timeout (100ms) low as we don't want to wait in non-EC2 environments
	client.Timeout = 100 * time.Millisecond

	const userTimeoutEnvVar = "AWS_METADATA_TIMEOUT"
	userTimeout := os.Getenv(userTimeoutEnvVar)
	if userTimeout != "" {
		newTimeout, err := time.ParseDuration(userTimeout)
		if err == nil {
			if newTimeout.Nanoseconds() > 0 {
				client.Timeout = newTimeout
			} else {
				log.Printf("[WARN] Non-positive value of %s (%s) is meaningless, ignoring", userTimeoutEnvVar, newTimeout.String())
			}
		} else {
			log.Printf("[WARN] Error converting %s to time.Duration: %s", userTimeoutEnvVar, err)
		}
	}

	log.Printf("[INFO] Setting AWS metadata API timeout to %s", client.Timeout.String())
	cfg := &aws.Config{
		HTTPClient: client,
	}
	usedEndpoint := setOptionalEndpoint(cfg)

	// Add the default AWS provider for ECS Task Roles if the relevant env variable is set
	if uri := os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"); len(uri) > 0 {
		providers = append(providers, defaults.RemoteCredProvider(*cfg, defaults.Handlers()))
		log.Print("[INFO] ECS container credentials detected, RemoteCredProvider added to auth chain")
	}

	// Real AWS should reply to a simple metadata request.
	// We check it actually does to ensure something else didn't just
	// happen to be listening on the same IP:Port
	metadataClient := ec2metadata.New(session.New(cfg))
	if metadataClient.Available() {
		providers = append(providers, &ec2rolecreds.EC2RoleProvider{
			Client: metadataClient,
		})
		log.Print("[INFO] AWS EC2 instance detected via default metadata" +
			" API endpoint, EC2RoleProvider added to the auth chain")
	} else {
		if usedEndpoint == "" {
			usedEndpoint = "default location"
		}
		log.Printf("[INFO] Ignoring AWS metadata API endpoint at %s "+
			"as it doesn't return any instance-id", usedEndpoint)
	}

	return awsCredentials.NewChainCredentials(providers), nil
}

func setOptionalEndpoint(cfg *aws.Config) string {
	endpoint := os.Getenv("AWS_METADATA_URL")
	if endpoint != "" {
		log.Printf("[INFO] Setting custom metadata endpoint: %q", endpoint)
		cfg.Endpoint = aws.String(endpoint)
		return endpoint
	}
	return ""
}
