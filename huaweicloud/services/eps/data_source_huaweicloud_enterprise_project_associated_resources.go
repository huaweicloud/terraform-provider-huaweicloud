package eps

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS GET /v1.0/associated-resources/{resource_id}
func DataSourceAssociatedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssociatedResourcesRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Attributes.
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associated_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"errors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAssociatedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	showPath := client.Endpoint + buildAssociatedResourcesQueryParameter(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", showPath, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.PathSearch("id", respBody, "").(string))
	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("associated_resources", flattenListAssociatedResourcesResponseBody(respBody)),
		d.Set("errors", flattenListErrorsResponseBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAssociatedResourcesQueryParameter(d *schema.ResourceData) string {
	httpUrl := "v1.0/associated-resources/{resource_id}"

	httpUrl = strings.ReplaceAll(httpUrl, "{resource_id}", d.Get("resource_id").(string))
	httpUrl += fmt.Sprintf("?project_id=%s", d.Get("project_id"))
	httpUrl += fmt.Sprintf("&region_id=%s", d.Get("region_id"))
	httpUrl += fmt.Sprintf("&resource_type=%s", d.Get("resource_type"))

	return httpUrl
}

func flattenListAssociatedResourcesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("associated_resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"eip":           utils.PathSearch("eip", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}
	return res
}

func flattenListErrorsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("errors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"project_id":    utils.PathSearch("project_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"error_code":    utils.PathSearch("error_code", v, nil),
			"error_msg":     utils.PathSearch("error_msg", v, nil),
		})
	}
	return res
}
