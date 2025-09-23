// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"fmt"
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

// @API RDS GET /v3/{project_id}/storage-type/{db_type}
func DataSourceStoragetype() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceStoragetypeRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_types": {
				Type:     schema.TypeList,
				Elem:     StoragetypeStorageTypeSchema(),
				Computed: true,
			},
		},
	}
}

func StoragetypeStorageTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Storage type.`,
			},
			"az_status": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The status details of the AZs to which the specification belongs.`,
			},
			"support_compute_group_type": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Performance specifications.`,
			},
		},
	}
	return &sc
}

func resourceStoragetypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// listStorageTypes: Query the List of RDS storage-types.
	var (
		listStorageTypesHttpUrl = "v3/{project_id}/storage-type/{db_type}"
		listStorageTypesProduct = "rds"
	)
	listStorageTypesClient, err := config.NewServiceClient(listStorageTypesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Storagetype Client: %s", err)
	}

	listStorageTypesPath := listStorageTypesClient.Endpoint + listStorageTypesHttpUrl
	listStorageTypesPath = strings.ReplaceAll(listStorageTypesPath, "{project_id}", listStorageTypesClient.ProjectID)
	listStorageTypesPath = strings.ReplaceAll(listStorageTypesPath, "{db_type}",
		fmt.Sprintf("%v", d.Get("db_type")))

	listStorageTypesqueryParams := buildListStorageTypesQueryParams(d)
	listStorageTypesPath += listStorageTypesqueryParams

	listStorageTypesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	listStorageTypesResp, err := listStorageTypesClient.Request("GET", listStorageTypesPath, &listStorageTypesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Storagetype")
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
		d.Set("storage_types", flattenListStorageTypesstorageType(listStorageTypesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListStorageTypesstorageType(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("storage_type", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":                       utils.PathSearch("name", v, nil),
			"az_status":                  utils.PathSearch("az_status", v, nil),
			"support_compute_group_type": utils.PathSearch("support_compute_group_type", v, nil),
		})
	}
	return rst
}

func buildListStorageTypesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("db_version"); ok {
		res = fmt.Sprintf("%s&version_name=%v", res, v)
	}

	if v, ok := d.GetOk("instance_mode"); ok {
		res = fmt.Sprintf("%s&ha_mode=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
