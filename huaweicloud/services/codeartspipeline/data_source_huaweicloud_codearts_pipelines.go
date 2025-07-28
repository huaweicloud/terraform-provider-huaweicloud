package codeartspipeline

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/list
func DataSourceCodeArtsPipelines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsPipelinesRead,

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
				Description: `Specifies the CodeArts project ID.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the component ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline name.`,
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the status.`,
			},
			"is_publish": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether the pipeline is a change pipeline.`,
			},
			"creator_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the creator ID list.`,
			},
			"executor_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the executor ID list.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting field name.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting rule.`,
			},
			"group_path_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group ID path.`,
			},
			"by_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether to query by group or not.`,
			},
			"is_banned": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether the pipeline is banned.`,
			},
			"pipelines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the pipeline list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline name.`,
						},
						"component_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the component ID.`,
						},
						"is_publish": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the pipeline is a change pipeline.`,
						},
						"is_collect": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the pipeline is collected.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"manifest_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline version.`,
						},
						"latest_run": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the latest running information.`,
							Elem:        schemaPipelineRunRecords(),
						},
						"convert_sign": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the sign of converting an old version to a new version.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsPipelinesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipelines/list"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildCodeArtsPipelinesQueryParams(d, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipelines: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error retrieving pipelines: %s", err)
		}

		pipelines := utils.PathSearch("pipelines", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(pipelines) == 0 {
			break
		}
		for _, pipeline := range pipelines {
			rst = append(rst, map[string]interface{}{
				"id":               utils.PathSearch("pipeline_id", pipeline, nil),
				"name":             utils.PathSearch("name", pipeline, nil),
				"component_id":     utils.PathSearch("component_id", pipeline, nil),
				"is_publish":       utils.PathSearch("is_publish", pipeline, nil),
				"is_collect":       utils.PathSearch("is_collect", pipeline, nil),
				"create_time":      utils.PathSearch("create_time", pipeline, nil),
				"manifest_version": utils.PathSearch("manifest_version", pipeline, nil),
				"latest_run":       []interface{}{flattenPipelineRunRecords(utils.PathSearch("latest_run", pipeline, nil))},
				"convert_sign":     utils.PathSearch("convert_sign", pipeline, nil),
			})
		}

		offset += len(pipelines)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pipelines", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeArtsPipelinesQueryParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":    d.Get("project_id"),
		"component_id":  utils.ValueIgnoreEmpty(d.Get("component_id")),
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"status":        utils.ValueIgnoreEmpty(d.Get("status")),
		"creator_ids":   utils.ValueIgnoreEmpty(d.Get("creator_ids")),
		"executor_ids":  utils.ValueIgnoreEmpty(d.Get("executor_ids")),
		"start_time":    utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":      utils.ValueIgnoreEmpty(d.Get("end_time")),
		"sort_key":      utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"sort_dir":      utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"group_path_id": utils.ValueIgnoreEmpty(d.Get("group_path_id")),
		"limit":         50,
		"offset":        offset,
	}

	temp := []string{"is_publish", "by_group", "is_banned"}
	for _, k := range temp {
		if v, ok := d.GetOk(k); ok {
			if val, err := strconv.ParseBool(v.(string)); err == nil {
				bodyParams[k] = val
			} else {
				log.Printf("[DEBUG] error parsing param %s string into bool: %s", k, err)
			}
		}
	}

	return bodyParams
}
