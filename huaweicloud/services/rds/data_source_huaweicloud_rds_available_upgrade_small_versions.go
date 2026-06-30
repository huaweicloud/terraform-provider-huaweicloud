package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/datastores/{database_name}/small-version
func DataSourceAvailableUpgradeSmallVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailableUpgradeSmallVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     availableUpgradeSmallVersionsDataStoresSchema(),
			},
		},
	}
}

func availableUpgradeSmallVersionsDataStoresSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"favored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAvailableUpgradeSmallVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/datastores/{database_name}/small-version"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{database_name}", d.Get("database_name").(string))
	listPath += buildGetAvailableUpgradeSmallVersionsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving RDS small versions: %s", err)
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

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("data_stores", flattenGetAvailableUpgradeSmallDataStoresBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAvailableUpgradeSmallVersionsQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?version=%v", d.Get("version"))
}

func flattenGetAvailableUpgradeSmallDataStoresBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("data_stores", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":      utils.PathSearch("id", v, nil),
			"name":    utils.PathSearch("name", v, nil),
			"favored": utils.PathSearch("favored", v, nil),
		})
	}
	return res
}
