package drs

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

// @API DRS POST /v3/{project_id}/jobs/batch-struct-process
func DataSourceDrsBatchStructProcess() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsBatchStructProcessRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: resultsSchema(),
				},
			},
		},
	}
}

func resultsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"job_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"struct_process": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: structProcessSchema(),
			},
		},
		"error_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"error_message": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"database_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: databaseInfoSchema(),
			},
		},
	}
}

func structProcessSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"result": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resultSchema(),
			},
		},
		"create_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resultSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"src_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dst_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"start_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"end_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func databaseInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_database": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"disaster_recovery_database": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceDrsBatchStructProcessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/batch-struct-process"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	jobIds := d.Get("job_ids").([]interface{})
	jobsArray := make([]string, len(jobIds))
	for i, id := range jobIds {
		jobsArray[i] = id.(string)
	}

	requestBody := map[string]interface{}{
		"jobs": jobsArray,
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: requestBody,
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS batch struct process: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("results", flattenBatchStructProcessResults(results)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchStructProcessResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	resultList := make([]interface{}, 0, len(results))
	for _, item := range results {
		structProcess := utils.PathSearch("struct_process", item, nil)
		var structProcessList []interface{}
		if structProcess != nil {
			structProcessList = []interface{}{flattenStructProcess(structProcess)}
		}

		databaseInfo := utils.PathSearch("database_info", item, nil)
		var databaseInfoList []interface{}
		if databaseInfo != nil {
			databaseInfoList = []interface{}{flattenBatchStructDatabaseInfo(databaseInfo)}
		}

		resultList = append(resultList, map[string]interface{}{
			"job_id":         utils.PathSearch("job_id", item, nil),
			"struct_process": structProcessList,
			"error_code":     utils.PathSearch("error_code", item, nil),
			"error_message":  utils.PathSearch("error_message", item, nil),
			"database_info":  databaseInfoList,
		})
	}

	return resultList
}

func flattenStructProcess(structProcess interface{}) map[string]interface{} {
	result := utils.PathSearch("result", structProcess, make([]interface{}, 0)).([]interface{})

	resultList := make([]interface{}, 0, len(result))
	for _, item := range result {
		resultList = append(resultList, map[string]interface{}{
			"type":       utils.PathSearch("type", item, nil),
			"status":     utils.PathSearch("status", item, nil),
			"src_count":  utils.PathSearch("src_count", item, nil),
			"dst_count":  utils.PathSearch("dst_count", item, nil),
			"start_time": utils.PathSearch("start_time", item, nil),
			"end_time":   utils.PathSearch("end_time", item, nil),
		})
	}

	return map[string]interface{}{
		"result":      resultList,
		"create_time": utils.PathSearch("create_time", structProcess, nil),
	}
}

func flattenBatchStructDatabaseInfo(databaseInfo interface{}) map[string]interface{} {
	return map[string]interface{}{
		"service_database":           utils.PathSearch("service_database", databaseInfo, nil),
		"disaster_recovery_database": utils.PathSearch("disaster_recovery_database", databaseInfo, nil),
	}
}
