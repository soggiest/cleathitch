package connector

import (
	"fmt"

	"github.com/soggiest/cleathitch/config"
	ldap "gopkg.in/ldap.v2"
)

func GetGroups(cfg config.Config, username string) []string {

	var groups []string
	var userID string
	var groupsAttr string

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

	if len(cfg.UserSearch.IDAttr) < 1 {
		userID = "sAMAccount"
	} else {
		userID = cfg.UserSearch.IDAttr
	}

	if len(cfg.UserSearch.GroupsAttr) < 1 {
		groupsAttr = "memberOf"
	} else {
		groupsAttr = cfg.UserSearch.GroupsAttr
	}

	searchRequest := ldap.NewSearchRequest(
		cfg.UserSearch.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(%s)(%s=%s))", cfg.UserSearch.Filter, userID, username),
		[]string{groupsAttr},
		nil,
	)

	userResult, err := ldapC.Search(searchRequest)
	if err != nil {
		fmt.Printf("Error searching LDAP: %v\n", err.Error())
	}

	userAttributes := userResult.Entries[0].GetAttributeValues(groupsAttr)

	for _, group := range userAttributes {
		groupSearchRequest := ldap.NewSearchRequest(
			group,
			ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
			cfg.GroupSearch.Filter,
			[]string{cfg.GroupSearch.NameAttr},
			nil,
		)
		groupResult, err := ldapC.Search(groupSearchRequest)
		if err != nil {
			fmt.Printf("Error searching LDAP for group: %v\n", err.Error())
		}

		groupname := groupResult.Entries[0].GetAttributeValue(cfg.GroupSearch.NameAttr)
		groups = append(groups, groupname)
	}
	/*
		for _, group := range groupResult.Entries {
			for _, groupAttribute := range group.Attributes {
				for _, groupValues := range groupAttribute.Values {
					if strings.Contains(groupValues, "cn=") {

					}
				}
				groups = append(groups, groupAttribute.Values)
			}
		}
	*/
	return groups
}
