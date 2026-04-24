package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/biz/catalogs
func DataSourceArchitectureProcess() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureProcessRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the process architectures are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace where the process architectures are located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name or code of the process architecture for fuzzy matching.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parent directory ID.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the process architecture.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The owner of the process architecture.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The left boundary of time filtering, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The right boundary of time filtering, in RFC3339 format.`,
			},

			// Attributes.
			"processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        processSchema(),
				Description: `The list of process architectures that match the filter parameters.`,
			},
		},
	}
}

func processSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the process architecture.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the process architecture.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the process architecture.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the process architecture.`,
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset ID corresponding to the process architecture.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The owner of the process architecture.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent directory ID of the process architecture.`,
			},
			"prev_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The previous node ID of the process architecture.`,
			},
			"next_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The next node ID of the process architecture.`,
			},
			"qualified_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication ID of the process architecture, automatically generated.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the process architecture.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the process architecture.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the process architecture, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the process architecture, in RFC3339 format.`,
			},
			"bizmetric_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of business metrics owned by the process architecture.`,
			},
			"children_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of child processes owned by the process architecture.`,
			},
			"children": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The list of child directories under the process architecture, in JSON format.`,
			},
		},
	}
}

func buildArchitectureProcessQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}
	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}
	if v, ok := d.GetOk("owner"); ok {
		res = fmt.Sprintf("%s&owner=%v", res, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	return res
}

func listArchitectureProcesses(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl     = "v2/{project_id}/design/biz/catalogs?limit={limit}"
		limit       = 100
		offset      = 0
		workspaceId = d.Get("workspace_id").(string)
		result      = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", fmt.Sprintf("%d", limit))
	listPathWithLimit += buildArchitectureProcessQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving DataArts process architectures: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		processes := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(processes) < 1 {
			break
		}
		result = append(result, processes...)
		offset += len(processes)
	}

	return result, nil
}

func flattenProcesses(processes []interface{}) []interface{} {
	if len(processes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(processes))
	for _, process := range processes {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", process, nil),
			"name":          utils.PathSearch("name", process, nil),
			"name_en":       utils.PathSearch("name_en", process, nil),
			"description":   utils.PathSearch("description", process, nil),
			"guid":          utils.PathSearch("guid", process, nil),
			"owner":         utils.PathSearch("owner", process, nil),
			"parent_id":     utils.PathSearch("parent_id", process, nil),
			"prev_id":       utils.PathSearch("prev_id", process, nil),
			"next_id":       utils.PathSearch("next_id", process, nil),
			"qualified_id":  utils.PathSearch("qualified_id", process, nil),
			"create_by":     utils.PathSearch("create_by", process, nil),
			"update_by":     utils.PathSearch("update_by", process, nil),
			"create_time":   utils.PathSearch("create_time", process, nil),
			"update_time":   utils.PathSearch("update_time", process, nil),
			"bizmetric_num": utils.PathSearch("bizmetric_num", process, 0),
			"children_num":  utils.PathSearch("children_num", process, 0),
			"children":      utils.JsonToString(utils.PathSearch("children", process, nil)),
		})
	}

	return result
}

func dataSourceArchitectureProcessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	processes, err := listArchitectureProcesses(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("processes", flattenProcesses(processes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
