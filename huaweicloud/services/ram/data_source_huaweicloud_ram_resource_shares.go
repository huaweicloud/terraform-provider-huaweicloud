package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/resource-shares/search
func DataSourceRAMShares() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRAMSharesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the owner type of resource sharing instance.`,
			},
			"resource_share_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of resource share IDs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the resource share.`,
			},
			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags attached to the resource share.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the identifier or name of the tag key.`,
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the list of values for the tag key.`,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the resource share.`,
			},
			"permission_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the permission ID.`,
			},
			"resource_shares": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of details about resource shares.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource share.`,
						},
						"allow_external_principals": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether resources can be shared with any accounts outside the organization.`,
						},
						"tags": common.TagsComputedSchema(),
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the resource share was created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the resource share was last updated.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource share.`,
						},
						"owning_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource owner in a resource share.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Status of the resource share.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the resource share.`,
						},
					},
				},
			},
		},
	}
}

func buildRAMSharesTagFilterRequestBodyParams(tagFilterRaw interface{}) []interface{} {
	if tagFilterRaw == nil {
		return nil
	}

	tagFilterArray := tagFilterRaw.([]interface{})
	rst := make([]interface{}, 0, len(tagFilterArray))
	for _, v := range tagFilterArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}
	return rst
}

// There is no value range marked limit in openapi, so the limit value is not configured here.
func buildRAMSharesRequestBodyParams(d *schema.ResourceData, nextMarker string) map[string]interface{} {
	requestParams := map[string]interface{}{
		"resource_owner":        d.Get("resource_owner"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"permission_id":         utils.ValueIgnoreEmpty(d.Get("permission_id")),
		"resource_share_ids":    utils.ValueIgnoreEmpty(d.Get("resource_share_ids")),
		"resource_share_status": utils.ValueIgnoreEmpty(d.Get("status")),
		"tag_filters":           utils.ValueIgnoreEmpty(buildRAMSharesTagFilterRequestBodyParams(d.Get("tag_filters"))),
	}

	if nextMarker != "" {
		requestParams["marker"] = nextMarker
	}
	return requestParams
}

func dataSourceRAMSharesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		mErr                *multierror.Error
		nextMarker          string
		httpUrl             = "v1/resource-shares/search"
		product             = "ram"
		totalResourceShares []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestOpt.JSONBody = utils.RemoveNil(buildRAMSharesRequestBodyParams(d, nextMarker))
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource shares: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceShares := utils.PathSearch("resource_shares", respBody, make([]interface{}, 0)).([]interface{})
		if len(resourceShares) > 0 {
			totalResourceShares = append(totalResourceShares, resourceShares...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resource_shares", flattenRAMResourceShares(totalResourceShares)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRAMResourceShares(totalResourceShares []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(totalResourceShares))
	for _, v := range totalResourceShares {
		rst = append(rst, map[string]interface{}{
			"name":                      utils.PathSearch("name", v, nil),
			"allow_external_principals": utils.PathSearch("allow_external_principals", v, nil),
			"tags":                      utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"created_at":                utils.PathSearch("created_at", v, nil),
			"updated_at":                utils.PathSearch("updated_at", v, nil),
			"id":                        utils.PathSearch("id", v, nil),
			"owning_account_id":         utils.PathSearch("owning_account_id", v, nil),
			"status":                    utils.PathSearch("status", v, nil),
			"description":               utils.PathSearch("description", v, nil),
		})
	}
	return rst
}
