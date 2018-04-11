package config

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	// The host and optional port of the LDAP server. If port isn't supplied, it will be
	// guessed based on the TLS configuration. 389 or 636.
	Host string `yaml:"host" envconfig:"host"`

	Port string `yaml:"port" envconfig:"port"`

	// Required if LDAP host does not use TLS.
	InsecureNoSSL bool `yaml:"insecure_no_ssl" envconfig:"insecure_no_ssl"`

	// Don't verify the CA.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify" envconfig:"insecure_skip_verify"`

	// Connect to the insecure port then issue a StartTLS command to negotiate a
	// secure connection. If unsupplied secure connections will use the LDAPS
	// protocol.
	StartTLS bool `yaml:"start_tls" envconfig:"start_tls"`

	// Path to a trusted root certificate file.
	RootCA string `yaml:"root_ca" envconfig:"root_ca"`

	// Base64 encoded PEM data containing root CAs.
	RootCAData []byte `yaml:"root_ca_data" envconfig:"root_ca_data"`

	// BindDN and BindPW for an application service account. The connector uses these
	// credentials to search for users and groups.
	BindDN string `yaml:"bind_dn" envconfig:"bind_dn"`
	BindPW string `yaml:"bind_pw" envconfig:"bind_pw"`

	// UsernamePrompt allows users to override the username attribute (displayed
	// in the username/password prompt). If unset, the handler will use
	// "Username".
	UsernamePrompt string `yaml:"usernamePrompt"`

	// User entry search configuration.
	UserSearch struct {
		// BaseDN to start the search from. For example "cn=users,dc=example,dc=com"
		BaseDN string `yaml:"baseDN"`

		// Optional filter to apply when searching the directory. For example "(objectClass=person)"
		Filter string `yaml:"filter"`

		// Attribute to match against the inputted username. This will be translated and combined
		// with the other filter as "(<attr>=<username>)".
		Username string `yaml:"username"`

		// Can either be:
		// * "sub" - search the whole sub tree
		// * "one" - only search one level
		Scope string `yaml:"scope"`

		// A mapping of attributes on the user entry to claims.
		IDAttr    string `yaml:"idAttr"`    // Defaults to "uid"
		EmailAttr string `yaml:"emailAttr"` // Defaults to "mail"
		NameAttr  string `yaml:"nameAttr"`  // No default.

	} `yaml:"userSearch"`

	// Group search configuration.
	GroupSearch struct {
		// BaseDN to start the search from. For example "cn=groups,dc=example,dc=com"
		BaseDN string `yaml:"baseDN"`

		// Optional filter to apply when searching the directory. For example "(objectClass=posixGroup)"
		Filter string `yaml:"filter"`

		Scope string `yaml:"scope"` // Defaults to "sub"

		// These two fields are use to match a user to a group.
		//
		// It adds an additional requirement to the filter that an attribute in the group
		// match the user's attribute value. For example that the "members" attribute of
		// a group matches the "uid" of the user. The exact filter being added is:
		//
		//   (<groupAttr>=<userAttr value>)
		//
		UserAttr  string `yaml:"userAttr"`
		GroupAttr string `yaml:"groupAttr"`

		// The attribute of the group that represents its name.
		NameAttr string `yaml:"nameAttr"`
	} `yaml:"groupSearch"`
}

func ReadConfig(configFile string) Config {
	var config Config
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", configFile)
	}

	file, _ := os.Open(configFile)
	fileBytes, err := ioutil.ReadAll(file)
	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalln("Config unmarshalling error:", err)

	}
	return config
}
