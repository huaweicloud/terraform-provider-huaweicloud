package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityV5PolicyVersions
// @API IAM GET /v5/policies/{policy_id}/versions
// @API IAM GET /v5/policies/{policy_id}/versions/{version_id}
func DataSourceIdentityV5PolicyVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5PolicyVersionsRead,

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5PolicyVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allVersions []interface{}
	var marker string
	policyId := d.Get("policy_id").(string)
	versionId := d.Get("version_id").(string)
	reqOpt := &golangsdk.RequestOpts{KeepResponseBody: true}
	if versionId != "" {
		// Get specific policy version
		path := client.Endpoint + "v5/policies/{policy_id}/versions/{version_id}"
		path = strings.ReplaceAll(path, "{policy_id}", policyId)
		path = strings.ReplaceAll(path, "{version_id}", versionId)

		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving policy version: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		version := utils.PathSearch("policy_version", resp, map[string]interface{}{}).(map[string]interface{})
		versions := map[string]interface{}{
			"versions": []interface{}{version},
		}
		allVersions = append(allVersions, flattenListPolicyVersionsV5Response(versions)...)
	} else {
		// List policy versions
		for {
			path := client.Endpoint + "v5/policies/{policy_id}/versions" + buildListPolicyVersionsV5Params(marker)
			path = strings.ReplaceAll(path, "{policy_id}", policyId)

			r, err := client.Request("GET", path, reqOpt)
			if err != nil {
				return diag.Errorf("error retrieving policy versions: %s", err)
			}
			resp, err := utils.FlattenResponse(r)
			if err != nil {
				return diag.FromErr(err)
			}
			versions := flattenListPolicyVersionsV5Response(resp)
			allVersions = append(allVersions, versions...)

			marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
			if marker == "" {
				break
			}
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("versions", allVersions),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPolicyVersionsV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	versions := utils.PathSearch("versions", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(versions))
	for i, version := range versions {
		result[i] = map[string]interface{}{
			"version_id": utils.PathSearch("version_id", version, nil),
			"is_default": utils.PathSearch("is_default", version, false),
			"created_at": utils.PathSearch("created_at", version, nil),
			"document":   utils.PathSearch("document", version, nil),
		}
	}
	return result
}

func buildListPolicyVersionsV5Params(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
