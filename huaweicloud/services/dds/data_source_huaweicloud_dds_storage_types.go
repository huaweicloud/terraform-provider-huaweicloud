package dds

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

// @API DDS GET /v3/{project_id}/storage-type
func DataSourceDdsStorageTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdsStorageTypesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_types": {
				Type:     schema.TypeList,
				Elem:     ddsStorageTypesSchema(),
				Computed: true,
			},
		},
	}
}

func ddsStorageTypesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDdsStorageTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getDdsStorageTypes: Query the List of DDS storage types.
	var (
		listDdsStorageTypesHttpUrl = "v3/{project_id}/storage-type"
		listDdsStorageTypesProduct = "dds"
	)
	listDdsStorageTypesClient, err := conf.NewServiceClient(listDdsStorageTypesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	listDdsStorageTypesPath := listDdsStorageTypesClient.Endpoint + listDdsStorageTypesHttpUrl
	listDdsStorageTypesPath = strings.ReplaceAll(listDdsStorageTypesPath, "{project_id}", listDdsStorageTypesClient.ProjectID)

	listDdsStorageTypesQueryParams := buildListDdsStorageTypesQueryParams(d)
	listDdsStorageTypesPath += listDdsStorageTypesQueryParams

	listStorageTypesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listStorageTypesResp, err := listDdsStorageTypesClient.Request("GET", listDdsStorageTypesPath, &listStorageTypesOpt)

	if err != nil {
		return diag.Errorf("error retrieving DDS Storage types %s", err)
	}

	listStorageTypesRespBody, err := utils.FlattenResponse(listStorageTypesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("storage_types", flattenGetDdsStorageTypesResponseBody(listStorageTypesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDdsStorageTypesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("storage_type", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":      utils.PathSearch("name", v, nil),
			"az_status": utils.PathSearch("az_status", v, nil),
		})
	}
	return rst
}

func buildListDdsStorageTypesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("engine_name"); ok {
		res = fmt.Sprintf("%s&engine_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
