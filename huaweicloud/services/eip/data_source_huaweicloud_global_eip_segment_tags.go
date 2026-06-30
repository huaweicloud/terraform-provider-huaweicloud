package eip

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/global-eip-segment/tags
func DataSourceGlobalEipSegmentTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipSegmentTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGlobalEipSegmentTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/global-eip-segment/tags"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving global EIP segment tags: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}
	d.SetId(generateId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tags", flattenGlobalEipSegmentTags(
			utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipSegmentTags(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}

	return rst
}
