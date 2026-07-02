package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/configurations/functions
func DataSourceConfigurationFunctions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationFunctionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_postpaid_mgmt": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_large_screen_mgmt": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_purchase_label_mgmt": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"billing_type_mgmt": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConfigurationFunctionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/configurations/functions"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster configuration functions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("support_postpaid_mgmt", utils.PathSearch("support_postpaid_mgmt", respBody, false)),
		d.Set("support_large_screen_mgmt", utils.PathSearch("support_large_screen_mgmt", respBody, false)),
		d.Set("support_purchase_label_mgmt", utils.PathSearch("support_purchase_label_mgmt", respBody, false)),
		d.Set("billing_type_mgmt", utils.PathSearch("billing_type_mgmt", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
