// SPDX-FileCopyrightText: 2025 The dpservice Authors
// SPDX-License-Identifier: Apache-2.0

package clientv2_test

import (
	"context"
	"net/netip"

	"github.com/ironcore-dev/dpservice/go/dpservice-go/api"
	clientv2 "github.com/ironcore-dev/dpservice/go/dpservice-go/clientv2"
	dpdkproto "github.com/ironcore-dev/dpservice/go/dpservice-go/proto"
)

func ExampleNewFromProto() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)
	_ = v2
	_ = ctx
	// Output:
}

func ExampleClient_LoadBalancers() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)

	// Core LB CRUD
	_, _ = v2.LoadBalancers().Get(ctx, "lb-1")
	_, _ = v2.LoadBalancers().List(ctx)

	_, _ = v2.LoadBalancers().Create(ctx, &api.LoadBalancer{
		LoadBalancerMeta: api.LoadBalancerMeta{ID: "lb-1"},
	})
	_, _ = v2.LoadBalancers().Delete(ctx, "lb-1")

	// LB sub-resources: prefixes
	_, _ = v2.LoadBalancers().Prefixes().List(ctx, "iface-1", clientv2.WithIgnoredCodes(1001))
	_, _ = v2.LoadBalancers().Prefixes().Create(ctx, &api.LoadBalancerPrefix{
		LoadBalancerPrefixMeta: api.LoadBalancerPrefixMeta{InterfaceID: "iface-1"},
	})
	var pfx netip.Prefix
	_, _ = v2.LoadBalancers().Prefixes().Delete(ctx, "iface-1", &pfx)

	// LB sub-resources: targets
	_, _ = v2.LoadBalancers().Targets().List(ctx, "lb-1")
	_, _ = v2.LoadBalancers().Targets().Create(ctx, &api.LoadBalancerTarget{
		LoadBalancerTargetMeta: api.LoadBalancerTargetMeta{LoadbalancerID: "lb-1"},
	})
	ip := netip.Addr{}
	_, _ = v2.LoadBalancers().Targets().Delete(ctx, "lb-1", &ip)
	// Output:
}

func ExampleClient_Interfaces() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)

	_, _ = v2.Interfaces().Get(ctx, "iface-1")
	_, _ = v2.Interfaces().List(ctx)
	_, _ = v2.Interfaces().Create(ctx, &api.Interface{InterfaceMeta: api.InterfaceMeta{ID: "iface-1"}})
	_, _ = v2.Interfaces().Delete(ctx, "iface-1")

	// VIP
	_, _ = v2.Interfaces().VIP().Get(ctx, "iface-1")
	_, _ = v2.Interfaces().VIP().Create(ctx, &api.VirtualIP{VirtualIPMeta: api.VirtualIPMeta{InterfaceID: "iface-1"}})
	_, _ = v2.Interfaces().VIP().Delete(ctx, "iface-1")

	// Interface prefixes
	_, _ = v2.Interfaces().Prefixes().List(ctx, "iface-1")
	_, _ = v2.Interfaces().Prefixes().Create(ctx, &api.Prefix{PrefixMeta: api.PrefixMeta{InterfaceID: "iface-1"}})
	var ip netip.Addr
	p := netip.PrefixFrom(ip, 24)
	_, _ = v2.Interfaces().Prefixes().Delete(ctx, "iface-1", &p)

	// Firewall rules (interface-scoped)
	_, _ = v2.Interfaces().Firewall().List(ctx, "iface-1")
	_, _ = v2.Interfaces().Firewall().Create(ctx, &api.FirewallRule{FirewallRuleMeta: api.FirewallRuleMeta{InterfaceID: "iface-1"}})
	_, _ = v2.Interfaces().Firewall().Get(ctx, "iface-1", "rule-1")
	_, _ = v2.Interfaces().Firewall().Delete(ctx, "iface-1", "rule-1")
	// Output:
}

func ExampleClient_Routes() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)
	_, _ = v2.Routes().List(ctx, 42)
	_, _ = v2.Routes().Create(ctx, &api.Route{RouteMeta: api.RouteMeta{VNI: 42}})
	var ip netip.Addr
	p := netip.PrefixFrom(ip, 24)
	_, _ = v2.Routes().Delete(ctx, 42, &p)
	// Output:
}

func ExampleClient_NATs() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)
	_, _ = v2.NATs().Get(ctx, "iface-1")
	_, _ = v2.NATs().Create(ctx, &api.Nat{NatMeta: api.NatMeta{InterfaceID: "iface-1"}})
	_, _ = v2.NATs().Delete(ctx, "iface-1")

	var natIP netip.Addr
	_, _ = v2.NATs().ListAny(ctx, &natIP)
	_, _ = v2.NATs().ListLocal(ctx, &natIP)
	_, _ = v2.NATs().ListNeighbors(ctx, &natIP)
	_, _ = v2.NATs().CreateNeighbor(ctx, &api.NeighborNat{})
	_, _ = v2.NATs().DeleteNeighbor(ctx, &api.NeighborNat{})
	// Output:
}

func ExampleClient_System() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)
	_, _ = v2.System().CheckInitialized(ctx)
	_, _ = v2.System().Initialize(ctx)
	_, _ = v2.System().GetVni(ctx, 42, 1)
	_, _ = v2.System().ResetVni(ctx, 42, 1)
	_, _ = v2.System().GetVersion(ctx, &api.Version{})
	// Output:
}

func ExampleClient_Capture() {
	var rpc dpdkproto.DPDKironcoreClient
	ctx := context.TODO()

	v2 := clientv2.NewFromProto(rpc)
	_, _ = v2.Capture().Start(ctx, &api.CaptureStart{})
	_, _ = v2.Capture().Status(ctx)
	_, _ = v2.Capture().Stop(ctx)
	// Output:
}
