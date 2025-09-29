package hss

import (
	"context"
	"fmt"
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

var ciCdConfigurationNonUpdatableParams = []string{
	"cicd_name",
	"enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/image/cicd/configurations
// @API HSS GET /v5/{project_id}/image/cicd/configurations/{cicd_id}
// @API HSS GET /v5/{project_id}/image/cicd/configurations
// @API HSS PUT /v5/{project_id}/image/cicd/configurations/{cicd_id}
// @API HSS POST /v5/{project_id}/image/cicd/configurations/batch-delete
func ResourceCiCdConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCiCdConfigurationCreate,
		ReadContext:   resourceCiCdConfigurationRead,
		UpdateContext: resourceCiCdConfigurationUpdate,
		DeleteContext: resourceCiCdConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(ciCdConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cicd_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vulnerability_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vulnerability_blocklist": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"associated_images_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCiCdConfigurationQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildCreateCiCdConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cicd_name":               d.Get("cicd_name"),
		"vulnerability_whitelist": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("vulnerability_whitelist").([]interface{}))),
		"vulnerability_blocklist": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("vulnerability_blocklist").([]interface{}))),
		"image_whitelist":         utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("image_whitelist").([]interface{}))),
	}

	return bodyParams
}

func resourceCiCdConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/image/cicd/configurations"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCiCdConfigurationQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCiCdConfigurationBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating HSS CiCd configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("cicd_id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating HSS CiCd configuration: ID is not found in API response")
	}

	d.SetId(id)

	return resourceCiCdConfigurationRead(ctx, d, meta)
}

func buildListCiCdConfigurationQueryParams(id, epsId string) string {
	req := fmt.Sprintf("?cicd_id=%v", id)
	if epsId != "" {
		return fmt.Sprintf("%s&enterprise_project_id=%v", req, epsId)
	}

	return req
}

func ListCiCdConfiguration(client *golangsdk.ServiceClient, id, epsId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/image/cicd/configurations"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListCiCdConfigurationQueryParams(id, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS CiCd configurations: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data_list[0]", respBody, nil), nil
}

func resourceCiCdConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/image/cicd/configurations/{cicd_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cicd_id}", d.Id())
	requestPath += buildCiCdConfigurationQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the error information returned by the get API is inaccurate and cannot
		// determine whether the resource exists.
		// At this point, it is necessary to call the query list API to determine whether the resource exists.
		listRespBody, listErr := ListCiCdConfiguration(client, d.Id(), epsId)
		if listErr != nil {
			return diag.Errorf("get API failed (%s), and list API check also failed (%s)", err, listErr)
		}

		if listRespBody == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
		}

		return diag.Errorf("error retrieving HSS CiCd configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cicd_name", utils.PathSearch("cicd_name", respBody, nil)),
		d.Set("vulnerability_whitelist", utils.ExpandToStringList(
			utils.PathSearch("vulnerability_whitelist", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("vulnerability_blocklist", utils.ExpandToStringList(
			utils.PathSearch("vulnerability_blocklist", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("image_whitelist", utils.ExpandToStringList(
			utils.PathSearch("image_whitelist", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("associated_images_num", utils.PathSearch("associated_images_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateCiCdConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// We can pass an empty array to call the API.
		"vulnerability_whitelist": utils.ExpandToStringList(d.Get("vulnerability_whitelist").([]interface{})),
		"vulnerability_blocklist": utils.ExpandToStringList(d.Get("vulnerability_blocklist").([]interface{})),
		"image_whitelist":         utils.ExpandToStringList(d.Get("image_whitelist").([]interface{})),
	}

	return bodyParams
}

func resourceCiCdConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/image/cicd/configurations/{cicd_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cicd_id}", d.Id())
	requestPath += buildCiCdConfigurationQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateCiCdConfigurationBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating HSS CiCd configuration: %s", err)
	}

	return resourceCiCdConfigurationRead(ctx, d, meta)
}

func buildDeleteCiCdConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"cicd_id_list": []string{d.Id()},
	}
}

func resourceCiCdConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + "v5/{project_id}/image/cicd/configurations/batch-delete"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildCiCdConfigurationQueryParams(epsId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteCiCdConfigurationBodyParams(d),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS CiCd configuration: %s", err)
	}

	return nil
}
