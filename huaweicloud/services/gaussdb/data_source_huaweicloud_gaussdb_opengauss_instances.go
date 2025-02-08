package gaussdb

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances
func DataSourceOpenGaussInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenGaussInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"switch_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintenance_window": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"backup_strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keep_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ha": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replication_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"coordinator_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sharding_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replica_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mysql_compatibility_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
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
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"public_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"nodes": {
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
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"availability_zone": {
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

func dataSourceOpenGaussInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.OpenGaussV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listOpts := instances.ListGaussDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to list instances: %s", err)
	}

	allInstances, err := instances.ExtractGaussDBInstances(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve instances: %s", err)
	}

	var instancesToSet []map[string]interface{}
	var instancesIds []string

	for _, instanceInAll := range allInstances.Instances {
		instanceToSet := map[string]interface{}{
			"id":                       instanceInAll.Id,
			"region":                   region,
			"name":                     instanceInAll.Name,
			"status":                   instanceInAll.Status,
			"type":                     instanceInAll.Type,
			"vpc_id":                   instanceInAll.VpcId,
			"subnet_id":                instanceInAll.SubnetId,
			"security_group_id":        instanceInAll.SecurityGroupId,
			"enterprise_project_id":    instanceInAll.EnterpriseProjectId,
			"db_user_name":             instanceInAll.DbUserName,
			"time_zone":                instanceInAll.TimeZone,
			"flavor":                   instanceInAll.FlavorRef,
			"port":                     instanceInAll.Port,
			"switch_strategy":          instanceInAll.SwitchStrategy,
			"maintenance_window":       instanceInAll.MaintenanceWindow,
			"public_ips":               instanceInAll.PublicIps,
			"mysql_compatibility_port": instanceInAll.MysqlCompatibility.Port,
		}

		instanceID := instanceInAll.Id
		instancesIds = append(instancesIds, instanceID)

		if len(instanceInAll.PrivateIps) > 0 {
			privateIps := instanceInAll.PrivateIps[0]
			ipList := strings.Split(privateIps, "/")
			for i := 0; i < len(ipList); i++ {
				ipList[i] = strings.Trim(ipList[i], " ")
			}
			instanceToSet["private_ips"] = ipList
		}

		// set data store
		dbList := make([]map[string]interface{}, 1)
		db := map[string]interface{}{
			"version": instanceInAll.DataStore.Version,
			"engine":  instanceInAll.DataStore.Type,
		}
		dbList[0] = db
		instanceToSet["datastore"] = dbList

		// set nodes
		var dnNum int
		shardingNum := 0
		coordinatorNum := 0
		azList := []string{}
		nodesList := make([]map[string]interface{}, 0, 1)
		for _, raw := range instanceInAll.Nodes {
			node := map[string]interface{}{
				"id":                raw.Id,
				"name":              raw.Name,
				"status":            raw.Status,
				"role":              raw.Role,
				"availability_zone": raw.AvailabilityZone,
			}
			nodesList = append(nodesList, node)
			azList = append(azList, raw.AvailabilityZone)
			if strings.Contains(raw.Name, "_gaussdbv5cn") {
				coordinatorNum++
			} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
				shardingNum++
			}
		}

		if shardingNum > 0 && coordinatorNum > 0 {
			dnNum = shardingNum / instanceInAll.ReplicaNum
			instanceToSet["nodes"] = nodesList
			instanceToSet["coordinator_num"] = coordinatorNum
			instanceToSet["sharding_num"] = dnNum
		} else {
			// If the HA mode is centralized, the HA structure of API response is nil.
			dnNum = instanceInAll.ReplicaNum + 1
			instanceToSet["nodes"] = nodesList
			instanceToSet["replica_num"] = instanceInAll.ReplicaNum
		}

		// remove duplicate az
		azList = utils.RemoveDuplicateElem(azList)
		sort.Strings(azList)
		instanceToSet["availability_zone"] = strings.Join(azList, ",")

		// set backup_strategy
		backupStrategyList := make([]map[string]interface{}, 1)
		backupStrategy := map[string]interface{}{
			"start_time": instanceInAll.BackupStrategy.StartTime,
			"keep_days":  instanceInAll.BackupStrategy.KeepDays,
		}
		backupStrategyList[0] = backupStrategy
		instanceToSet["backup_strategy"] = backupStrategyList

		// set ha
		haList := make([]map[string]interface{}, 1)
		ha := map[string]interface{}{
			"replication_mode": instanceInAll.Ha.ReplicationMode,
		}
		haList[0] = ha
		instanceToSet["ha"] = haList

		// set volume
		volumeList := make([]map[string]interface{}, 1)
		volume := map[string]interface{}{
			"type": instanceInAll.Volume.Type,
			"size": instanceInAll.Volume.Size / dnNum,
		}
		volumeList[0] = volume
		instanceToSet["volume"] = volumeList

		instancesToSet = append(instancesToSet, instanceToSet)
	}

	d.SetId(hashcode.Strings(instancesIds))
	var mErr *multierror.Error
	mErr = multierror.Append(mErr,
		d.Set("instances", instancesToSet),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
