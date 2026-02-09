package ccm

import (
	"context"
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

// @API CCM POST /v3/scm/csr
// @API CCM GET /v3/scm/csr/{id}
// @API CCM PUT /v3/scm/csr/{id}
// @API CCM DELETE /v3/scm/csr/{id}
func ResourceCsr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCsrCreate,
		UpdateContext: resourceCsrUpdate,
		ReadContext:   resourceCsrRead,
		DeleteContext: resourceCsrDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"domain_name",
			"private_key_algo",
			"usage",
			"sans",
			"company_country",
			"company_province",
			"company_city",
			"company_name",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key_algo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"usage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sans": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"company_country": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"company_province": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"company_city": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"company_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"csr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateCsrBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"domain_name":      d.Get("domain_name"),
		"private_key_algo": d.Get("private_key_algo"),
		"usage":            d.Get("usage"),
		"sans":             utils.ValueIgnoreEmpty(d.Get("sans")),
		"company_country":  utils.ValueIgnoreEmpty(d.Get("company_country")),
		"company_province": utils.ValueIgnoreEmpty(d.Get("company_province")),
		"company_city":     utils.ValueIgnoreEmpty(d.Get("company_city")),
		"company_name":     utils.ValueIgnoreEmpty(d.Get("company_name")),
	}
}

func resourceCsrCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "scm"
		httpUrl = "v3/scm/csr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCsrBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CCM SSL CSR: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CCM SSL CSR ID from the API response")
	}
	d.SetId(id)

	return resourceCsrRead(ctx, d, meta)
}

func GetCsr(client *golangsdk.ServiceClient, csrId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/scm/csr/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{id}", csrId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCsrRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "scm"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	respBody, err := GetCsr(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SCM.0707"),
			"error retrieving CCM SSL CSR",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("csr", utils.PathSearch("csr", respBody, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", respBody, nil)),
		d.Set("sans", utils.PathSearch("sans", respBody, nil)),
		d.Set("private_key_algo", utils.PathSearch("private_key_algo", respBody, nil)),
		d.Set("usage", utils.PathSearch("usage", respBody, nil)),
		d.Set("company_country", utils.PathSearch("company_country", respBody, nil)),
		d.Set("company_province", utils.PathSearch("company_province", respBody, nil)),
		d.Set("company_city", utils.PathSearch("company_city", respBody, nil)),
		d.Set("company_name", utils.PathSearch("company_name", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCsrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "scm"
		httpUrl = "v3/scm/csr/{id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": d.Get("name"),
		},
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating CCM SSL CSR: %s", err)
	}

	return resourceCsrRead(ctx, d, meta)
}

func resourceCsrDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "scm"
		httpUrl = "v3/scm/csr/{id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CCM SSL CSR: %s", err)
	}

	return nil
}
