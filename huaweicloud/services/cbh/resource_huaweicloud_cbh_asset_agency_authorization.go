package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH POST /v2/{project_id}/cbs/agency/authorization
// @API CBH GET /v2/{project_id}/cbs/agency/authorization
func ResourceAssetAgencyAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetAgencyAuthorizationCreate,
		UpdateContext: resourceAssetAgencyAuthorizationUpdate,
		ReadContext:   resourceAssetAgencyAuthorizationRead,
		DeleteContext: resourceAssetAgencyAuthorizationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"csms": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"kms": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func assetAgencyAuthorizationOperation(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	basePath := client.Endpoint + "v2/{project_id}/cbs/agency/authorization"
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	baseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	baseOpt.JSONBody = map[string]interface{}{
		"authorization": map[string]interface{}{
			"csms": d.Get("csms").(bool),
			"kms":  d.Get("kms").(bool),
		},
	}

	_, err := client.Request("POST", basePath, &baseOpt)

	return err
}

func resourceAssetAgencyAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	err = assetAgencyAuthorizationOperation(client, d)
	if err != nil {
		return diag.Errorf("error creating CBH asset agency authorization: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(id)

	return resourceAssetAgencyAuthorizationRead(ctx, d, meta)
}

func resourceAssetAgencyAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	basePath := client.Endpoint + "v2/{project_id}/cbs/agency/authorization"
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	baseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", basePath, &baseOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH asset agency authorization: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		// This resource is an operational resource, the create method supports enabling and disabling authorization.
		// The API for querying will not report errors and has no parameters, and will always be successful.
		// So the logic of checkDeleted is not added.
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("csms", utils.PathSearch("authorization.csms", getRespBody, false)),
		d.Set("kms", utils.PathSearch("authorization.kms", getRespBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAssetAgencyAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	err = assetAgencyAuthorizationOperation(client, d)
	if err != nil {
		return diag.Errorf("error updating CBH asset agency authorization: %s", err)
	}

	return resourceAssetAgencyAuthorizationRead(ctx, d, meta)
}

func resourceAssetAgencyAuthorizationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
