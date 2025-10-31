package rgc

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

// @API RGC GET /v1/managed-organization/managed-accounts/{managed_account_id}/blueprint
func DataSourceBluePrint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBluePrintRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_impl_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_impl_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBluePrintRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getBluePrintProduct = "rgc"
	getBluePrintClient, err := cfg.NewServiceClient(getBluePrintProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getBluePrintRespBody, err := getBluePrint(getBluePrintClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC blueprint: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("manage_account_id", utils.PathSearch("manage_account_id", getBluePrintRespBody, nil)),
		d.Set("account_id", utils.PathSearch("account_id", getBluePrintRespBody, nil)),
		d.Set("account_name", utils.PathSearch("account_name", getBluePrintRespBody, nil)),
		d.Set("blueprint_product_id", utils.PathSearch("blueprint_product_id", getBluePrintRespBody, nil)),
		d.Set("blueprint_product_name", utils.PathSearch("blueprint_product_name", getBluePrintRespBody, nil)),
		d.Set("blueprint_product_version", utils.PathSearch("blueprint_product_version", getBluePrintRespBody, nil)),
		d.Set("blueprint_status", utils.PathSearch("blueprint_status", getBluePrintRespBody, nil)),
		d.Set("blueprint_product_impl_type", utils.PathSearch("blueprint_product_impl_type", getBluePrintRespBody, nil)),
		d.Set("blueprint_product_impl_detail", utils.PathSearch("blueprint_product_impl_detail", getBluePrintRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getBluePrint(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	accountId := d.Get("managed_account_id").(string)

	var (
		getBluePrintHttpUrl = "v1/managed-organization/managed-accounts/{managed_account_id}/blueprint"
	)
	getBluePrintHttpPath := client.Endpoint + getBluePrintHttpUrl
	getBluePrintHttpPath = strings.ReplaceAll(getBluePrintHttpPath, "{managed_account_id}", accountId)

	getBluePrintHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBluePrintHttpResp, err := client.Request("GET", getBluePrintHttpPath, &getBluePrintHttpOpt)
	if err != nil {
		return nil, err
	}
	getBluePrintRespBody, err := utils.FlattenResponse(getBluePrintHttpResp)
	if err != nil {
		return nil, err
	}
	return getBluePrintRespBody, nil
}
