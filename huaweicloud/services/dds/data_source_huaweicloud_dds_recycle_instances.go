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

// @API DDS GET /v3/{project_id}/recycle-instances
func DataSourceDdsRecycleInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsRecycleInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the instances.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance name.`,
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance mode.`,
						},
						"backup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the backup ID.`,
						},
						"datastore": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the database information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the database version.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the database type.`,
									},
								},
							},
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the charging mode.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the enterprise project ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"deleted_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deletion time.`,
						},
						"retained_until": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the retention end time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsRecycleInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/recycle-instances"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getPath += fmt.Sprintf("?limit=%d", pageLimit)
	currentTotal := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", currentTotal)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving DDS recycle insatnces: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}

		instances := utils.PathSearch("instances", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, instance := range instances {
			rst = append(rst, map[string]interface{}{
				"id":                    utils.PathSearch("id", instance, nil),
				"name":                  utils.PathSearch("name", instance, nil),
				"mode":                  utils.PathSearch("mode", instance, nil),
				"backup_id":             utils.PathSearch("backup_id", instance, nil),
				"datastore":             flatteRecycleInstancesResponseDatastore(instance),
				"charging_mode":         parseChargingMode(utils.PathSearch("pay_model", instance, "").(string)),
				"enterprise_project_id": utils.PathSearch("enterprise_project_id", instance, nil),
				"created_at":            utils.PathSearch("create_at", instance, nil),
				"deleted_at":            utils.PathSearch("deleted_at", instance, nil),
				"retained_until":        utils.PathSearch("retained_until", instance, nil),
			})
		}

		// `total_count` means the number of all `instances`, and type is float64.
		currentTotal += len(instances)
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0))
		if int(totalCount.(float64)) == currentTotal {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flatteRecycleInstancesResponseDatastore(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("data_store", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("type", curJson, nil),
			"version": utils.PathSearch("version", curJson, nil),
		},
	}
	return rst
}

func parseChargingMode(v string) string {
	switch v {
	case "0":
		return "postPaid"
	case "1":
		return "prePaid"
	default:
		return v
	}
}
