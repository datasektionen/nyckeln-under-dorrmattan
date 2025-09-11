package sso

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type storage struct {
	lock         sync.Mutex
	authRequests map[string]*authRequest
	codes        map[string]string
	tokens       map[string]*accessToken
	signingKey   signingKey

	dao *dao.Dao
}

type accessToken struct {
	kthid  string
	scopes []string
}

// AuthRequestByCode implements op.Storage.
func (s *storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	request, ok := s.authRequests[id]
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	return request, nil
}

// AuthRequestByCode implements the op.Storage interface
func (s *storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	requestID, ok := func() (string, bool) {
		s.lock.Lock()
		defer s.lock.Unlock()
		requestID, ok := s.codes[code]
		return requestID, ok
	}()
	if !ok {
		return nil, fmt.Errorf("code invalid or expired")
	}
	return s.AuthRequestByID(ctx, requestID)
}

// AuthorizeClientIDSecret implements the op.Storage interface
func (s *storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	client, err := s.dao.GetClient(clientID)
	if err != nil {
		return err
	}
	if client.Secret != clientSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

// CreateAccessAndRefreshTokens implements op.Storage.
func (s *storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAccessToken implements op.Storage.
func (s *storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	tokenID := uuid.New()
	s.accessToken(tokenID.String(), request.GetSubject(), request.GetScopes())

	return tokenID.String(), time.Now().Add(time.Minute), nil
}

// CreateAuthRequest implements the op.Storage interface
func (s *storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id := uuid.NewString()
	request := authRequest{id: id, authCode: "", inner: authReq}
	s.authRequests[id] = &request

	return request, nil
}

// DeleteAuthRequest implements op.Storage.
func (s *storage) DeleteAuthRequest(ctx context.Context, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.authRequests, id)
	for code, requestID := range s.codes {
		if id == requestID {
			delete(s.codes, code)
			return nil
		}
	}
	return nil
}

// GetClientByClientID implements the op.Storage interface
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	c, err := s.dao.GetClient(clientID)
	if err != nil {
		return nil, err
	}

	return client{
		id:           c.Id,
		redirectURIs: c.RedirectURIs,
	}, nil
}

func (s *storage) accessToken(id string, subject string, scopes []string) *accessToken {
	s.lock.Lock()
	defer s.lock.Unlock()
	token := &accessToken{
		kthid:  subject,
		scopes: scopes,
	}

	s.tokens[id] = token
	return token
}

// CheckLogin implements authentication for the login page.
func (s *storage) CheckLogin(kthid, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	request, ok := s.authRequests[id]
	if !ok {
		return fmt.Errorf("request not found")
	}

	user, err := s.dao.GetUser(kthid)
	if err != nil {
		return err
	}
	if user.KTHID != kthid {
		return fmt.Errorf("invalid kthid")
	}

	request.kthid = kthid
	return nil
}

// GetKeyByIDAndClientID implements op.Storage.
func (s *storage) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented")
}

// GetPrivateClaimsFromScopes implements op.Storage.
func (s *storage) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented")
}

// GetRefreshTokenInfo implements op.Storage.
func (s *storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented")
}

// Health implements op.Storage.
func (s *storage) Health(ctx context.Context) error {
	panic("unimplemented")
}

// RevokeToken implements op.Storage.
func (s *storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented")
}

// SaveAuthCode implements op.Storage.
func (s *storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.codes[code] = id
	return nil
}

// SetIntrospectionFromToken implements op.Storage.
func (s *storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented")
}

// SetUserinfoFromScopes implements op.Storage.
func (s *storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, kthid string, clientID string, scopes []string) error {
	user, err := s.dao.GetUser(kthid)
	if err != nil {
		return err
	}

	if err := s.setUserinfo(ctx, userinfo, user, scopes); err != nil {
		return err
	}
	return nil
}

// SetUserinfoFromToken implements op.Storage.
func (s *storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, kthid, origin string) error {
	token, ok := func() (*accessToken, bool) {
		s.lock.Lock()
		defer s.lock.Unlock()
		token, ok := s.tokens[tokenID]
		return token, ok
	}()
	if !ok {
		return fmt.Errorf("token is invalid or has expired")
	}

	if token.kthid != kthid {
		return fmt.Errorf("You're asking to get info about a different user than who the token is for")
	}

	user, err := s.dao.GetUser(kthid)
	if err != nil {
		return err
	}

	if err := s.setUserinfo(ctx, userinfo, user, token.scopes); err != nil {
		return err
	}

	return nil
}

func (s *storage) setUserinfo(ctx context.Context, userinfo *oidc.UserInfo, user *dao.User, scopes []string) error {
	if userinfo.Claims == nil {
		userinfo.Claims = make(map[string]any)
	}

	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userinfo.Subject = user.KTHID

		case oidc.ScopeProfile:
			userinfo.Name = user.FirstName + " " + user.FamilyName
			userinfo.GivenName = user.FirstName
			userinfo.FamilyName = user.FamilyName

		case oidc.ScopeEmail:
			userinfo.Email = user.Email
			userinfo.EmailVerified = true

		case "year":
			userinfo.Claims["year_tag"] = user.YearTag

		default:
			if group, ok := strings.CutPrefix(scope, "pls_"); ok {
				perms := s.dao.GetUserPermissionsForGroup(user.KTHID, group)
				userinfo.Claims[scope] = perms
			}
		}
	}

	return nil
}

// SignatureAlgorithms implements op.Storage.
func (s *storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{jose.RS256}, nil
}

type signingKey struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *signingKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKey) Key() any {
	return s.key
}

func (s *signingKey) ID() string {
	return s.id
}

type publicKey struct {
	signingKey
}

func (s *publicKey) ID() string {
	return s.id
}

func (s *publicKey) Algorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *publicKey) Use() string {
	return "sig"
}

func (s *publicKey) Key() any {
	return &s.key.PublicKey
}

// KeySet implements the op.Storage interface
func (s *storage) KeySet(ctx context.Context) ([]op.Key, error) {
	return []op.Key{&publicKey{s.signingKey}}, nil
}

// SigningKey implements op.Storage.
func (s *storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	return &s.signingKey, nil
}

// TerminateSession implements op.Storage.
func (s *storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented")
}

// TokenRequestByRefreshToken implements op.Storage.
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	panic("unimplemented")
}

// ValidateJWTProfileScopes implements op.Storage.
func (s *storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented")
}
