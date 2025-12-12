package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/resource-shares/resource-instances/count
func DataSourceResourceInstancesCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInstancesCountRead,
		Schema: map[string]*schema.Schema{
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceResourceInstancesCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getResourceInstancesCountProduct = "ram"
	getResourceInstancesCountClient, err := cfg.NewServiceClient(getResourceInstancesCountProduct, region)
	if err != nil {
		return diag.Errorf("Error creating RAM client: %s", err)
	}

	getResourceInstancesCountRespBody, err := getResourceInstancesCount(getResourceInstancesCountClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RAM resource instances count: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("total_count", utils.PathSearch("total_count", getResourceInstancesCountRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getResourceInstancesCount(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		getResourceInstancesCountHttpUrl = "v1/resource-shares/resource-instances/count"
	)
	getResourceInstancesCountHttpPath := client.Endpoint + getResourceInstancesCountHttpUrl

	getResourceInstancesCountHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildResourceInstancesCountFilterBody(d),
	}
	getResourceInstancesCountHttpResp, err := client.Request("POST", getResourceInstancesCountHttpPath, &getResourceInstancesCountHttpOpt)
	if err != nil {
		return nil, err
	}
	getResourceInstancesCountRespBody, err := utils.FlattenResponse(getResourceInstancesCountHttpResp)
	if err != nil {
		return nil, err
	}
	return getResourceInstancesCountRespBody, nil
}

func buildResourceInstancesCountFilterBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{}

	if v, ok := d.GetOk("without_any_tag"); ok {
		params["without_any_tag"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		params["tags"] = v
	}

	if v, ok := d.GetOk("matches"); ok {
		params["matches"] = v
	}

	return params
}
