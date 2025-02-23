package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/forbearing/golib/config"
	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
)

func Init() (err error) {
	return nil
}

func New(cfg config.LdapConfig) (*ldap.Conn, error) {
	return newConn(cfg)
}

// SearchUser search user in ldap.
// If the length of attribute is 0, attribute default to []string{"cn"}.
func SearchUser(username string, attribute []string) (*ldap.Conn, []*ldap.Entry, error) {
	conn, err := newConn(config.App.LdapConfig)
	if err != nil {
		return nil, nil, err
	}
	if len(attribute) == 0 {
		attribute = []string{"cn"}
	}
	cfg := config.App.LdapConfig
	if err = conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
		conn.Close()
		return nil, nil, err
	}
	searchReq := ldap.NewSearchRequest(cfg.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, fmt.Sprintf(cfg.SearchFilter, username), attribute, nil)
	sr, err := conn.Search(searchReq)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, sr.Entries, nil
}

// AuthUser checks that username and password matched.
func AuthUser(username, password string) bool {
	if username == config.App.AuthConfig.NoneExpireUsername && password == config.App.AuthConfig.NoneExpirePassword {
		return true
	}

	conn, entry, err := SearchUser(username, nil)
	if err != nil {
		zap.S().Error(err)
		return false
	}
	defer conn.Close()
	if len(entry) != 1 {
		zap.S().Error("ldap entry length not equal 1")
		return false
	}
	userDN := entry[0].DN
	if err = conn.Bind(userDN, password); err != nil {
		zap.S().Error(err)
		return false
	}
	return true
}

// newConn creates a ldap connection, and it's your responsebility to close the connection.
func newConn(cfg config.LdapConfig) (conn *ldap.Conn, err error) {
	schema := "ldap"
	if cfg.UseSsl {
		schema = "ldaps"
	}
	return ldap.DialURL(fmt.Sprintf("%s://%s:%d", schema, cfg.Host, cfg.Port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
}
