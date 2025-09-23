package dew

import (
	"context"
	"errors"
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

var aliasNonUpdatableParams = []string{"key_id", "alias"}

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

		CustomizeDiff: config.FlexibleForceNew(aliasNonUpdatableParams),

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
				Description: `Specifies the key ID used to bind the alias.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alias of the key, it can only be prefixed with **alias/**.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
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
		JSONBody:         buildCreateAliasBodyParams(d),
	}
	_, err = createAliasClient.Request("POST", createAliasPath, &createAliasOpt)
	if err != nil {
		return diag.Errorf("error creating KMS alias: %s", err)
	}

	d.SetId(fmt.Sprintf("%s?%s", d.Get("key_id"), d.Get("alias")))

	return resourceKmsAliasRead(ctx, d, meta)
}

func buildCreateAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id": d.Get("key_id"),
		"alias":  d.Get("alias"),
	}
	return bodyParams
}

func resourceKmsAliasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		getAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		getAliasProduct = "kms"
	)
	getAliasClient, err := cfg.NewServiceClient(getAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	basePath := getAliasClient.Endpoint + getAliasHttpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", getAliasClient.ProjectID)

	getAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	basePath += fmt.Sprintf("?key_id=%s&limit=50", d.Get("key_id").(string))
	for {
		getAliasPath := basePath
		if marker != "" {
			getAliasPath = basePath + fmt.Sprintf("&marker=%s", marker)
		}

		getAliasResp, err := getAliasClient.Request("GET", getAliasPath, &getAliasOpt)
		if err != nil {
			return diag.Errorf("error retrieve KMS alias response error: %s", err)
		}

		getAliasRespBody, err := utils.FlattenResponse(getAliasResp)
		if err != nil {
			return diag.Errorf("flatten KMS alias response error: %s", err)
		}
		aliases := utils.PathSearch("aliases", getAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) == 0 {
			break
		}

		allAliases = append(allAliases, aliases...)

		marker = utils.PathSearch("page_info.next_marker", getAliasRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	searchPath := fmt.Sprintf("[?alias=='%s']|[0]", d.Get("alias").(string))
	aliasDetail := utils.PathSearch(searchPath, allAliases, nil)
	if aliasDetail == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("key_id", utils.PathSearch("key_id", aliasDetail, nil)),
		d.Set("alias", utils.PathSearch("alias", aliasDetail, nil)),
		d.Set("alias_urn", utils.PathSearch("alias_urn", aliasDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", aliasDetail, nil)),
		d.Set("update_time", utils.PathSearch("update_time", aliasDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsAliasUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceKmsAliasDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		deleteAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		deleteAliasProduct = "kms"
	)
	deleteAliasClient, err := cfg.NewServiceClient(deleteAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	deleteAliasPath := deleteAliasClient.Endpoint + deleteAliasHttpUrl
	deleteAliasPath = strings.ReplaceAll(deleteAliasPath, "{project_id}", deleteAliasClient.ProjectID)

	deleteAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteAliasBodyParams(d),
	}
	_, err = deleteAliasClient.Request("DELETE", deleteAliasPath, &deleteAliasOpt)
	if err != nil {
		return diag.Errorf("error deleting KMS alias: %s", err)
	}

	return nil
}

func buildDeleteAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":  d.Get("key_id"),
		"aliases": []string{d.Get("alias").(string)},
	}
	return bodyParams
}

func resourceKmsAliasImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "?", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <key_id>?<alias>")
	}

	mErr := multierror.Append(
		d.Set("key_id", parts[0]),
		d.Set("alias", parts[1]),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import KMS alias：%s", err)
	}

	return []*schema.ResourceData{d}, nil
}
