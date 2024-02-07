// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DNS
// ---------------------------------------------------------------

package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS GET /v2/zones/{zone_id}/recordsets
// @API DNS GET /v2.1/zones/{zone_id}/recordsets
func DataSourceRecordsets() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceRecordsetsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the zone ID.`,
			},
			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resolution line ID.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource tag.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the recordset to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the recordset type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the recordset to be queried. Fuzzy matching will work.`,
			},
			"recordset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the recordset to be queried. Fuzzy matching will work.`,
			},
			"search_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the query criteria search mode.`,
			},
			"recordsets": {
				Type:        schema.TypeList,
				Elem:        recordsetSchema(),
				Computed:    true,
				Description: `The list of recordsets.`,
			},
		},
	}
}

func recordsetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recordset ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recordset name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recordset description.`,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone ID of the recordset.`,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone name of the recordset.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recordset type.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The recordset caching duration (in seconds) on a local DNS server.`,
			},
			"records": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The recordset values.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recordset status.`,
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the record set is created by default. A default record set cannot be deleted.`,
			},
			"line_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resolution line ID.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The weight of the recordset.`,
			},
		},
	}
	return &sc
}

func resourceRecordsetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	zoneID := d.Get("zone_id").(string)
	client, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// The private zone can only use v2 version API. The public zone use v2.1 version API
	version := getApiVersionByZoneType(zoneType)
	listHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets", version)
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{zone_id}", d.Get("zone_id").(string))
	listPath += buildListRecordsetsQueryParams(d, zoneType)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DNS recordsets, %s", err)
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
		d.Set("recordsets", flattenListRecordsets(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListRecordsets(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("recordsets", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"zone_id":     utils.PathSearch("zone_id", v, nil),
			"zone_name":   utils.PathSearch("zone_name", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"ttl":         utils.PathSearch("ttl", v, nil),
			"records":     utils.PathSearch("records", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"default":     utils.PathSearch("default", v, nil),
			"line_id":     utils.PathSearch("line", v, nil),
			"weight":      utils.PathSearch("weight", v, nil),
		}
	}
	return rst
}

func buildListRecordsetsQueryParams(d *schema.ResourceData, zoneType string) string {
	queryParam := ""
	if v, ok := d.GetOk("line_id"); ok && zoneType == "public" {
		queryParam = fmt.Sprintf("%s&line_id=%v", queryParam, v)
	}

	if v, ok := d.GetOk("tags"); ok {
		queryParam = fmt.Sprintf("%s&tags=%v", queryParam, v)
	}

	if v, ok := d.GetOk("status"); ok {
		queryParam = fmt.Sprintf("%s&status=%v", queryParam, v)
	}

	if v, ok := d.GetOk("type"); ok {
		queryParam = fmt.Sprintf("%s&type=%v", queryParam, v)
	}

	if v, ok := d.GetOk("name"); ok {
		queryParam = fmt.Sprintf("%s&name=%v", queryParam, v)
	}

	if v, ok := d.GetOk("recordset_id"); ok {
		queryParam = fmt.Sprintf("%s&id=%v", queryParam, v)
	}

	if v, ok := d.GetOk("search_mode"); ok {
		queryParam = fmt.Sprintf("%s&search_mode=%v", queryParam, v)
	}

	if queryParam != "" {
		queryParam = "?" + queryParam[1:]
	}
	return queryParam
}
