package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v3/{project_id}/vpc/firewalls
// @API VPC GET /v3/{project_id}/vpc/firewalls/{id}
func DataSourceNetworkAcls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkAclsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ingress_rules": {
							Type:     schema.TypeList,
							Elem:     networkAclsRuleSchema(),
							Computed: true,
						},
						"egress_rules": {
							Type:     schema.TypeList,
							Elem:     networkAclsRuleSchema(),
							Computed: true,
						},
						"associated_subnets": {
							Type:     schema.TypeSet,
							Elem:     networkAclsSubnetSchema(),
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

func networkAclsRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
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
			"source_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip_address_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_ip_address_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func networkAclsSubnetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceNetworkAclsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getNetworkAclsHttpUrl := "v3/{project_id}/vpc/firewalls"
	getNetworkAclsPath := client.Endpoint + getNetworkAclsHttpUrl
	getNetworkAclsPath = strings.ReplaceAll(getNetworkAclsPath, "{project_id}", client.ProjectID)

	getNetworkAclsQueryParams := buildNetworkAclsQueryParams(d, cfg)
	getNetworkAclsPath += getNetworkAclsQueryParams

	getNetworkAclsResp, err := pagination.ListAllItems(
		client,
		"marker",
		getNetworkAclsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.FromErr(err)
	}

	getNetworkAclsRespJson, err := json.Marshal(getNetworkAclsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getNetworkAclsRespBody interface{}
	err = json.Unmarshal(getNetworkAclsRespJson, &getNetworkAclsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	ids := utils.PathSearch("firewalls[*].id", getNetworkAclsRespBody, []interface{}{})
	networkAcls := make([]map[string]interface{}, len(ids.([]interface{})))

	for i, id := range ids.([]interface{}) {
		getNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls/" + id.(string)
		getNetworkAclPath := client.Endpoint + getNetworkAclHttpUrl
		getNetworkAclPath = strings.ReplaceAll(getNetworkAclPath, "{project_id}", client.ProjectID)

		getNetworkAclOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		getNetworkAclResp, err := client.Request("GET", getNetworkAclPath, &getNetworkAclOpt)
		if err != nil {
			return diag.Errorf("error retrieving VPC network ACL(%s): %s", id, err)
		}

		getNetworkAclRespBody, err := utils.FlattenResponse(getNetworkAclResp)
		if err != nil {
			return diag.FromErr(err)
		}

		networkAcls[i] = flattenNetworkAcl(getNetworkAclRespBody)
	}

	mErr := multierror.Append(
		nil,
		d.Set("network_acls", networkAcls),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildNetworkAclsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := fmt.Sprintf("?enterprise_project_id=%v", cfg.GetEnterpriseProjectID(d, "all_granted_eps"))

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("network_acl_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("enabled"); ok {
		res = fmt.Sprintf("%s&admin_state_up=%v", res, v)
	}

	return res
}

func flattenNetworkAcl(getNetworkAclRespBody interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":                  utils.PathSearch("firewall.name", getNetworkAclRespBody, nil),
		"id":                    utils.PathSearch("firewall.id", getNetworkAclRespBody, nil),
		"description":           utils.PathSearch("firewall.description", getNetworkAclRespBody, nil),
		"enterprise_project_id": utils.PathSearch("firewall.enterprise_project_id", getNetworkAclRespBody, nil),
		"enabled":               utils.PathSearch("firewall.admin_state_up", getNetworkAclRespBody, nil),
		"ingress_rules":         flattenRules(utils.PathSearch("firewall.ingress_rules", getNetworkAclRespBody, nil)),
		"egress_rules":          flattenRules(utils.PathSearch("firewall.egress_rules", getNetworkAclRespBody, nil)),
		"associated_subnets":    flattenSubnets(utils.PathSearch("firewall.associations", getNetworkAclRespBody, nil)),
		"status":                utils.PathSearch("firewall.status", getNetworkAclRespBody, nil),
		"created_at":            utils.PathSearch("firewall.created_at", getNetworkAclRespBody, nil),
		"updated_at":            utils.PathSearch("firewall.updated_at", getNetworkAclRespBody, nil),
	}
}
