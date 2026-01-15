package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI POST /v1/agency
// @API CCI GET /v1/agency
func ResourceAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgencyCreate,
		ReadContext:   resourceAgencyRead,
		DeleteContext: resourceAgencyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trust_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trust_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"need_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/agency"
		product = "cci"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CCI agency: %s", err)
	}

	d.SetId(cfg.DomainID)

	return resourceAgencyRead(ctx, d, meta)
}

func resourceAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/agency"
		product = "cci"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	getPath := client.Endpoint + httpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCI agency")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("[0].name", getRespBody, nil)),
		d.Set("domain_id", utils.PathSearch("[0].domain_id", getRespBody, nil)),
		d.Set("trust_domain_id", utils.PathSearch("[0].trust_domain_id", getRespBody, nil)),
		d.Set("trust_domain_name", utils.PathSearch("[0].trust_domain_name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("[0].description", getRespBody, nil)),
		d.Set("duration", utils.PathSearch("[0].duration", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("[0].create_time", getRespBody, nil)),
		d.Set("need_update", utils.PathSearch("[0].need_update", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAgencyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CCI agency resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
