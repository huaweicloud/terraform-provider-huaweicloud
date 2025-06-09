package cbr

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

// @API CBR GET /v3/{project_id}/storage_usage
func DataSourceStorageUsages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStorageUsagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region where the CBR storage usage is located.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the resource to filter the storage usage.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the resource to filter the storage usage.`,
			},
			"storage_usages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource.`,
						},
						"backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of backups.`,
						},
						"backup_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The size of backups in bytes.`,
						},
						"backup_size_multiaz": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The size of multi-AZ backups in bytes.`,
						},
					},
				},
			},
		},
	}
}

func buildStorageUsageQueryParams(d *schema.ResourceData, offset int) string {
	res := ""

	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&resource_id=%v", res, v)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		res = fmt.Sprintf("%s&resource_type=%v", res, v)
	}

	if offset != 0 {
		res = fmt.Sprintf("%s&offset=%v", res, offset)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func dataSourceStorageUsagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/storage_usage"
		product    = "cbr"
		totalItems []interface{}
		offset     = 0
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		eachRequestPath := requestPath + buildStorageUsageQueryParams(d, offset)
		resp, err := client.Request("GET", eachRequestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CBR storage usage: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		items := utils.PathSearch("storage_usage", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) == 0 {
			break
		}

		totalItems = append(totalItems, items...)
		offset += len(items)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("storage_usages", flattenStorageUsage(totalItems)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStorageUsage(totalItems []interface{}) []interface{} {
	if len(totalItems) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(totalItems))
	for _, v := range totalItems {
		result = append(result, map[string]interface{}{
			"resource_id":         utils.PathSearch("resource_id", v, nil),
			"resource_name":       utils.PathSearch("resource_name", v, nil),
			"resource_type":       utils.PathSearch("resource_type", v, nil),
			"backup_count":        utils.PathSearch("backup_count", v, nil),
			"backup_size":         utils.PathSearch("backup_size", v, nil),
			"backup_size_multiaz": utils.PathSearch("backup_size_multiaz", v, nil),
		})
	}

	return result
}
