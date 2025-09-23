package as

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group_tag/{id}/tags
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group
func DataSourceASGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceASGroupRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the AS groups are located.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AS group name used to query group list.",
			},
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AS group configuration id used to query group list.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AS group status used to query group list.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Enterprise Project id used to query group list.",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group name of the AS group.",
						},
						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group id of the AS group.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group status of the AS group.",
						},
						"scaling_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration id of the AS group.",
						},
						"scaling_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration name of the AS group.",
						},
						"current_instance_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of current instances in the AS group.",
						},
						"desire_instance_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The expected number of instances in the AS group.",
						},
						"min_instance_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum number of instances in the AS group.",
						},
						"max_instance_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of instances in the AS group.",
						},
						"cool_down_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The cooling duration of the AS group.",
						},
						"lbaas_listeners": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        groupDataSourceLBaasListenersSchema(),
							Description: "The enhanced load balancers of the AS group.",
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The AZ information of the AS group.",
						},
						"networks": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        groupDataSourceNetworksSchema(),
							Description: "The network information of the AS group.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        groupDataSourceSecurityGroupsSchema(),
							Description: "The security group information of the AS group.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when an AS group was created.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VPC to which the AS group belongs.",
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Details about the AS group.",
						},
						"is_scaling": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The scaling flag of the AS group.",
						},
						"health_periodic_audit_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health check method.",
						},
						"health_periodic_audit_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check interval.",
						},
						"health_periodic_audit_grace_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The grace period for health check.",
						},
						"instance_terminate_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance removal policy.",
						},
						"delete_publicip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether to delete the EIP bound to the ECS when deleting the ECS.",
						},
						"delete_volume": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: "Specifies whether to delete the data disks attached to the ECS when " +
								"deleting the ECS.",
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise project ID of the AS group.",
						},
						"activity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the AS action.",
						},
						"multi_az_scaling_policy": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "The priority policy used to select target AZs when adjusting the number of" +
								" instances in an AS group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the AS group.",
						},
						"iam_agency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The agency name.",
						},
						"tags": common.TagsComputedSchema(),
						"instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The scaling group instances ids.",
						},
					},
				},
				Description: "A list of AS groups",
			},
		},
	}
}

func groupDataSourceLBaasListenersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backend ECS group ID.",
			},
			"protocol_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The backend protocol ID.",
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The weight, which determines the portion of requests a backend " +
					"ECS processes compared to other backend ECSs added to the same listener.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the listener associate with the ELB.",
			},
			"protocol_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of IP addresses of backend servers to be bound with the ELB.",
			},
		},
	}
}

func groupDataSourceNetworksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet ID.",
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether to support IPv6 addresses.",
			},
			"ipv6_bandwidth_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the shared bandwidth of an IPv6 address.",
			},
			"source_dest_check": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func groupDataSourceSecurityGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the security group.",
			},
		},
	}
}

func flattenDataSourceLBaaSListeners(listeners []groups.LBaaSListener) []map[string]interface{} {
	if len(listeners) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(listeners))
	for i, item := range listeners {
		res[i] = map[string]interface{}{
			"pool_id":          item.PoolID,
			"protocol_port":    item.ProtocolPort,
			"weight":           item.Weight,
			"protocol_version": item.ProtocolVersion,
			"listener_id":      item.ListenerID,
		}
	}
	return res
}

func flattenDataSourceNetworks(networks []groups.Network) []map[string]interface{} {
	res := make([]map[string]interface{}, len(networks))
	for i, item := range networks {
		res[i] = map[string]interface{}{
			"id":                item.ID,
			"ipv6_enable":       item.IPv6Enable,
			"ipv6_bandwidth_id": item.IPv6BandWidth.ID,
			"source_dest_check": len(item.AllowedAddressPairs) == 0,
		}
	}
	return res
}

func flattenDataSourceSecurityGroups(sgs []groups.SecurityGroup) []map[string]interface{} {
	res := make([]map[string]interface{}, len(sgs))
	for i, item := range sgs {
		res[i] = map[string]interface{}{
			"id": item.ID,
		}
	}
	return res
}

// getDataSourceInstancesIDs using to collecting total instance IDs in AS group.
// When the query API reports an error, only the failure log is printed and the program is not terminated.
func getDataSourceInstancesIDs(asClient *golangsdk.ServiceClient, groupID string) []string {
	allIns, err := getAllInstancesInGroup(asClient, groupID)
	if err != nil {
		log.Printf("[WARN] Error fetching instances in AS group (%s): %s", groupID, err)
		return nil
	}

	allIDs := make([]string, 0, len(allIns))
	for _, ins := range allIns {
		if instanceID := utils.PathSearch("instance_id", ins, "").(string); instanceID != "" {
			allIDs = append(allIDs, instanceID)
		}
	}

	return allIDs
}

func getDataSourceGroupTags(asClient *golangsdk.ServiceClient, groupID string) map[string]string {
	resourceTags, err := tags.Get(asClient, groupID).Extract()
	if err == nil {
		tagMap := make(map[string]string)
		for _, val := range resourceTags.Tags {
			tagMap[val.Key] = val.Value
		}
		return tagMap
	}
	log.Printf("[WARN] Error fetching tags of AS group (%s): %s", groupID, err)
	return nil
}

func flattenDataSourceGroups(asClient *golangsdk.ServiceClient, groupList []groups.Group) ([]string, []map[string]interface{}) {
	ids := make([]string, 0, len(groupList))
	elements := make([]map[string]interface{}, 0, len(groupList))
	for _, group := range groupList {
		groupID := group.ID

		groupMap := map[string]interface{}{
			"scaling_group_name":                 group.Name,
			"scaling_group_id":                   groupID,
			"status":                             group.Status,
			"scaling_configuration_id":           group.ConfigurationID,
			"scaling_configuration_name":         group.ConfigurationName,
			"current_instance_number":            group.ActualInstanceNumber,
			"desire_instance_number":             group.DesireInstanceNumber,
			"min_instance_number":                group.MinInstanceNumber,
			"max_instance_number":                group.MaxInstanceNumber,
			"cool_down_time":                     group.CoolDownTime,
			"lbaas_listeners":                    flattenDataSourceLBaaSListeners(group.LBaaSListeners),
			"availability_zones":                 group.AvailableZones,
			"networks":                           flattenDataSourceNetworks(group.Networks),
			"security_groups":                    flattenDataSourceSecurityGroups(group.SecurityGroups),
			"created_at":                         group.CreateTime,
			"vpc_id":                             group.VpcID,
			"detail":                             group.Detail,
			"is_scaling":                         group.IsScaling,
			"health_periodic_audit_method":       group.HealthPeriodicAuditMethod,
			"health_periodic_audit_time":         group.HealthPeriodicAuditTime,
			"health_periodic_audit_grace_period": group.HealthPeriodicAuditGrace,
			"instance_terminate_policy":          group.InstanceTerminatePolicy,
			"delete_publicip":                    group.DeletePublicip,
			"delete_volume":                      group.DeleteVolume,
			"enterprise_project_id":              group.EnterpriseProjectID,
			"activity_type":                      group.ActivityType,
			"multi_az_scaling_policy":            group.MultiAZPriorityPolicy,
			"description":                        group.Description,
			"iam_agency_name":                    group.IamAgencyName,
			"instances":                          getDataSourceInstancesIDs(asClient, groupID),
			"tags":                               getDataSourceGroupTags(asClient, groupID),
		}
		elements = append(elements, groupMap)
		ids = append(ids, groupID)
	}
	return ids, elements
}

func buildDataSourceGroupOpts(d *schema.ResourceData) groups.ListOpts {
	return groups.ListOpts{
		Name:                d.Get("name").(string),
		ConfigurationID:     d.Get("scaling_configuration_id").(string),
		Status:              d.Get("status").(string),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
	}
}

func dataSourceASGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf   = meta.(*config.Config)
		region = conf.GetRegion(d)
		opts   = buildDataSourceGroupOpts(d)
	)

	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	pages, err := groups.List(asClient, opts).AllPages()
	if err != nil {
		return diag.Errorf("error getting AS group list: %s", err)
	}
	groupList, err := pages.(groups.GroupPage).Extract()
	if err != nil {
		return diag.Errorf("error extract to AS group list: %s", err)
	}

	ids, elements := flattenDataSourceGroups(asClient, groupList)
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("groups", elements),
		d.Set("region", region),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting AS group fields: %s", mErr)
	}
	return nil
}
