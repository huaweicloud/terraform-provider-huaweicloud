package huaweicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/instances"
)

func dataSourceGaussDBMysqlInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGaussDBMysqlInstancesRead,

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
						"private_write_ip": {
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
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.gaussdbV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	listOpts := instances.ListTaurusDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return err
	}

	allInstances, err := instances.ExtractTaurusDBInstances(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve instances: %s", err)
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
			return err
		}
		log.Printf("[DEBUG] Retrieved Instance %s: %+v", instance.Id, instance)

		instanceToSet["configuration_id"] = instance.ConfigurationId
		instanceToSet["availability_zone_mode"] = instance.AZMode
		instanceToSet["master_availability_zone"] = instance.MasterAZ

		if len(instance.PrivateIps) > 0 {
			instanceToSet["private_write_ip"] = instance.PrivateIps[0]
		}

		flavor := ""
		slave_count := 0
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
				slave_count += 1
			}
			if flavor == "" {
				flavor = raw.Flavor
			}
		}

		instanceToSet["nodes"] = nodesList
		instanceToSet["read_replicas"] = slave_count
		if flavor != "" {
			log.Printf("[DEBUG] Node Flavor: %s", flavor)
			instanceToSet["flavor"] = flavor
		}

		instancesToSet = append(instancesToSet, instanceToSet)
	}

	d.SetId(hashcode.Strings(instancesIds))
	d.Set("instances", instancesToSet)

	return nil
}
