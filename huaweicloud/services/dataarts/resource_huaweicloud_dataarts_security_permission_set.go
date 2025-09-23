package dataarts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/security/permission-sets
// @API DataArtsStudio DELETE /v1/{project_id}/security/permission-sets/{permission_set_id}
// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/{permission_set_id}
// @API DataArtsStudio PUT /v1/{project_id}/security/permission-sets/{permission_set_id}
func ResourceSecurityPermissionSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPermissionSetCreate,
		ReadContext:   resourceSecurityPermissionSetRead,
		UpdateContext: resourceSecurityPermissionSetUpdate,
		DeleteContext: resourceSecurityPermissionSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityPermissionSetImportState,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manager_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"manager_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_type": {
				Type:     schema.TypeString,
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
	}
}

func resourceSecurityPermissionSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	// createDataArtsStudioPermissionSet: create a permission set.
	var (
		createPermissionSetUrl     = "v1/{project_id}/security/permission-sets"
		createPermissionSetProduct = "dataarts"
	)

	createPermissionSetClient, err := conf.NewServiceClient(createPermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	createPermissionSetPath := createPermissionSetClient.Endpoint + createPermissionSetUrl
	createPermissionSetPath = strings.ReplaceAll(createPermissionSetPath, "{project_id}", createPermissionSetClient.ProjectID)

	createPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	createPermissionSetOpt.JSONBody = utils.RemoveNil(buildCreatePermissionSetBodyParams(d))
	createPermissionSetResp, err := createPermissionSetClient.Request("POST", createPermissionSetPath, &createPermissionSetOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Security permissions set: %s", err)
	}

	createPermissionSetRespBody, err := utils.FlattenResponse(createPermissionSetResp)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Security permission set: %s", err)
	}

	setId := utils.PathSearch("id", createPermissionSetRespBody, "").(string)
	if setId == "" {
		return diag.Errorf("unable to find the DataArts Security permission set ID from the API response")
	}

	d.SetId(setId)
	return resourceSecurityPermissionSetRead(ctx, d, meta)
}

func buildCreatePermissionSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"parent_id":    d.Get("parent_id"),
		"manager_id":   d.Get("manager_id"),
		"manager_name": utils.ValueIgnoreEmpty(d.Get("manager_name")),
		"manager_type": utils.ValueIgnoreEmpty(d.Get("manager_type")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceSecurityPermissionSetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	// getDataArtsPermissionSet: Query the DataArts permission set detail.
	var (
		getPermissionSetHttpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}"
		getPermissionSetProduct = "dataarts"
	)

	getPermissionSetClient, err := conf.NewServiceClient(getPermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPermissionSetPath := getPermissionSetClient.Endpoint + getPermissionSetHttpUrl
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{project_id}", getPermissionSetClient.ProjectID)
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{permission_set_id}", d.Id())

	getPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getPermissionSetResp, err := getPermissionSetClient.Request("GET", getPermissionSetPath, &getPermissionSetOpt)
	if err != nil {
		if hasErrorCode(err, "DLS.3027") {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "error retrieving DataArts Security permission set")
	}

	getPermissionSetRespBody, err := utils.FlattenResponse(getPermissionSetResp)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Security permission set: %s", err)
	}

	createAt := utils.PathSearch("create_time", getPermissionSetRespBody, float64(0)).(float64)
	updateAt := utils.PathSearch("update_time", getPermissionSetRespBody, float64(0)).(float64)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("name", utils.PathSearch("name", getPermissionSetRespBody, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", getPermissionSetRespBody, nil)),
		d.Set("manager_id", utils.PathSearch("manager_id", getPermissionSetRespBody, nil)),
		d.Set("manager_name", utils.PathSearch("manager_name", getPermissionSetRespBody, nil)),
		d.Set("manager_type", utils.PathSearch("manager_type", getPermissionSetRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getPermissionSetRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getPermissionSetRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("instance_id", getPermissionSetRespBody, nil)),
		d.Set("datasource_type", utils.PathSearch("datasource_type", getPermissionSetRespBody, nil)),
		d.Set("created_by", utils.PathSearch("create_user", getPermissionSetRespBody, nil)),
		d.Set("updated_by", utils.PathSearch("update_user", getPermissionSetRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createAt), false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(updateAt), false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSecurityPermissionSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// updateDataArtsPermissionSet: Update the DataArts permission set detail.
	var (
		updatePermissionSetHttpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}"
		updatePermissionSetProduct = "dataarts"
	)

	updatePermissionSetClient, err := conf.NewServiceClient(updatePermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePermissionSetPath := updatePermissionSetClient.Endpoint + updatePermissionSetHttpUrl
	updatePermissionSetPath = strings.ReplaceAll(updatePermissionSetPath, "{project_id}", updatePermissionSetClient.ProjectID)
	updatePermissionSetPath = strings.ReplaceAll(updatePermissionSetPath, "{permission_set_id}", d.Id())

	updatePermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	updatePermissionSetOpt.JSONBody = utils.RemoveNil(buildUpdatePermissionSetBodyParams(d))
	_, err = updatePermissionSetClient.Request("PUT", updatePermissionSetPath, &updatePermissionSetOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts Security permission set: %s", err)
	}

	return resourceSecurityPermissionSetRead(ctx, d, meta)
}

func buildUpdatePermissionSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"manager_id":   d.Get("manager_id"),
		"manager_name": utils.ValueIgnoreEmpty(d.Get("manager_name")),
		"manager_type": utils.ValueIgnoreEmpty(d.Get("manager_type")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceSecurityPermissionSetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// delateDataArtsPermissionSet: Delete the DataArts permission set detail.
	var (
		deletePermissionSetHttpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}"
		deletePermissionSetProduct = "dataarts"
	)

	deletePermissionSetClient, err := conf.NewServiceClient(deletePermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePermissionSetPath := deletePermissionSetClient.Endpoint + deletePermissionSetHttpUrl
	deletePermissionSetPath = strings.ReplaceAll(deletePermissionSetPath, "{project_id}", deletePermissionSetClient.ProjectID)
	deletePermissionSetPath = strings.ReplaceAll(deletePermissionSetPath, "{permission_set_id}", d.Id())

	deletePermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	_, err = deletePermissionSetClient.Request("DELETE", deletePermissionSetPath, &deletePermissionSetOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Security permission set: %s", err)
	}

	return nil
}

func resourceSecurityPermissionSetImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<id>")
	}

	d.Set("workspace_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("error_code", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse error_code from response body: %s", parseErr)
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}
