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

// @API HSS GET /v5/{project_id}/kubernetes/container/detail
func DataSourceKubernetesContainerDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesContainerDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_password": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_port_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"enable_simulate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extra": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"openvpn": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"outside_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"outside_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"linux": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"os": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"rdp": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proto_env": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"mysql": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_path": {
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

func buildKubernetesContainerDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?container_id=%v", d.Get("container_id").(string))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceKubernetesContainerDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/kubernetes/container/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildKubernetesContainerDetailQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS kubernetes container detail: %s", err)
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
		d.Set("service_name", utils.PathSearch("service_name", respBody, nil)),
		d.Set("service_username", utils.PathSearch("service_username", respBody, nil)),
		d.Set("service_password", utils.PathSearch("service_password", respBody, nil)),
		d.Set("service_port_list", flattenKubernetesContainerServicePortList(
			utils.PathSearch("service_port_list", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("enable_simulate", utils.PathSearch("enable_simulate", respBody, nil)),
		d.Set("hosts", utils.PathSearch("hosts", respBody, nil)),
		d.Set("extra", flattenKubernetesContainerExtraInfo(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKubernetesContainerServicePortList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"desc":      utils.PathSearch("desc", v, nil),
			"type":      utils.PathSearch("type", v, nil),
			"protocol":  utils.PathSearch("protocol", v, nil),
			"user_port": utils.PathSearch("user_port", v, nil),
			"port":      utils.PathSearch("port", v, nil),
		})
	}

	return rst
}

func flattenKubernetesContainerExtraInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	extraInfo := utils.PathSearch("extra", resp, nil)
	if extraInfo == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	item := map[string]interface{}{
		"openvpn": flattenOpenvpnInfo(extraInfo),
		"linux":   flattenLinuxInfo(extraInfo),
		"rdp":     flattenRdpInfo(extraInfo),
		"mysql":   flattenMysqlInfo(extraInfo),
	}
	rst = append(rst, item)

	return rst
}

func flattenOpenvpnInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	openvpnInfo := utils.PathSearch("openvpn", resp, nil)
	if openvpnInfo == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"outside_ip":   utils.PathSearch("outside_ip", openvpnInfo, nil),
		"outside_port": utils.PathSearch("outside_port", openvpnInfo, nil),
	})

	return rst
}

func flattenLinuxInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	linuxInfo := utils.PathSearch("linux", resp, nil)
	if linuxInfo == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"os": utils.PathSearch("os", linuxInfo, nil),
	})

	return rst
}

func flattenRdpInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rdpInfo := utils.PathSearch("rdp", resp, nil)
	if rdpInfo == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"proto_env": utils.PathSearch("proto_env", rdpInfo, nil),
		"system":    utils.PathSearch("system", rdpInfo, nil),
	})

	return rst
}

func flattenMysqlInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	mysqlInfo := utils.PathSearch("mysql", resp, nil)
	if mysqlInfo == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"custom_path": utils.PathSearch("custom_path", mysqlInfo, nil),
	})

	return rst
}
