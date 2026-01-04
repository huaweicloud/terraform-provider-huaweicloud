package codeartsdeploy

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

// @API CodeArtsDeploy POST /v1/projects/{project_id}/applications/groups
// @API CodeArtsDeploy PUT /v1/projects/{project_id}/applications/groups/{group_id}
// @API CodeArtsDeploy DELETE /v1/projects/{project_id}/applications/groups/{group_id}
// @API CodeArtsDeploy GET /v1/projects/{project_id}/applications/groups
func ResourceDeployApplicationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployApplicationGroupCreate,
		ReadContext:   resourceDeployApplicationGroupRead,
		UpdateContext: resourceDeployApplicationGroupUpdate,
		DeleteContext: resourceDeployApplicationGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployApplicationGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application group name.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the parent application group ID.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the group path.`,
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the group sorting field.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the group creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the user who last updates the group.`,
			},
			"application_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of applications in the group.`,
			},
			"children": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the child group name list.`,
			},
		},
	}
}

func resourceDeployApplicationGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	createHttpUrl := "v1/projects/{project_id}/applications/groups"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildCreateDeployApplicationGroupBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating application group: %s", err)
	}

	// get groups list
	groups, err := getDeployApplicationGroup(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// filter by its name and parent ID
	groupId := filterGroupByNameAndParentId(groups.([]interface{}), d)
	if groupId == "" {
		return diag.Errorf("unable to find group ID from API response")
	}

	d.SetId(groupId)

	return resourceDeployApplicationGroupRead(ctx, d, meta)
}

func buildCreateDeployApplicationGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"parent_id": utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}

	return bodyParams
}

func resourceDeployApplicationGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	// get groups list
	groups, err := getDeployApplicationGroup(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving application groups")
	}

	// filter group by group ID
	group := filterGroupByGroupId(groups.([]interface{}), d.Id())
	if group == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving application group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", group, nil)),
		d.Set("project_id", utils.PathSearch("project_id", group, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", group, nil)),
		d.Set("path", utils.PathSearch("path", group, nil)),
		d.Set("ordinal", utils.PathSearch("ordinal", group, nil)),
		d.Set("created_by", utils.PathSearch("create_user_id", group, nil)),
		d.Set("updated_by", utils.PathSearch("last_update_user_id", group, nil)),
		d.Set("application_count", utils.PathSearch("count", group, nil)),
		d.Set("children", utils.PathSearch("children[*].name", group, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDeployApplicationGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v1/projects/{project_id}/applications/groups"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", d.Get("project_id").(string))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	groups := utils.PathSearch("result", listRespBody, make([]interface{}, 0)).([]interface{})
	if len(groups) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/projects/{project_id}/applications/groups",
				RequestId: "NONE",
				Body:      []byte("the application groups does not exist, the result field is empty"),
			},
		}
	}

	return groups, nil
}

func filterGroupByNameAndParentId(currentList []interface{}, d *schema.ResourceData) string {
	// search first level
	searchPath := fmt.Sprintf("[?name=='%s']|[0].id", d.Get("name").(string))
	parentId, ok := d.GetOk("parent_id")
	if ok {
		searchPath = fmt.Sprintf("[?name=='%s'&&parent_id=='%s']|[0].id", d.Get("name").(string), parentId)
	}
	groupId := utils.PathSearch(searchPath, currentList, "").(string)

	if groupId == "" && ok {
		// search next level
		children := utils.PathSearch("[*].children[]", currentList, make([]interface{}, 0)).([]interface{})
		if len(children) != 0 {
			groupId = filterGroupByNameAndParentId(children, d)
		}
	}

	return groupId
}

func filterGroupByGroupId(currentList []interface{}, id string) interface{} {
	// search first level
	searchPath := fmt.Sprintf(`[?id=='%s']|[0]`, id)
	group := utils.PathSearch(searchPath, currentList, nil)

	if group == nil {
		// search next level
		children := utils.PathSearch("[*].children[]", currentList, make([]interface{}, 0)).([]interface{})
		if len(children) != 0 {
			group = filterGroupByGroupId(children, id)
		}
	}

	return group
}

func resourceDeployApplicationGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	updateHttpUrl := "v1/projects/{project_id}/applications/groups/{group_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", d.Get("project_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: map[string]interface{}{
			"name": d.Get("name"),
		},
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating application group: %s", err)
	}

	return resourceDeployApplicationGroupRead(ctx, d, meta)
}

func resourceDeployApplicationGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	deleteHttpUrl := "v1/projects/{project_id}/applications/groups/{group_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", d.Get("project_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting application group")
	}

	return nil
}

func resourceDeployApplicationGroupImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<id>', but got '%s'",
			d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("project_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving project ID: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
