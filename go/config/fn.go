// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Fn struct {
	input string
	// output of the last function,
	// might be useful when chaining multiple functions?
	output string
}

// .............................................................................
// Methods to set function input

// FnDomain sets the function input to the value of APP_DOMAIN
func (c *Config) FnDomain() *Fn {
	fn := Fn{}
	fn.input = c.domain
	fn.output = ""
	return &fn
}

// FnDomainHosts sets the function input to the value of APP_DOMAIN_HOSTS
func (c *Config) FnDomainHosts() *Fn {
	fn := Fn{}
	fn.input = c.domainHosts
	fn.output = ""
	return &fn
}

// FnInstanceId sets the function input to the value of APP_INSTANCE_ID
func (c *Config) FnInstanceId() *Fn {
	fn := Fn{}
	fn.input = c.instanceId
	fn.output = ""
	return &fn
}

// FnListen sets the function input to the value of APP_LISTEN
func (c *Config) FnListen() *Fn {
	fn := Fn{}
	fn.input = c.listen
	fn.output = ""
	return &fn
}

// FnPortApi sets the function input to the value of APP_PORT_API
func (c *Config) FnPortApi() *Fn {
	fn := Fn{}
	fn.input = c.portApi
	fn.output = ""
	return &fn
}

// FnPortCaddy sets the function input to the value of APP_PORT_CADDY
func (c *Config) FnPortCaddy() *Fn {
	fn := Fn{}
	fn.input = c.portCaddy
	fn.output = ""
	return &fn
}

// FnTemplateDomainDir sets the function input to the value of APP_TEMPLATE_DOMAIN_DIR
func (c *Config) FnTemplateDomainDir() *Fn {
	fn := Fn{}
	fn.input = c.templateDomainDir
	fn.output = ""
	return &fn
}

// FnAwsProfile sets the function input to the value of AWS_PROFILE
func (c *Config) FnAwsProfile() *Fn {
	fn := Fn{}
	fn.input = c.awsProfile
	fn.output = ""
	return &fn
}

// FnDir sets the function input to the value of APP_DIR
func (c *Config) FnDir() *Fn {
	fn := Fn{}
	fn.input = c.dir
	fn.output = ""
	return &fn
}

// .............................................................................
// Type conversion functions

// Bool parses a bool from the value or returns an error.
// Valid values are "1", "0", "true", or "false".
// The value is not case-sensitive
func (fn *Fn) Bool() (bool, error) {
	v := strings.ToLower(fn.input)
	if v == "1" || v == "true" {
		return true, nil
	}
	if v == "0" || v == "false" {
		return false, nil
	}
	return false, fmt.Errorf("invalid value %s", fn.input)
}

// Float64 parses a float64 from the value or returns an error
func (fn *Fn) Float64() (float64, error) {
	f, err := strconv.ParseFloat(fn.input, 64)
	if err != nil {
		return f, err
	}
	return f, nil
}

// Int64 parses an int64 from the value or returns an error
func (fn *Fn) Int64() (int64, error) {
	i, err := strconv.ParseInt(fn.input, 10, 64)
	if err != nil {
		return i, err
	}
	return i, nil
}

// String returns the input as is
func (fn *Fn) String() string {
	return fn.input
}
