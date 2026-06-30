package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/biz/catalogs
func DataSourceArchitectureProcesses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureProcessesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the processes are located.",
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace to which the processes belong.",
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of processes.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent ID of processes.",
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The creator of the processes.",
			},

			// Attributes.
			"processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureProcessesElem(),
				Description: "The list of processes that matched filter parameters.",
			},
		},
	}
}

func dataArchitectureProcessesElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the process, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the process.",
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English name of the process.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the process.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the process.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent ID of process.",
			},
			"prev_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The previous ID of process.",
			},
			"next_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The next ID of process.",
			},
			"qualified_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The qualified ID of process.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the process, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the process, in RFC3339 format.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the process.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last editor of the process.",
			},
		},
	}
	return &sc
}

func buildArchitectureProcessesQueryParams(d *schema.ResourceData) string {
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

	return res
}

func buildArchitectureMoreHeaders(workspaceId string) map[string]string {
	result := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		result["workspace"] = workspaceId
	}

	return result
}

func listArchitectureProcesses(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/biz/catalogs?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureProcessesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		processes := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, processes...)

		if len(processes) < limit {
			break
		}
		offset += len(processes)
	}

	return result, nil
}

func flattenArchitectureProcesses(processes []interface{}) []map[string]interface{} {
	if len(processes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(processes))
	for _, process := range processes {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", process, nil),
			"name":         utils.PathSearch("name", process, nil),
			"name_en":      utils.PathSearch("name_en", process, nil),
			"description":  utils.PathSearch("description", process, nil),
			"owner":        utils.PathSearch("owner", process, nil),
			"parent_id":    utils.PathSearch("parent_id", process, nil),
			"prev_id":      utils.PathSearch("prev_id", process, nil),
			"next_id":      utils.PathSearch("next_id", process, nil),
			"qualified_id": utils.PathSearch("qualified_id", process, nil),
			"created_at":   utils.PathSearch("create_time", process, nil),
			"updated_at":   utils.PathSearch("update_time", process, nil),
			"created_by":   utils.PathSearch("create_by", process, nil),
			"updated_by":   utils.PathSearch("update_by", process, nil),
		})
	}
	return result
}

func dataSourceArchitectureProcessesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	processes, err := listArchitectureProcesses(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture processes: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("processes", flattenArchitectureProcesses(processes)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
