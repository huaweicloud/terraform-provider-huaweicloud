package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/app-servers/access-agent/upgrade-record
func DataSourceAppHdaUpgradeRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppHdaUpgradeRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the HDA upgrade records are located.`,
			},

			// Attributes.
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of HDA upgrade records that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the server.`,
						},
						"machine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The machine name of the server.`,
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server.`,
						},
						"server_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server group.`,
						},
						"sid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SID of the server.`,
						},
						"current_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current version of the access agent.`,
						},
						"target_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The target version of the access agent.`,
						},
						"upgrade_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The HDA upgrade status.`,
						},
						"upgrade_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The upgrade time.`,
						},
					},
				},
			},
		},
	}
}

func listAppHdaUpgradeRecords(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-servers/access-agent/upgrade-record?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}

		offset += len(items)
	}

	return result, nil
}

func flattenAppHdaUpgradeRecords(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"server_id":         utils.PathSearch("server_id", item, nil),
			"machine_name":      utils.PathSearch("machine_name", item, nil),
			"server_name":       utils.PathSearch("server_name", item, nil),
			"server_group_name": utils.PathSearch("server_group_name", item, nil),
			"sid":               utils.PathSearch("sid", item, nil),
			"current_version":   utils.PathSearch("current_version", item, nil),
			"target_version":    utils.PathSearch("target_version", item, nil),
			"upgrade_status":    utils.PathSearch("upgrade_status", item, nil),
			"upgrade_time": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("upgrade_time", item, "").(string))/1000,
				false,
			),
		}
	}

	return result
}

func dataSourceAppHdaUpgradeRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	records, err := listAppHdaUpgradeRecords(client)
	if err != nil {
		return diag.Errorf("error querying HDA upgrade records: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenAppHdaUpgradeRecords(records)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
