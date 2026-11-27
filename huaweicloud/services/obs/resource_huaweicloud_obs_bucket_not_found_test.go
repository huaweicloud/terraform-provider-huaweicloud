package obs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// newNotFoundOBSConfig returns a provider configuration whose OBS endpoint answers every request with
// 404 NoSuchBucket, as OBS does after the bucket has been deleted outside Terraform.
func newNotFoundOBSConfig(t *testing.T) *config.Config {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("x-obs-error-code", "NoSuchBucket")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>NoSuchBucket</Code>` +
			`<Message>The specified bucket does not exist</Message></Error>`))
	}))
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse the server URL: %s", err)
	}

	domainClient := &golangsdk.ProviderClient{}
	domainClient.HTTPClient = http.Client{
		Transport: &config.LogRoundTripper{Rt: &http.Transport{}},
	}

	return &config.Config{
		AccessKey: "access-key",
		SecretKey: "secret-key",
		DomainID:  "domain-id",
		Region:    "region-1",
		// "localhost" has no dots, so the endpoint is used as it is.
		Endpoints:    map[string]string{"obs": "http://localhost:" + serverURL.Port()},
		DomainClient: domainClient,
	}
}

func TestOBSBucketSubResources_bucketNotFound(t *testing.T) {
	cases := []struct {
		name     string
		resource *schema.Resource
		raw      map[string]interface{}
		id       string
		call     func(*schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
	}{
		{
			name:     "huaweicloud_obs_bucket_acl read",
			resource: ResourceOBSBucketAcl(),
			raw:      map[string]interface{}{"bucket": "test-bucket"},
			id:       "test-bucket",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.ReadContext
			},
		},
		{
			name:     "huaweicloud_obs_bucket_acl delete",
			resource: ResourceOBSBucketAcl(),
			raw:      map[string]interface{}{"bucket": "test-bucket"},
			id:       "test-bucket",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.DeleteContext
			},
		},
		{
			name:     "huaweicloud_obs_bucket_object_acl read",
			resource: ResourceOBSBucketObjectAcl(),
			raw:      map[string]interface{}{"bucket": "test-bucket", "key": "test-object"},
			id:       "test-object",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.ReadContext
			},
		},
		{
			name:     "huaweicloud_obs_bucket_object_acl delete",
			resource: ResourceOBSBucketObjectAcl(),
			raw:      map[string]interface{}{"bucket": "test-bucket", "key": "test-object"},
			id:       "test-object",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.DeleteContext
			},
		},
		{
			name:     "huaweicloud_obs_bucket_replication read",
			resource: ResourceObsBucketReplication(),
			raw:      map[string]interface{}{"bucket": "test-bucket", "destination_bucket": "dest-bucket", "agency": "agency"},
			id:       "test-bucket",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.ReadContext
			},
		},
		{
			name:     "huaweicloud_obs_bucket_replication delete",
			resource: ResourceObsBucketReplication(),
			raw:      map[string]interface{}{"bucket": "test-bucket", "destination_bucket": "dest-bucket", "agency": "agency"},
			id:       "test-bucket",
			call: func(r *schema.Resource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return r.DeleteContext
			},
		},
	}

	cfg := newNotFoundOBSConfig(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tc.resource.Schema, tc.raw)
			d.SetId(tc.id)

			diags := tc.call(tc.resource)(context.Background(), d, cfg)
			if diags.HasError() {
				t.Fatalf("expected the resource to be removed from the state, got error: %v", diags)
			}
			if d.Id() != "" {
				t.Errorf("expected an empty resource ID, got %q", d.Id())
			}
		})
	}
}
