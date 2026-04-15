package dataarts

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const studioWorkspaceUserTypeUser = 0

var studioWorkspaceUserNonUpdatableParams = []string{
	"workspace_id",
	"user_id",
}

// @API IAM GET /v3.0/OS-USER/users/{user_id}
// @API DataArtsStudio POST /v2/{project_id}/{workspace_id}/users
// @API DataArtsStudio GET /v2/{project_id}/{workspace_id}/users
// @API DataArtsStudio PUT /v2/{project_id}/{workspace_id}/users/{user_id}
// @API DataArtsStudio POST /v2/{project_id}/{workspace_id}/delete-users
func ResourceStudioWorkspaceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudioWorkspaceUserCreate,
		ReadContext:   resourceStudioWorkspaceUserRead,
		UpdateContext: resourceStudioWorkspaceUserUpdate,
		DeleteContext: resourceStudioWorkspaceUserDelete,

		CustomizeDiff: config.FlexibleForceNew(studioWorkspaceUserNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceStudioWorkspaceUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workspace user is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the user belongs.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM user to which the workspace user correspond.`,
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The role ID.`,
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role code.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role description.`,
						},
					},
				},
				Description: `The role list of the workspace user.`,
			},

			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workspace user, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the workspace user, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildStudioWorkspaceUserRoleIds(roleIds []interface{}) []map[string]interface{} {
	if len(roleIds) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(roleIds))
	for _, v := range roleIds {
		result = append(result, map[string]interface{}{
			"role_id": utils.PathSearch("id", v, nil),
		})
	}
	return result
}

func getIdentityUserById(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	httpUrl := "v3.0/OS-USER/users/{user_id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{user_id}", userId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", path, &opts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

// The client is the IAM client, before building the body params, we need to get the user detail using IAM client.
func buildStudioWorkspaceUserBodyParams(client *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]interface{}, error) {
	userId := d.Get("user_id").(string)
	user, err := getIdentityUserById(client, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":      studioWorkspaceUserTypeUser,
		"roles_ids": buildStudioWorkspaceUserRoleIds(d.Get("roles").(*schema.Set).List()),
		"user_ids": []map[string]interface{}{
			{
				"user_id":         userId,
				"user_name":       utils.PathSearch("user.name", user, "").(string),
				"domain_id":       utils.PathSearch("user.domain_id", user, "").(string),
				"is_domain_owner": utils.PathSearch("user.is_domain_owner", user, false).(bool),
			},
		},
	}, nil
}

func createStudioWorkspaceUser(dataartsClient, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/{workspace_id}/users"
	path := dataartsClient.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", dataartsClient.ProjectID)
	path = strings.ReplaceAll(path, "{workspace_id}", d.Get("workspace_id").(string))

	bodyParams, err := buildStudioWorkspaceUserBodyParams(iamClient, d)
	if err != nil {
		return err
	}
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(bodyParams),
	}

	requestResp, err := dataartsClient.Request("POST", path, &opts)
	if err != nil {
		return err
	}
	resBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	if !utils.PathSearch("is_success", resBody, false).(bool) {
		return errors.New(utils.PathSearch("message", resBody, "").(string))
	}
	return nil
}

func resourceStudioWorkspaceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		userId = d.Get("user_id").(string)
	)

	dataartsClient, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	iamClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if err := createStudioWorkspaceUser(dataartsClient, iamClient, d); err != nil {
		return diag.Errorf("error adding workspace user (%s): %s", userId, err)
	}

	d.SetId(userId)
	return resourceStudioWorkspaceUserRead(ctx, d, meta)
}

func listStudioWorkspaceUsers(client *golangsdk.ServiceClient, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/{workspace_id}/users?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opts)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		users := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		offset += len(users)
	}

	return result, nil
}

func GetStudioWorkspaceUserById(client *golangsdk.ServiceClient, workspaceId, userId string) (interface{}, error) {
	users, err := listStudioWorkspaceUsers(client, workspaceId)
	if err != nil {
		return nil, err
	}

	user := utils.PathSearch(fmt.Sprintf("[?user_id=='%s']|[0]", userId), users, nil)
	if user == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/{workspace_id}/users",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the workspace user (%s) does not exist", userId)),
			},
		}
	}
	return user, nil
}

func resourceStudioWorkspaceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		userId      = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	user, err := GetStudioWorkspaceUserById(client, workspaceId, userId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving workspace user")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("roles", flattenStudioWorkspaceUserRoles(utils.PathSearch("roles", user, make([]interface{}, 0)).([]interface{}))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", user, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", user, float64(0)).(float64))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func updateStudioWorkspaceUser(dataartsClient, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v2/{project_id}/{workspace_id}/users/{user_id}"
		workspaceId = d.Get("workspace_id").(string)
		userId      = d.Get("user_id").(string)
	)

	path := dataartsClient.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", dataartsClient.ProjectID)
	path = strings.ReplaceAll(path, "{workspace_id}", workspaceId)
	path = strings.ReplaceAll(path, "{user_id}", userId)

	bodyParams, err := buildStudioWorkspaceUserBodyParams(iamClient, d)
	if err != nil {
		return err
	}
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(bodyParams),
	}

	requestResp, err := dataartsClient.Request("PUT", path, &opts)
	if err != nil {
		return err
	}
	resBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", resBody, false).(bool) {
		return errors.New(utils.PathSearch("message", resBody, "").(string))
	}
	return nil
}

func resourceStudioWorkspaceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		userId = d.Id()
	)

	dataartsClient, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	iamClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChange("roles") {
		if err := updateStudioWorkspaceUser(dataartsClient, iamClient, d); err != nil {
			return diag.Errorf("error updating workspace user (%s): %s", userId, err)
		}
	}

	return resourceStudioWorkspaceUserRead(ctx, d, meta)
}

func deleteStudioWorkspaceUser(dataartsClient *golangsdk.ServiceClient, workspaceId, userId string) error {
	httpUrl := "v2/{project_id}/{workspace_id}/delete-users"
	path := dataartsClient.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", dataartsClient.ProjectID)
	path = strings.ReplaceAll(path, "{workspace_id}", workspaceId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"user_ids": []interface{}{userId},
		},
	}

	requestResp, err := dataartsClient.Request("POST", path, &opts)
	if err != nil {
		return err
	}
	resBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", resBody, false).(bool) {
		return golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/{workspace_id}/users",
				RequestId: "NONE",
				Body:      []byte(utils.PathSearch("message", resBody, "").(string)),
			},
		}
	}
	return nil
}

func resourceStudioWorkspaceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		userId      = d.Get("user_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := deleteStudioWorkspaceUser(client, workspaceId, userId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting workspace user (%s)", userId))
	}
	return nil
}

func resourceStudioWorkspaceUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<user_id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
		d.Set("user_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
