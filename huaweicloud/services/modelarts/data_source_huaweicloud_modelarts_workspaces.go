// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/workspaces
func DataSourceWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Workspace name. Fuzzy match is supported.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID to which the workspace belongs.`,
			},
			"filter_accessible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to filter that the current user does not have permission to access.`,
			},
			"workspaces": {
				Type:        schema.TypeList,
				Elem:        workspacesWorkspacesSchema(),
				Computed:    true,
				Description: `The list of workspaces.`,
			},
		},
	}
}

func workspacesWorkspacesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Workspace ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Workspace name.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Authorization type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the worksapce.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Account name of the owner.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID to which the workspace belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Workspace status.`,
			},
			"status_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status details.`,
			},
		},
	}
	return &sc
}

func resourceWorkspacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listWorkspacesHttpUrl = "v1/{project_id}/workspaces"
		listWorkspacesProduct = "modelarts"
	)
	listWorkspacesClient, err := cfg.NewServiceClient(listWorkspacesProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	listWorkspacesPath := listWorkspacesClient.Endpoint + listWorkspacesHttpUrl
	listWorkspacesPath = strings.ReplaceAll(listWorkspacesPath, "{project_id}", listWorkspacesClient.ProjectID)

	listWorkspacesqueryParams := buildListWorkspacesQueryParams(d)
	listWorkspacesPath += listWorkspacesqueryParams

	listWorkspacesResp, err := pagination.ListAllItems(
		listWorkspacesClient,
		"offset",
		listWorkspacesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ModelArts workspaces")
	}

	listWorkspacesRespJson, err := json.Marshal(listWorkspacesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listWorkspacesRespBody interface{}
	err = json.Unmarshal(listWorkspacesRespJson, &listWorkspacesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspaces", flattenListWorkspacesWorkspaces(listWorkspacesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListWorkspacesWorkspaces(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("workspaces", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"auth_type":             utils.PathSearch("auth_type", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"owner":                 utils.PathSearch("owner", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"status_info":           utils.PathSearch("status_info", v, nil),
		})
	}
	return rst
}

func buildListWorkspacesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	if v, ok := d.GetOk("filter_accessible"); ok {
		res = fmt.Sprintf("%s&filter_accessible=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
