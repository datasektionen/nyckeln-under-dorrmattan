package sso

import (
	"log/slog"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type client struct {
	id           string
	secret       string
	redirectURIs []string
}

var _ op.Client = client{}

// AccessTokenType implements op.Client.
func (c client) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeBearer
}

// ApplicationType implements op.Client.
func (c client) ApplicationType() op.ApplicationType {
	slog.Warn("oidcprovider.client.ApplicationType")
	return op.ApplicationTypeWeb
}

// AuthMethod implements op.Client.
func (c client) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethodBasic
}

// ClockSkew implements op.Client.
func (c client) ClockSkew() time.Duration {
	return time.Second * 10
}

// DevMode implements op.Client.
func (c client) DevMode() bool {
	return true
}

// GetID implements op.Client.
func (c client) GetID() string {
	return c.id
}

// GrantTypes implements op.Client.
func (c client) GrantTypes() []oidc.GrantType {
	slog.Warn("oidcprovider.client.GrantTypes")
	return []oidc.GrantType{oidc.GrantTypeCode}
}

// IDTokenLifetime implements op.Client.
func (c client) IDTokenLifetime() time.Duration {
	slog.Warn("oidcprovider.client.IDTokenLifetime")
	return time.Hour * 24
}

// IDTokenUserinfoClaimsAssertion implements op.Client.
func (c client) IDTokenUserinfoClaimsAssertion() bool {
	return true
}

// IsScopeAllowed implements op.Client.
func (c client) IsScopeAllowed(scope string) bool {
	slog.Warn("oidcprovider.client.IsScopeAllowed", "scope", scope)
	if after, ok := strings.CutPrefix(scope, "pls_"); ok {
		return len(after) > 0
	}
	return false
}

// LoginURL implements op.Client.
func (c client) LoginURL(authRequestID string) string {
	return "/login?authRequestID=" + authRequestID
}

// PostLogoutRedirectURIs implements op.Client.
func (c client) PostLogoutRedirectURIs() []string {
	slog.Warn("oidcprovider.client.PostLogoutRedirectURIs")
	panic("unimplemented")
}

// RedirectURIs implements op.Client.
func (c client) RedirectURIs() []string {
	return c.redirectURIs
}

// ResponseTypes implements op.Client.
func (c client) ResponseTypes() []oidc.ResponseType {
	return []oidc.ResponseType{"code"}
}

// RestrictAdditionalAccessTokenScopes implements op.Client.
func (c client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	slog.Warn("oidcprovider.client.RestrictAdditionalAccessTokenScopes")
	panic("unimplemented")
}

// RestrictAdditionalIdTokenScopes implements op.Client.
func (c client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	slog.Warn("oidcprovider.client.RestrictAdditionalIdTokenScopes")
	return func(scopes []string) []string { return scopes }
}
