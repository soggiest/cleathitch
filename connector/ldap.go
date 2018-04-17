package connector

import (
	"fmt"

	ldap "gopkg.in/ldap.v2"

	"github.com/soggiest/cleathitch/config"
)

func GetGroups(cfg config.Config, username string) []string {

	var groups []string

	ldapC, err := ldap.Dial(cfg.Protocol, fmt.Sprintf("%s:%s", cfg.LDAPHost, cfg.LDAPPort))
	if err != nil {
		fmt.Printf("Error contacting LDAP server: %v\n", err.Error())
	}
	err = ldapC.Bind(cfg.BindDN, cfg.BindPW)
	if err != nil {
		fmt.Printf("Error Binding: %v\n", err.Error())
	}

	//	filter := fmt.Sprintf("(&(cn=*) (memberUid=%s))", username)

	//fmt.Printf("FILTER: %v", filter)
	searchRequest := ldap.NewSearchRequest(
		cfg.GroupSearch.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(cn=*)(memberUid=%s))", username),
		[]string{"cn"},
		nil,
	)
	groupResult, err := ldapC.Search(searchRequest)
	if err != nil {
		fmt.Printf("Error searching LDAP: %v\n", err.Error())
	}

	for _, group := range groupResult.Entries {
		groups = append(groups, group.Attributes[0].Values[0])
	}

	return groups
}
