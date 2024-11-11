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

// @API DDS GET /v3/{project_id}/configurations/{config_id}/histories
func DataSourceDdsPtModificationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsPtModificationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the parameter template.`,
			},
			"histories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the modification records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter name.`,
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the new value.`,
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the old value.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time, in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsPtModificationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getTasksHttpUrl := "v3/{project_id}/configurations/{config_id}/histories"
	getTasksPath := client.Endpoint + getTasksHttpUrl
	getTasksPath = strings.ReplaceAll(getTasksPath, "{project_id}", client.ProjectID)
	getTasksPath = strings.ReplaceAll(getTasksPath, "{config_id}", d.Get("configuration_id").(string))
	getTasksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// pagelimit is `10`
	getTasksPath += fmt.Sprintf("?limit=%v", pageLimit)
	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getTasksPath + fmt.Sprintf("&offset=%d", currentTotal)
		getTasksResp, err := client.Request("GET", currentPath, &getTasksOpt)
		if err != nil {
			return diag.Errorf("error retrieving records: %s", err)
		}
		getTasksRespBody, err := utils.FlattenResponse(getTasksResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}

		records := utils.PathSearch("histories", getTasksRespBody, make([]interface{}, 0)).([]interface{})
		for _, record := range records {
			results = append(results, map[string]interface{}{
				"parameter_name": utils.PathSearch("parameter_name", record, nil),
				"new_value":      utils.PathSearch("new_value", record, nil),
				"old_value":      utils.PathSearch("old_value", record, nil),
				"updated_at":     utils.PathSearch("updated_at", record, nil),
			})
		}

		if len(records) < pageLimit {
			break
		}
		currentTotal += len(records)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("histories", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
