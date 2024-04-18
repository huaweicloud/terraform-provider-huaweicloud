package gaussdb

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances
func DataSourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenGaussInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
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
	}
}

func dataSourceOpenGaussInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	allInstances, err := instances.ExtractGaussDBInstances(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve instances: %s", err)
	}

	if allInstances.TotalCount < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if allInstances.TotalCount > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	instance := allInstances.Instances[0]

	log.Printf("[DEBUG] retrieved instance %s: %+v", instance.Id, instance)
	d.SetId(instance.Id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("type", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("enterprise_project_id", instance.EnterpriseProjectId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("flavor", instance.FlavorRef),
		d.Set("port", instance.Port),
		d.Set("switch_strategy", instance.SwitchStrategy),
		d.Set("maintenance_window", instance.MaintenanceWindow),
		d.Set("public_ips", instance.PublicIps),
	)

	if len(instance.PrivateIps) > 0 {
		privateIps := instance.PrivateIps[0]
		ipList := strings.Split(privateIps, "/")
		for i := 0; i < len(ipList); i++ {
			ipList[i] = strings.Trim(ipList[i], " ")
		}
		d.Set("private_ips", ipList)
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
	var dnNum int
	shardingNum := 0
	coordinatorNum := 0
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
			coordinatorNum++
		} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
			shardingNum++
		}
	}

	if shardingNum > 0 && coordinatorNum > 0 {
		dnNum = shardingNum / instance.ReplicaNum
		d.Set("nodes", nodesList)
		d.Set("sharding_num", dnNum)
		d.Set("coordinator_num", coordinatorNum)
	} else {
		// If the HA mode is centralized, the HA structure of API response is nil.
		dnNum = instance.ReplicaNum + 1
		d.Set("nodes", nodesList)
		d.Set("replica_num", instance.ReplicaNum)
	}

	// remove duplicate az
	azList = utils.RemoveDuplicateElem(azList)
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
	volumeList := make([]map[string]interface{}, 1)
	volume := map[string]interface{}{
		"type": instance.Volume.Type,
		"size": instance.Volume.Size / dnNum,
	}
	volumeList[0] = volume
	d.Set("volume", volumeList)

	return diag.FromErr(mErr.ErrorOrNil())
}
