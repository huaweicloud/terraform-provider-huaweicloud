package dataarts

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

// @API DataArtsStudio POST /v1/{project_id}/security/permission-sets/{permission_set_id}/members
// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/{permission_set_id}/members
// @API DataArtsStudio POST /v1/{project_id}/security/permission-sets/{permission_set_id}/members/batch-delete
func ResourceSecurityPermissionSetMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPermissionSetMemberCreate,
		ReadContext:   resourceSecurityPermissionSetMemberRead,
		DeleteContext: resourceSecurityPermissionSetMemberDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityPermissionSetMemberImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecurityPermissionSetMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		workspaceId     = d.Get("workspace_id").(string)
		permissionSetId = d.Get("permission_set_id").(string)
		httpUrl         = "v1/{project_id}/security/permission-sets/{permission_set_id}/members"
		product         = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{permission_set_id}", permissionSetId)
	createPermissionSetMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceId},
		JSONBody: map[string]interface{}{
			"member_id":   d.Get("object_id"),
			"member_name": d.Get("name"),
			"member_type": d.Get("type"),
		},
	}
	resp, err := client.Request("POST", createPath, &createPermissionSetMemberOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Security permission set member: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Security permission set member: %s", err)
	}

	memberId := utils.PathSearch("member_id", respBody, "").(string)
	if memberId == "" {
		return diag.Errorf("unable to find the member ID of the DataArts Security permission set from the API response")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceId, permissionSetId, memberId))

	return resourceSecurityPermissionSetMemberRead(ctx, d, meta)
}

// GetMemberByObjectId is a method used to query the specified member information.
func GetMemberByObjectId(cfg *config.Config, region, workspaceId, permissionSetId, objectId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/security/permission-sets/{permission_set_id}/members"
	product := "dataarts"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{permission_set_id}", permissionSetId)
	getPermissionSetMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceId},
	}
	resp, err := client.Request("GET", getPath, &getPermissionSetMemberOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if member := utils.PathSearch(fmt.Sprintf("permission_set_members|[?member_id=='%s']|[0]", objectId), respBody, nil); member != nil {
		return member, nil
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/security/permission-sets/{permission_set_id}/members",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("unable to find the member using object ID (%s) in the permission set", objectId)),
		},
	}
}

func resourceSecurityPermissionSetMemberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		workspaceId     = d.Get("workspace_id").(string)
		permissionSetId = d.Get("permission_set_id").(string)
	)
	resp, err := GetMemberByObjectId(cfg, region, workspaceId, permissionSetId, d.Get("object_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DLS.6036"),
			"error retrieving DataArts Security permission set member")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("object_id", utils.PathSearch("member_id", resp, nil)),
		d.Set("name", utils.PathSearch("member_name", resp, nil)),
		d.Set("type", utils.PathSearch("member_type", resp, nil)),
		d.Set("member_id", utils.PathSearch("id", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSecurityPermissionSetMemberDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		permissionSetId = d.Get("permission_set_id").(string)
		httpUrl         = "v1/{project_id}/security/permission-sets/{permission_set_id}/members/batch-delete"
		product         = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{permission_set_id}", permissionSetId)
	deletePermissionSetMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody: map[string]interface{}{
			"ids": []string{d.Get("member_id").(string)},
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("POST", deletePath, &deletePermissionSetMemberOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Security permission set member: %s", err)
	}

	return nil
}

func resourceSecurityPermissionSetMemberImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<permission_set_id>/<member_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
		d.Set("permission_set_id", parts[1]),
		d.Set("object_id", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
