package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live PUT /v1/{project_id}/guard/referer-chain
// @API Live GET /v1/{project_id}/guard/referer-chain
// @API Live DELETE /v1/{project_id}/guard/referer-chain
func ResourceRefererValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRefererValidationCreate,
		ReadContext:   resourceRefererValidationRead,
		UpdateContext: resourceRefererValidationUpdate,
		DeleteContext: resourceRefererValidationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRefererValidationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the streaming domain name to which the referer validation belongs.`,
			},
			"referer_config_empty": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies whether the referer header is included.`,
			},
			"referer_white_list": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies whether the referer is in the trustlist.`,
			},
			"referer_auth_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the domain name list.`,
			},
		},
	}
}

func resourceRefererValidationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateRefererValidation(client, d)
	if err != nil {
		return diag.Errorf("error creating Live referer validation: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceRefererValidationRead(ctx, d, meta)
}

func createOrUpdateRefererValidation(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	validationHttpUrl := "v1/{project_id}/guard/referer-chain"
	validationPath := client.Endpoint + validationHttpUrl
	validationPath = strings.ReplaceAll(validationPath, "{project_id}", client.ProjectID)

	validationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRefererValidationBodyParams(d)),
	}

	_, err := client.Request("PUT", validationPath, &validationOpt)
	return err
}

func buildRefererValidationBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"guard_switch":         "true",
		"domain":               d.Get("domain_name"),
		"referer_config_empty": d.Get("referer_config_empty"),
		"referer_white_list":   d.Get("referer_white_list"),
		"referer_auth_list":    utils.ExpandToStringList(d.Get("referer_auth_list").([]interface{})),
	}

	return params
}

func resourceRefererValidationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		getHttpUrl = "v1/{project_id}/guard/referer-chain"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%v", getPath, domainName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error retrieving referer validation")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	refererAuthList := utils.PathSearch("referer_auth_list", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(refererAuthList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "referer validation")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("referer_config_empty", utils.PathSearch("referer_config_empty", getRespBody, nil)),
		d.Set("referer_white_list", utils.PathSearch("referer_white_list", getRespBody, nil)),
		d.Set("referer_auth_list", utils.PathSearch("referer_auth_list", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRefererValidationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateRefererValidation(client, d)
	if err != nil {
		return diag.Errorf("error updating referer validation: %s", err)
	}

	return resourceRefererValidationRead(ctx, d, meta)
}

func resourceRefererValidationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deleteHttpUrl = "v1/{project_id}/guard/referer-chain"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?domain=%v", deletePath, d.Get("domain_name"))
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error deleting referer validation")
	}

	return nil
}

func resourceRefererValidationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	if importedId == "" {
		return nil, fmt.Errorf("invalid format specified for import ID, 'domain_name' is empty")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("domain_name", importedId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
