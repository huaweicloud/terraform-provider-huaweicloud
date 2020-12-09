package huaweicloud

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/opengauss/v3/instances"
)

func dataSourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpenGaussInstanceRead,

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
	}
}

func dataSourceOpenGaussInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.openGaussV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	listOpts := instances.ListGaussDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return err
	}

	allInstances, err := instances.ExtractGaussDBInstances(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve instances: %s", err)
	}

	if allInstances.TotalCount < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if allInstances.TotalCount > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	instance := allInstances.Instances[0]

	log.Printf("[DEBUG] Retrieved Instance %s: %+v", instance.Id, instance)
	d.SetId(instance.Id)

	d.Set("region", region)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)
	d.Set("db_user_name", instance.DbUserName)
	d.Set("time_zone", instance.TimeZone)
	d.Set("flavor", instance.FlavorRef)
	d.Set("port", instance.Port)
	d.Set("switch_strategy", instance.SwitchStrategy)
	d.Set("maintenance_window", instance.MaintenanceWindow)

	if len(instance.PrivateIps) > 0 {
		private_ips := instance.PrivateIps[0]
		ip_list := strings.Split(private_ips, "/")
		for i := 0; i < len(ip_list); i++ {
			ip_list[i] = strings.Trim(ip_list[i], " ")
		}
		d.Set("private_ips", ip_list)
	}

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": instance.DataStore.Version,
		"engine":  instance.DataStore.Type,
	}
	dbList[0] = db
	d.Set("datastore", dbList)

	// set nodes
	sharding_num := 0
	coordinator_num := 0
	azList := []string{}
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
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
			coordinator_num += 1
		} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
			sharding_num += 1
		}
	}
	d.Set("nodes", nodesList)
	d.Set("coordinator_num", coordinator_num)

	dn_num := sharding_num / 3
	d.Set("sharding_num", dn_num)

	//remove duplicate az
	azList = removeDuplicateElem(azList)
	sort.Strings(azList)
	d.Set("availability_zone", strings.Join(azList, ","))

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)

	// set ha
	haList := make([]map[string]interface{}, 1)
	ha := map[string]interface{}{
		"replication_mode": instance.Ha.ReplicationMode,
	}
	haList[0] = ha
	d.Set("ha", haList)

	// set volume
	volume_size := instance.Volume.Size
	dn_size := volume_size / dn_num
	volumeList := make([]map[string]interface{}, 1)
	volume := map[string]interface{}{
		"type": instance.Volume.Type,
		"size": dn_size,
	}
	volumeList[0] = volume
	d.Set("volume", volumeList)

	return nil
}
