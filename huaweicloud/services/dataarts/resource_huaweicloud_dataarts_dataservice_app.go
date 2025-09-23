// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArtsStudio
// ---------------------------------------------------------------

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

var appNotFoundErrors = []string{
	"DLM.4001", // Workspace not found
	"DLM.4063", // Application not found
}

// @API DataArtsStudio POST /v1/{project_id}/service/apps
// @API DataArtsStudio GET /v1/{project_id}/service/apps/{app_id}
// @API DataArtsStudio PUT /v1/{project_id}/service/apps/{app_id}
// @API DataArtsStudio DELETE /v1/{project_id}/service/apps/{app_id}
func ResourceDataServiceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceAppCreate,
		UpdateContext: resourceDataServiceAppUpdate,
		ReadContext:   resourceDataServiceAppRead,
		DeleteContext: resourceDataServiceAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAppImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of DLM engine.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The type of the application.`,
			},
			"app_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the application.`,
			},
			"app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secret of the application.`,
			},
		},
	}
}

func resourceDataServiceAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createApp: create an app
	var (
		createAppHttpUrl = "v1/{project_id}/service/apps"
		createAppProduct = "dataarts"
	)
	createAppClient, err := cfg.NewServiceClient(createAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createAppPath := createAppClient.Endpoint + createAppHttpUrl
	createAppPath = strings.ReplaceAll(createAppPath, "{project_id}", createAppClient.ProjectID)

	createAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	createAppOpt.JSONBody = utils.RemoveNil(buildCreateAppBodyParams(d))
	createAppResp, err := createAppClient.Request("POST", createAppPath, &createAppOpt)
	if err != nil {
		return diag.Errorf("error creating app: %s", err)
	}

	createAppRespBody, err := utils.FlattenResponse(createAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("id", createAppRespBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the DataArts DataService APP ID from the API response")
	}
	d.SetId(appId)

	return resourceDataServiceAppRead(ctx, d, meta)
}

func buildCreateAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"app_type":    utils.ValueIgnoreEmpty(d.Get("app_type")),
	}
	return bodyParams
}

func resourceDataServiceAppRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getApp: Query the app
	var (
		getAppHttpUrl = "v1/{project_id}/service/apps/{id}"
		getAppProduct = "dataarts"
	)
	getAppClient, err := cfg.NewServiceClient(getAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getAppPath := getAppClient.Endpoint + getAppHttpUrl
	getAppPath = strings.ReplaceAll(getAppPath, "{project_id}", getAppClient.ProjectID)
	getAppPath = strings.ReplaceAll(getAppPath, "{id}", d.Id())

	getAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	getAppResp, err := getAppClient.Request("GET", getAppPath, &getAppOpt)

	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", appNotFoundErrors...),
			"error retrieving app")
	}

	getAppRespBody, err := utils.FlattenResponse(getAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getAppRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getAppRespBody, nil)),
		d.Set("app_key", utils.PathSearch("app_key", getAppRespBody, nil)),
		d.Set("app_secret", utils.PathSearch("app_secret", getAppRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDataServiceAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAppChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateAppChanges...) {
		// updateApp: update the App
		var (
			updateAppHttpUrl = "v1/{project_id}/service/apps/{id}"
			updateAppProduct = "dataarts"
		)
		updateAppClient, err := cfg.NewServiceClient(updateAppProduct, region)
		if err != nil {
			return diag.Errorf("error creating DataArts Studio client: %s", err)
		}

		updateAppPath := updateAppClient.Endpoint + updateAppHttpUrl
		updateAppPath = strings.ReplaceAll(updateAppPath, "{project_id}", updateAppClient.ProjectID)
		updateAppPath = strings.ReplaceAll(updateAppPath, "{id}", d.Id())

		updateAppOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"workspace":    d.Get("workspace_id").(string),
				"dlm-type":     d.Get("dlm_type").(string),
			},
		}

		updateAppOpt.JSONBody = utils.RemoveNil(buildUpdateAppBodyParams(d))
		_, err = updateAppClient.Request("PUT", updateAppPath, &updateAppOpt)
		if err != nil {
			return diag.Errorf("error updating app: %s", err)
		}
	}
	return resourceDataServiceAppRead(ctx, d, meta)
}

func buildUpdateAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceDataServiceAppDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteApp: delete the app
	var (
		deleteAppHttpUrl = "v1/{project_id}/service/apps/{id}"
		deleteAppProduct = "dataarts"
	)
	deleteAppClient, err := cfg.NewServiceClient(deleteAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deleteAppPath := deleteAppClient.Endpoint + deleteAppHttpUrl
	deleteAppPath = strings.ReplaceAll(deleteAppPath, "{project_id}", deleteAppClient.ProjectID)
	deleteAppPath = strings.ReplaceAll(deleteAppPath, "{id}", d.Id())

	deleteAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	_, err = deleteAppClient.Request("DELETE", deleteAppPath, &deleteAppOpt)
	if err != nil {
		return diag.Errorf("error deleting app: %s", err)
	}

	return nil
}

func resourceAppImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<dlm_type>/<id>")
	}

	d.Set("workspace_id", parts[0])
	d.Set("dlm_type", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
