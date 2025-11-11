package aom

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v2/{project_id}/aom/dashboards-folder
// @API AOM PUT /v2/{project_id}/aom/dashboards-folder/{folder_id}
// @API AOM DELETE /v2/{project_id}/aom/dashboards-folder/{folder_id}
// @API AOM GET /v2/{project_id}/aom/dashboards-folder
func ResourceDashboardsFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardsFolderCreate,
		ReadContext:   resourceDashboardsFolderRead,
		UpdateContext: resourceDashboardsFolderUpdate,
		DeleteContext: resourceDashboardsFolderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"folder_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"delete_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dashboard_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_template": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDashboardsFolderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/aom/dashboards-folder"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildCreateDashboardsFolderBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM dashboards folder: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating dashboards folder response: %s", err)
	}

	id := utils.PathSearch("folder_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating AOM dashboards folder: can not found folder_id in return")
	}

	d.SetId(id)

	return resourceDashboardsFolderRead(ctx, d, meta)
}

func buildCreateDashboardsFolderBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"folder_title": d.Get("folder_title"),
	}

	return bodyParams
}

func resourceDashboardsFolderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	folder, err := getDashboardsFolder(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving dashboards folder")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("folder_title", utils.PathSearch("folder_title", folder, nil)),
		d.Set("created_by", utils.PathSearch("created_by", folder, nil)),
		d.Set("dashboard_ids", utils.PathSearch("dashboard_ids", folder, nil)),
		d.Set("is_template", utils.PathSearch("is_template", folder, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", folder, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDashboardsFolder(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/aom/dashboards-folder"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards folder: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards folder: %s", err)
	}

	jsonPath := fmt.Sprintf("[?folder_id=='%s']|[0]", d.Id())
	folder := utils.PathSearch(jsonPath, listRespBody, nil)
	if folder == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return folder, nil
}

func resourceDashboardsFolderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	if d.HasChange("folder_title") {
		updateHttpUrl := "v2/{project_id}/aom/dashboards-folder/{folder_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{folder_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
			JSONBody:         buildUpdateDashboardsFolderBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating dashboards folder: %s", err)
		}
	}

	return resourceDashboardsFolderRead(ctx, d, meta)
}

func buildUpdateDashboardsFolderBodyParams(d *schema.ResourceData) map[string]interface{} {
	oldRaw, newRaw := d.GetChange("folder_title")
	bodyParams := map[string]interface{}{
		"folder_title": newRaw,
		"old_title":    oldRaw,
	}

	return bodyParams
}

func resourceDashboardsFolderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// DELETE will return 200 even deleting a non exist folder
	_, err = getDashboardsFolder(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting dashboards folder")
	}

	deleteHttpUrl := "v2/{project_id}/aom/dashboards-folder/{folder_id}?delete_all={delete_all}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{folder_id}", d.Id())
	deletePath = strings.ReplaceAll(deletePath, "{delete_all}", strconv.FormatBool(d.Get("delete_all").(bool)))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting dashboards folder")
	}

	return nil
}
