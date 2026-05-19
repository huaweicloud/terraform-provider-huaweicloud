package taurusdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3.1/{project_id}/instances
// @API TaurusDB GET /v3.1/{project_id}/instances/{instance_id}
func DataSourceTaurusDBInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBInstanceRead,

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

func dataSourceTaurusDBInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listBody := map[string]interface{}{
		"name":      d.Get("name"),
		"vpc_id":    d.Get("vpc_id"),
		"subnet_id": d.Get("subnet_id"),
	}

	allInstances, err := getInstancesList(client, utils.RemoveNil(listBody))
	if err != nil {
		return diag.Errorf("error retrieving instances: %s", err)
	}

	if len(allInstances) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if len(allInstances) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	instanceInList := allInstances[0].(map[string]interface{})
	instanceID := instanceInList["id"].(string)
	instanceDetail, err := GetInstanceDetail(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	instance := utils.PathSearch("instance", instanceDetail, nil)

	d.SetId(instanceID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("mode", utils.PathSearch("type", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("configuration_id", utils.PathSearch("configuration_id", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("db_user_name", utils.PathSearch("db_user_name", instance, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", instance, nil)),
		d.Set("availability_zone_mode", utils.PathSearch("az_mode", instance, nil)),
		d.Set("master_availability_zone", utils.PathSearch("master_az_code", instance, nil)),
	)

	// Set port number with safe type assertion
	if portRaw := utils.PathSearch("port", instanceInList, nil); portRaw != nil {
		if portStr, ok := portRaw.(string); ok && portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				mErr = multierror.Append(mErr, d.Set("port", port))
			}
		}
	}

	// Set private_write_ip with safe type assertion
	if privateWriteIpsRaw := utils.PathSearch("private_write_ips", instance, nil); privateWriteIpsRaw != nil {
		if privateWriteIps, ok := privateWriteIpsRaw.([]interface{}); ok && len(privateWriteIps) > 0 {
			mErr = multierror.Append(mErr, d.Set("private_write_ip", privateWriteIps[0]))
		}
	}

	// Set datastore with safe type assertion
	var engine, version string
	if engineRaw := utils.PathSearch("datastore.type", instance, nil); engineRaw != nil {
		if e, ok := engineRaw.(string); ok {
			engine = e
		}
	}
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	if versionRaw := utils.PathSearch("datastore.version", instance, nil); versionRaw != nil {
		if v, ok := versionRaw.(string); ok {
			version = v
		}
	}
	datastore := []map[string]interface{}{
		{
			"engine":  engine,
			"version": version,
		},
	}
	mErr = multierror.Append(mErr, d.Set("datastore", datastore))

	// Set nodes, read_replicas and flavor
	nodesList, slaveCount, flavor := extractNodesInfo(instance)

	mErr = multierror.Append(mErr,
		d.Set("nodes", nodesList),
		d.Set("read_replicas", slaveCount),
	)

	if flavor != "" {
		mErr = multierror.Append(mErr, d.Set("flavor", flavor))
	}

	// Set backup_strategy
	backupStrategy := extractBackupStrategy(instance)
	if backupStrategy != nil {
		mErr = multierror.Append(mErr, d.Set("backup_strategy", []map[string]interface{}{backupStrategy}))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

// GetInstanceDetail gets the details of a TaurusDB instance
func GetInstanceDetail(client *golangsdk.ServiceClient, instanceID string) (interface{}, error) {
	httpUrl := "v3.1/{project_id}/instances/{instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error gettting TaurusDB instance (%s): %s", instanceID, err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving TaurusDB instance (%s): %s", instanceID, err)
	}
	return utils.PathSearch("instance", respBody, nil), nil
}

// extractBackupStrategy extracts backup strategy information from instance data
func extractBackupStrategy(instance interface{}) map[string]interface{} {
	backupStrategyRaw := utils.PathSearch("backup_strategy", instance, nil)
	if backupStrategyRaw == nil {
		return nil
	}

	var startTime string
	if startTimeRaw := utils.PathSearch("start_time", backupStrategyRaw, nil); startTimeRaw != nil {
		if s, ok := startTimeRaw.(string); ok {
			startTime = s
		}
	}

	backupStrategyMap := map[string]interface{}{
		"start_time": startTime,
	}

	extractKeepDays(backupStrategyRaw, backupStrategyMap)

	return backupStrategyMap
}

// extractKeepDays extracts keep_days from backup strategy
func extractKeepDays(backupStrategyRaw interface{}, backupStrategyMap map[string]interface{}) {
	if keepDaysRaw := utils.PathSearch("keep_days", backupStrategyRaw, nil); keepDaysRaw != nil {
		if keepDaysStr, ok := keepDaysRaw.(string); ok && keepDaysStr != "" {
			if keepDays, err := strconv.Atoi(keepDaysStr); err == nil {
				backupStrategyMap["keep_days"] = keepDays
			}
		}
	}
}
