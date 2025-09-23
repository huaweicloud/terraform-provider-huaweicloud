// Some response parameters in the API documentation display abnormally, causing the PMS platform to be unable to
// recognize the response parameters. Therefore, this dataSource is written in an automatically generated format.
package ims

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

// @API IMS GET /v1/cloudimages/os_version
func DataSourceOsVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOsVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"platform": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_version_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_bit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOsVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.ImageV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v1 client: %s", err)
	}

	getPath := client.Endpoint + "v1/cloudimages/os_version"
	if v, ok := d.GetOk("tag"); ok && v.(string) != "" {
		result := strings.Split(v.(string), ",")
		for i, tag := range result {
			if i == 0 {
				getPath += "?tag=" + tag
			} else {
				getPath += "&tag=" + tag
			}
		}
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving IMS OS versions, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("os_versions", flattenOsVersions(getRespBody.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOsVersions(osVersions []interface{}) []interface{} {
	result := make([]interface{}, 0, len(osVersions))

	for _, osVersion := range osVersions {
		result = append(result, map[string]interface{}{
			"platform": utils.PathSearch("platform", osVersion, nil),
			"versions": flattenVersions(utils.PathSearch("version_list", osVersion, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenVersions(versions []interface{}) []interface{} {
	result := make([]interface{}, 0, len(versions))

	for _, version := range versions {
		result = append(result, map[string]interface{}{
			"platform":       utils.PathSearch("platform", version, nil),
			"os_version_key": utils.PathSearch("os_version_key", version, nil),
			"os_version":     utils.PathSearch("os_version", version, nil),
			"os_bit":         utils.PathSearch("os_bit", version, nil),
			"os_type":        utils.PathSearch("os_type", version, nil),
		})
	}

	return result
}
