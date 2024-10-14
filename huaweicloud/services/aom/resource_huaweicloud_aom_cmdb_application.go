// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AOM
// ---------------------------------------------------------------

package aom

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	AppNotExistsCode = "AOM.30004003"
	ComNotExistsCode = "AOM.30004203"
	EnvNotExistsCode = "AOM.30004303"
)

// @API AOM POST /v1/applications
// @API AOM DELETE /v1/applications/{application_id}
// @API AOM GET /v1/applications/{application_id}
// @API AOM PUT /v1/applications/{application_id}
func ResourceCmdbApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCmdbApplicationCreate,
		ReadContext:   resourceCmdbApplicationRead,
		UpdateContext: resourceCmdbApplicationUpdate,
		DeleteContext: resourceCmdbApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// attributes
			"register_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateApplicationBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"display_name": utils.ValueIgnoreEmpty(d.Get("display_name")),
		"eps_id":       utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func resourceCmdbApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createApplicationHttpUrl := "v1/applications"
	createApplicationPath := client.Endpoint + createApplicationHttpUrl

	createApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createApplicationOpt.JSONBody = utils.RemoveNil(buildCreateApplicationBodyParams(d, cfg))
	createApplicationResp, err := client.Request("POST", createApplicationPath, &createApplicationOpt)
	if err != nil {
		return diag.Errorf("error creating CMDB application: %s", err)
	}
	createApplicationRespBody, err := utils.FlattenResponse(createApplicationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createApplicationRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find application ID from the API response")
	}

	d.SetId(id)
	return resourceCmdbApplicationRead(ctx, d, meta)
}

func resourceCmdbApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	getApplicationHttpUrl := "v1/applications/{application_id}"
	getApplicationPath := client.Endpoint + getApplicationHttpUrl
	getApplicationPath = strings.ReplaceAll(getApplicationPath, "{application_id}", d.Id())

	getApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getApplicationResp, err := client.Request("GET", getApplicationPath, &getApplicationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", AppNotExistsCode),
			"error retrieving CMDB application")
	}

	getApplicationRespBody, err := utils.FlattenResponse(getApplicationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", getApplicationRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getApplicationRespBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", getApplicationRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("eps_id", getApplicationRespBody, nil)),
		d.Set("register_type", utils.PathSearch("register_type", getApplicationRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getApplicationRespBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CMDB application fields: %s", err)
	}

	return nil
}

func buildUpdateApplicationBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"description":  d.Get("description"),
		"display_name": utils.ValueIgnoreEmpty(d.Get("display_name")),
		"eps_id":       utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func resourceCmdbApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateApplicationHttpUrl := "v1/applications/{application_id}"
	updateApplicationPath := client.Endpoint + updateApplicationHttpUrl
	updateApplicationPath = strings.ReplaceAll(updateApplicationPath, "{application_id}", d.Id())

	updateApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateApplicationOpt.JSONBody = utils.RemoveNil(buildUpdateApplicationBodyParams(d, cfg))
	_, err = client.Request("PUT", updateApplicationPath, &updateApplicationOpt)
	if err != nil {
		return diag.Errorf("error updating CMDB application: %s", err)
	}

	return resourceCmdbApplicationRead(ctx, d, meta)
}

func resourceCmdbApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteApplicationHttpUrl := "v1/applications/{application_id}"
	deleteApplicationPath := client.Endpoint + deleteApplicationHttpUrl
	deleteApplicationPath = strings.ReplaceAll(deleteApplicationPath, "{application_id}", d.Id())

	deleteApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteApplicationPath, &deleteApplicationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", AppNotExistsCode),
			"error deleting CMDB application")
	}

	return nil
}
