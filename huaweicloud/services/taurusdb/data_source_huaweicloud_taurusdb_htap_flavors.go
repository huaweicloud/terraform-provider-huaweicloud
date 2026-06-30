package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/htap/flavors/{engine_name}
func DataSourceTaurusDBHtapFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHtapFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapFlavorsSchema(),
			},
		},
	}
}

func htapFlavorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceHtapFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/htap/flavors/{engine_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{engine_name}", d.Get("engine_name").(string))
	getPath += buildGetHtapFlavorsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB HTAP flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flattenHtapFlavorsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetHtapFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("availability_zone_mode"); ok {
		res = fmt.Sprintf("%s&availability_zone_mode=%v", res, v)
	}
	if v, ok := d.GetOk("spec_code"); ok {
		res = fmt.Sprintf("%s&spec_code=%v", res, v)
	}
	if v, ok := d.GetOk("version_name"); ok {
		res = fmt.Sprintf("%s&version_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenHtapFlavorsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"type":          utils.PathSearch("type", v, nil),
			"spec_code":     utils.PathSearch("spec_code", v, nil),
			"version_name":  utils.PathSearch("version_name", v, nil),
			"instance_mode": utils.PathSearch("instance_mode", v, nil),
			"vcpus":         utils.PathSearch("vcpus", v, nil),
			"ram":           utils.PathSearch("ram", v, nil),
			"az_status":     utils.PathSearch("az_status", v, nil),
		})
	}
	return res
}
