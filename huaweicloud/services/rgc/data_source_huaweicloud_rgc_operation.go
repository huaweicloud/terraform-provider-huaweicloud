package rgc

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

// @API RGC GET /v1/managed-organization
func DataSourceOperation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperationRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"percentage_complete": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"percentage_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"percentage_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"percentage_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOperationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getOperationProduct = "rgc"
	getOperationClient, err := cfg.NewServiceClient(getOperationProduct, region)
	if err != nil {
		return diag.Errorf("Error creating RGC client: %s", err)
	}

	getOperationRespBody, err := getOperation(getOperationClient, d)
	if err != nil {
		return diag.Errorf("error retrieving RGC operation: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("operation_id", utils.PathSearch("operation_id", getOperationRespBody, nil)),
		d.Set("percentage_complete", utils.PathSearch("percentage_complete", getOperationRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getOperationRespBody, nil)),
		d.Set("percentage_details", utils.PathSearch("percentage_details", getOperationRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getOperation(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		getOperationHttpUrl = "v1/managed-organization"
	)
	getOperationHttpPath := client.Endpoint + getOperationHttpUrl
	getOperationHttpPath += buildOperationQueryParams(d)

	getOperationHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOperationHttpResp, err := client.Request("GET", getOperationHttpPath, &getOperationHttpOpt)
	if err != nil {
		return nil, err
	}
	getOperationRespBody, err := utils.FlattenResponse(getOperationHttpResp)
	if err != nil {
		return nil, err
	}
	return getOperationRespBody, nil
}

func buildOperationQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("account_id"); ok {
		res = fmt.Sprintf("%s&account_id=%v", res, v)
	}

	if v, ok := d.GetOk("organizational_unit_id"); ok {
		res = fmt.Sprintf("%s&organizational_unit_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
