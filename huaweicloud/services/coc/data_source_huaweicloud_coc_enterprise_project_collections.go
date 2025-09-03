package coc

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

// @API COC GET /v1/enterprise-project-collect
func DataSourceCocEnterpriseProjectCollections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocEnterpriseProjectCollectionsRead,

		Schema: map[string]*schema.Schema{
			"unique_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ep_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceCocEnterpriseProjectCollectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/enterprise-project-collect"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var marker string
	res := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + buildGetEnterpriseProjectCollectionsParams(d, marker)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving COC enterprise project collections: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		views, nextMarker := flattenCocGetEnterpriseProjectCollections(getRespBody)
		if len(views) < 1 {
			break
		}
		res = append(res, views...)
		marker = nextMarker
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetEnterpriseProjectCollectionsParams(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if uniqueID, ok := d.GetOk("unique_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, uniqueID)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenCocGetEnterpriseProjectCollections(resp interface{}) ([]map[string]interface{}, string) {
	if resp == nil {
		return nil, ""
	}
	viewsJson := utils.PathSearch("data", resp, make([]interface{}, 0))
	viewsArray := viewsJson.([]interface{})
	if len(viewsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(viewsArray))
	marker := ""
	for _, view := range viewsArray {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", view, nil),
			"user_id":    utils.PathSearch("user_id", view, nil),
			"ep_id_list": utils.PathSearch("ep_id_list", view, nil),
		})
		marker = utils.PathSearch("id", view, "").(string)
	}
	return result, marker
}
