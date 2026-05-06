package eip

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{project_id}/geip/bindings
func DataSourceGlobalEipBindings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipBindingsRead,

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
			"geip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"geip_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `gcbandwidth.id`, here it is named `gcbandwidth_id`.
			"gcbandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `gcbandwidth.admin_status`,
			// here it is named `gcbandwidth_admin_status`.
			"gcbandwidth_admin_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `gcbandwidth.size`,
			// here it is named `gcbandwidth_size`.
			"gcbandwidth_size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `gcbandwidth.sla_level`,
			// here it is named `gcbandwidth_sla_level`.
			"gcbandwidth_sla_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `gcbandwidth.dscp`,
			// here it is named `gcbandwidth_dscp`.
			"gcbandwidth_dscp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.private_ip_address`,
			// here it is named `vnic_private_ip_address`.
			"vnic_private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.vpc_id`,
			// here it is named `vnic_vpc_id`.
			"vnic_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.port_id`,
			// here it is named `vnic_port_id`.
			"vnic_port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.device_id`,
			// here it is named `vnic_device_id`.
			"vnic_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.device_owner`,
			// here it is named `vnic_device_owner`.
			"vnic_device_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.device_owner_prefixlike`,
			// here it is named `vnic_device_owner_prefixlike`.
			"vnic_device_owner_prefixlike": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.instance_type`,
			// here it is named `vnic_instance_type`.
			"vnic_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field name in the API document is `vnic.instance_id`,
			// here it is named `vnic_instance_id`.
			"vnic_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"geip_bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geip_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_border_group": {
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
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"gcbandwidth": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"admin_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"short_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sla_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dscp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"vnic": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device_owner": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vtep": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vni": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port_profile": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vn_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"forward_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hash_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vni": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"nexthops": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_address": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mode": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
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
				},
			},
		},
	}
}

func buildGlobalGeipBindingsQueryParams(d *schema.ResourceData) string {
	values := url.Values{}
	// The maximum `limit` is `2000`.
	values.Set("limit", "2000")

	if raw, ok := d.Get("fields").([]interface{}); ok {
		for _, v := range raw {
			values.Add("fields", fmt.Sprint(v))
		}
	}

	type param struct {
		key   string
		query string
	}

	params := []param{
		{"geip_id", "geip_id"},
		{"geip_ip_address", "geip_ip_address"},
		{"public_border_group", "public_border_group"},
		{"instance_type", "instance_type"},
		{"instance_id", "instance_id"},
		{"instance_vpc_id", "instance_vpc_id"},
		{"gcbandwidth_id", "gcbandwidth.id"},
		{"gcbandwidth_admin_status", "gcbandwidth.admin_status"},
		{"gcbandwidth_size", "gcbandwidth.size"},
		{"gcbandwidth_sla_level", "gcbandwidth.sla_level"},
		{"gcbandwidth_dscp", "gcbandwidth.dscp"},
		{"vnic_private_ip_address", "vnic.private_ip_address"},
		{"vnic_vpc_id", "vnic.vpc_id"},
		{"vnic_port_id", "vnic.port_id"},
		{"vnic_device_id", "vnic.device_id"},
		{"vnic_device_owner", "vnic.device_owner"},
		{"vnic_device_owner_prefixlike", "vnic.device_owner_prefixlike"},
		{"vnic_instance_type", "vnic.instance_type"},
		{"vnic_instance_id", "vnic.instance_id"},
		{"sort_key", "sort_key"},
		{"sort_dir", "sort_dir"},
	}

	for _, p := range params {
		if v, ok := d.GetOk(p.key); ok {
			values.Set(p.query, fmt.Sprint(v))
		}
	}

	enc := values.Encode()
	if enc == "" {
		return ""
	}

	return "?" + enc
}

func dataSourceGlobalEipBindingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/geip/bindings"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 networking client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildGlobalGeipBindingsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global eip bindings: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		bindingsResp := utils.PathSearch("geip_bindings", respBody, make([]interface{}, 0)).([]interface{})
		if len(bindingsResp) == 0 {
			break
		}

		result = append(result, bindingsResp...)
		offset += len(bindingsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("geip_bindings", flattenGlobalGeipBindings(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalGeipBindings(bindings []interface{}) []interface{} {
	if len(bindings) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(bindings))
	for _, b := range bindings {
		rst = append(rst, map[string]interface{}{
			"geip_id":             utils.PathSearch("geip_id", b, nil),
			"geip_ip_address":     utils.PathSearch("geip_ip_address", b, nil),
			"public_border_group": utils.PathSearch("public_border_group", b, nil),
			"created_at":          utils.PathSearch("created_at", b, nil),
			"updated_at":          utils.PathSearch("updated_at", b, nil),
			"instance_type":       utils.PathSearch("instance_type", b, nil),
			"instance_id":         utils.PathSearch("instance_id", b, nil),
			"version":             utils.PathSearch("version", b, nil),
			"gcbandwidth": flattenGlobalGeipBindingGcbandwidth(
				utils.PathSearch("gcbandwidth", b, nil)),
			"vnic": flattenGlobalGeipBindingVnic(utils.PathSearch("vnic", b, nil)),
			"vn_list": flattenGlobalGeipBindingVnList(
				utils.PathSearch("vn_list", b, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenGlobalGeipBindingGcbandwidth(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":           utils.PathSearch("id", raw, nil),
			"admin_status": utils.PathSearch("admin_status", raw, nil),
			"size":         utils.PathSearch("size", raw, nil),
			"short_id":     utils.PathSearch("short_id", raw, nil),
			"sla_level":    utils.PathSearch("sla_level", raw, nil),
			"dscp":         utils.PathSearch("dscp", raw, nil),
		},
	}
}

func flattenGlobalGeipBindingVnic(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"private_ip_address": utils.PathSearch("private_ip_address", raw, nil),
			"device_id":          utils.PathSearch("device_id", raw, nil),
			"device_owner":       utils.PathSearch("device_owner", raw, nil),
			"vpc_id":             utils.PathSearch("vpc_id", raw, nil),
			"port_id":            utils.PathSearch("port_id", raw, nil),
			"mac":                utils.PathSearch("mac", raw, nil),
			"vtep":               utils.PathSearch("vtep", raw, nil),
			"vni":                utils.PathSearch("vni", raw, nil),
			"instance_id":        utils.PathSearch("instance_id", raw, nil),
			"instance_type":      utils.PathSearch("instance_type", raw, nil),
			"port_profile":       utils.PathSearch("port_profile", raw, nil),
		},
	}
}

func flattenGlobalGeipBindingVnList(raw []interface{}) []interface{} {
	if len(raw) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(raw))
	for _, item := range raw {
		rst = append(rst, map[string]interface{}{
			"id":           utils.PathSearch("id", item, nil),
			"owner":        utils.PathSearch("owner", item, nil),
			"location":     utils.PathSearch("location", item, nil),
			"forward_mode": utils.PathSearch("forward_mode", item, nil),
			"cluster_id":   utils.PathSearch("cluster_id", item, nil),
			"hash_mode":    utils.PathSearch("hash_mode", item, nil),
			"type":         utils.PathSearch("type", item, nil),
			"vni":          utils.PathSearch("vni", item, nil),
			"nexthops": flattenGlobalGeipBindingNexthops(
				utils.PathSearch("nexthops", item, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.PathSearch("created_at", item, nil),
			"updated_at": utils.PathSearch("updated_at", item, nil),
		})
	}

	return rst
}

func flattenGlobalGeipBindingNexthops(raw []interface{}) []interface{} {
	if len(raw) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(raw))
	for _, item := range raw {
		rst = append(rst, map[string]interface{}{
			"ip_address": utils.PathSearch("ip_address", item, nil),
			"mode":       utils.PathSearch("mode", item, nil),
		})
	}

	return rst
}
