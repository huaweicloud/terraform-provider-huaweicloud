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

// @API AAD GET /v2/aad/policies/waf/blackwhite-list
func DataSourcePolicyBlackWhiteLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyBlackWhiteListsRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"overseas_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"black": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     blackWhiteSchema(),
			},
			"white": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     blackWhiteSchema(),
			},
		},
	}
}

func blackWhiteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyBlackWhiteListsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/policies/waf/blackwhite-list"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = fmt.Sprintf("%s?domain_name=%v", requestPath, d.Get("domain_name"))
	requestPath = fmt.Sprintf("%s&overseas_type=%v", requestPath, d.Get("overseas_type"))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD policy black white lists: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("black", flattenPolicyBlackOrWhiteListsAttribute(
			utils.PathSearch("black", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("white", flattenPolicyBlackOrWhiteListsAttribute(
			utils.PathSearch("white", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolicyBlackOrWhiteListsAttribute(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":        utils.PathSearch("id", v, nil),
			"type":      utils.PathSearch("type", v, nil),
			"ip":        utils.PathSearch("ip", v, nil),
			"domain_id": utils.PathSearch("domain_id", v, nil),
		})
	}

	return rst
}
