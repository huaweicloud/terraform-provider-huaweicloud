package eip

import (
	"context"
	"fmt"
	"net/url"
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

// @API EIP GET /v3/{domain_id}/global-eip-segments
func DataSourceGlobalEipSegments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipSegmentsRead,

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
			// Named as `id` in the API documentation, now changed to `segment_ids`.
			"segment_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"internet_bandwidth_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_site": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"geip_pool_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cidr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cidr_v6": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"freezen": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"internet_bandwidth_is_null": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_region": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_instance_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_public_border_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_instance_site": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_instance_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_project_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_service_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"associate_instance_service_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// There is no explanation of the field structure for response in the API documentation.
			// The response field in the current schema is referenced from the response example section in the API
			// documentation.
			// The structure of some fields is unclear, so they were discarded in the schema, such as `tags`.
			"global_eip_segments": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geip_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_v6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"freezen": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"freezen_info": {
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
						"internet_bandwidth": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"is_pre_paid": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_charged": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipSegmentsQueryParams(d *schema.ResourceData, marker string) string {
	values := url.Values{}
	values.Set("limit", "2000")

	add := func(name string, raw interface{}, useMulti bool) {
		if !useMulti {
			values.Set(name, fmt.Sprint(raw))
			return
		}
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				values.Add(name, fmt.Sprint(v))
			}
		}
	}

	if marker != "" {
		values.Set("marker", marker)
	}

	type param struct {
		key      string
		query    string
		useMulti bool
	}
	params := []param{
		{"page_reverse", "page_reverse", false},
		{"fields", "fields", true},
		{"sort_key", "sort_key", false},
		{"sort_dir", "sort_dir", false},
		{"segment_ids", "id", true},
		{"internet_bandwidth_id", "internet_bandwidth_id", true},
		{"name", "name", true},
		{"name_like", "name_like", false},
		{"access_site", "access_site", true},
		{"geip_pool_name", "geip_pool_name", true},
		{"isp", "isp", true},
		{"ip_address", "ip_address", true},
		{"ipv6_address", "ipv6_address", true},
		{"ip_version", "ip_version", true},
		{"cidr", "cidr", true},
		{"cidr_v6", "cidr_v6", true},
		{"freezen", "freezen", true},
		{"internet_bandwidth_is_null", "internet_bandwidth_is_null", true},
		{"status", "status", true},
		{"associate_instance_region", "associate_instance.region", true},
		{"associate_instance_instance_type", "associate_instance.instance_type", true},
		{"associate_instance_public_border_group", "associate_instance.public_border_group", true},
		{"associate_instance_instance_site", "associate_instance.instance_site", true},
		{"associate_instance_instance_id", "associate_instance.instance_id", true},
		{"associate_instance_project_id", "associate_instance.project_id", true},
		{"associate_instance_service_id", "associate_instance.service_id", true},
		{"associate_instance_service_type", "associate_instance.service_type", true},
		{"enterprise_project_id", "enterprise_project_id", true},
	}

	for _, p := range params {
		if raw, ok := d.GetOk(p.key); ok {
			add(p.query, raw, p.useMulti)
		}
	}

	return "?" + values.Encode()
}

func dataSourceGlobalEipSegmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "geip"
		httpUrl    = "v3/{domain_id}/global-eip-segments"
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
		requestPathWithMarker := requestPath + buildGlobalEipSegmentsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global EIP segments: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		segmentsResp := utils.PathSearch("global_eip_segments", respBody, make([]interface{}, 0)).([]interface{})
		if len(segmentsResp) == 0 {
			break
		}

		result = append(result, segmentsResp...)

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
		d.Set("global_eip_segments", flattenGlobalEipSegments(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipSegments(segments []interface{}) []interface{} {
	if len(segments) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(segments))
	for _, seg := range segments {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", seg, nil),
			"name":           utils.PathSearch("name", seg, nil),
			"description":    utils.PathSearch("description", seg, nil),
			"domain_id":      utils.PathSearch("domain_id", seg, nil),
			"access_site":    utils.PathSearch("access_site", seg, nil),
			"geip_pool_name": utils.PathSearch("geip_pool_name", seg, nil),
			"isp":            utils.PathSearch("isp", seg, nil),
			"ip_version":     int(utils.PathSearch("ip_version", seg, float64(0)).(float64)),
			"cidr":           utils.PathSearch("cidr", seg, nil),
			"cidr_v6":        utils.PathSearch("cidr_v6", seg, nil),
			"freezen":        utils.PathSearch("freezen", seg, false),
			"freezen_info":   utils.PathSearch("freezen_info", seg, nil),
			"status":         utils.PathSearch("status", seg, nil),
			"created_at":     utils.PathSearch("created_at", seg, nil),
			"updated_at":     utils.PathSearch("updated_at", seg, nil),
			"internet_bandwidth": flattenGlobalEipSegmentsInternetBandwidth(
				utils.PathSearch("internet_bandwidth", seg, nil)),
			"is_pre_paid":           utils.PathSearch("is_pre_paid", seg, false),
			"is_charged":            utils.PathSearch("is_charged", seg, false),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", seg, nil),
		})
	}

	return rst
}

func flattenGlobalEipSegmentsInternetBandwidth(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", raw, nil),
			"size": int(utils.PathSearch("size", raw, float64(0)).(float64)),
		},
	}
}
