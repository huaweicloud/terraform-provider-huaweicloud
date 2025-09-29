package ecs

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

// @API ECS GET /v1/{project_id}/cloudservers/limits
func DataSourceEcsComputeQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEcsComputeQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"absolute": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the tenant quotas.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_security_groups": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_server_group_members": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_server_groups": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_floating_ips": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_spot_instances": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_spot_ram_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_spot_cores_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_personality": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_security_group_rules": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_security_groups_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_spot_ram_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_server_meta": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_instances": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_ram_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_cores_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_instances_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_ram_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_image_meta": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_personality_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_keypairs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_floating_ips_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_server_groups_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_total_spot_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_spot_instances_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_cluster_server_group_members": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_fault_domain_members": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEcsComputeQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/cloudservers/limits"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ECS quotas: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("absolute", flattenComputeAbsoluteResponseBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComputeAbsoluteResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("absolute", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"max_security_groups":              utils.PathSearch("maxSecurityGroups", curJson, nil),
			"max_server_group_members":         utils.PathSearch("maxServerGroupMembers", curJson, nil),
			"max_server_groups":                utils.PathSearch("maxServerGroups", curJson, nil),
			"max_total_cores":                  utils.PathSearch("maxTotalCores", curJson, nil),
			"max_total_floating_ips":           utils.PathSearch("maxTotalFloatingIps", curJson, nil),
			"max_total_spot_instances":         utils.PathSearch("maxTotalSpotInstances", curJson, nil),
			"max_total_spot_ram_size":          utils.PathSearch("maxTotalSpotRAMSize", curJson, nil),
			"total_spot_cores_used":            utils.PathSearch("totalSpotCoresUsed", curJson, nil),
			"max_personality":                  utils.PathSearch("maxPersonality", curJson, nil),
			"max_security_group_rules":         utils.PathSearch("maxSecurityGroupRules", curJson, nil),
			"total_security_groups_used":       utils.PathSearch("totalSecurityGroupsUsed", curJson, nil),
			"total_spot_ram_used":              utils.PathSearch("totalSpotRAMUsed", curJson, nil),
			"max_server_meta":                  utils.PathSearch("maxServerMeta", curJson, nil),
			"max_total_instances":              utils.PathSearch("maxTotalInstances", curJson, nil),
			"max_total_ram_size":               utils.PathSearch("maxTotalRAMSize", curJson, nil),
			"total_cores_used":                 utils.PathSearch("totalCoresUsed", curJson, nil),
			"total_instances_used":             utils.PathSearch("totalInstancesUsed", curJson, nil),
			"total_ram_used":                   utils.PathSearch("totalRAMUsed", curJson, nil),
			"max_image_meta":                   utils.PathSearch("maxImageMeta", curJson, nil),
			"max_personality_size":             utils.PathSearch("maxPersonalitySize", curJson, nil),
			"max_total_keypairs":               utils.PathSearch("maxTotalKeypairs", curJson, nil),
			"total_floating_ips_used":          utils.PathSearch("totalFloatingIpsUsed", curJson, nil),
			"total_server_groups_used":         utils.PathSearch("totalServerGroupsUsed", curJson, nil),
			"max_total_spot_cores":             utils.PathSearch("maxTotalSpotCores", curJson, nil),
			"total_spot_instances_used":        utils.PathSearch("totalSpotInstancesUsed", curJson, nil),
			"max_cluster_server_group_members": utils.PathSearch("maxClusterServerGroupMembers", curJson, nil),
			"max_fault_domain_members":         utils.PathSearch("maxFaultDomainMembers", curJson, nil),
		},
	}
	return res
}
