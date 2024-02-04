// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations/delegated-administrators/register
// @API Organizations GET /v1/organizations/delegated-administrators
// @API Organizations POST /v1/organizations/delegated-administrators/deregister
func ResourceDelegatedAdministrator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDelegatedAdministratorCreate,
		ReadContext:   resourceDelegatedAdministratorRead,
		DeleteContext: resourceDelegatedAdministratorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the unique ID of an account.`,
			},
			"service_principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the service principal.`,
			},
		},
	}
}

func resourceDelegatedAdministratorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDelegatedAdministrator: create Organizations delegated administrator
	var (
		createDelegatedAdministratorHttpUrl = "v1/organizations/delegated-administrators/register"
		createDelegatedAdministratorProduct = "organizations"
	)
	createDelegatedAdministratorClient, err := cfg.NewServiceClient(createDelegatedAdministratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createDelegatedAdministratorPath := createDelegatedAdministratorClient.Endpoint + createDelegatedAdministratorHttpUrl

	createDelegatedAdministratorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createDelegatedAdministratorOpt.JSONBody = utils.RemoveNil(deleteDelegatedAdministratorBodyParams(d))
	_, err = createDelegatedAdministratorClient.Request("POST", createDelegatedAdministratorPath,
		&createDelegatedAdministratorOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations delegated administrator: %s", err)
	}

	accountID := d.Get("account_id")
	servicePrincipal := d.Get("service_principal")

	d.SetId(fmt.Sprintf("%s/%s", accountID, servicePrincipal))

	return resourceDelegatedAdministratorRead(ctx, d, meta)
}

func resourceDelegatedAdministratorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDelegatedAdministrator: Query Organizations delegated administrator
	var (
		getDelegatedAdministratorHttpUrl = "v1/organizations/delegated-administrators"
		getDelegatedAdministratorProduct = "organizations"
	)
	getDelegatedAdministratorClient, err := cfg.NewServiceClient(getDelegatedAdministratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <account_id>/<service_principal>")
	}
	accountID := parts[0]
	servicePrincipal := parts[1]

	getDelegatedAdministratorPath := getDelegatedAdministratorClient.Endpoint + getDelegatedAdministratorHttpUrl

	getDelegatedAdministratorQueryParams := buildGetDelegatedAdministratorQueryParams(servicePrincipal)
	getDelegatedAdministratorPath += getDelegatedAdministratorQueryParams

	getDelegatedAdministratorResp, err := pagination.ListAllItems(
		getDelegatedAdministratorClient,
		"marker",
		getDelegatedAdministratorPath,
		&pagination.QueryOpts{MarkerField: "account_id"})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations delegated administrator")
	}

	getDelegatedAdministratorRespJson, err := json.Marshal(getDelegatedAdministratorResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDelegatedAdministratorRespBody interface{}
	err = json.Unmarshal(getDelegatedAdministratorRespJson, &getDelegatedAdministratorRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	delegatedAdministrator := utils.PathSearch(fmt.Sprintf("delegated_administrators|[?account_id=='%s']|[0]",
		accountID), getDelegatedAdministratorRespBody, nil)
	if delegatedAdministrator == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("account_id", accountID),
		d.Set("service_principal", servicePrincipal),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDelegatedAdministratorQueryParams(servicePrincipal string) string {
	return fmt.Sprintf("?service_principal=%v", servicePrincipal)
}

func resourceDelegatedAdministratorDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDelegatedAdministrator: Delete Organizations delegated administrator
	var (
		deleteDelegatedAdministratorHttpUrl = "v1/organizations/delegated-administrators/deregister"
		deleteDelegatedAdministratorProduct = "organizations"
	)
	deleteDelegatedAdministratorClient, err := cfg.NewServiceClient(deleteDelegatedAdministratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deleteDelegatedAdministratorPath := deleteDelegatedAdministratorClient.Endpoint + deleteDelegatedAdministratorHttpUrl

	deleteDelegatedAdministratorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteDelegatedAdministratorOpt.JSONBody = utils.RemoveNil(deleteDelegatedAdministratorBodyParams(d))
	_, err = deleteDelegatedAdministratorClient.Request("POST", deleteDelegatedAdministratorPath,
		&deleteDelegatedAdministratorOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations delegated administrator: %s", err)
	}

	return nil
}

func deleteDelegatedAdministratorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"account_id":        d.Get("account_id"),
		"service_principal": d.Get("service_principal"),
	}
	return bodyParams
}
