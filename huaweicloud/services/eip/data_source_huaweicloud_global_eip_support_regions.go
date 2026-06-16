package eip

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{domain_id}/geip/support-regions
func DataSourceGlobalEipSupportRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipSupportRegionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field is of boolean type in the API documentation,
			// here it is modified to string to meet different scenarios.
			"page_reverse": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The `sort_key` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `sort_dir` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Named as `id` in the API documentation, now changed to `support_region_ids`.
			"support_region_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_border_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_site": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remote_endpoint": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"support_regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_border_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipSupportRegionsQueryParams(d *schema.ResourceData, marker string) string {
	values := url.Values{}
	values.Set("limit", "2000")

	if raw, ok := d.GetOk("page_reverse"); ok {
		if b, err := strconv.ParseBool(raw.(string)); err == nil {
			values.Set("page_reverse", fmt.Sprint(b))
		}
	}

	type param struct {
		key      string
		query    string
		useMulti bool
	}

	params := []param{
		{"fields", "fields", true},
		{"sort_key", "sort_key", false},
		{"sort_dir", "sort_dir", false},
		{"support_region_ids", "id", true},
		{"instance_type", "instance_type", true},
		{"public_border_group", "public_border_group", true},
		{"access_site", "access_site", true},
		{"region_id", "region_id", true},
		{"remote_endpoint", "remote_endpoint", true},
		{"status", "status", true},
	}

	for _, p := range params {
		if raw, ok := d.GetOk(p.key); ok {
			if !p.useMulti {
				values.Set(p.query, fmt.Sprint(raw))
				continue
			}

			arr, ok := raw.([]interface{})
			if !ok {
				values.Add(p.query, fmt.Sprint(raw))
				continue
			}
			for _, v := range arr {
				values.Add(p.query, fmt.Sprint(v))
			}
		}
	}

	if marker != "" {
		values.Set("marker", marker)
	}

	enc := values.Encode()
	if enc == "" {
		return ""
	}

	return "?" + enc
}

func dataSourceGlobalEipSupportRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "geip"
		httpUrl    = "v3/{domain_id}/geip/support-regions"
		result     = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildGlobalEipSupportRegionsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global EIP support regions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		regionsResp := utils.PathSearch("support_regions", respBody, make([]interface{}, 0)).([]interface{})
		if len(regionsResp) == 0 {
			break
		}

		result = append(result, regionsResp...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("support_regions", flattenGlobalEipSupportRegions(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipSupportRegions(regions []interface{}) []interface{} {
	if len(regions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(regions))
	for _, r := range regions {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", r, nil),
			"instance_type":       utils.PathSearch("instance_type", r, nil),
			"access_site":         utils.PathSearch("access_site", r, nil),
			"region_id":           utils.PathSearch("region_id", r, nil),
			"public_border_group": utils.PathSearch("public_border_group", r, nil),
			"remote_endpoint":     utils.PathSearch("remote_endpoint", r, nil),
			"status":              utils.PathSearch("status", r, nil),
			"created_at":          utils.PathSearch("created_at", r, nil),
			"updated_at":          utils.PathSearch("updated_at", r, nil),
		})
	}

	return rst
}
