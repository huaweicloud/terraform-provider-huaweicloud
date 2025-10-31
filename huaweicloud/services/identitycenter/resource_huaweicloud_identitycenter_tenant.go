package identitycenter

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

var identityCenterTenantNonUpdateParams = []string{"identity_store_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/provision-tenant
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/provision-tenant
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}
func ResourceIdentityCenterTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterTenantCreate,
		UpdateContext: resourceIdentityCenterTenantUpdate,
		ReadContext:   resourceIdentityCenterTenantRead,
		DeleteContext: resourceIdentityCenterTenantDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterTenantImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterTenantNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scim_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterTenantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v1/identity-stores/{identity_store_id}/provision-tenant"
		createProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center tenant: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantId := utils.PathSearch("tenant_id", createRespBody, "").(string)
	if tenantId == "" {
		return diag.Errorf("unable to find the Identity Center tenant ID from the API response")
	}
	d.SetId(tenantId)

	return resourceIdentityCenterTenantRead(ctx, d, meta)
}

func resourceIdentityCenterTenantRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/provision-tenant"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center provisioning tenant.")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tenant := utils.PathSearch(fmt.Sprintf("provisioning_tenants|[?tenant_id =='%s']|[0]", d.Id()), listRespBody, nil)
	if tenant == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no tenant found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("creation_time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("creation_time", tenant, float64(0)).(float64))/1000, false)),
		d.Set("scim_endpoint", utils.PathSearch("scim_endpoint", tenant, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterTenantUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterTenantDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}"
		deleteProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deletePath = strings.ReplaceAll(deletePath, "{tenant_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center Tenant: %s", err)
	}

	return nil
}

func resourceIdentityCenterTenantImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid id format, must be <identity_store_id>/<tenant_id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
