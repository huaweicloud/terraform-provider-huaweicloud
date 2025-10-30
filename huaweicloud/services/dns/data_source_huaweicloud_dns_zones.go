package dns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS GET /v2/zones
func DataSourceZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceZonesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The region where the zones are located.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
			"zone_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The zone type. The value can be **public** or **private**.`,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the zone.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource tag.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The zone name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The zone status.`,
			},
			"search_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The query criteria search mode.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID which the zone associated.`,
			},
			"router_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the VPC associated with the private zone.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting filed for the list of the zones.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting mode for the list of the zones.`,
			},
			"zones": {
				Type:        schema.TypeList,
				Elem:        zoneSchema(),
				Computed:    true,
				Description: `The list of zones.`,
			},
		},
	}
}

func zoneSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone description.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address of the administrator managing the zone.`,
			},
			"zone_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone type.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time to live (TTL) of the zone.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone status.`,
			},
			"record_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of record sets in the zone.`,
			},
			"masters": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The master DNS servers, from which the slave servers get DNS information.`,
			},
			"routers": {
				Type:        schema.TypeList,
				Elem:        zoneRouterSchema(),
				Computed:    true,
				Description: `The list of VPCs associated with the zone.`,
			},
			"tags": common.TagsComputedSchema(),
			"proxy_pattern": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recursive resolution proxy mode for subdomains of the private zone.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the zone, in RFC3339 format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the zone, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the zone, in RFC3339 format.`,
			},
		},
	}
	return &sc
}

func zoneRouterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC associated with the zone.`,
			},
			"router_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region of the VPC.`,
			},
		},
	}
	return &sc
}

func resourceZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = ""
		mErr        *multierror.Error
		listHttpUrl = "v2/zones"
		dnsProduct  = "dns"
	)
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
		cfg.RegionClient = true
	}

	client, err := cfg.NewServiceClient(dnsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath += buildListZonesQueryParams(d)
	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DNS zones, %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("zones", flattenListZones(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListZones(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("zones", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"email":                 utils.PathSearch("email", v, nil),
			"zone_type":             utils.PathSearch("zone_type", v, nil),
			"ttl":                   utils.PathSearch("ttl", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"record_num":            utils.PathSearch("record_num", v, nil),
			"masters":               utils.PathSearch("masters", v, nil),
			"routers":               flattenZonesRouters(v),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"proxy_pattern":         utils.PathSearch("proxy_pattern", v, nil),
			"pool_id":               utils.PathSearch("pool_id", v, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
				v, "").(string), "2006-01-02T15:04:05")/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
				v, "").(string), "2006-01-02T15:04:05")/1000, false),
		}
	}
	return rst
}

func flattenZonesRouters(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("routers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"router_id":     utils.PathSearch("router_id", v, nil),
			"router_region": utils.PathSearch("router_region", v, nil),
		}
	}
	return rst
}

func buildListZonesQueryParams(d *schema.ResourceData) string {
	queryParam := fmt.Sprintf("?type=%v", d.Get("zone_type"))
	if v, ok := d.GetOk("zone_id"); ok {
		queryParam = fmt.Sprintf("%s&id=%v", queryParam, v)
	}

	if v, ok := d.GetOk("tags"); ok {
		queryParam = fmt.Sprintf("%s&tags=%v", queryParam, v)
	}

	if v, ok := d.GetOk("name"); ok {
		queryParam = fmt.Sprintf("%s&name=%v", queryParam, v)
	}

	if v, ok := d.GetOk("status"); ok {
		queryParam = fmt.Sprintf("%s&status=%v", queryParam, v)
	}

	if v, ok := d.GetOk("search_mode"); ok {
		queryParam = fmt.Sprintf("%s&search_mode=%v", queryParam, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryParam = fmt.Sprintf("%s&enterprise_project_id=%v", queryParam, v)
	}

	if v, ok := d.GetOk("router_id"); ok {
		queryParam = fmt.Sprintf("%s&router_id=%v", queryParam, v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		queryParam = fmt.Sprintf("%s&sort_key=%v", queryParam, v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		queryParam = fmt.Sprintf("%s&sort_dir=%v", queryParam, v)
	}
	return queryParam
}
