// SPDX-FileCopyrightText: 2025 The dpservice Authors
// SPDX-License-Identifier: Apache-2.0

package clientv2

import (
	"context"
	"net/netip"

	"github.com/ironcore-dev/dpservice/go/dpservice-go/api"
	legacy "github.com/ironcore-dev/dpservice/go/dpservice-go/client"
	dpdkproto "github.com/ironcore-dev/dpservice/go/dpservice-go/proto"
)

// CallOption allows customizing client call behavior.
type CallOption func(*callOptions)

type callOptions struct {
	ignoredCodes []uint32
}

// WithIgnoredCodes configures error codes that should be treated as non-fatal.
func WithIgnoredCodes(codes ...uint32) CallOption {
	return func(o *callOptions) {
		o.ignoredCodes = append(o.ignoredCodes, codes...)
	}
}

func buildCallOptions(opts ...CallOption) callOptions {
	var o callOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}

// toLegacyIgnored converts CallOptions to the legacy variadic []uint32 form.
func toLegacyIgnored(opts ...CallOption) [][]uint32 {
	o := buildCallOptions(opts...)
	if len(o.ignoredCodes) == 0 {
		return nil
	}
	return [][]uint32{o.ignoredCodes}
}

// Client is the root v2 client exposing domain-specific sub-clients.
type Client interface {
	LoadBalancers() LoadBalancers
	Interfaces() Interfaces
	Routes() Routes
	NATs() NATs
	Firewall() Firewall
	System() System
	Capture() Capture
}

// NewFromProto builds a v2 Client from a grpc/proto client.
func NewFromProto(rpc dpdkproto.DPDKironcoreClient) Client {
	return &rootAdapter{legacy: legacy.NewClient(rpc)}
}

// AsV2 adapts an existing legacy client to the v2 Client.
func AsV2(c legacy.Client) Client {
	return &rootAdapter{legacy: c}
}

// rootAdapter implements Client by delegating to the legacy client.
type rootAdapter struct {
	legacy legacy.Client
}

func (r *rootAdapter) LoadBalancers() LoadBalancers { return &lbClient{legacy: r.legacy} }
func (r *rootAdapter) Interfaces() Interfaces       { return &ifaceClient{legacy: r.legacy} }
func (r *rootAdapter) Routes() Routes               { return &routeClient{legacy: r.legacy} }
func (r *rootAdapter) NATs() NATs                   { return &natClient{legacy: r.legacy} }
func (r *rootAdapter) Firewall() Firewall           { return &fwClient{legacy: r.legacy} }
func (r *rootAdapter) System() System               { return &systemClient{legacy: r.legacy} }
func (r *rootAdapter) Capture() Capture             { return &captureClient{legacy: r.legacy} }

//
// Load Balancers
//

type LoadBalancers interface {
	Get(ctx context.Context, id string, opts ...CallOption) (*api.LoadBalancer, error)
	List(ctx context.Context, opts ...CallOption) (*api.LoadBalancerList, error)
	Create(ctx context.Context, lb *api.LoadBalancer, opts ...CallOption) (*api.LoadBalancer, error)
	Delete(ctx context.Context, id string, opts ...CallOption) (*api.LoadBalancer, error)

	Prefixes() LoadBalancerPrefixes
	Targets() LoadBalancerTargets
}

type LoadBalancerPrefixes interface {
	List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.PrefixList, error)
	Create(ctx context.Context, prefix *api.LoadBalancerPrefix, opts ...CallOption) (*api.LoadBalancerPrefix, error)
	Delete(ctx context.Context, interfaceID string, prefix *netip.Prefix, opts ...CallOption) (*api.LoadBalancerPrefix, error)
}

type LoadBalancerTargets interface {
	List(ctx context.Context, loadBalancerID string, opts ...CallOption) (*api.LoadBalancerTargetList, error)
	Create(ctx context.Context, target *api.LoadBalancerTarget, opts ...CallOption) (*api.LoadBalancerTarget, error)
	Delete(ctx context.Context, lbID string, targetIP *netip.Addr, opts ...CallOption) (*api.LoadBalancerTarget, error)
}

type lbClient struct{ legacy legacy.Client }

func (c *lbClient) Get(ctx context.Context, id string, opts ...CallOption) (*api.LoadBalancer, error) {
	return c.legacy.GetLoadBalancer(ctx, id, toLegacyIgnored(opts...)...)
}
func (c *lbClient) List(ctx context.Context, opts ...CallOption) (*api.LoadBalancerList, error) {
	return c.legacy.ListLoadBalancers(ctx, toLegacyIgnored(opts...)...)
}
func (c *lbClient) Create(ctx context.Context, lb *api.LoadBalancer, opts ...CallOption) (*api.LoadBalancer, error) {
	return c.legacy.CreateLoadBalancer(ctx, lb, toLegacyIgnored(opts...)...)
}
func (c *lbClient) Delete(ctx context.Context, id string, opts ...CallOption) (*api.LoadBalancer, error) {
	return c.legacy.DeleteLoadBalancer(ctx, id, toLegacyIgnored(opts...)...)
}
func (c *lbClient) Prefixes() LoadBalancerPrefixes { return &lbPrefixesClient{legacy: c.legacy} }
func (c *lbClient) Targets() LoadBalancerTargets   { return &lbTargetsClient{legacy: c.legacy} }

type lbPrefixesClient struct{ legacy legacy.Client }

func (c *lbPrefixesClient) List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.PrefixList, error) {
	return c.legacy.ListLoadBalancerPrefixes(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *lbPrefixesClient) Create(ctx context.Context, prefix *api.LoadBalancerPrefix, opts ...CallOption) (*api.LoadBalancerPrefix, error) {
	return c.legacy.CreateLoadBalancerPrefix(ctx, prefix, toLegacyIgnored(opts...)...)
}
func (c *lbPrefixesClient) Delete(ctx context.Context, interfaceID string, prefix *netip.Prefix, opts ...CallOption) (*api.LoadBalancerPrefix, error) {
	return c.legacy.DeleteLoadBalancerPrefix(ctx, interfaceID, prefix, toLegacyIgnored(opts...)...)
}

type lbTargetsClient struct{ legacy legacy.Client }

func (c *lbTargetsClient) List(ctx context.Context, loadBalancerID string, opts ...CallOption) (*api.LoadBalancerTargetList, error) {
	return c.legacy.ListLoadBalancerTargets(ctx, loadBalancerID, toLegacyIgnored(opts...)...)
}
func (c *lbTargetsClient) Create(ctx context.Context, target *api.LoadBalancerTarget, opts ...CallOption) (*api.LoadBalancerTarget, error) {
	return c.legacy.CreateLoadBalancerTarget(ctx, target, toLegacyIgnored(opts...)...)
}
func (c *lbTargetsClient) Delete(ctx context.Context, lbID string, targetIP *netip.Addr, opts ...CallOption) (*api.LoadBalancerTarget, error) {
	return c.legacy.DeleteLoadBalancerTarget(ctx, lbID, targetIP, toLegacyIgnored(opts...)...)
}

//
// Interfaces and sub-resources
//

type Interfaces interface {
	Get(ctx context.Context, id string, opts ...CallOption) (*api.Interface, error)
	List(ctx context.Context, opts ...CallOption) (*api.InterfaceList, error)
	Create(ctx context.Context, iface *api.Interface, opts ...CallOption) (*api.Interface, error)
	Delete(ctx context.Context, id string, opts ...CallOption) (*api.Interface, error)

	VIP() VirtualIPs
	Prefixes() InterfacePrefixes
	Firewall() Firewall
}

type VirtualIPs interface {
	Get(ctx context.Context, interfaceID string, opts ...CallOption) (*api.VirtualIP, error)
	Create(ctx context.Context, vip *api.VirtualIP, opts ...CallOption) (*api.VirtualIP, error)
	Delete(ctx context.Context, interfaceID string, opts ...CallOption) (*api.VirtualIP, error)
}

type InterfacePrefixes interface {
	List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.PrefixList, error)
	Create(ctx context.Context, prefix *api.Prefix, opts ...CallOption) (*api.Prefix, error)
	Delete(ctx context.Context, interfaceID string, prefix *netip.Prefix, opts ...CallOption) (*api.Prefix, error)
}

type ifaceClient struct{ legacy legacy.Client }

func (c *ifaceClient) Get(ctx context.Context, id string, opts ...CallOption) (*api.Interface, error) {
	return c.legacy.GetInterface(ctx, id, toLegacyIgnored(opts...)...)
}
func (c *ifaceClient) List(ctx context.Context, opts ...CallOption) (*api.InterfaceList, error) {
	return c.legacy.ListInterfaces(ctx, toLegacyIgnored(opts...)...)
}
func (c *ifaceClient) Create(ctx context.Context, iface *api.Interface, opts ...CallOption) (*api.Interface, error) {
	return c.legacy.CreateInterface(ctx, iface, toLegacyIgnored(opts...)...)
}
func (c *ifaceClient) Delete(ctx context.Context, id string, opts ...CallOption) (*api.Interface, error) {
	return c.legacy.DeleteInterface(ctx, id, toLegacyIgnored(opts...)...)
}
func (c *ifaceClient) VIP() VirtualIPs             { return &vipClient{legacy: c.legacy} }
func (c *ifaceClient) Prefixes() InterfacePrefixes { return &ifacePrefixesClient{legacy: c.legacy} }
func (c *ifaceClient) Firewall() Firewall          { return &fwClient{legacy: c.legacy} }

type vipClient struct{ legacy legacy.Client }

func (c *vipClient) Get(ctx context.Context, interfaceID string, opts ...CallOption) (*api.VirtualIP, error) {
	return c.legacy.GetVirtualIP(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *vipClient) Create(ctx context.Context, vip *api.VirtualIP, opts ...CallOption) (*api.VirtualIP, error) {
	return c.legacy.CreateVirtualIP(ctx, vip, toLegacyIgnored(opts...)...)
}
func (c *vipClient) Delete(ctx context.Context, interfaceID string, opts ...CallOption) (*api.VirtualIP, error) {
	return c.legacy.DeleteVirtualIP(ctx, interfaceID, toLegacyIgnored(opts...)...)
}

type ifacePrefixesClient struct{ legacy legacy.Client }

func (c *ifacePrefixesClient) List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.PrefixList, error) {
	return c.legacy.ListPrefixes(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *ifacePrefixesClient) Create(ctx context.Context, prefix *api.Prefix, opts ...CallOption) (*api.Prefix, error) {
	return c.legacy.CreatePrefix(ctx, prefix, toLegacyIgnored(opts...)...)
}
func (c *ifacePrefixesClient) Delete(ctx context.Context, interfaceID string, prefix *netip.Prefix, opts ...CallOption) (*api.Prefix, error) {
	return c.legacy.DeletePrefix(ctx, interfaceID, prefix, toLegacyIgnored(opts...)...)
}

//
// Routes
//

type Routes interface {
	List(ctx context.Context, vni uint32, opts ...CallOption) (*api.RouteList, error)
	Create(ctx context.Context, route *api.Route, opts ...CallOption) (*api.Route, error)
	Delete(ctx context.Context, vni uint32, prefix *netip.Prefix, opts ...CallOption) (*api.Route, error)
}

type routeClient struct{ legacy legacy.Client }

func (c *routeClient) List(ctx context.Context, vni uint32, opts ...CallOption) (*api.RouteList, error) {
	return c.legacy.ListRoutes(ctx, vni, toLegacyIgnored(opts...)...)
}
func (c *routeClient) Create(ctx context.Context, route *api.Route, opts ...CallOption) (*api.Route, error) {
	return c.legacy.CreateRoute(ctx, route, toLegacyIgnored(opts...)...)
}
func (c *routeClient) Delete(ctx context.Context, vni uint32, prefix *netip.Prefix, opts ...CallOption) (*api.Route, error) {
	return c.legacy.DeleteRoute(ctx, vni, prefix, toLegacyIgnored(opts...)...)
}

//
// NATs
//

type NATs interface {
	Get(ctx context.Context, interfaceID string, opts ...CallOption) (*api.Nat, error)
	Create(ctx context.Context, nat *api.Nat, opts ...CallOption) (*api.Nat, error)
	Delete(ctx context.Context, interfaceID string, opts ...CallOption) (*api.Nat, error)

	ListAny(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error)
	ListLocal(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error)
	ListNeighbors(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error)
	CreateNeighbor(ctx context.Context, n *api.NeighborNat, opts ...CallOption) (*api.NeighborNat, error)
	DeleteNeighbor(ctx context.Context, n *api.NeighborNat, opts ...CallOption) (*api.NeighborNat, error)
}

type natClient struct{ legacy legacy.Client }

func (c *natClient) Get(ctx context.Context, interfaceID string, opts ...CallOption) (*api.Nat, error) {
	return c.legacy.GetNat(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *natClient) Create(ctx context.Context, nat *api.Nat, opts ...CallOption) (*api.Nat, error) {
	return c.legacy.CreateNat(ctx, nat, toLegacyIgnored(opts...)...)
}
func (c *natClient) Delete(ctx context.Context, interfaceID string, opts ...CallOption) (*api.Nat, error) {
	return c.legacy.DeleteNat(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *natClient) ListAny(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error) {
	return c.legacy.ListNats(ctx, natIP, "any", toLegacyIgnored(opts...)...)
}
func (c *natClient) ListLocal(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error) {
	return c.legacy.ListLocalNats(ctx, natIP, toLegacyIgnored(opts...)...)
}
func (c *natClient) ListNeighbors(ctx context.Context, natIP *netip.Addr, opts ...CallOption) (*api.NatList, error) {
	return c.legacy.ListNeighborNats(ctx, natIP, toLegacyIgnored(opts...)...)
}
func (c *natClient) CreateNeighbor(ctx context.Context, n *api.NeighborNat, opts ...CallOption) (*api.NeighborNat, error) {
	return c.legacy.CreateNeighborNat(ctx, n, toLegacyIgnored(opts...)...)
}
func (c *natClient) DeleteNeighbor(ctx context.Context, n *api.NeighborNat, opts ...CallOption) (*api.NeighborNat, error) {
	return c.legacy.DeleteNeighborNat(ctx, n, toLegacyIgnored(opts...)...)
}

//
// Firewall
//

type Firewall interface {
	List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.FirewallRuleList, error)
	Get(ctx context.Context, interfaceID string, ruleID string, opts ...CallOption) (*api.FirewallRule, error)
	Create(ctx context.Context, rule *api.FirewallRule, opts ...CallOption) (*api.FirewallRule, error)
	Delete(ctx context.Context, interfaceID string, ruleID string, opts ...CallOption) (*api.FirewallRule, error)
}

type fwClient struct{ legacy legacy.Client }

func (c *fwClient) List(ctx context.Context, interfaceID string, opts ...CallOption) (*api.FirewallRuleList, error) {
	return c.legacy.ListFirewallRules(ctx, interfaceID, toLegacyIgnored(opts...)...)
}
func (c *fwClient) Get(ctx context.Context, interfaceID string, ruleID string, opts ...CallOption) (*api.FirewallRule, error) {
	return c.legacy.GetFirewallRule(ctx, interfaceID, ruleID, toLegacyIgnored(opts...)...)
}
func (c *fwClient) Create(ctx context.Context, rule *api.FirewallRule, opts ...CallOption) (*api.FirewallRule, error) {
	return c.legacy.CreateFirewallRule(ctx, rule, toLegacyIgnored(opts...)...)
}
func (c *fwClient) Delete(ctx context.Context, interfaceID string, ruleID string, opts ...CallOption) (*api.FirewallRule, error) {
	return c.legacy.DeleteFirewallRule(ctx, interfaceID, ruleID, toLegacyIgnored(opts...)...)
}

//
// System
//

type System interface {
	CheckInitialized(ctx context.Context, opts ...CallOption) (*api.Initialized, error)
	Initialize(ctx context.Context, opts ...CallOption) (*api.Initialized, error)
	GetVni(ctx context.Context, vni uint32, vniType uint8, opts ...CallOption) (*api.Vni, error)
	ResetVni(ctx context.Context, vni uint32, vniType uint8, opts ...CallOption) (*api.Vni, error)
	GetVersion(ctx context.Context, version *api.Version, opts ...CallOption) (*api.Version, error)
}

type systemClient struct{ legacy legacy.Client }

func (c *systemClient) CheckInitialized(ctx context.Context, opts ...CallOption) (*api.Initialized, error) {
	return c.legacy.CheckInitialized(ctx, toLegacyIgnored(opts...)...)
}
func (c *systemClient) Initialize(ctx context.Context, opts ...CallOption) (*api.Initialized, error) {
	return c.legacy.Initialize(ctx, toLegacyIgnored(opts...)...)
}
func (c *systemClient) GetVni(ctx context.Context, vni uint32, vniType uint8, opts ...CallOption) (*api.Vni, error) {
	return c.legacy.GetVni(ctx, vni, vniType, toLegacyIgnored(opts...)...)
}
func (c *systemClient) ResetVni(ctx context.Context, vni uint32, vniType uint8, opts ...CallOption) (*api.Vni, error) {
	return c.legacy.ResetVni(ctx, vni, vniType, toLegacyIgnored(opts...)...)
}
func (c *systemClient) GetVersion(ctx context.Context, version *api.Version, opts ...CallOption) (*api.Version, error) {
	return c.legacy.GetVersion(ctx, version, toLegacyIgnored(opts...)...)
}

//
// Capture
//

type Capture interface {
	Start(ctx context.Context, capture *api.CaptureStart, opts ...CallOption) (*api.CaptureStart, error)
	Stop(ctx context.Context, opts ...CallOption) (*api.CaptureStop, error)
	Status(ctx context.Context, opts ...CallOption) (*api.CaptureStatus, error)
}

type captureClient struct{ legacy legacy.Client }

func (c *captureClient) Start(ctx context.Context, capture *api.CaptureStart, opts ...CallOption) (*api.CaptureStart, error) {
	return c.legacy.CaptureStart(ctx, capture, toLegacyIgnored(opts...)...)
}
func (c *captureClient) Stop(ctx context.Context, opts ...CallOption) (*api.CaptureStop, error) {
	return c.legacy.CaptureStop(ctx, toLegacyIgnored(opts...)...)
}
func (c *captureClient) Status(ctx context.Context, opts ...CallOption) (*api.CaptureStatus, error) {
	return c.legacy.CaptureStatus(ctx, toLegacyIgnored(opts...)...)
}
