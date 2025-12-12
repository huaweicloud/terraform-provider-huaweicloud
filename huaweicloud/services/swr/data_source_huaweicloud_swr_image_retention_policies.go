package swr

import (
	"context"
	"encoding/json"
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

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/retentions
func DataSourceImageRetentionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRetentionPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the data source. If omitted, the provider-level region will be used.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the organization to which the image belongs.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the repository to which the image belongs.`,
			},
			"retention_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All retention policies that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The image retention policy ID.`,
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image retention policy matching rule.`,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        retentionPoliciesRulesSchema(),
							Description: `The rules of the image retention policy.`,
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The reserved field.`,
						},
					},
				},
			},
		},
	}
}

func retentionPoliciesRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The template of the image retention policy.`,
			},
			"params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The params of matching template.`,
			},
			"tag_selectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The exception images.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The matching rule. The value can be **label** or **regexp**.`,
						},
						"pattern": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The pattern of the matching kind.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceImageRetentionPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	// Query the list of SWR image retention policies.
	resp, err := getImageRetentionPolicies(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("retention_policies", flattenImageRetentionPoliciesResponse(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getImageRetentionPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	listRetentionPoliciesHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/retentions"
	listRetentionPoliciesPath := client.Endpoint + listRetentionPoliciesHttpUrl
	listRetentionPoliciesPath = strings.ReplaceAll(listRetentionPoliciesPath, "{namespace}", organization)
	listRetentionPoliciesPath = strings.ReplaceAll(listRetentionPoliciesPath, "{repository}", repository)

	listRetentionPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listRetentionPoliciesResp, err := client.Request("GET", listRetentionPoliciesPath, &listRetentionPoliciesOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying SWR image retention policies: %s", err)
	}
	listRetentionPoliciesRespBody, err := utils.FlattenResponse(listRetentionPoliciesResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening SWR image retention policies: %s", err)
	}
	retentionPolicies := listRetentionPoliciesRespBody.([]interface{})

	return retentionPolicies, nil
}

func flattenImageRetentionPoliciesResponse(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	results := make([]interface{}, len(curArray))
	for i, v := range curArray {
		results[i] = map[string]interface{}{
			"id":        int(utils.PathSearch("id", v, float64(0)).(float64)),
			"algorithm": utils.PathSearch("algorithm", v, nil),
			"rules":     flattenImageRetentionPoliciesRulesResponse(utils.PathSearch("rules", v, nil)),
			"scope":     utils.PathSearch("scope", v, nil),
		}
	}
	return results
}

func flattenImageRetentionPoliciesRulesResponse(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	results := make([]interface{}, len(curArray))
	for i, v := range curArray {
		// If `template` is set to **date_rule**, the `params` to **{"days": "xxx"}**.
		// If `template` is set to **tag_rule**, set `params` to **{"num": "xxx"}**.
		// So `params` use string type, convert `params` object to json object string.
		paramsJson, _ := json.Marshal(utils.PathSearch("params", v, nil))
		results[i] = map[string]interface{}{
			"template":      utils.PathSearch("template", v, nil),
			"params":        string(paramsJson),
			"tag_selectors": flattenTagSelectorsResponse(utils.PathSearch("tag_selectors", v, nil)),
		}
	}
	return results
}

func flattenTagSelectorsResponse(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	results := make([]interface{}, len(curArray))
	for i, v := range curArray {
		results[i] = map[string]interface{}{
			"kind":    utils.PathSearch("kind", v, nil),
			"pattern": utils.PathSearch("pattern", v, nil),
		}
	}
	return results
}
