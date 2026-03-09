package organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	delegatedAdministratorNotFoundErrCodes = []string{
		"Organizations.1500", // The delegated administrator not found.
	}
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
				Description: `The unique ID of an account.`,
			},
			"service_principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the service principal.`,
			},
		},
	}
}

func resourceDelegatedAdministratorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/delegated-administrators/register"
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDelegatedAdministratorBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations delegated administrator: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("account_id"), d.Get("service_principal")))

	return resourceDelegatedAdministratorRead(ctx, d, meta)
}

func listDelegatedAdministrators(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/delegated-administrators"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		delegatedAdministrators := utils.PathSearch("delegated_administrators", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, delegatedAdministrators...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func GetDelegatedAdministrator(client *golangsdk.ServiceClient, accountId, servicePrincipal string) (interface{}, error) {
	delegatedAdministrators, err := listDelegatedAdministrators(client, fmt.Sprintf("?service_principal=%v", servicePrincipal))
	if err != nil {
		return nil, err
	}

	delegatedAdministrator := utils.PathSearch(fmt.Sprintf("[?account_id=='%s']|[0]", accountId), delegatedAdministrators, nil)
	if delegatedAdministrator == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/delegated-administrators",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the delegated administrator (%s) does not exist for account (%s)", servicePrincipal, accountId)),
			},
		}
	}

	return delegatedAdministrator, nil
}

func resourceDelegatedAdministratorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<account_id>/<service_principal>', but got '%s'", d.Id())
	}

	servicePrincipal := parts[1]
	delegatedAdministrator, err := GetDelegatedAdministrator(client, parts[0], servicePrincipal)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving Organizations delegated administrator",
		)
	}

	mErr := multierror.Append(
		d.Set("account_id", utils.PathSearch("account_id", delegatedAdministrator, nil)),
		d.Set("service_principal", servicePrincipal),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDelegatedAdministratorDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/delegated-administrators/deregister"
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDelegatedAdministratorBodyParams(d),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
				"error_code", delegatedAdministratorNotFoundErrCodes...),
			"error deleting Organizations delegated administrator",
		)
	}

	return nil
}

func buildDelegatedAdministratorBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"account_id":        d.Get("account_id"),
		"service_principal": d.Get("service_principal"),
	}
}
