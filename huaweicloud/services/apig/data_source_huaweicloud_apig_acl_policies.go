package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/acls
func DataSourceAclPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAclPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the ACL policies belong.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the ACL policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the ACL policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the ACL policy.`,
			},
			"entity_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the entity type of the ACL policy.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All ACL policies that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the ACL policy.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the ACL policy.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the ACL policy.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the ACL policy.`,
						},
						"bind_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of bound APIs.`,
						},
						"entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The entity type of the ACL policy.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the policy.`,
						},
					},
				},
			},
		},
	}
}

func buildListAclPoliciesParams(d *schema.ResourceData) string {
	res := ""
	if policyId, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, policyId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&acl_name=%v", res, name)
	}
	if aclType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&acl_type=%v", res, aclType)
	}
	if entityType, ok := d.GetOk("entity_type"); ok {
		res = fmt.Sprintf("%s&entity_type=%v", res, entityType)
	}
	return res
}

func queryAclPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/acls?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListAclPoliciesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving ACL policies under specified "+
				"dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		policies := utils.PathSearch("acls", respBody, make([]interface{}, 0)).([]interface{})
		if len(policies) < 1 {
			break
		}
		result = append(result, policies...)
		offset += len(policies)
	}
	return result, nil
}

func dataSourceAclPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	policies, err := queryAclPolicies(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenAclPolicies(policies)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAclPolicies(policies []interface{}) []interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", policy, nil),
			"name":        utils.PathSearch("acl_name", policy, nil),
			"type":        utils.PathSearch("acl_type", policy, nil),
			"value":       utils.PathSearch("acl_value", policy, nil),
			"bind_num":    utils.PathSearch("bind_num", policy, nil),
			"entity_type": utils.PathSearch("entity_type", policy, nil),
			"updated_at":  utils.PathSearch("update_time", policy, nil),
		})
	}
	return result
}
