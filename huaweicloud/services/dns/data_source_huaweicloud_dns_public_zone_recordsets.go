package dns

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

// @API DNS GET /v2.1/zones/{zone_id}/email-recordsets
// @API DNS GET /v2.1/zones/{zone_id}/website-recordsets
func DataSourcePublicZoneRecordsets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicZoneRecordsetsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the public zone.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the domain name to be queried.`,
			},
			"recordsets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the recordset.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the recordset.`,
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the zone to which the recordset belongs.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the recordset.`,
						},
						"default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the recordset is default.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the recordset, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the recordset, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of recordsets.`,
			},
		},
	}
}

func listPublicZoneRecordsets(client *golangsdk.ServiceClient, zoneId, queryType string) ([]interface{}, error) {
	var (
		httpUrl = "v2.1/zones/{zone_id}/{type}-recordsets"
		limit   = 500
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{type}", queryType)
	listPath = strings.ReplaceAll(listPath, "{zone_id}", zoneId)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		recordsets := utils.PathSearch("recordsets", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, recordsets...)

		if len(recordsets) < limit {
			break
		}
		offset += len(recordsets)
	}

	return result, nil
}

func dataSourcePublicZoneRecordsetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)
	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	zoneId := d.Get("zone_id").(string)
	queryType := d.Get("type").(string)
	recordsets, err := listPublicZoneRecordsets(client, zoneId, queryType)
	if err != nil {
		return diag.Errorf("error querying public zone %s recordsets: %s", queryType, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("recordsets", flattenPublicZoneRecordsets(recordsets)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPublicZoneRecordsets(recordsets []interface{}) []map[string]interface{} {
	if len(recordsets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(recordsets))
	for _, item := range recordsets {
		result = append(result, map[string]interface{}{
			"id":      utils.PathSearch("id", item, nil),
			"name":    utils.PathSearch("name", item, nil),
			"zone_id": utils.PathSearch("zone_id", item, nil),
			"type":    utils.PathSearch("type", item, nil),
			"default": utils.PathSearch("default", item, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_at",
				item, "").(string), "2006-01-02T15:04:05")/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_at",
				item, "").(string), "2006-01-02T15:04:05")/1000, false),
		})
	}

	return result
}
