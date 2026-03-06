package css

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections
func DataSourceVpcepserviceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_session": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpcep_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpcep_ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpcep_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpcep_update_switch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcepConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		httpUrl           = "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections?limit={limit}"
		clusterId         = d.Get("cluster_id").(string)
		offset            = 0
		limit             = 100
		result            = make([]interface{}, 0)
		permissions       = make([]interface{}, 0)
		vpcServiceName    = ""
		vpcepUpdateSwitch bool
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving endpoint connections: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		vpcServiceName = utils.PathSearch("vpcServiceName", getRespBody, "").(string)
		vpcepUpdateSwitch = utils.PathSearch("vpcepUpdateSwitch", getRespBody, false).(bool)
		permissions = utils.PathSearch("permissions", getRespBody, make([]interface{}, 0)).([]interface{})
		connections := utils.PathSearch("connections", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, connections...)
		if len(connections) < limit {
			break
		}

		offset += len(connections)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_service_name", vpcServiceName),
		d.Set("vpcep_update_switch", vpcepUpdateSwitch),
		d.Set("connections", flattenVpcepserviceConnections(result)),
		d.Set("permissions", flattenVpcepservicePermissions(permissions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVpcepserviceConnections(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"max_session":        utils.PathSearch("maxSession", v, nil),
			"specification_name": utils.PathSearch("specificationName", v, nil),
			"created_at":         utils.PathSearch("created_at", v, nil),
			"update_at":          utils.PathSearch("update_at", v, nil),
			"domain_id":          utils.PathSearch("domain_id", v, nil),
			"vpcep_ip":           utils.PathSearch("vpcepIp", v, nil),
			"vpcep_ipv6_address": utils.PathSearch("vpcepIpv6Address", v, nil),
			"vpcep_dns_name":     utils.PathSearch("vpcepDnsName", v, nil),
		})
	}

	return rst
}

func flattenVpcepservicePermissions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"permission":      utils.PathSearch("permission", v, nil),
			"permission_type": utils.PathSearch("permission_type", v, nil),
			"created_at":      utils.PathSearch("created_at", v, nil),
		})
	}

	return rst
}
