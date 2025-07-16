package workspace

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

// @API Workspace GET /v2/{project_id}/subnets/{subnet_id}/available-ip
func DataSourceAvailableIpNumber() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailableIpNumberRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the subnet is located.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the subnet to be queried.`,
			},
			"available_ip": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of available IPs in the subnet.`,
			},
		},
	}
}

func queryAvailableIpNumber(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl  = "v2/{project_id}/subnets/{subnet_id}/available-ip"
		subnetId = d.Get("subnet_id").(string)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{subnet_id}", subnetId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	availableIpNumber := int(utils.PathSearch("available_ip", respBody, float64(0)).(float64))

	return availableIpNumber, nil
}

func dataSourceAvailableIpNumberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	availableIpNumber, err := queryAvailableIpNumber(client, d)
	if err != nil {
		return diag.Errorf("error getting available IP number: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("available_ip", availableIpNumber),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
