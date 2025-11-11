package urlbuilder

import (
	"fmt"
	"net/url"
)

type URLBuilder struct {
	scheme   string
	user     string
	password string
	host     string
	port     string
	path     string
	query    url.Values
}

func New() *URLBuilder {
	return &URLBuilder{query: url.Values{}}
}

func (b *URLBuilder) Scheme(s string) *URLBuilder {
	b.scheme = s

	return b
}

func (b *URLBuilder) User(user, pass string) *URLBuilder {
	b.user, b.password = user, pass

	return b
}

func (b *URLBuilder) Host(host string) *URLBuilder {
	b.host = host

	return b
}

func (b *URLBuilder) Port(port string) *URLBuilder {
	b.port = port

	return b
}

func (b *URLBuilder) Path(p string) *URLBuilder {
	b.path = p

	return b
}

func (b *URLBuilder) AddQuery(key, value string) *URLBuilder {
	b.query.Add(key, value)

	return b
}

func (b *URLBuilder) Build() string {
	u := &url.URL{
		Scheme: b.scheme,
		Host:   fmt.Sprintf("%s:%s", b.host, b.port),
		Path:   b.path,
	}
	if b.user != "" {
		u.User = url.UserPassword(b.user, b.password)
	}

	if len(b.query) > 0 {
		u.RawQuery = b.query.Encode()
	}

	return u.String()
}
