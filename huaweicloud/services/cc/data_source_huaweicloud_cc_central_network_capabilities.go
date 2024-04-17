package cc

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC GET /v3/{domain_id}/gcn/capabilities
func DataSourceCcCentralNetworkCapabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkCapabilitiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"capability": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the capability of the central network.`,
			},
			"capabilities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Central network capability list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account that the central network belongs to.`,
						},
						"capability": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The capability of the central network.`,
						},
						"specifications": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specifications of the central network capability.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCcCentralNetworkCapabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)

	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	httpUrl := "v3/{domain_id}/gcn/capabilities"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)

	params := buildCentralNetworkCapabilitiesQueryParams(d)
	path += params

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return diag.Errorf("error retrieving central network capabilities: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	rawCapabilities := utils.PathSearch("capabilities", respBody, make([]interface{}, 0))
	capabilities, err := flattenCentralNetworkCapabilities(rawCapabilities.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("capabilities", capabilities),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCentralNetworkCapabilitiesQueryParams(d *schema.ResourceData) string {
	var params string
	if capability, ok := d.GetOk("capability"); ok {
		params += "?capability=" + capability.(string)
	}
	return params
}

func flattenCentralNetworkCapabilities(capabilities []interface{}) ([]interface{}, error) {
	if capabilities == nil {
		return nil, nil
	}

	rst := make([]interface{}, 0, len(capabilities))

	for _, v := range capabilities {
		specifications := utils.PathSearch("specifications", v, nil)

		specJson, err := json.Marshal(specifications)
		if err != nil {
			return nil, err
		}

		rst = append(rst, map[string]interface{}{
			"domain_id":      utils.PathSearch("domain_id", v, nil),
			"capability":     utils.PathSearch("capability", v, nil),
			"specifications": string(specJson),
		})
	}

	return rst, nil
}
