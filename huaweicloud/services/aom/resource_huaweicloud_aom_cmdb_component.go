// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AOM
// ---------------------------------------------------------------

package aom

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v1/components
// @API AOM PUT /v1/components/{component_id}
// @API AOM DELETE /v1/components/{component_id}
// @API AOM GET /v1/components/{component_id}
func ResourceCmdbComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCmdbComponentCreate,
		ReadContext:   resourceCmdbComponentRead,
		UpdateContext: resourceCmdbComponentUpdate,
		DeleteContext: resourceCmdbComponentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// attributes
			"register_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_app_id": {
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

func buildCreateComponentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"model_id":    d.Get("model_id"),
		"model_type":  d.Get("model_type"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceCmdbComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createComponentHttpUrl := "v1/components"
	createComponentPath := client.Endpoint + createComponentHttpUrl

	createComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createComponentOpt.JSONBody = utils.RemoveNil(buildCreateComponentBodyParams(d))
	createComponentResp, err := client.Request("POST", createComponentPath, &createComponentOpt)
	if err != nil {
		return diag.Errorf("error creating CMDB component: %s", err)
	}
	createComponentRespBody, err := utils.FlattenResponse(createComponentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createComponentRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find component ID from the API response")
	}

	d.SetId(id)
	return resourceCmdbComponentRead(ctx, d, meta)
}

func resourceCmdbComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	getComponentHttpUrl := "v1/components/{component_id}"
	getComponentPath := client.Endpoint + getComponentHttpUrl
	getComponentPath = strings.ReplaceAll(getComponentPath, "{component_id}", d.Id())

	getComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getComponentResp, err := client.Request("GET", getComponentPath, &getComponentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ComNotExistsCode),
			"error retrieving CMDB component")
	}

	getComponentRespBody, err := utils.FlattenResponse(getComponentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var modelID, modelType string
	appID := utils.PathSearch("app_id", getComponentRespBody, "")
	subAppID := utils.PathSearch("sub_app_id", getComponentRespBody, "")

	switch {
	case subAppID != "":
		modelID = subAppID.(string)
		modelType = "SUB_APPLICATION"
	case appID != "":
		modelID = appID.(string)
		modelType = "APPLICATION"
	default:
		log.Printf("[WARN] both app_id and sub_app_id do not exist in API response")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", getComponentRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getComponentRespBody, nil)),
		d.Set("register_type", utils.PathSearch("register_type", getComponentRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getComponentRespBody, nil)),
		d.Set("app_id", appID),
		d.Set("sub_app_id", subAppID),
		d.Set("model_id", modelID),
		d.Set("model_type", modelType),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CMDB component fields: %s", err)
	}

	return nil
}

func buildUpdateComponentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceCmdbComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateComponentHttpUrl := "v1/components/{component_id}"
	updateComponentPath := client.Endpoint + updateComponentHttpUrl
	updateComponentPath = strings.ReplaceAll(updateComponentPath, "{component_id}", d.Id())

	updateComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateComponentOpt.JSONBody = utils.RemoveNil(buildUpdateComponentBodyParams(d))
	_, err = client.Request("PUT", updateComponentPath, &updateComponentOpt)
	if err != nil {
		return diag.Errorf("error updating CMDB component: %s", err)
	}

	return resourceCmdbComponentRead(ctx, d, meta)
}

func resourceCmdbComponentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteComponentHttpUrl := "v1/components/{component_id}"
	deleteComponentPath := client.Endpoint + deleteComponentHttpUrl
	deleteComponentPath = strings.ReplaceAll(deleteComponentPath, "{component_id}", d.Id())

	deleteComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteComponentPath, &deleteComponentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ComNotExistsCode),
			"error deleting CMDB component")
	}

	return nil
}
