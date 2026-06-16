package geminidb

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

// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/configuration-histories
func DataSourceGeminiDBInstanceParametersHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBInstanceParametersHistoriesRead,

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
			"parameter_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBInstanceParametersHistorySchema(),
			},
		},
	}
}

func geminiDBInstanceParametersHistorySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"parameter_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"old_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applied": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applied_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGeminiDBInstanceParametersHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/configuration-histories"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGeminiDBInstanceParametersHistoriesQueryParams(d)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB instance parameters histories: %s", err)
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

	histories := flattenListGeminiDBInstanceParametersHistories(getRespBody)
	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("histories", histories),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGeminiDBInstanceParametersHistoriesQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("parameter_name"); ok {
		queryParams = fmt.Sprintf("%s&parameter_name=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenListGeminiDBInstanceParametersHistories(resp interface{}) []map[string]interface{} {
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
			"parameter_name": utils.PathSearch("parameter_name", historyRaw, nil),
			"old_value":      utils.PathSearch("old_value", historyRaw, nil),
			"new_value":      utils.PathSearch("new_value", historyRaw, nil),
			"update_result":  utils.PathSearch("update_result", historyRaw, nil),
			"applied":        utils.PathSearch("applied", historyRaw, nil),
			"updated_at":     utils.PathSearch("updated_at", historyRaw, nil),
			"applied_at":     utils.PathSearch("applied_at", historyRaw, nil),
		}
		histories = append(histories, historyMap)
	}

	return histories
}
