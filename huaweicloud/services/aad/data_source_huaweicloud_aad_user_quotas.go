package aad

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v2/aad/user/quotas
func DataSourceUserQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserQuotasRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"overseas_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_port_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cc_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"geo_ip": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"white_ip": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAadUserQuotasQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?type=%v", d.Get("type"))

	if v, ok := d.GetOk("overseas_type"); ok {
		queryParams = fmt.Sprintf("%s&overseas_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("ip"); ok {
		queryParams = fmt.Sprintf("%s&ip=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceUserQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/user/quotas"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildAadUserQuotasQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving AAD user quotas: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("domain", utils.PathSearch("domain", respBody, nil)),
		d.Set("instance", utils.PathSearch("instance", respBody, nil)),
		d.Set("port", utils.PathSearch("port", respBody, nil)),
		d.Set("domain_port_quota", utils.PathSearch("domain_port_quota", respBody, nil)),
		d.Set("cc_quota", utils.PathSearch("cc_quota", respBody, nil)),
		d.Set("custom", utils.PathSearch("custom", respBody, nil)),
		d.Set("geo_ip", utils.PathSearch("geo_ip", respBody, nil)),
		d.Set("white_ip", utils.PathSearch("white_ip", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
