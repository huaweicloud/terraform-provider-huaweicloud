package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/cluster-protect/protection-item
func DataSourceClusterProtectProtectionItems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterProtectProtectionItemsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vuls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"baselines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"baseline_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_index": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"malwares": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"malware_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_ns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cluster_labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildClusterProtectProtectionItemsQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceClusterProtectProtectionItemsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/cluster-protect/protection-item"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildClusterProtectProtectionItemsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS cluster protect all protection items: %s", err)
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
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, 0)),
		d.Set("vuls", utils.ExpandToStringList(utils.PathSearch("vuls", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("baselines", flattenProtectionItemsBaselines(utils.PathSearch("baselines", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("malwares", flattenProtectionItemsMalwares(utils.PathSearch("malwares", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("images", flattenProtectionItemsImages(utils.PathSearch("images", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("clusters", flattenProtectionItemsClusters(utils.PathSearch("clusters", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProtectionItemsBaselines(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"baseline_desc":  utils.PathSearch("baseline_desc", v, nil),
			"baseline_index": utils.PathSearch("baseline_index", v, nil),
			"baseline_type":  utils.PathSearch("baseline_type", v, nil),
		})
	}

	return rst
}

func flattenProtectionItemsMalwares(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"malware_type": utils.PathSearch("malware_type", v, nil),
		})
	}

	return rst
}

func flattenProtectionItemsImages(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"image_name":    utils.PathSearch("image_name", v, nil),
			"image_version": utils.PathSearch("image_version", v, nil),
			"id":            utils.PathSearch("id", v, nil),
		})
	}

	return rst
}

func flattenProtectionItemsClusters(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"cluster_ns":     utils.ExpandToStringList(utils.PathSearch("cluster_ns", v, make([]interface{}, 0)).([]interface{})),
			"cluster_labels": utils.ExpandToStringList(utils.PathSearch("cluster_labels", v, make([]interface{}, 0)).([]interface{})),
			"protect_status": utils.PathSearch("protect_status", v, nil),
		})
	}

	return rst
}
