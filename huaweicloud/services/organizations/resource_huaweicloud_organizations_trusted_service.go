package organizations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	trustedServiceNotFoundErrCodes = []string{
		"Organizations.1900", // The specified trusted service is disabled.
	}

	organizationNotFoundErrCodes = []string{
		"Organizations.1015", // The organization is not found.
	}
)

// @API Organizations GET /v1/organizations/trusted-services
// @API Organizations POST /v1/organizations/trusted-services/disable
// @API Organizations POST /v1/organizations/trusted-services/enable
func ResourceTrustedService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrustedServiceCreate,
		ReadContext:   resourceTrustedServiceRead,
		DeleteContext: resourceTrustedServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the trusted service principal.`,
			},
			"enabled_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the trusted service was integrated with Organizations.`,
			},
		},
	}
}

func resourceTrustedServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/trusted-services/enable"
	)

	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTrustedServiceBodyParams(d)),
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error enabling trusted service: %s", err)
	}

	d.SetId(d.Get("service").(string))

	return resourceTrustedServiceRead(ctx, d, meta)
}

func buildTrustedServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"service_principal": utils.ValueIgnoreEmpty(d.Get("service")),
	}
}

func listTrustedServices(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/trusted-services"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		trustedServices := utils.PathSearch("trusted_services", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, trustedServices...)
		if len(trustedServices) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func GetTrustedService(client *golangsdk.ServiceClient, service string) (interface{}, error) {
	trustedServices, err := listTrustedServices(client)
	if err != nil {
		return nil, err
	}

	trustedService := utils.PathSearch(fmt.Sprintf("[?service_principal=='%s']|[0]", service), trustedServices, nil)
	if trustedService == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/trusted-services",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the trusted service (%s) does not exist", service)),
			},
		}
	}

	return trustedService, nil
}

func resourceTrustedServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	trustedService, err := GetTrustedService(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving trusted service",
		)
	}

	mErr := multierror.Append(
		d.Set("service", utils.PathSearch("service_principal", trustedService, nil)),
		d.Set("enabled_at", utils.PathSearch("enabled_at", trustedService, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTrustedServiceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/trusted-services/disable"
	)

	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTrustedServiceBodyParams(d)),
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(
				common.ConvertExpected400ErrInto404Err(err, "error_code", trustedServiceNotFoundErrCodes...),
				"error_code",
				organizationNotFoundErrCodes...,
			),
			fmt.Sprintf("error deleting trusted service (%s)", d.Id()))
	}

	return nil
}
