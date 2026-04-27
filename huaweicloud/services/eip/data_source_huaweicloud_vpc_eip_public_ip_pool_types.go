package eip

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

// @API EIP GET /v3/{project_id}/eip/publicip-pool-types
func DataSourceVpcEipPublicIpPoolTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipPublicIpPoolTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The `fields` field is of type string in the API documentation,
			// but according to its field description, it should be a list.
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `id`, here it is named `type_id`.
			"type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `publicip-pool-types`, object type.
			// But the actual return is `publicip_pool_types`, list type, here it is named `public_ip_pool_types`.
			"public_ip_pool_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildVpcEipPublicIpPoolTypesQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	rawArray, ok := d.Get("fields").([]interface{})
	if ok && len(rawArray) > 0 {
		for _, v := range rawArray {
			rst += fmt.Sprintf("&fields=%s", v.(string))
		}
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%s", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v)
	}

	if v, ok := d.GetOk("type_id"); ok {
		rst += fmt.Sprintf("&id=%s", v)
	}

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&name=%s", v)
	}

	if v, ok := d.GetOk("size"); ok {
		rst += fmt.Sprintf("&size=%d", v)
	}

	if v, ok := d.GetOk("status"); ok {
		rst += fmt.Sprintf("&status=%s", v)
	}

	if v, ok := d.GetOk("type"); ok {
		rst += fmt.Sprintf("&type=%s", v)
	}

	if v, ok := d.GetOk("description"); ok {
		rst += fmt.Sprintf("&description=%s", v)
	}

	if v, ok := d.GetOk("public_border_group"); ok {
		rst += fmt.Sprintf("&public_border_group=%s", v)
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceVpcEipPublicIpPoolTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/eip/publicip-pool-types"
		result     = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating EIP v3 client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildVpcEipPublicIpPoolTypesQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EIP public IP pool types: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		types := utils.PathSearch("publicip_pool_types", respBody, make([]interface{}, 0)).([]interface{})
		if len(types) == 0 {
			break
		}

		result = append(result, types...)

		// There is no return for the `page_info` field in the API documentation, but it actually exists.
		nextMarker = utils.PathSearch("page_info.previous_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("public_ip_pool_types", flattenVpcEipPublicIpPoolTypes(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVpcEipPublicIpPoolTypes(typesResp []interface{}) []interface{} {
	if len(typesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(typesResp))
	for _, v := range typesResp {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"type": utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
