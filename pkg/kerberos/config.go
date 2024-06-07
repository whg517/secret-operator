package kerberos

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"
	"strings"
)

/*
Example krb5.conf file:

# [logging]
#  default = FILE:/tmp/krb5libs.log
#  kdc = FILE:/tmp/krb5kdc.log
#  admin_server = FILE:/tmp/kadmind.log

[libdefaults]
 default_realm = EXAMPLE.COM
#  dns_lookup_realm = false
#  dns_lookup_kdc = true
#  rdns = false
#  ticket_lifetime = 24h
#  forwardable = true
#  udp_preference_limit = 0

[realms]
 EXAMPLE.COM = {
  kdc = fodera.example.cn:88
  # master_kdc = fodera.example.cn:88
  # kpasswd_server = fodera.example.cn:464
  admin_server = fodera.example.cn:749
  # default_domain = example.cn
}

# [domain_realm]
#  .example.cn = EXAMPLE.COM
#  example.cn = EXAMPLE.COM

*/

type Krb5Config struct {
	Realm       string
	AdminServer string
	KDC         string

	hashed string
}

func (c *Krb5Config) GetRealm() string {
	return strings.ToUpper(c.Realm)
}

// ref: https://web.mit.edu/kerberos/krb5-latest/doc/admin/conf_files/krb5_conf.html#sample-krb5-conf-file
func (c *Krb5Config) Content() string {
	content := `# krb5.conf generated by secrets.zncdata.dev
# It will be overwritten by the secret-operator

[libdefaults]
  default_realm = ` + c.GetRealm() + `
  dns_lookup_realm = false
  dns_lookup_kdc = true
  rdns = false
  udp_preference_limit = 1

[realms]
  ` + c.GetRealm() + ` = {
    kdc = ` + c.KDC + `
    admin_server = ` + c.AdminServer + `
  }

[domain_realm]
  cluster.local = ` + c.GetRealm() + `
  .cluster.local = ` + c.GetRealm() + `

`
	return content
}

// Save write krb5.conf file

// Default krb5.conf in Linux is /etc/krb5.conf, if you want to use custom krb5.conf file, you can set KRB5_CONFIG env.
func (c *Krb5Config) Save(path string) error {
	return os.WriteFile(path, []byte(c.Content()), 0644)
}

func (c *Krb5Config) CheckSum() string {
	if c.hashed == "" {
		hash := sha256.Sum256([]byte(c.Content()))
		c.hashed = hex.EncodeToString(hash[:])[:23]
	}
	return c.hashed
}

func (c *Krb5Config) GetTempPath() (string, error) {
	absFilename := path.Join(os.TempDir(), "krb5-"+c.CheckSum()+".conf")

	if _, err := os.Stat(absFilename); os.IsNotExist(err) {
		if err := c.Save(absFilename); err != nil {
			return "", err
		}
	}

	return absFilename, nil
}
