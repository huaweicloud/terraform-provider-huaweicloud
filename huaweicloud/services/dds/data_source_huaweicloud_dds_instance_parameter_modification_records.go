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

// @API DDS GET /v3/{project_id}/instances/{instance_id}/configuration-histories
func DataSourceDdsInstanceParameterModificationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsInstanceParameterModificationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the entity ID.`,
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
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the old value.`,
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the new value.`,
						},
						"update_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update result.`,
						},
						"applied": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the parameter is applied.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time, in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"applied_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the apply time, in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsInstanceParameterModificationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	instId := d.Get("instance_id").(string)
	entityId := d.Get("entity_id").(string)
	if entityId == "" {
		entityId = instId
	}

	getHttpUrl := "v3/{project_id}/instances/{instance_id}/configuration-histories?entity_id={entity_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instId)
	getPath = strings.ReplaceAll(getPath, "{entity_id}", entityId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	getPath += fmt.Sprintf("&limit=%v", pageLimit)
	currentTotal := 0

	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", currentTotal)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		records := utils.PathSearch("histories", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, record := range records {
			results = append(results, map[string]interface{}{
				"parameter_name": utils.PathSearch("parameter_name", record, nil),
				"old_value":      utils.PathSearch("old_value", record, nil),
				"new_value":      utils.PathSearch("new_value", record, nil),
				"update_result":  utils.PathSearch("update_result", record, nil),
				"applied":        utils.PathSearch("applied", record, nil),
				"updated_at":     utils.PathSearch("updated_at", record, nil),
				"applied_at":     utils.PathSearch("applied_at", record, nil),
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("histories", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
