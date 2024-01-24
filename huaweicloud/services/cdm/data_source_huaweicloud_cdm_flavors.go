package cdm

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDM GET /v1.1/{project_id}/datastores
// @API CDM GET /v1.1/{project_id}/datastores/{id}/flavors
func DataSourceCdmFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCdmFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCdmFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "cdm"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDM client: %s", err)
	}

	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getDataStoreIDHttpUrl := "v1.1/{project_id}/datastores"
	getDataStoreIDPath := client.Endpoint + getDataStoreIDHttpUrl
	getDataStoreIDPath = strings.ReplaceAll(getDataStoreIDPath, "{project_id}", client.ProjectID)
	getDataStoreIDResp, err := client.Request("GET", getDataStoreIDPath, &requestOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataStore ID")
	}
	getDataStoreIDBody, err := utils.FlattenResponse(getDataStoreIDResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataStoreID := utils.PathSearch("datastores|[0].id", getDataStoreIDBody, "").(string)
	if dataStoreID == "" {
		return common.CheckDeletedDiag(d, err, "error retrieving DataStore ID")
	}

	getDataStoreFlavorsHttpUrl := "v1.1/{project_id}/datastores/{id}/flavors"
	getDataStoreFlavorsPath := client.Endpoint + getDataStoreFlavorsHttpUrl
	getDataStoreFlavorsPath = strings.ReplaceAll(getDataStoreFlavorsPath, "{project_id}", client.ProjectID)
	getDataStoreFlavorsPath = strings.ReplaceAll(getDataStoreFlavorsPath, "{id}", dataStoreID)
	getDataStoreFlavorsResp, err := client.Request("GET", getDataStoreFlavorsPath, &requestOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataStore Flavors")
	}
	getDataStoreFlavorsBody, err := utils.FlattenResponse(getDataStoreFlavorsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	cdmFlavor := utils.PathSearch("versions|[0]", getDataStoreFlavorsBody, nil)
	if cdmFlavor == nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataStore Flavors")
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("version", utils.PathSearch("name", cdmFlavor, nil)),
		d.Set("flavors", flattenFlavors(utils.PathSearch("flavors", cdmFlavor, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlavors(curJson interface{}) []interface{} {
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("str_id", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"cpu":    strconv.FormatFloat(utils.PathSearch("cpu", v, float64(0)).(float64), 'f', -1, 64),
			"memory": strconv.FormatFloat(utils.PathSearch("ram", v, float64(0)).(float64), 'f', -1, 64),
		})
	}
	return rst
}
