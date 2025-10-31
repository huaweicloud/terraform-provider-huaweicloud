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

var identityCenterBearerTokenNonUpdateParams = []string{"identity_store_id", "tenant_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token/{token_id}
func ResourceIdentityCenterBearerToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterBearerTokenCreate,
		UpdateContext: resourceIdentityCenterBearerTokenUpdate,
		ReadContext:   resourceIdentityCenterBearerTokenRead,
		DeleteContext: resourceIdentityCenterBearerTokenDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterBearerTokenImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterBearerTokenNonUpdateParams),

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
			"tenant_id": {
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
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterBearerTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token"
		createProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	createPath = strings.ReplaceAll(createPath, "{tenant_id}", fmt.Sprintf("%v", d.Get("tenant_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center bearer token: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenId := utils.PathSearch("token_id", createRespBody, "").(string)
	if tokenId == "" {
		return diag.Errorf("unable to find the Identity Center bearer token ID from the API response")
	}
	d.SetId(tokenId)

	return resourceIdentityCenterBearerTokenRead(ctx, d, meta)
}

func resourceIdentityCenterBearerTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	listPath = strings.ReplaceAll(listPath, "{tenant_id}", fmt.Sprintf("%v", d.Get("tenant_id")))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center bearer token.")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	token := utils.PathSearch(fmt.Sprintf("bearer_tokens[?token_id =='%s']|[0]", d.Id()), listRespBody, nil)
	if token == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no bearer token found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("creation_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("creation_time", token, float64(0)).(float64))/1000, false)),
		d.Set("expiration_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expiration_time", token, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterBearerTokenUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterBearerTokenDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token/{token_id}"
		deleteProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deletePath = strings.ReplaceAll(deletePath, "{tenant_id}", fmt.Sprintf("%v", d.Get("tenant_id")))
	deletePath = strings.ReplaceAll(deletePath, "{token_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center bearer token: %s", err)
	}

	return nil
}

func resourceIdentityCenterBearerTokenImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, errors.New("invalid id format, must be <identity_store_id>/<tenant_id>/<token_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
		d.Set("tenant_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
