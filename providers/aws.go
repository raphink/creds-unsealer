package providers

import (
	"fmt"

	"github.com/raphink/narcissus"
	"honnef.co/go/augeas"
)

// AWS stores data used to manage AWS credentials
type AWS struct {
	*BaseProvider
}

// AWSConfig represents an entry in the AWS's configuration file
type AWSConfig struct {
	AccessKeyID     string `yaml:"aws_access_key_id,omitempty" narcissus:"aws_access_key_id"`
	SecretAccessKey string `yaml:"aws_secret_access_key,omitempty" narcissus:"aws_secret_access_key"`
	Region          string `yaml:"region,omitempty" narcissus:"region"`
}

// AWSConfigs is used by Narcissus to manage entries on the AWS's configuration file
type AWSConfigs struct {
	augeasFile string
	augeasLens string `default:"IniFile.lns_loose"`
	augeasPath string
	Configs    map[string]AWSConfig `narcissus:"section"`
}

// GetName returns the provider's name
func (a *AWS) GetName() string {
	return "AWS"
}

// Unseal unseals a secret from the backend and add it to the config file
func (a *AWS) Unseal(cred string) (err error) {
	var secret AWSConfig
	err = a.backend.GetSecret(a.inputPath+"/"+cred, &secret)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}

	err = a.writeSecret(cred, secret)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %s", err)
	}

	return
}

func (a *AWS) writeSecret(name string, config AWSConfig) (err error) {
	aug, err := augeas.New("/", "", augeas.NoModlAutoload)
	defer aug.Close()
	if err != nil {
		return fmt.Errorf("failed to initialize Augeas: %s", err)
	}

	n := narcissus.New(&aug)
	configs := AWSConfigs{
		augeasFile: a.outputPath,
		augeasPath: "/files" + a.outputPath,
	}
	configs.Configs = make(map[string]AWSConfig)
	configs.Configs[a.outputKeyPrefix+name] = config

	err = n.Write(&configs)
	if err != nil {
		return fmt.Errorf("failed to write configs: %s", err)
	}

	return
}
