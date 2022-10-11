package integration

import (
	"context"
	"github.com/linode/linodego"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache_RegionList(t *testing.T) {
	validateResult := func(r []linodego.Region, err error) {
		if err != nil {
			t.Fatal(err)
		}

		if len(r) == 0 {
			t.Fatalf("expected a list of regions - %v", r)
		}
	}

	client, teardown := createTestClient(t, "fixtures/TestCache_RegionList")
	defer teardown()

	// Collect request number
	totalRequests := int64(0)

	client.OnBeforeRequest(func(request *linodego.Request) error {
		if !strings.Contains(request.URL, "regions") {
			return nil
		}

		atomic.AddInt64(&totalRequests, 1)
		return nil
	})

	// First request (no cache)
	validateResult(client.ListRegions(context.Background(), nil))

	// Second request (cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Clear cache
	client.InvalidateCache()

	// Third request (non-cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Invalidate the region response
	if err := client.InvalidateCacheEndpoint("/regions"); err != nil {
		t.Fatal(err)
	}

	// Fourth request (non-cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Fifth request (cache disabled)
	client.UseCache(false)
	validateResult(client.ListRegions(context.Background(), nil))

	// Sixth request (cache disabled)
	validateResult(client.ListRegions(context.Background(), nil))

	// Validate request count
	if totalRequests != 4 {
		t.Fatalf("expected 4 requests, got %d", totalRequests)
	}
}

func TestCache_Expiration(t *testing.T) {
	validateResult := func(r []linodego.Region, err error) {
		if err != nil {
			t.Fatal(err)
		}

		if len(r) == 0 {
			t.Fatalf("expected a list of regions - %v", r)
		}
	}

	client, teardown := createTestClient(t, "fixtures/TestCache_Expiration")
	defer teardown()

	// Collect request number
	totalRequests := int64(0)

	client.OnBeforeRequest(func(request *linodego.Request) error {
		if !strings.Contains(request.URL, "regions") {
			return nil
		}

		atomic.AddInt64(&totalRequests, 1)
		return nil
	})

	// First request (no cache)
	validateResult(client.ListRegions(context.Background(), nil))

	// Second request (cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Entries should expire immediately
	client.SetCacheExpiration(0)

	// Third request (non-cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Entries shouldn't expire
	client.SetCacheExpiration(time.Hour)

	// Fourth request (cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Validate request count
	if totalRequests != 2 {
		t.Fatalf("expected 2 requests, got %d", totalRequests)
	}
}
