package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestCredentialAttributesAreSensitive verifies that the attributes carrying the kubeconfig and the client
// credentials are marked sensitive in the schema Terraform core receives.
//
// The terraform-plugin-sdk v2 derives the core schema of a computed list of objects with ImpliedType, which
// discards the Sensitive flag of the nested attributes (see helper/schema/core_schema.go, coreConfigSchemaType).
// The nested flags alone therefore do not hide anything in plan output, so the attributes must be marked at the
// top level.
func TestCredentialAttributesAreSensitive(t *testing.T) {
	cases := []struct {
		name       string
		resource   *schema.Resource
		attributes []string
	}{
		{
			name:       "huaweicloud_cce_cluster resource",
			resource:   ResourceCluster(),
			attributes: []string{"kube_config_raw", "certificate_users"},
		},
		{
			name:       "huaweicloud_cce_cluster data source",
			resource:   DataSourceCCEClusterV3(),
			attributes: []string{"kube_config_raw", "certificate_users"},
		},
		{
			name:       "huaweicloud_cce_cluster_certificate data source",
			resource:   DataSourceCCEClusterCertificate(),
			attributes: []string{"kube_config_raw", "users"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			block := tc.resource.CoreConfigSchema()
			for _, name := range tc.attributes {
				attr, ok := block.Attributes[name]
				if !ok {
					t.Fatalf("attribute %q is missing from the core schema", name)
				}
				if !attr.Sensitive {
					t.Errorf("attribute %q is not marked sensitive in the core schema", name)
				}
			}
		})
	}
}
