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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/binded-acls
func DataSourceApiAssociatedAclPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiAssociatedAclPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the ACL policies belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API bound to the ACL policy.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the ACL policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the ACL policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the ACL policy.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment where the API is published.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"entity_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The entity type of the ACL policy.`,
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
							Description: `One or more objects from which the access will be controlled.`,
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment where the API is published.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment where the API is published.`,
						},
						"entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The entity type of the ACL policy.`,
						},
						"bind_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bind ID.`,
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time that the ACL policy is bound to the API.`,
						},
					},
				},
			},
		},
	}
}

func buildListApiAssociatedAclPoliciesParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	if v, ok := d.GetOk("env_name"); ok {
		res = fmt.Sprintf("%s&env_name=%v", res, v)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&acl_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&acl_name=%v", res, v)
	}
	return res
}

func queryApiAssociatedAclPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/binded-acls?api_id={api_id}"
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)

	queryParams := buildListApiAssociatedAclPoliciesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving associated ACL policies (bound to the API: %s) under specified "+
				"dedicated instance (%s): %s", apiId, instanceId, err)
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

func dataSourceApiAssociatedAclPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	policies, err := queryApiAssociatedAclPolicies(client, d)
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
		d.Set("policies", filterAssociatedAclPolicies(flattenAssociatedAclPolicies(policies), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterAssociatedAclPolicies(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("entity_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("entity_type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenAssociatedAclPolicies(policies []interface{}) []interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("acl_id", policy, nil),
			"name":        utils.PathSearch("acl_name", policy, nil),
			"type":        utils.PathSearch("acl_type", policy, nil),
			"value":       utils.PathSearch("acl_value", policy, nil),
			"entity_type": utils.PathSearch("entity_type", policy, nil),
			"env_id":      utils.PathSearch("env_id", policy, nil),
			"env_name":    utils.PathSearch("env_name", policy, nil),
			"bind_id":     utils.PathSearch("bind_id", policy, nil),
			"bind_time":   utils.PathSearch("bind_time", policy, nil),
		})
	}
	return result
}
