// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SWR
// ---------------------------------------------------------------

package swr

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

// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo
// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo
// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo
func ResourceSwrImageAutoSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrImageAutoSyncCreate,
		ReadContext:   resourceSwrImageAutoSyncRead,
		DeleteContext: resourceSwrImageAutoSyncDelete,
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
			"target_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target region name.`,
			},
			"target_organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target organization name.`,
			},
			"override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to overwrite.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
}

func resourceSwrImageAutoSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSwrImageAutoSync: create SWR image auto sync
	var (
		createSwrImageAutoSyncHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo"
		createSwrImageAutoSyncProduct = "swr"
	)
	createSwrImageAutoSyncClient, err := cfg.NewServiceClient(createSwrImageAutoSyncProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	targetRegion := d.Get("target_region").(string)
	targetOrganization := d.Get("target_organization").(string)
	createSwrImageAutoSyncPath := createSwrImageAutoSyncClient.Endpoint + createSwrImageAutoSyncHttpUrl
	createSwrImageAutoSyncPath = strings.ReplaceAll(createSwrImageAutoSyncPath, "{namespace}", organization)
	createSwrImageAutoSyncPath = strings.ReplaceAll(createSwrImageAutoSyncPath, "{repository}", repository)

	createSwrImageAutoSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createSwrImageAutoSyncOpt.JSONBody = utils.RemoveNil(buildCreateSwrImageAutoSyncBodyParams(d))
	_, err = createSwrImageAutoSyncClient.Request("POST", createSwrImageAutoSyncPath,
		&createSwrImageAutoSyncOpt)
	if err != nil {
		return diag.Errorf("error creating SWR image auto sync: %s", err)
	}

	id := organization + "/" + repository + "/" + targetRegion + "/" + targetOrganization
	d.SetId(id)

	return resourceSwrImageAutoSyncRead(ctx, d, meta)
}

func buildCreateSwrImageAutoSyncBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"remoteRegionId":  utils.ValueIgnoreEmpty(d.Get("target_region")),
		"remoteNamespace": utils.ValueIgnoreEmpty(d.Get("target_organization")),
		"syncAuto":        true,
		"override":        utils.ValueIgnoreEmpty(d.Get("override")),
	}
	return bodyParams
}

func resourceSwrImageAutoSyncRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSwrImageAutoSync: Query SWR image auto sync
	var (
		getSwrImageAutoSyncHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo"
		getSwrImageAutoSyncProduct = "swr"
	)
	getSwrImageAutoSyncClient, err := cfg.NewServiceClient(getSwrImageAutoSyncProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 4)
	if len(parts) != 4 {
		return diag.Errorf("invalid id format, must be " +
			"<organization_name>/<repository_name>/<target_region>/<target_organization>")
	}
	organization := parts[0]
	repository := parts[1]
	targetRegion := parts[2]
	targetOrganization := parts[3]

	getSwrImageAutoSyncPath := getSwrImageAutoSyncClient.Endpoint + getSwrImageAutoSyncHttpUrl
	getSwrImageAutoSyncPath = strings.ReplaceAll(getSwrImageAutoSyncPath, "{namespace}", organization)
	getSwrImageAutoSyncPath = strings.ReplaceAll(getSwrImageAutoSyncPath, "{repository}", repository)

	getSwrImageAutoSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageAutoSyncResp, err := getSwrImageAutoSyncClient.Request("GET",
		getSwrImageAutoSyncPath, &getSwrImageAutoSyncOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR image auto sync")
	}

	getSwrImageAutoSyncRespBody, err := utils.FlattenResponse(getSwrImageAutoSyncResp)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, res := range getSwrImageAutoSyncRespBody.([]interface{}) {
		resTargetRegion := utils.PathSearch("remoteRegionId", res, "").(string)
		resTargetOrganization := utils.PathSearch("remoteNamespace", res, "").(string)
		if resTargetRegion == targetRegion && resTargetOrganization == targetOrganization {
			mErr = multierror.Append(
				mErr,
				d.Set("region", region),
				d.Set("organization", utils.PathSearch("namespace", res, nil)),
				d.Set("override", utils.PathSearch("override", res, nil)),
				d.Set("target_region", resTargetRegion),
				d.Set("target_organization", resTargetOrganization),
				d.Set("repository", utils.PathSearch("repoName", res, nil)),
				d.Set("created_at", utils.PathSearch("createdAt", res, nil)),
				d.Set("updated_at", utils.PathSearch("updatedAt", res, nil)),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
}

func resourceSwrImageAutoSyncDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSwrImageAutoSync: Delete SWR image auto sync
	var (
		deleteSwrImageAutoSyncHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo"
		deleteSwrImageAutoSyncProduct = "swr"
	)
	deleteSwrImageAutoSyncClient, err := cfg.NewServiceClient(deleteSwrImageAutoSyncProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteSwrImageAutoSyncPath := deleteSwrImageAutoSyncClient.Endpoint + deleteSwrImageAutoSyncHttpUrl
	deleteSwrImageAutoSyncPath = strings.ReplaceAll(deleteSwrImageAutoSyncPath, "{namespace}",
		fmt.Sprintf("%v", d.Get("organization")))
	deleteSwrImageAutoSyncPath = strings.ReplaceAll(deleteSwrImageAutoSyncPath, "{repository}",
		fmt.Sprintf("%v", d.Get("repository")))

	deleteSwrImageAutoSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteSwrImageAutoSyncOpt.JSONBody = utils.RemoveNil(buildDeleteSwrImageAutoSyncBodyParams(d))
	_, err = deleteSwrImageAutoSyncClient.Request("DELETE", deleteSwrImageAutoSyncPath,
		&deleteSwrImageAutoSyncOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].errorCode", "SVCSTG.SWR.4001158"),
			"error deleting SWR image auto sync")
	}

	return nil
}

func buildDeleteSwrImageAutoSyncBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"remoteRegionId":  utils.ValueIgnoreEmpty(d.Get("target_region")),
		"remoteNamespace": utils.ValueIgnoreEmpty(d.Get("target_organization")),
	}
	return bodyParams
}
