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

const (
	domainNameNotExistsCode = "LIVE.100011001"
)

// @API Live PUT /v1/{project_id}/guard/key-chain
// @API Live GET /v1/{project_id}/guard/key-chain
// @API Live DELETE /v1/{project_id}/guard/key-chain
func ResourceUrlValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUrlValidationCreate,
		ReadContext:   resourceUrlValidationRead,
		UpdateContext: resourceUrlValidationUpdate,
		DeleteContext: resourceUrlValidationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUrlValidationImportState,
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
				Description: `Specifies the domain name to which the URL validation belongs.`,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the URL validation key value.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signing method of the URL validation.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the timeout interval of URL validation.`,
			},
		},
	}
}

func resourceUrlValidationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateUrlValidation(client, d, domainName)
	if err != nil {
		return diag.Errorf("error creating Live URL validation: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceUrlValidationRead(ctx, d, meta)
}

func createOrUpdateUrlValidation(client *golangsdk.ServiceClient, d *schema.ResourceData, domainName string) error {
	validationHttpUrl := "v1/{project_id}/guard/key-chain"
	validationPath := client.Endpoint + validationHttpUrl
	validationPath = strings.ReplaceAll(validationPath, "{project_id}", client.ProjectID)
	validationPath = fmt.Sprintf("%s?domain=%v", validationPath, domainName)

	validationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUrlValidationBodyParams(d)),
	}

	_, err := client.Request("PUT", validationPath, &validationOpt)
	return err
}

func buildUrlValidationBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"key":       d.Get("key"),
		"auth_type": d.Get("auth_type"),
		"timeout":   d.Get("timeout"),
	}

	return params
}

func resourceUrlValidationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		getHttpUrl = "v1/{project_id}/guard/key-chain"
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
			"error retrieving URL validation")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	key := utils.PathSearch("key", getRespBody, "").(string)
	if key == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "URL validation")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("key", utils.PathSearch("key", getRespBody, nil)),
		d.Set("auth_type", utils.PathSearch("auth_type", getRespBody, nil)),
		d.Set("timeout", utils.PathSearch("timeout", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceUrlValidationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateUrlValidation(client, d, domainName)
	if err != nil {
		return diag.Errorf("error updating Live URL validation: %s", err)
	}

	return resourceUrlValidationRead(ctx, d, meta)
}

func resourceUrlValidationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deleteHttpUrl = "v1/{project_id}/guard/key-chain"
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
			"error deleting URL validation")
	}

	return nil
}

func resourceUrlValidationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
