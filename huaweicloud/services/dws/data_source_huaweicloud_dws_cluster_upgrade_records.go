package dws

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

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/upgrade-management/records
func DataSourceClusterUpgradeRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterUpgradeRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster upgrade records are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the upgrade records belong.`,
			},

			// Attributes.
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        clusterUpgradeRecordSchema(),
				Description: `The list of cluster upgrade records that matched filter parameters.`,
			},
		},
	}
}

func clusterUpgradeRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the upgrade record.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the upgrade record.`,
			},
			"record_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the upgrade record.`,
			},
			"from_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source version before the upgrade.`,
			},
			"to_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The target version after the upgrade.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the upgrade task, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the upgrade task, in RFC3339 format.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the upgrade job.`,
			},
			"failed_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason why the upgrade failed.`,
			},
		},
	}
}

func listClusterUpgradeRecords(client *golangsdk.ServiceClient, clusterID string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/upgrade-management/records?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{cluster_id}", clusterID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, records...)
		if len(records) < limit {
			break
		}
		offset += len(records)
	}

	return result, nil
}

func flattenClusterUpgradeRecords(records []interface{}) []map[string]interface{} {
	if len(records) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(records))
	for i, record := range records {
		result[i] = map[string]interface{}{
			"id":            utils.PathSearch("item_id", record, nil),
			"status":        utils.PathSearch("status", record, nil),
			"record_type":   utils.PathSearch("record_type", record, nil),
			"from_version":  utils.PathSearch("from_version", record, nil),
			"to_version":    utils.PathSearch("to_version", record, nil),
			"start_time":    utils.PathSearch("start_time", record, nil),
			"end_time":      utils.PathSearch("end_time", record, nil),
			"job_id":        utils.PathSearch("job_id", record, nil),
			"failed_reason": utils.PathSearch("failed_reason", record, nil),
		}
	}

	return result
}

func dataSourceClusterUpgradeRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterID = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	records, err := listClusterUpgradeRecords(client, clusterID)
	if err != nil {
		return diag.Errorf("error querying cluster upgrade records: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenClusterUpgradeRecords(records)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
