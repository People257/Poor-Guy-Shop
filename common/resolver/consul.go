package resolver

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jpillora/backoff"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"time"
)

// resolvr implements resolver.Resolver from the gRPC package.
// It watches for endpoints changes and pushes them to the underlying gRPC connection.
type resolvr struct {
	cancelFunc context.CancelFunc
	cc         resolver.ClientConn
	target     *target
}

// ResolveNow will be skipped due unnecessary in this case
func (r *resolvr) ResolveNow(resolver.ResolveNowOptions) {}

// Close closes the resolver.
func (r *resolvr) Close() {
	r.cancelFunc()
}

func (r *resolvr) watchConsulService(ctx context.Context, s *api.Health) {
	tgt := r.target

	bck := &backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    10 * time.Millisecond,
		Max:    tgt.MaxBackoff,
	}
	var lastIndex uint64
	for {
		ss, meta, err := s.ServiceMultipleTags(tgt.Service, tgt.TagsSlice(), tgt.Healthy, &api.QueryOptions{
			WaitIndex:         lastIndex,
			Near:              tgt.Near,
			WaitTime:          tgt.Wait,
			Datacenter:        tgt.Dc,
			AllowStale:        tgt.AllowStale,
			RequireConsistent: tgt.RequireConsistent,
		})
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				zap.L().Error("[Consul resolver] Couldn't fetch endpoints", zap.String("target", tgt.String()), zap.Error(err))
				time.Sleep(bck.Duration())
				continue
			}
		}
		bck.Reset()
		lastIndex = meta.LastIndex
		zap.L().Debug("[Consul resolver] Endpoints fetched",
			zap.Int("count", len(ss)),
			zap.String("target", tgt.String()),
			zap.Duration("wait", meta.RequestTime))

		cc := make([]string, 0, len(ss))
		for _, s := range ss {
			address := s.Service.Address
			if s.Service.Address == "" {
				address = s.Node.Address
			}
			cc = append(cc, fmt.Sprintf("%s:%d", address, s.Service.Port))
		}

		if tgt.Limit != 0 && len(cc) > tgt.Limit {
			cc = cc[:tgt.Limit]
		}

		connsSet := make(map[string]struct{}, len(cc))
		for _, c := range cc {
			connsSet[c] = struct{}{}
		}
		conns := make([]resolver.Address, 0, len(connsSet))
		for c := range connsSet {
			conns = append(conns, resolver.Address{Addr: c})
		}
		err = r.cc.UpdateState(resolver.State{Addresses: conns})
		if err != nil {
			zap.L().Error("[Consul resolver] Couldn't update client connection", zap.String("service", tgt.Service), zap.Error(err))
			continue
		}
	}
}
