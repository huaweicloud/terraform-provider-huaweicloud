// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

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

// @API WAF GET /v1/{project_id}/waf/config/lts
// @API WAF PUT /v1/{project_id}/waf/config/lts/{ltsconfig_id}
func ResourceWAFAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAFAccessCreate,
		UpdateContext: resourceWAFAccessUpdate,
		ReadContext:   resourceWAFAccessRead,
		DeleteContext: resourceWAFAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFAccessImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lts_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log group ID.`,
			},
			"lts_attack_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log stream ID for attack logs.`,
			},
			"lts_access_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log stream ID for access logs.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID.`,
			},
		},
	}
}

func resourceWAFAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	detailRespBody, err := queryWAFAccessDetail(client, d, cfg)
	if err != nil {
		return diag.Errorf("error creating LTS access WAF logs configuration (failed to query detail): %s", err)
	}

	configId := utils.PathSearch("id", detailRespBody, "").(string)
	if configId == "" {
		return diag.Errorf("unable to find the configuration ID of the LTS access WAF logs from the API response")
	}

	if err := modifyWAFAccessConfiguration(client, d, cfg, configId); err != nil {
		return diag.Errorf("error creating LTS access WAF logs configuration: %s", err)
	}
	d.SetId(configId)

	return resourceWAFAccessRead(ctx, d, meta)
}

func modifyWAFAccessConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config, id string) error {
	modifyPath := client.Endpoint + "v1/{project_id}/waf/config/lts/{ltsconfig_id}"
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{ltsconfig_id}", id)
	modifyPath += buildAccessQueryParams(d, cfg)
	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildModifyWAFAccessBodyParams(d)),
	}

	_, err := client.Request("PUT", modifyPath, &modifyOpt)
	return err
}

func buildAccessQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildModifyWAFAccessBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enabled":   true,
		"ltsIdInfo": buildWAFAccessLTSInfoBodyParams(d),
	}
	return bodyParams
}

func buildWAFAccessLTSInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ltsGroupId":        d.Get("lts_group_id"),
		"ltsAttackStreamID": utils.ValueIgnoreEmpty(d.Get("lts_attack_stream_id")),
		"ltsAccessStreamID": utils.ValueIgnoreEmpty(d.Get("lts_access_stream_id")),
	}
}

func queryWAFAccessDetail(client *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/waf/config/lts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAccessQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func resourceWAFAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	detailRespBody, err := queryWAFAccessDetail(client, d, cfg)
	if err != nil {
		return diag.Errorf("error retrieving LTS access WAF logs configuration: %s", err)
	}

	enabled := utils.PathSearch("enabled", detailRespBody, false).(bool)
	if !enabled {
		// the WAF access is not exist
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			"error retrieving LTS access WAF logs configuration")
	}

	// update the resource ID for import operation
	id := utils.PathSearch("id", detailRespBody, "")
	if id != "" {
		d.SetId(id.(string))
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("lts_group_id", utils.PathSearch("ltsIdInfo.ltsGroupId", detailRespBody, nil)),
		d.Set("lts_attack_stream_id", utils.PathSearch("ltsIdInfo.ltsAttackStreamID",
			detailRespBody, nil)),
		d.Set("lts_access_stream_id", utils.PathSearch("ltsIdInfo.ltsAccessStreamID",
			detailRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWAFAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if err := modifyWAFAccessConfiguration(client, d, cfg, d.Id()); err != nil {
		return diag.Errorf("error updating LTS access WAF logs configuration: %s", err)
	}

	return resourceWAFAccessRead(ctx, d, meta)
}

func resourceWAFAccessDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/config/lts/{ltsconfig_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{ltsconfig_id}", d.Id())
	deletePath += buildAccessQueryParams(d, cfg)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteWAFAccessBodyParams(),
	}

	_, err = client.Request("PUT", deletePath, &deleteOpt)
	if err != nil {
		diag.Errorf("error deleting LTS access WAF logs configuration: %s", err)
	}

	return nil
}

func buildDeleteWAFAccessBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enabled":   false,
		"ltsIdInfo": make(map[string]interface{}),
	}
	return bodyParams
}

func resourceWAFAccessImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	var mErr *multierror.Error
	importId := d.Id()
	// `0` means default enterprise project.
	// The non-default enterprise project ID format is a 32-bit UUID with hyphens.
	// The resource ID format is a 32-bit UUID without hyphens.
	if utils.IsUUIDWithHyphens(importId) || importId == "0" {
		mErr = multierror.Append(d.Set("enterprise_project_id", d.Id()))
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
