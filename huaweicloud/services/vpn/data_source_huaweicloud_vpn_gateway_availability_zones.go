// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/vpn-gateways/availability-zones
func DataSourceVpnGatewayAZs() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpnGatewayAZsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the flavor name.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vpc",
				Description: `Specifies the attachment type.`,
			},
			"names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The names of the availability zones.`,
			},
		},
	}
}

func resourceVpnGatewayAZsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGatewayAZs: Query VPN gateway AZs
	var (
		getGatewayAZsHttpUrl = "v5/{project_id}/vpn-gateways/availability-zones"
		getGatewayAZsProduct = "vpn"
	)
	getGatewayAZsClient, err := cfg.NewServiceClient(getGatewayAZsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getGatewayAZsPath := getGatewayAZsClient.Endpoint + getGatewayAZsHttpUrl
	getGatewayAZsPath = strings.ReplaceAll(getGatewayAZsPath, "{project_id}", getGatewayAZsClient.ProjectID)

	getGatewayAZsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getGatewayAZsResp, err := getGatewayAZsClient.Request("GET", getGatewayAZsPath, &getGatewayAZsOpt)

	if err != nil {
		return diag.Errorf("error retrieving VPN gateway AZs: %s", err)
	}

	getGatewayAZsRespBody, err := utils.FlattenResponse(getGatewayAZsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	flavor := strings.ToLower(d.Get("flavor").(string))
	attachmentType := d.Get("attachment_type").(string)

	azPath := fmt.Sprintf("availability_zones.%s.%s", flavor, attachmentType)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("names", utils.PathSearch(azPath, getGatewayAZsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
