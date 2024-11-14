// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/mozey/config/pkg/share"
	"github.com/pkg/errors"
)

// KeyPrefix is not made publicly available on this package,
// users must use the getter or setter methods.
// This package must not change the config file

// APP_DOMAIN
var domain string

// APP_DOMAIN_HOSTS
var domainHosts string

// APP_INSTANCE_ID
var instanceId string

// APP_LISTEN
var listen string

// APP_PORT_API
var portApi string

// APP_PORT_CADDY
var portCaddy string

// APP_TEMPLATE_DOMAIN_DIR
var templateDomainDir string

// AWS_PROFILE
var awsProfile string

// APP_DIR
var dir string

// Config fields correspond to config file keys less the prefix
type Config struct {
	domain            string // APP_DOMAIN
	domainHosts       string // APP_DOMAIN_HOSTS
	instanceId        string // APP_INSTANCE_ID
	listen            string // APP_LISTEN
	portApi           string // APP_PORT_API
	portCaddy         string // APP_PORT_CADDY
	templateDomainDir string // APP_TEMPLATE_DOMAIN_DIR
	awsProfile        string // AWS_PROFILE
	dir               string // APP_DIR
}

// Domain is APP_DOMAIN
func (c *Config) Domain() string {
	return c.domain
}

// DomainHosts is APP_DOMAIN_HOSTS
func (c *Config) DomainHosts() string {
	return c.domainHosts
}

// InstanceId is APP_INSTANCE_ID
func (c *Config) InstanceId() string {
	return c.instanceId
}

// Listen is APP_LISTEN
func (c *Config) Listen() string {
	return c.listen
}

// PortApi is APP_PORT_API
func (c *Config) PortApi() string {
	return c.portApi
}

// PortCaddy is APP_PORT_CADDY
func (c *Config) PortCaddy() string {
	return c.portCaddy
}

// TemplateDomainDir is APP_TEMPLATE_DOMAIN_DIR
func (c *Config) TemplateDomainDir() string {
	return c.templateDomainDir
}

// AwsProfile is AWS_PROFILE
func (c *Config) AwsProfile() string {
	return c.awsProfile
}

// Dir is APP_DIR
func (c *Config) Dir() string {
	return c.dir
}

// SetDomain overrides the value of domain
func (c *Config) SetDomain(v string) {
	c.domain = v
}

// SetDomainHosts overrides the value of domainHosts
func (c *Config) SetDomainHosts(v string) {
	c.domainHosts = v
}

// SetInstanceId overrides the value of instanceId
func (c *Config) SetInstanceId(v string) {
	c.instanceId = v
}

// SetListen overrides the value of listen
func (c *Config) SetListen(v string) {
	c.listen = v
}

// SetPortApi overrides the value of portApi
func (c *Config) SetPortApi(v string) {
	c.portApi = v
}

// SetPortCaddy overrides the value of portCaddy
func (c *Config) SetPortCaddy(v string) {
	c.portCaddy = v
}

// SetTemplateDomainDir overrides the value of templateDomainDir
func (c *Config) SetTemplateDomainDir(v string) {
	c.templateDomainDir = v
}

// SetAwsProfile overrides the value of awsProfile
func (c *Config) SetAwsProfile(v string) {
	c.awsProfile = v
}

// SetDir overrides the value of dir
func (c *Config) SetDir(v string) {
	c.dir = v
}

// New creates an instance of Config.
// Build with ldflags to set the package vars.
// Env overrides package vars.
// Fields correspond to the config file keys less the prefix.
// The config file must have a flat structure
func New() *Config {
	conf := &Config{}
	SetVars(conf)
	SetEnv(conf)
	return conf
}

// SetVars sets non-empty package vars on Config
func SetVars(conf *Config) {

	if domain != "" {
		conf.domain = domain
	}

	if domainHosts != "" {
		conf.domainHosts = domainHosts
	}

	if instanceId != "" {
		conf.instanceId = instanceId
	}

	if listen != "" {
		conf.listen = listen
	}

	if portApi != "" {
		conf.portApi = portApi
	}

	if portCaddy != "" {
		conf.portCaddy = portCaddy
	}

	if templateDomainDir != "" {
		conf.templateDomainDir = templateDomainDir
	}

	if awsProfile != "" {
		conf.awsProfile = awsProfile
	}

	if dir != "" {
		conf.dir = dir
	}

}

// SetEnv sets non-empty env vars on Config
func SetEnv(conf *Config) {
	var v string

	v = os.Getenv("APP_DOMAIN")
	if v != "" {
		conf.domain = v
	}

	v = os.Getenv("APP_DOMAIN_HOSTS")
	if v != "" {
		conf.domainHosts = v
	}

	v = os.Getenv("APP_INSTANCE_ID")
	if v != "" {
		conf.instanceId = v
	}

	v = os.Getenv("APP_LISTEN")
	if v != "" {
		conf.listen = v
	}

	v = os.Getenv("APP_PORT_API")
	if v != "" {
		conf.portApi = v
	}

	v = os.Getenv("APP_PORT_CADDY")
	if v != "" {
		conf.portCaddy = v
	}

	v = os.Getenv("APP_TEMPLATE_DOMAIN_DIR")
	if v != "" {
		conf.templateDomainDir = v
	}

	v = os.Getenv("AWS_PROFILE")
	if v != "" {
		conf.awsProfile = v
	}

	v = os.Getenv("APP_DIR")
	if v != "" {
		conf.dir = v
	}

}

// GetMap of all env vars
func (c *Config) GetMap() map[string]string {
	m := make(map[string]string)

	m["APP_DOMAIN"] = c.domain

	m["APP_DOMAIN_HOSTS"] = c.domainHosts

	m["APP_INSTANCE_ID"] = c.instanceId

	m["APP_LISTEN"] = c.listen

	m["APP_PORT_API"] = c.portApi

	m["APP_PORT_CADDY"] = c.portCaddy

	m["APP_TEMPLATE_DOMAIN_DIR"] = c.templateDomainDir

	m["AWS_PROFILE"] = c.awsProfile

	m["APP_DIR"] = c.dir

	return m
}

// LoadMap sets the env from a map and returns a new instance of Config
func LoadMap(configMap map[string]string) (conf *Config) {
	for key, val := range configMap {
		_ = os.Setenv(key, val)
	}
	return New()
}

// SetEnvBase64 decodes and sets env from the given base64 string
func SetEnvBase64(configBase64 string) (err error) {
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(configBase64)
	if err != nil {
		return errors.WithStack(err)
	}
	// UnMarshall json
	configMap := make(map[string]string)
	err = json.Unmarshal(decoded, &configMap)
	if err != nil {
		return errors.WithStack(err)
	}
	// Set config
	for key, value := range configMap {
		err = os.Setenv(key, value)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// LoadFile sets the env from file and returns a new instance of Config
func LoadFile(env string) (conf *Config, err error) {
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		// Use current working dir
		appDir, err = os.Getwd()
		if err != nil {
			return conf, errors.WithStack(err)
		}
	}

	var configPath string
	filePaths, err := share.GetConfigFilePaths(appDir, env)
	if err != nil {
		return conf, err
	}
	for _, configPath = range filePaths {
		_, err := os.Stat(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Path does not exist
				continue
			}
			return conf, errors.WithStack(err)
		}
		// Path exists
		break
	}
	if configPath == "" {
		return conf, errors.Errorf("config file not found in %s", appDir)
	}

	b, err := os.ReadFile(configPath)
	if err != nil {
		return conf, errors.WithStack(err)
	}

	configMap, err := share.UnmarshalConfig(configPath, b)
	if err != nil {
		return conf, err
	}
	for key, val := range configMap {
		_ = os.Setenv(key, val)
	}
	return New(), nil
}
