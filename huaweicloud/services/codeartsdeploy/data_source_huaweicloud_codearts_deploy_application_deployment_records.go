package codeartsdeploy

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

// @API CodeArtsDeploy GET /v2/{project_id}/task/{id}/history
func DataSourceCodeartsDeployApplicationDeploymentRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsDeployApplicationDeploymentRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the deployment task ID.`,
			},
			"start_date": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time.`,
			},
			"end_date": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the record list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the record ID.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time of application deployment.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time of application deployment.`,
						},
						"duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deployment duration.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application status.`,
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operator user name.`,
						},
						"release_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the deployment record sequence number.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deployment type.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsDeployApplicationDeploymentRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/task/{id}/history"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{id}", d.Get("task_id").(string))
	getPath += buildApplicationDeploymentRecordsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	// pageSize is `10`
	getPath += fmt.Sprintf("&size=%v", pageSize)
	pageIndex := 1

	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&page=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		records := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, record := range records {
			rst = append(rst, map[string]interface{}{
				"id":         utils.PathSearch("execution_id", record, nil),
				"start_time": utils.PathSearch("start_time", record, nil),
				"end_time":   utils.PathSearch("end_time", record, nil),
				"duration":   utils.PathSearch("duration", record, nil),
				"state":      utils.PathSearch("state", record, nil),
				"operator":   utils.PathSearch("operator", record, nil),
				"release_id": utils.PathSearch("release_id", record, nil),
				"type":       utils.PathSearch("type", record, nil),
			})
		}

		total := utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		if pageSize*(pageIndex-1)+len(records) >= int(total) {
			break
		}
		pageIndex++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildApplicationDeploymentRecordsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?start_date=%v", d.Get("start_date"))
	res += fmt.Sprintf("&end_date=%v", d.Get("end_date"))

	return res
}
