// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SWR
// ---------------------------------------------------------------

package swr

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/access
// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/access
// @API SWR PATCH /v2/manage/namespaces/{namespace}/repos/{repository}/access
// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/access
func ResourceSwrImagePermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrImagePermissionsCreate,
		UpdateContext: resourceSwrImagePermissionsUpdate,
		ReadContext:   resourceSwrImagePermissionsRead,
		DeleteContext: resourceSwrImagePermissionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrImagePermissionsImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the organization.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the repository.`,
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        ImagePermissionsUserSchema(),
				Required:    true,
				Description: `Specifies the users to access to the image (repository).`,
			},
			"self_permission": {
				Type:        schema.TypeList,
				Elem:        ImagePermissionsSelfPermissionSchema(),
				Computed:    true,
				Description: `Indicates the permission information of current user.`,
			},
		},
	}
}

func ImagePermissionsUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the existing HuaweiCloud user.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the existing HuaweiCloud user.`,
			},
			"permission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the user permission of the existing HuaweiCloud user.`,
			},
		},
	}
	return &sc
}

func ImagePermissionsSelfPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the ID of the existing HuaweiCloud user.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the name of the existing HuaweiCloud user.`,
			},
			"permission": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the user permission of the existing HuaweiCloud user.`,
			},
		},
	}
	return &sc
}

func resourceSwrImagePermissionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createSwrImagePermissionsProduct = "swr"
	)
	createSwrImagePermissionsClient, err := cfg.NewServiceClient(createSwrImagePermissionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR Client: %s", err)
	}
	err = addSwrImagePermissions(d, cfg, createSwrImagePermissionsClient, d.Get("users").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)

	d.SetId(organization + "/" + repository)

	return resourceSwrImagePermissionsRead(ctx, d, meta)
}

func resourceSwrImagePermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSwrImagePermissions: Query SWR image permissions
	var (
		getSwrImagePermissionsHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/access"
		getSwrImagePermissionsProduct = "swr"
	)
	getSwrImagePermissionsClient, err := cfg.NewServiceClient(getSwrImagePermissionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR Client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	getSwrImagePermissionsPath := getSwrImagePermissionsClient.Endpoint + getSwrImagePermissionsHttpUrl
	getSwrImagePermissionsPath = strings.ReplaceAll(getSwrImagePermissionsPath, "{namespace}",
		organization)
	getSwrImagePermissionsPath = strings.ReplaceAll(getSwrImagePermissionsPath, "{repository}",
		repository)

	getSwrImagePermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImagePermissionsResp, err := getSwrImagePermissionsClient.Request("GET",
		getSwrImagePermissionsPath, &getSwrImagePermissionsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR image permissions")
	}

	getSwrImagePermissionsRespBody, err := utils.FlattenResponse(getSwrImagePermissionsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("organization", organization),
		d.Set("repository", utils.PathSearch("name", getSwrImagePermissionsRespBody, nil)),
		d.Set("self_permission", flattenGetImagePermissionsSelfPermissionResponseBody(getSwrImagePermissionsRespBody)),
		d.Set("users", flattenGetImagePermissionsResponseBodyUser(getSwrImagePermissionsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetImagePermissionsSelfPermissionResponseBody(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("self_auth", resp, nil)
	if curJson == nil {
		return rst
	}

	permission := utils.PathSearch("auth", curJson, float64(0)).(float64)
	rst = []interface{}{
		map[string]interface{}{
			"user_id":    utils.PathSearch("user_id", curJson, nil),
			"user_name":  utils.PathSearch("user_name", curJson, nil),
			"permission": resourceSWRAuthToPermission(int(permission)),
		},
	}
	return rst
}

func flattenGetImagePermissionsResponseBodyUser(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("others_auths", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0)
	for _, v := range curArray {
		permission := utils.PathSearch("auth", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"user_id":    utils.PathSearch("user_id", v, nil),
			"user_name":  utils.PathSearch("user_name", v, nil),
			"permission": resourceSWRAuthToPermission(int(permission)),
		})
	}
	return rst
}

func resourceSwrImagePermissionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSwrImagePermissionsHasChanges := []string{
		"users",
	}

	if d.HasChanges(updateSwrImagePermissionsHasChanges...) {
		oUsersRaw, nUsersRaw := d.GetChange("users")

		var (
			updateSwrImagePermissionsProduct = "swr"
		)
		updateSwrImagePermissionsClient, err := cfg.NewServiceClient(updateSwrImagePermissionsProduct, region)
		if err != nil {
			return diag.Errorf("error creating SWR Client: %s", err)
		}

		err = deleteSwrImagePermissions(d, updateSwrImagePermissionsClient, oUsersRaw.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		addUsers := dealAddUsers(nUsersRaw.([]interface{}))

		err = addSwrImagePermissions(d, cfg, updateSwrImagePermissionsClient, addUsers)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceSwrImagePermissionsRead(ctx, d, meta)
}

// if the user_name is empty and the user_id is changed, the user_id and user_name do not match
// so it is need to delete the user_name and retrieve by the user_id
func dealAddUsers(addUsers []interface{}) []interface{} {
	res := make([]interface{}, 0, len(addUsers))
	for _, addUser := range addUsers {
		user := addUser.(map[string]interface{})
		user["user_name"] = ""
		res = append(res, user)
	}
	return res
}

func resourceSwrImagePermissionsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSwrImagePermissions: Delete SWR image permissions
	var (
		deleteSwrImagePermissionsProduct = "swr"
	)
	deleteSwrImagePermissionsClient, err := cfg.NewServiceClient(deleteSwrImagePermissionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR Client: %s", err)
	}

	err = deleteSwrImagePermissions(d, deleteSwrImagePermissionsClient, d.Get("users").([]interface{}))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR image permissions")
	}

	return nil
}

func addSwrImagePermissions(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient,
	addUsers []interface{}) error {
	// createSwrImagePermissions: create SWR image permissions
	var (
		createSwrImagePermissionsHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/access"
	)

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")
	createSwrImagePermissionsPath := client.Endpoint + createSwrImagePermissionsHttpUrl
	createSwrImagePermissionsPath = strings.ReplaceAll(createSwrImagePermissionsPath, "{namespace}",
		organization)
	createSwrImagePermissionsPath = strings.ReplaceAll(createSwrImagePermissionsPath, "{repository}",
		repository)

	createSwrImagePermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	params, err := buildSwrImagePermissionsUsersChildBody(d, cfg, addUsers)
	if err != nil {
		return err
	}
	createSwrImagePermissionsOpt.JSONBody = params
	_, err = client.Request("POST", createSwrImagePermissionsPath, &createSwrImagePermissionsOpt)
	if err != nil {
		return fmt.Errorf("error creating SWR image permissions: %s", err)
	}
	return nil
}

func deleteSwrImagePermissions(d *schema.ResourceData, client *golangsdk.ServiceClient,
	deleteUsers []interface{}) error {
	// deleteSwrImagePermissions: Delete SWR image permissions
	var (
		deleteSwrImagePermissionsHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/access"
	)

	deleteSwrImagePermissionsPath := client.Endpoint + deleteSwrImagePermissionsHttpUrl
	deleteSwrImagePermissionsPath = strings.ReplaceAll(deleteSwrImagePermissionsPath, "{namespace}",
		fmt.Sprintf("%v", d.Get("organization")))
	deleteSwrImagePermissionsPath = strings.ReplaceAll(deleteSwrImagePermissionsPath, "{repository}",
		strings.ReplaceAll(d.Get("repository").(string), "/", "$"))

	deleteSwrImagePermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	userIDs := make([]string, 0, len(deleteUsers))

	for _, deleteUser := range deleteUsers {
		user := deleteUser.(map[string]interface{})
		userIDs = append(userIDs, user["user_id"].(string))
	}
	deleteSwrImagePermissionsOpt.JSONBody = userIDs
	_, err := client.Request("DELETE", deleteSwrImagePermissionsPath,
		&deleteSwrImagePermissionsOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildSwrImagePermissionsUsersChildBody(d *schema.ResourceData, cfg *config.Config,
	rawParams []interface{}) ([]map[string]interface{}, error) {
	if len(rawParams) == 0 {
		return nil, nil
	}
	params := make([]map[string]interface{}, 0)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	for _, rawParam := range rawParams {
		raw := rawParam.(map[string]interface{})
		userId := utils.ValueIgnoreEmpty(raw["user_id"]).(string)
		userName := raw["user_name"].(string)
		if userName == "" {
			user, err := users.Get(iamClient, userId).Extract()
			if err != nil {
				return nil, fmt.Errorf("error retrieving user(%s): %s", userId, err)
			}
			userName = user.Name
		}
		param := map[string]interface{}{
			"user_id":   userId,
			"user_name": userName,
			"auth":      resourceSWRPermissionToAuth(raw["permission"].(string)),
		}
		params = append(params, param)
	}

	return params, nil
}

func resourceSwrImagePermissionsImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 2 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 2 {
			return nil, errors.New("invalid id format, must be <organization_name>/<repository_name> or " +
				"<organization_name>,<repository_name>")
		}
	} else {
		// reform ID to be separated by a slash
		id := fmt.Sprintf("%s/%s", parts[0], parts[1])
		d.SetId(id)
	}

	d.Set("organization", parts[0])
	d.Set("repository", parts[1])

	return []*schema.ResourceData{d}, nil
}
