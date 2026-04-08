package dcs

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/offline/key-analysis
func DataSourceOfflineKeyAnalyses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOfflineKeyAnalysesRead,

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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     offlineKeyAnalysesRecordSchema(),
			},
		},
	}
}

func offlineKeyAnalysesRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOfflineKeyAnalysesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/offline/key-analysis"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildListOfflineKeyAnalysesQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DCS offline key analyses: %s", err)
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

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenOfflineKeyAnalysesRecords(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListOfflineKeyAnalysesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("status"); ok {
		res = res + "&status=" + v.(string)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenOfflineKeyAnalysesRecords(resp interface{}) []map[string]interface{} {
	curArray := utils.PathSearch("records", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, len(curArray))

	for i, item := range curArray {
		result[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", item, nil),
			"status":      utils.PathSearch("status", item, nil),
			"created_at":  utils.PathSearch("created_at", item, nil),
			"started_at":  utils.PathSearch("started_at", item, nil),
			"finished_at": utils.PathSearch("finished_at", item, nil),
		}
	}

	return result
}
