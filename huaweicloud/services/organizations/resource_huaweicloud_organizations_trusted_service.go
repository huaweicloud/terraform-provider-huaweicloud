// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

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
				Description: `Specifies the name of the trusted service principal.`,
			},
			"enabled_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the date when the trusted service was integrated with Organizations.`,
			},
		},
	}
}

func resourceTrustedServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createTrustedService: create Organizations trusted service
	var (
		createTrustedServiceHttpUrl = "v1/organizations/trusted-services/enable"
		createTrustedServiceProduct = "organizations"
	)
	createTrustedServiceClient, err := cfg.NewServiceClient(createTrustedServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createTrustedServicePath := createTrustedServiceClient.Endpoint + createTrustedServiceHttpUrl

	createTrustedServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createTrustedServiceOpt.JSONBody = utils.RemoveNil(buildTrustedServiceBodyParams(d))
	_, err = createTrustedServiceClient.Request("POST", createTrustedServicePath, &createTrustedServiceOpt)
	if err != nil {
		return diag.Errorf("error creating TrustedService: %s", err)
	}

	serviceName := d.Get("service").(string)
	d.SetId(serviceName)

	return resourceTrustedServiceRead(ctx, d, meta)
}

func buildTrustedServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"service_principal": utils.ValueIgnoreEmpty(d.Get("service")),
	}
	return bodyParams
}

func resourceTrustedServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getTrustedService: Query Organizations trusted service
	var (
		getTrustedServiceHttpUrl = "v1/organizations/trusted-services"
		getTrustedServiceProduct = "organizations"
	)
	getTrustedServiceClient, err := cfg.NewServiceClient(getTrustedServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getTrustedServiceBasePath := getTrustedServiceClient.Endpoint + getTrustedServiceHttpUrl

	getTrustedServicePath := getTrustedServiceBasePath + buildGetTrustedServiceQueryParams("")

	getTrustedServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var serviceName string
	var enabledAt string
getTrustedServicesLoop:
	for {
		getTrustedServiceResp, err := getTrustedServiceClient.Request("GET", getTrustedServicePath,
			&getTrustedServiceOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving TrustedService")
		}
		getTrustedServiceRespBody, err := utils.FlattenResponse(getTrustedServiceResp)
		if err != nil {
			return diag.FromErr(err)
		}

		trustedServices := utils.PathSearch("trusted_services", getTrustedServiceRespBody, nil)
		if trustedServices == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
		}

		for _, trustedService := range trustedServices.([]interface{}) {
			servicePrincipal := utils.PathSearch("service_principal", trustedService, "").(string)
			if servicePrincipal == d.Id() {
				serviceName = servicePrincipal
				enabledAt = utils.PathSearch("enabled_at", trustedService, "").(string)
				break getTrustedServicesLoop
			}
		}
		marker := utils.PathSearch("page_info.next_marker", getTrustedServiceRespBody, nil)
		if marker == nil {
			break
		}
		getTrustedServicePath = getTrustedServiceBasePath + buildGetTrustedServiceQueryParams(marker.(string))
	}

	if serviceName == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("service", serviceName),
		d.Set("enabled_at", enabledAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTrustedServiceQueryParams(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%s", res, marker)
	}
	return res
}

func resourceTrustedServiceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteTrustedService: Delete Organizations trusted service
	var (
		deleteTrustedServiceHttpUrl = "v1/organizations/trusted-services/disable"
		deleteTrustedServiceProduct = "organizations"
	)
	deleteTrustedServiceClient, err := cfg.NewServiceClient(deleteTrustedServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	deleteTrustedServicePath := deleteTrustedServiceClient.Endpoint + deleteTrustedServiceHttpUrl

	deleteTrustedServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteTrustedServiceOpt.JSONBody = utils.RemoveNil(buildTrustedServiceBodyParams(d))
	_, err = deleteTrustedServiceClient.Request("POST", deleteTrustedServicePath, &deleteTrustedServiceOpt)
	if err != nil {
		return diag.Errorf("error deleting TrustedService: %s", err)
	}

	return nil
}
