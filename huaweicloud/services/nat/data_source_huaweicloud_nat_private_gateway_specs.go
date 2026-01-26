package nat

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

// @API NAT GET /v3/{project_id}/private-nat/specs
func DataSourcePrivateGatewaySpecs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateGatewaySpecsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"specs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cbc_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sess_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bps_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pps_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qps_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenPrivateGatewaySpecs(specs []interface{}) []map[string]interface{} {
	if len(specs) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(specs))
	for _, v := range specs {
		rst = append(rst, map[string]interface{}{
			"name":     utils.PathSearch("name", v, nil),
			"code":     utils.PathSearch("code", v, nil),
			"cbc_code": utils.PathSearch("cbc_code", v, nil),
			"rule_max": utils.PathSearch("rule_max", v, nil),
			"sess_max": utils.PathSearch("sess_max", v, nil),
			"bps_max":  utils.PathSearch("bps_max", v, nil),
			"pps_max":  utils.PathSearch("pps_max", v, nil),
			"qps_max":  utils.PathSearch("qps_max", v, nil),
		})
	}

	return rst
}

func dataSourcePrivateGatewaySpecsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v3/{project_id}/private-nat/specs"
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving NAT private gateway specs %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("specs", flattenPrivateGatewaySpecs(utils.PathSearch("specs", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
