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

// @API TaurusDB GET /v3/{project_id}/htap/storage-type/{database}
func DataSourceTaurusDBHtapStorageTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStorageTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapStorageTypesSchema(),
			},
		},
	}
}

func htapStorageTypesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"min_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStorageTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/htap/storage-type/{database}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{database}", d.Get("database").(string))
	getPath = fmt.Sprintf("%s?version_name=%s", getPath, d.Get("version_name").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB HTAP storage types: %s", err)
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
		d.Set("storage_type", flattenTaurusDBHtapStorageTypesBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTaurusDBHtapStorageTypesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("storage_type", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":            utils.PathSearch("name", v, nil),
			"az_status":       utils.PathSearch("az_status", v, nil),
			"min_volume_size": utils.PathSearch("min_volume_size", v, nil),
		})
	}
	return res
}
