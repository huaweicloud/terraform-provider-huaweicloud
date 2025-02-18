package codeartsinspector

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const pageLimit = 10

// @API VSS POST /v3/{project_id}/hostscan/groups
// @API VSS GET /v3/{project_id}/hostscan/groups
// @API VSS DELETE /v3/{project_id}/hostscan/groups/{group_id}
func ResourceInspectorHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInspectorHostGroupCreate,
		ReadContext:   resourceInspectorHostGroupRead,
		DeleteContext: resourceInspectorHostGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the host group name.`,
			},
		},
	}
}

func resourceInspectorHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vss", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts inspector client: %s", err)
	}

	createHttpUrl := "v3/{project_id}/hostscan/groups"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": d.Get("name"),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating host group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("id", createRespBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable find host group ID from the API response")
	}
	d.SetId(groupId)

	return resourceInspectorHostGroupRead(ctx, d, meta)
}

func resourceInspectorHostGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vss", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts inspector client: %s", err)
	}

	group, err := GetInspectorHostGroup(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving host group")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", group, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetInspectorHostGroup(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	listGroupsHttpUrl := "v3/{project_id}/hostscan/groups"
	listGroupsPath := client.Endpoint + listGroupsHttpUrl
	listGroupsPath = strings.ReplaceAll(listGroupsPath, "{project_id}", client.ProjectID)
	listGroupsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pagelimit is `10`
	listGroupsPath += fmt.Sprintf("?limit=%v", pageLimit)

	currentTotal := 0
	for {
		currentPath := listGroupsPath + fmt.Sprintf("&offset=%d", currentTotal)
		listGroupsResp, err := client.Request("GET", currentPath, &listGroupsOpt)
		if err != nil {
			return nil, err
		}
		listGroupsRespBody, err := utils.FlattenResponse(listGroupsResp)
		if err != nil {
			return nil, err
		}

		searchPath := fmt.Sprintf("items[?id=='%s']|[0]", id)
		group := utils.PathSearch(searchPath, listGroupsRespBody, nil)
		if group != nil {
			return group, nil
		}

		groups := utils.PathSearch("items", listGroupsRespBody, make([]interface{}, 0)).([]interface{})
		currentTotal += len(groups)
		totalCount := utils.PathSearch("total", listGroupsRespBody, float64(0))
		if int(totalCount.(float64)) <= currentTotal {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceInspectorHostGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vss", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts inspector client: %s", err)
	}

	deleteHttpUrl := "v3/{project_id}/hostscan/groups/{group_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting host group")
	}

	return nil
}
