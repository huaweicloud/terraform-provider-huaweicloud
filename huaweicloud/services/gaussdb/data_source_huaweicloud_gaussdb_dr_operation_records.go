package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/disaster-recovery/records
func DataSourceGaussDbDrOperationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbDrOperationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussBbDrOperationRecordsRecordSchema(),
			},
		},
	}
}

func gaussBbDrOperationRecordsRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbDrOperationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/disaster-recovery/records"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildDrOperationRecordsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB DR operation records: %s", err)
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

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenDrOperationRecords(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDrOperationRecordsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?entity_id=%s&entity_type=%s", d.Get("entity_id").(string), d.Get("entity_type").(string))

	return res
}

func flattenDrOperationRecords(resp interface{}) []interface{} {
	records := utils.PathSearch("records", resp, make([]interface{}, 0)).([]interface{})
	if len(records) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", record, nil),
			"action":      utils.PathSearch("action", record, nil),
			"status":      utils.PathSearch("status", record, nil),
			"message":     utils.PathSearch("message", record, nil),
			"entity_id":   utils.PathSearch("entity_id", record, nil),
			"entity_type": utils.PathSearch("entity_type", record, nil),
			"job_id":      utils.PathSearch("job_id", record, nil),
			"instance_id": utils.PathSearch("instance_id", record, nil),
			"created_at":  utils.PathSearch("created_at", record, nil),
			"updated_at":  utils.PathSearch("updated_at", record, nil),
		})
	}

	return result
}
