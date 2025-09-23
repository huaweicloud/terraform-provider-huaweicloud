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

var aliasAssociateNonUpdatableParams = []string{"target_key_id", "alias"}

// @API DEW POST /v1.0/{project_id}/kms/alias/associate
// @API DEW DELETE /v1.0/{project_id}/kms/aliases
// @API DEW GET /v1.0/{project_id}/kms/aliases
func ResourceKmsAliasAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsAliasAssociateCreate,
		ReadContext:   resourceKmsAliasAssociateRead,
		UpdateContext: resourceKmsAliasAssociateUpdate,
		DeleteContext: resourceKmsAliasAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKmsAliasAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(aliasAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"target_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the target key ID used to bind the alias.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alias used to bind the key, it can only be prefixed with **alias/**.`,
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

func resourceKmsAliasAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                         = meta.(*config.Config)
		region                      = cfg.GetRegion(d)
		createAssociateAliasHttpUrl = "v1.0/{project_id}/kms/alias/associate"
		createAssociateAliasProduct = "kms"
	)
	createAssociateAliasClient, err := cfg.NewServiceClient(createAssociateAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createAssociateAliasPath := createAssociateAliasClient.Endpoint + createAssociateAliasHttpUrl
	createAssociateAliasPath = strings.ReplaceAll(createAssociateAliasPath, "{project_id}", createAssociateAliasClient.ProjectID)

	createAssociateAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAssociateAliasBodyParams(d),
	}
	_, err = createAssociateAliasClient.Request("POST", createAssociateAliasPath, &createAssociateAliasOpt)
	if err != nil {
		return diag.Errorf("error creating KMS associate alias: %s", err)
	}

	d.SetId(fmt.Sprintf("%s?%s", d.Get("target_key_id"), d.Get("alias")))

	return resourceKmsAliasAssociateRead(ctx, d, meta)
}

func buildCreateAssociateAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_key_id": d.Get("target_key_id"),
		"alias":         d.Get("alias"),
	}
	return bodyParams
}

func resourceKmsAliasAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		getAssociateAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		getAssociateAliasProduct = "kms"
	)
	getAssociateAliasClient, err := cfg.NewServiceClient(getAssociateAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	basePath := getAssociateAliasClient.Endpoint + getAssociateAliasHttpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", getAssociateAliasClient.ProjectID)

	getAssociateAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	basePath += fmt.Sprintf("?key_id=%s&limit=50", d.Get("target_key_id").(string))
	for {
		getAssociateAliasPath := basePath
		if marker != "" {
			getAssociateAliasPath = basePath + fmt.Sprintf("&marker=%s", marker)
		}

		getAssociateAliasResp, err := getAssociateAliasClient.Request("GET", getAssociateAliasPath, &getAssociateAliasOpt)
		if err != nil {
			return diag.Errorf("error retrieving KMS associate alias: %s", err)
		}

		getAssociateAliasRespBody, err := utils.FlattenResponse(getAssociateAliasResp)
		if err != nil {
			return diag.Errorf("flatten KMS associate alias response error: %s", err)
		}
		aliases := utils.PathSearch("aliases", getAssociateAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) == 0 {
			break
		}

		allAliases = append(allAliases, aliases...)

		marker = utils.PathSearch("page_info.next_marker", getAssociateAliasRespBody, "").(string)
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
		d.Set("target_key_id", utils.PathSearch("key_id", aliasDetail, nil)),
		d.Set("alias", utils.PathSearch("alias", aliasDetail, nil)),
		d.Set("alias_urn", utils.PathSearch("alias_urn", aliasDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", aliasDetail, nil)),
		d.Set("update_time", utils.PathSearch("update_time", aliasDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsAliasAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceKmsAliasAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                         = meta.(*config.Config)
		region                      = cfg.GetRegion(d)
		deleteAssociateAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		deleteAssociateAliasProduct = "kms"
	)
	deleteAssociateAliasClient, err := cfg.NewServiceClient(deleteAssociateAliasProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	deleteAssociateAliasPath := deleteAssociateAliasClient.Endpoint + deleteAssociateAliasHttpUrl
	deleteAssociateAliasPath = strings.ReplaceAll(deleteAssociateAliasPath, "{project_id}", deleteAssociateAliasClient.ProjectID)

	deleteAssociateAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteAssociateAliasBodyParams(d),
	}
	_, err = deleteAssociateAliasClient.Request("DELETE", deleteAssociateAliasPath, &deleteAssociateAliasOpt)
	if err != nil {
		return diag.Errorf("error deleting KMS associate alias: %s", err)
	}

	return nil
}

func buildDeleteAssociateAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":  d.Get("target_key_id"),
		"aliases": []string{d.Get("alias").(string)},
	}
	return bodyParams
}

func resourceKmsAliasAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "?", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <target_key_id>?<alias>")
	}

	mErr := multierror.Append(
		d.Set("target_key_id", parts[0]),
		d.Set("alias", parts[1]),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import KMS associate alias：%s", err)
	}

	return []*schema.ResourceData{d}, nil
}
