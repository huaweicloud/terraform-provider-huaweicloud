package taurusdb

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDBforMySQL GET /v3/{project_id}/instances
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
func DataSourceGaussDBMysqlInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlInstanceRead,

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
	}
}

func dataSourceGaussDBMysqlInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if allInstances.TotalCount < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if allInstances.TotalCount > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	instanceID := allInstances.Instances[0].Id
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] retrieved instance %s: %+v", instance.Id, instance)
	d.SetId(instance.Id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("mode", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("configuration_id", instance.ConfigurationId),
		d.Set("enterprise_project_id", instance.EnterpriseProjectId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("availability_zone_mode", instance.AZMode),
		d.Set("master_availability_zone", instance.MasterAZ),
	)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}
	if len(instance.PrivateIps) > 0 {
		d.Set("private_write_ip", instance.PrivateIps[0])
	}

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": instance.DataStore.Version,
	}
	// normalize engine
	engine := instance.DataStore.Type
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	db["engine"] = engine
	dbList[0] = db
	d.Set("datastore", dbList)

	// set nodes
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
	d.Set("nodes", nodesList)
	d.Set("read_replicas", slaveCount)
	if flavor != "" {
		log.Printf("[DEBUG] node flavor: %s", flavor)
		d.Set("flavor", flavor)
	}

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
	}
	if days, err := strconv.Atoi(instance.BackupStrategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)

	return diag.FromErr(mErr.ErrorOrNil())
}
