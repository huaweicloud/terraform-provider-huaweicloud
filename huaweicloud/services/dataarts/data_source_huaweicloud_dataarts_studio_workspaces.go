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

// @API DataArtsStudio GET /v1/{project_id}/workspaces/{instance_id}
func DataSourceDataArtsStudioWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataArtsStudioWorkSpaceRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bad_record_location_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_log_location_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_num": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeFloat,
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
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDataArtsStudioWorkSpaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listWorkspaceHttpUrl := "v1/{project_id}/workspaces/{instance_id}"
	listWorkspaceProduct := "dataarts"

	listWorkspaceClient, err := cfg.NewServiceClient(listWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	listWorkspacePath := listWorkspaceClient.Endpoint + listWorkspaceHttpUrl
	listWorkspacePath = strings.ReplaceAll(listWorkspacePath, "{project_id}", listWorkspaceClient.ProjectID)
	listWorkspacePath = strings.ReplaceAll(listWorkspacePath, "{instance_id}", d.Get("instance_id").(string))

	listWorkSpaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	currentTotal := 0
	listWorkspacePath += fmt.Sprintf("?limit=10&offset=%v", currentTotal)
	results := make([]map[string]interface{}, 0)

	for {
		listWorkspaceResp, err := listWorkspaceClient.Request("GET", listWorkspacePath, &listWorkSpaceOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listWorkspaceRespBody, err := utils.FlattenResponse(listWorkspaceResp)
		if err != nil {
			return diag.FromErr(err)
		}

		workspaces := utils.PathSearch("data", listWorkspaceRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("count", listWorkspaceRespBody, 0)
		for _, workspace := range workspaces {
			// filter result
			workspaceId := utils.PathSearch("id", workspace, "").(string)
			name := utils.PathSearch("name", workspace, "").(string)
			epsId := utils.PathSearch("eps_id", workspace, "").(string)
			createdBy := utils.PathSearch("create_user", workspace, "").(string)
			if val, ok := d.GetOk("workspace_id"); ok && workspaceId != val {
				continue
			}
			if val, ok := d.GetOk("name"); ok && name != val {
				continue
			}
			if val, ok := d.GetOk("created_by"); ok && createdBy != val {
				continue
			}
			if val, ok := d.GetOk("enterprise_project_id"); ok && epsId != val {
				continue
			}

			results = append(results, map[string]interface{}{
				"name":                     utils.PathSearch("name", workspace, nil),
				"id":                       utils.PathSearch("id", workspace, nil),
				"enterprise_project_id":    utils.PathSearch("eps_id", workspace, nil),
				"bad_record_location_name": utils.PathSearch("bad_record_location_name", workspace, nil),
				"job_log_location_name":    utils.PathSearch("job_log_location_name", workspace, nil),
				"description":              utils.PathSearch("description", workspace, nil),
				"created_by":               utils.PathSearch("create_user", workspace, nil),
				"updated_by":               utils.PathSearch("update_user", workspace, nil),
				"member_num":               utils.PathSearch("member_num", workspace, float64(0)),
				"is_default":               utils.PathSearch("is_default", workspace, float64(0)),
				"created_at": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("create_time", workspace, float64(0)).(float64))/1000, false),
				"updated_at": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("update_time", workspace, float64(0)).(float64))/1000, false),
			})
		}
		currentTotal += len(workspaces)
		// type of `total` is float64
		if float64(currentTotal) == total {
			break
		}
		index := strings.Index(listWorkspacePath, "offset")
		listWorkspacePath = fmt.Sprintf("%soffset=%v", listWorkspacePath[:index], currentTotal)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspaces", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
