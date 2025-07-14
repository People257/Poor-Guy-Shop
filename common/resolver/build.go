package resolver

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"strings"
)

func init() {
	resolver.Register(NewBuilder())
}

// schemeName for the urls
// All target URLs like 'consul://.../...' will be resolved by this resolver
const schemeName = "consul"

// Builder implements resolver.Builder and use for constructing all consul resolvers
type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", url.URL.Host, url.URL.Path + "?" + url.URL.RawQuery}, "/")
	tgt, err := parseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("wrong consul URL: %w", err)
	}
	cli, err := api.NewClient(tgt.consulConfig())
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to the Consul API: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &resolvr{cancelFunc: cancel, cc: cc, target: &tgt}
	go r.watchConsulService(ctx, cli.Health())

	return r, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *Builder) Scheme() string {
	return schemeName
}
