package geminidb

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/configurations/{config_id}/applied-histories
func DataSourceGeminiDBPtApplicationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBPtApplicationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBPtApplicationRecordsSchema(),
			},
		},
	}
}

func geminiDBPtApplicationRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applied_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apply_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGeminiDBPtApplicationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/{config_id}/applied-histories"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", d.Get("config_id").(string))

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB applied histories: %s", err)
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("histories", flattenListGeminiDBPtApplicationRecords(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGeminiDBPtApplicationRecords(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	historiesRaw := utils.PathSearch("histories", resp, nil)
	if historiesRaw == nil {
		return nil
	}

	historiesSlice, ok := historiesRaw.([]interface{})
	if !ok {
		return nil
	}

	histories := make([]map[string]interface{}, 0, len(historiesSlice))
	for _, historyRaw := range historiesSlice {
		historyMap := map[string]interface{}{
			"instance_id":    utils.PathSearch("instance_id", historyRaw, nil),
			"instance_name":  utils.PathSearch("instance_name", historyRaw, nil),
			"applied_at":     utils.PathSearch("applied_at", historyRaw, nil),
			"apply_result":   utils.PathSearch("apply_result", historyRaw, nil),
			"failure_reason": utils.PathSearch("failure_reason", historyRaw, nil),
		}
		histories = append(histories, historyMap)
	}
	return histories
}
