// SPDX-FileCopyrightText: 2025 The dpservice Authors
// SPDX-License-Identifier: Apache-2.0

// Package clientv2 provides a small, domain-focused Go client for dpservice.
//
// It organizes operations by resource domains and sub-resources for better
// discoverability and testability. The v2 client can be constructed directly
// from a generated gRPC client, or adapted from the legacy client while you
// migrate.
//
// # Basic usage
//
//	v2 := clientv2.NewFromProto(rpc)
//
//	// Load balancers
//	lb, err := v2.LoadBalancers().Get(ctx, "lb-1")
//	_ = lb; _ = err
//
//	// LB sub-resources
//	_, _ = v2.LoadBalancers().Prefixes().List(ctx, "iface-1", clientv2.WithIgnoredCodes(1001))
//	_, _ = v2.LoadBalancers().Targets().List(ctx, "lb-1")
//
//	// Interfaces and sub-resources
//	_, _ = v2.Interfaces().Get(ctx, "iface-1")
//	_, _ = v2.Interfaces().VIP().Get(ctx, "iface-1")
//	_, _ = v2.Interfaces().Prefixes().List(ctx, "iface-1")
//	_, _ = v2.Interfaces().Firewall().List(ctx, "iface-1")
//
//	// Routes
//	_, _ = v2.Routes().List(ctx, 42)
//
//	// NATs
//	_, _ = v2.NATs().ListAny(ctx, &natIP)
//
//	// System
//	_, _ = v2.System().GetVersion(ctx, &api.Version{})
//
//	// Capture
//	_, _ = v2.Capture().Status(ctx)
//
// Migration from legacy
//
//	// If you already have a legacy client, adapt it without changing call sites
//	// that pass the legacy type around.
//	v2 := clientv2.AsV2(legacyClient)
//	_, _ = v2.Routes().List(ctx, 42)
package clientv2
