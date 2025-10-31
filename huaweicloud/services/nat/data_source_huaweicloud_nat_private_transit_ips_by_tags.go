package nat

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v3/{project_id}/transit-ips/resource_instances/action
func DataSourceNatPrivateTransitIpsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatPrivateTransitIpsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func flattenNatPrivateTransitIpsByTagsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": utils.PathSearch("resource_detail", v, nil),
			"tags":            flattenNatTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return resources
}

func buildNatPrivateTransitIpsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}
	if v, ok := d.GetOk("action"); ok {
		body["action"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		body["tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags_any"); ok {
		body["tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags"); ok {
		body["not_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags_any"); ok {
		body["not_tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("matches"); ok {
		body["matches"] = expandMatches(v.([]interface{}))
	}
	return body
}

func dataSourceNatPrivateTransitIpsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v3/{project_id}/transit-ips/resource_instances/action"
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	var allTrasitIps []interface{}
	offset := 0
	totalCount := 0

	for {
		requestBody := buildNatPrivateTransitIpsByTagsBodyParams(d)
		if requestBody["action"] == "filter" {
			requestBody["limit"] = "1000"
			requestBody["offset"] = strconv.Itoa(offset)
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         requestBody,
		}
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying Nat Private Transit Ips by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		transitIps := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if len(transitIps) == 0 {
			break
		}
		allTrasitIps = append(allTrasitIps, transitIps...)
		offset += len(transitIps)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenNatPrivateTransitIpsByTagsResponseBody(allTrasitIps)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
