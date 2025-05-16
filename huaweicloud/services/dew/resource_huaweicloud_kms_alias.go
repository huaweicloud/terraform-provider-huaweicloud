// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product KMS
// ---------------------------------------------------------------

package dew

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"

	"net/url"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{"key_id", "alias"}

// @API DEW POST /v1.0/{project_id}/kms/aliases
// @API DEW DELETE /v1.0/{project_id}/kms/aliases
// @API DEW GET /v1.0/{project_id}/kms/aliases
func ResourceKmsAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsAliasCreate,
		ReadContext:   resourceKmsAliasRead,
		UpdateContext: resourceKmsAliasUpdate,
		DeleteContext: resourceKmsAliasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKmsAliasImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key ID used to bind the alias, it cannot be updated.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alias of the key, It can only be prefixed with "alias\" and cannot be updated.`,
			},
			"alias_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias resource locator.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the alias. `,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the alias.`,
			},
		},
	}
}

func resourceKmsAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAlias: create a KMS Alias.
	var (
		createAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		createAliasProduct = "kms"
	)
	createAliasClient, err := cfg.NewServiceClient(createAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createAliasPath := createAliasClient.Endpoint + createAliasHttpUrl
	createAliasPath = strings.ReplaceAll(createAliasPath, "{project_id}", createAliasClient.ProjectID)

	createAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAliasOpt.JSONBody = utils.RemoveNil(buildCreateAliasBodyParams(d, cfg))
	createAliasResp, err := createAliasClient.Request("POST", createAliasPath, &createAliasOpt)
	if err != nil {
		return diag.Errorf("error creating KMS Alias: %s", err)
	}
	createAliasRespBody, err := utils.FlattenResponse(createAliasResp)
	if err != nil {
		return diag.FromErr(err)
	}

	alias := utils.PathSearch("alias", createAliasRespBody, "").(string)
	key_id := utils.PathSearch("key_id", createAliasRespBody, "").(string)
	if alias == "" {
		return diag.Errorf("unable to find the KMS alias ID from the API response")
	}
	d.SetId(fmt.Sprintf("%s?%s", key_id, alias))

	return resourceKmsAliasRead(ctx, d, meta)
}

func buildCreateAliasBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id": d.Get("key_id"),
		"alias":  d.Get("alias"),
	}
	return bodyParams
}

func resourceKmsAliasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		getAliasProduct = "kms"
	)
	getAliasClient, err := cfg.NewServiceClient(getAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS Client: %s", err)
	}

	getAliasPath := getAliasClient.Endpoint + getAliasHttpUrl
	getAliasPath = strings.ReplaceAll(getAliasPath, "{project_id}", getAliasClient.ProjectID)

	getAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	nextMarker := ""
	for {

		// 构建查询参数
		params := url.Values{}
		params.Add("key_id", d.Get("key_id").(string))
		params.Add("limit", "50")
		if marker != "" {
			params.Add("marker", marker)
		}
		// 生成完整 URL
		getAliasPath := fmt.Sprintf("%s?%s", getAliasPath, params.Encode())

		getAliasResp, err := getAliasClient.Request("GET", getAliasPath, &getAliasOpt)

		getAliasRespBody, err := utils.FlattenResponse(getAliasResp)
		if err != nil {
			return diag.FromErr(err)
		}
		aliases := utils.PathSearch("aliases", getAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) > 0 {
			allAliases = append(allAliases, aliases...)
		}
		nextMarker = utils.PathSearch("page_info.next_marker", getAliasRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	searchPath := fmt.Sprintf("[?alias=='%s']|[0]", d.Get("alias").(string))
	aliasDetail := utils.PathSearch(searchPath, allAliases, nil)
	if aliasDetail == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "KMS alias")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("key_id", utils.PathSearch("key_id", aliasDetail, nil)),
		d.Set("alias", utils.PathSearch("alias", aliasDetail, nil)),
		d.Set("alias_urn", utils.PathSearch("alias_urn", aliasDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", aliasDetail, nil)),
		d.Set("update_time", utils.PathSearch("update_time", aliasDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsAliasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceKmsAliasDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAlias: missing operation notes
	var (
		deleteAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		deleteAliasProduct = "kms"
	)
	deleteAliasClient, err := cfg.NewServiceClient(deleteAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS Client: %s", err)
	}

	deleteAliasPath := deleteAliasClient.Endpoint + deleteAliasHttpUrl
	deleteAliasPath = strings.ReplaceAll(deleteAliasPath, "{project_id}", deleteAliasClient.ProjectID)

	deleteAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteAliasOpt.JSONBody = utils.RemoveNil(buildDeleteAliasBodyParams(d, cfg))
	_, err = deleteAliasClient.Request("DELETE", deleteAliasPath, &deleteAliasOpt)
	if err != nil {
		return diag.Errorf("error deleting KMS alias: %s", err)
	}

	return nil
}

func buildDeleteAliasBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":  d.Get("key_id"),
		"aliases": []string{d.Get("alias").(string)},
	}
	return bodyParams
}

func resourceKmsAliasImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "?", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <key_id>/<alias>")
	}

	d.Set("key_id", parts[0])
	d.Set("alias", parts[1])
	d.SetId(d.Id())

	return []*schema.ResourceData{d}, nil
}
