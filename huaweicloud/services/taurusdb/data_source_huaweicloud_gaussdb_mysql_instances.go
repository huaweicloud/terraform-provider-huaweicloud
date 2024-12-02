package taurusdb

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API GaussDBforMySQL GET /v3/{project_id}/instances
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
func DataSourceGaussDBMysqlInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlInstancesRead,

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
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration_id": {
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
						"availability_zone_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_dns_name_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_write_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_begin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"read_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
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
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_read_ip": {
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

func dataSourceGaussDBMysqlInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GaussdbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listOpts := instances.ListTaurusDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allInstances, err := instances.ExtractTaurusDBInstances(pages)

	if err != nil {
		return diag.Errorf("unable to retrieve instances: %s", err)
	}

	var instancesToSet []map[string]interface{}
	var instancesIds []string

	for _, instanceInAll := range allInstances.Instances {
		instanceToSet := map[string]interface{}{
			"id":                    instanceInAll.Id,
			"region":                region,
			"name":                  instanceInAll.Name,
			"status":                instanceInAll.Status,
			"mode":                  instanceInAll.Type,
			"vpc_id":                instanceInAll.VpcId,
			"subnet_id":             instanceInAll.SubnetId,
			"security_group_id":     instanceInAll.SecurityGroupId,
			"enterprise_project_id": instanceInAll.EnterpriseProjectId,
			"db_user_name":          instanceInAll.DbUserName,
			"time_zone":             instanceInAll.TimeZone,
		}

		if dbPort, err := strconv.Atoi(instanceInAll.Port); err == nil {
			instanceToSet["port"] = dbPort
		}

		// set data store
		dbList := make([]map[string]interface{}, 1)
		db := map[string]interface{}{
			"version": instanceInAll.DataStore.Version,
		}
		// normalize engine
		engine := instanceInAll.DataStore.Type
		if engine == "GaussDB(for MySQL)" {
			engine = "gaussdb-mysql"
		}
		db["engine"] = engine
		dbList[0] = db
		instanceToSet["datastore"] = dbList

		// set backup_strategy
		backupStrategyList := make([]map[string]interface{}, 1)
		backupStrategy := map[string]interface{}{
			"start_time": instanceInAll.BackupStrategy.StartTime,
		}
		if days, err := strconv.Atoi(instanceInAll.BackupStrategy.KeepDays); err == nil {
			backupStrategy["keep_days"] = days
		}
		backupStrategyList[0] = backupStrategy
		instanceToSet["backup_strategy"] = backupStrategyList

		// set nodes, configuration_id, availability_zone_mode, master_availability_zone, private_write_ip
		instanceID := instanceInAll.Id
		instancesIds = append(instancesIds, instanceID)

		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] retrieved instance %s: %+v", instance.Id, instance)

		instanceToSet["configuration_id"] = instance.ConfigurationId
		instanceToSet["availability_zone_mode"] = instance.AZMode
		instanceToSet["master_availability_zone"] = instance.MasterAZ
		instanceToSet["description"] = instance.Alias
		instanceToSet["created_at"] = instance.Created
		instanceToSet["updated_at"] = instance.Updated

		if len(instance.PrivateIps) > 0 {
			instanceToSet["private_write_ip"] = instance.PrivateIps[0]
		}
		if len(instance.PrivateDnsNames) > 0 {
			instanceToSet["private_dns_name_prefix"] = strings.Split(instance.PrivateDnsNames[0], ".")[0]
			instanceToSet["private_dns_name"] = instance.PrivateDnsNames[0]
		}

		maintainWindow := strings.Split(instance.MaintenanceWindow, "-")
		if len(maintainWindow) == 2 {
			instanceToSet["maintain_begin"] = maintainWindow[0]
			instanceToSet["maintain_end"] = maintainWindow[1]
		}

		flavor := ""
		slaveCount := 0
		nodesList := make([]map[string]interface{}, 0, 1)
		for _, raw := range instance.Nodes {
			node := map[string]interface{}{
				"id":                raw.Id,
				"name":              raw.Name,
				"status":            raw.Status,
				"type":              raw.Type,
				"availability_zone": raw.AvailabilityZone,
			}
			if len(raw.PrivateIps) > 0 {
				node["private_read_ip"] = raw.PrivateIps[0]
			}
			nodesList = append(nodesList, node)
			if raw.Type == "slave" && raw.Status == "ACTIVE" {
				slaveCount++
			}
			if flavor == "" {
				flavor = raw.Flavor
			}
		}

		instanceToSet["nodes"] = nodesList
		instanceToSet["read_replicas"] = slaveCount
		if flavor != "" {
			log.Printf("[DEBUG] node flavor: %s", flavor)
			instanceToSet["flavor"] = flavor
		}

		instancesToSet = append(instancesToSet, instanceToSet)
	}

	d.SetId(hashcode.Strings(instancesIds))
	var mErr *multierror.Error
	mErr = multierror.Append(mErr,
		d.Set("instances", instancesToSet),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
