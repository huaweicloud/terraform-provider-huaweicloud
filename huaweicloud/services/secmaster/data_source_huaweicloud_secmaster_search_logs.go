package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/search/logs
func DataSourceSearchLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSearchLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"sort": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "desc",
			},
			"analysis_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     analysisResultsSchema(),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     searchResultSchema(),
			},
		},
	}
}

func analysisResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     analysisFieldSchema(),
			},
			"datarows": {
				// Convert field `datarows` to JSON string.
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func analysisFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func searchResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_source": {
				// Convert field `data_source` to JSON string.
				Type:     schema.TypeString,
				Computed: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSearchLogsBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dataspace_id": d.Get("dataspace_id"),
		"pipe_id":      d.Get("pipe_id"),
		"query":        d.Get("query"),
		"from":         d.Get("from"),
		"to":           d.Get("to"),
		"sort":         d.Get("sort"),
		"limit":        limit,
		"offset":       offset,
	}

	return bodyParams
}

func dataSourceSearchLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/search/logs"
		limit       = 500
		offset      = 0
		mErr        *multierror.Error
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	reqOpt.JSONBody = utils.RemoveNil(buildSearchLogsBodyParams(d, limit, offset))

	resp, err := client.Request("POST", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster search logs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("analysis_results", flattenAnalysisResults(utils.PathSearch("analysis_results", respBody, nil))),
		d.Set("results", flattenSearchResults(utils.PathSearch("results", respBody, make([]interface{}, 0)))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAnalysisResults(analysisResults interface{}) []interface{} {
	if analysisResults == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"size":     utils.PathSearch("size", analysisResults, nil),
			"total":    utils.PathSearch("total", analysisResults, nil),
			"schema":   flattenAnalysisFields(utils.PathSearch("schema", analysisResults, make([]interface{}, 0))),
			"datarows": utils.JsonToString(utils.PathSearch("datarows", analysisResults, nil)),
		},
	}
}

func flattenAnalysisFields(schemaList interface{}) []interface{} {
	if schemaList == nil {
		return nil
	}

	list, ok := schemaList.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", item, nil),
			"type":  utils.PathSearch("type", item, nil),
			"alias": utils.PathSearch("alias", item, nil),
		})
	}

	return result
}

func flattenSearchResults(results interface{}) []interface{} {
	if results == nil {
		return nil
	}

	list, ok := results.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		result = append(result, map[string]interface{}{
			"data_source": utils.JsonToString(utils.PathSearch("data_source", item, nil)),
			"timestamp":   utils.PathSearch("timestamp", item, nil),
		})
	}

	return result
}
