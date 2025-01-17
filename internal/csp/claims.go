// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	RoleOrgOwner         = "csp:org_owner"
	RolePlatformOperator = "csp:platform_operator"
)

type Claims struct {
	jwt.StandardClaims

	ContextName string   `json:"context_name,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Username    string   `json:"username,omitempty"`
	Perms       []string `json:"perms,omitempty"`

	Context         string `json:"context,omitempty"`
	AuthorizedParty string `json:"azp,omitempty"`

	// The token as a string, signed and ready to be put in an Authorization header
	Token string `json:"-"`
}

func (claims *Claims) GetQualifiedUsername() string {
	if !strings.Contains(claims.Username, "@") {
		return fmt.Sprintf("%s@%s", claims.Username, claims.Domain)
	}
	return claims.Username
}

func (claims *Claims) IsOrgOwner() bool {
	for _, p := range claims.Perms {
		if p == RoleOrgOwner {
			return true
		}
	}

	return false
}

func (claims *Claims) IsPlatformOperator() bool {
	for _, p := range claims.Perms {
		if p == RolePlatformOperator {
			return true
		}
	}

	return false
}
