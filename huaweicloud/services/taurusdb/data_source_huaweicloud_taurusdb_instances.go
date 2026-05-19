package taurusdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3.1/{project_id}/instances
// @API TaurusDB GET /v3.1/{project_id}/instances/{instance_id}
func DataSourceTaurusDBInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBInstancesRead,

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

func dataSourceTaurusDBInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	params := map[string]interface{}{
		"name":      d.Get("name"),
		"vpc_id":    d.Get("vpc_id"),
		"subnet_id": d.Get("subnet_id"),
	}

	allInstances, err := getInstancesList(client, utils.RemoveNil(params))
	if err != nil {
		return diag.Errorf("error retrieving instances: %s", err)
	}

	instancesToSet, err := extractInstances(client, allInstances)
	if err != nil {
		return diag.Errorf("error extracting instances: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", instancesToSet),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstancesList(client *golangsdk.ServiceClient, params map[string]interface{}) ([]interface{}, error) {
	var (
		allInstances = make([]interface{}, 0)
		httpUrl      = "v3.1/{project_id}/instances"
		offset       = 0
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = buildListInstantsQueryParams(listPath, params)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		// Add offset parameter to URL
		fullPath := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", fullPath, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("error getting TaurusDB instances: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("error retrieving TaurusDB instances: %s", err)
		}

		pageInstances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(pageInstances) == 0 {
			break
		}

		allInstances = append(allInstances, pageInstances...)

		totalCount := utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if int(totalCount) == len(allInstances) {
			break
		}

		offset += len(pageInstances)
	}
	return allInstances, nil
}

func buildListInstantsQueryParams(baseUrl string, params map[string]interface{}) string {
	res := fmt.Sprintf("%s?limit=100", baseUrl)
	for k, v := range params {
		if v != nil {
			if vStr := v.(string); vStr != "" {
				res += fmt.Sprintf("&%s=%s", k, vStr)
			}
		}
	}
	return res
}

func extractInstances(client *golangsdk.ServiceClient, allInstances []interface{}) ([]map[string]interface{}, error) {
	var instancesToSet []map[string]interface{}
	if len(allInstances) == 0 {
		return instancesToSet, nil
	}
	for _, instanceInList := range allInstances {
		instanceID := utils.PathSearch("id", instanceInList, nil).(string)
		instanceDetail, err := GetInstanceDetail(client, instanceID)
		if err != nil {
			return nil, err
		}
		extractedInstance := extractInstance(instanceInList, instanceDetail)

		instancesToSet = append(instancesToSet, extractedInstance)
	}
	return instancesToSet, nil
}

func extractInstance(instanceInList interface{}, instanceDetail interface{}) map[string]interface{} {
	instanceToSet := make(map[string]interface{})

	// Extract basic information from instance (list data)
	extractBasicInfo(instanceInList, instanceToSet)

	// Extract detailed information from instanceDetail
	extractDetailedInfo(instanceDetail, instanceToSet)

	// Set nodes, read_replicas and flavor
	nodesList, slaveCount, flavor := extractNodesInfo(instanceDetail)

	instanceToSet["nodes"] = nodesList
	instanceToSet["read_replicas"] = slaveCount
	if flavor != "" {
		instanceToSet["flavor"] = flavor
	}

	return instanceToSet
}

// extractBasicInfo extracts basic instance information from list data
func extractBasicInfo(instanceInList interface{}, instanceToSet map[string]interface{}) {
	instanceToSet["id"] = utils.PathSearch("id", instanceInList, nil)
	instanceToSet["region"] = utils.PathSearch("region", instanceInList, nil)
	instanceToSet["name"] = utils.PathSearch("name", instanceInList, nil)
	instanceToSet["status"] = utils.PathSearch("status", instanceInList, nil)
	instanceToSet["vpc_id"] = utils.PathSearch("vpc_id", instanceInList, nil)
	instanceToSet["subnet_id"] = utils.PathSearch("subnet_id", instanceInList, nil)
	instanceToSet["security_group_id"] = utils.PathSearch("security_group_id", instanceInList, nil)
	instanceToSet["enterprise_project_id"] = utils.PathSearch("enterprise_project_id", instanceInList, nil)
	instanceToSet["db_user_name"] = utils.PathSearch("db_user_name", instanceInList, nil)
	instanceToSet["time_zone"] = utils.PathSearch("time_zone", instanceInList, nil)
	instanceToSet["mode"] = utils.PathSearch("type", instanceInList, nil)

	// Handle port number with safe type assertion
	if portRaw := utils.PathSearch("port", instanceInList, nil); portRaw != nil {
		if portStr, ok := portRaw.(string); ok && portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				instanceToSet["port"] = port
			}
		}
	}
}

// extractDetailedInfo extracts detailed instance information from instance detail
func extractDetailedInfo(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	instanceToSet["description"] = utils.PathSearch("alias", instanceDetail, nil)
	instanceToSet["configuration_id"] = utils.PathSearch("configuration_id", instanceDetail, nil)
	instanceToSet["availability_zone_mode"] = utils.PathSearch("az_mode", instanceDetail, nil)
	instanceToSet["master_availability_zone"] = utils.PathSearch("master_az_code", instanceDetail, nil)
	instanceToSet["created_at"] = utils.PathSearch("created", instanceDetail, nil)
	instanceToSet["updated_at"] = utils.PathSearch("updated", instanceDetail, nil)

	// Set private_write_ip with safe type assertion
	setPrivateWriteIP(instanceDetail, instanceToSet)

	// Set private_dns_name and private_dns_name_prefix with safe type assertion
	setPrivateDNSNames(instanceDetail, instanceToSet)

	// Set maintain_begin and maintain_end with safe type assertion
	setMaintainWindow(instanceDetail, instanceToSet)

	// Set datastore with safe type assertion
	extractDatastore(instanceDetail, instanceToSet)

	// Set backup_strategy with safe type assertion
	extractBackupStrategyForList(instanceDetail, instanceToSet)
}

// setPrivateWriteIP sets the private write IP address
func setPrivateWriteIP(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	if privateWriteIpsRaw := utils.PathSearch("private_write_ips", instanceDetail, nil); privateWriteIpsRaw != nil {
		if privateWriteIps, ok := privateWriteIpsRaw.([]interface{}); ok && len(privateWriteIps) > 0 {
			instanceToSet["private_write_ip"] = privateWriteIps[0]
		}
	}
}

// setPrivateDNSNames sets the private DNS names
func setPrivateDNSNames(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	if privateDNSNamesRaw := utils.PathSearch("private_dns_names", instanceDetail, nil); privateDNSNamesRaw != nil {
		if privateDNSNames, ok := privateDNSNamesRaw.([]interface{}); ok && len(privateDNSNames) > 0 {
			if dnsName, ok := privateDNSNames[0].(string); ok && dnsName != "" {
				instanceToSet["private_dns_name"] = dnsName
				if parts := strings.Split(dnsName, "."); len(parts) > 0 && parts[0] != "" {
					instanceToSet["private_dns_name_prefix"] = parts[0]
				}
			}
		}
	}
}

// setMaintainWindow sets the maintenance window times
func setMaintainWindow(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	if maintainWindowRaw := utils.PathSearch("maintain_windows", instanceDetail, nil); maintainWindowRaw != nil {
		if maintainWindow, ok := maintainWindowRaw.(string); ok && maintainWindow != "" {
			if parts := strings.Split(maintainWindow, "-"); len(parts) == 2 {
				instanceToSet["maintain_begin"] = parts[0]
				instanceToSet["maintain_end"] = parts[1]
			}
		}
	}
}

// extractDatastore sets the datastore information
func extractDatastore(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	var engine, version string
	if engineRaw := utils.PathSearch("datastore.type", instanceDetail, nil); engineRaw != nil {
		if e, ok := engineRaw.(string); ok {
			engine = e
		}
	}
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	if versionRaw := utils.PathSearch("datastore.version", instanceDetail, nil); versionRaw != nil {
		if v, ok := versionRaw.(string); ok {
			version = v
		}
	}
	instanceToSet["datastore"] = []map[string]interface{}{
		{
			"engine":  engine,
			"version": version,
		},
	}
}

// extractBackupStrategyForList sets the backup strategy information
func extractBackupStrategyForList(instanceDetail interface{}, instanceToSet map[string]interface{}) {
	if backupStrategy := utils.PathSearch("backup_strategy", instanceDetail, nil); backupStrategy != nil {
		var startTime, keepDaysStr string

		if startTimeRaw := utils.PathSearch("start_time", backupStrategy, nil); startTimeRaw != nil {
			if s, ok := startTimeRaw.(string); ok {
				startTime = s
			}
		}

		if keepDaysRaw := utils.PathSearch("keep_days", backupStrategy, nil); keepDaysRaw != nil {
			if k, ok := keepDaysRaw.(string); ok {
				keepDaysStr = k
			}
		}

		backupStrategyMap := map[string]interface{}{
			"start_time": startTime,
		}

		if keepDaysStr != "" {
			if keepDays, err := strconv.Atoi(keepDaysStr); err == nil {
				backupStrategyMap["keep_days"] = keepDays
			}
		}

		instanceToSet["backup_strategy"] = []map[string]interface{}{backupStrategyMap}
	}
}

// extractNodesInfo extracts node information, slave count, and flavor from instance data
func extractNodesInfo(instanceDetail interface{}) ([]map[string]interface{}, int, string) {
	nodesList := make([]map[string]interface{}, 0)
	slaveCount := 0
	flavor := ""

	nodesRaw := utils.PathSearch("nodes", instanceDetail, nil)
	if nodesRaw == nil {
		return nodesList, slaveCount, flavor
	}

	nodes, ok := nodesRaw.([]interface{})
	if !ok || len(nodes) == 0 {
		return nodesList, slaveCount, flavor
	}

	for _, nodeRaw := range nodes {
		node, nodeType, nodeStatus, nodeFlavor := buildNodeInfo(nodeRaw)
		nodesList = append(nodesList, node)

		// Count slave nodes
		if nodeType == "slave" && nodeStatus == "normal" {
			slaveCount++
		}

		// Get flavor from nodes
		if flavor == "" && nodeFlavor != "" {
			flavor = nodeFlavor
		}
	}

	return nodesList, slaveCount, flavor
}

// buildNodeInfo builds a node map and returns node type and status for counting
func buildNodeInfo(nodeRaw interface{}) (node map[string]interface{}, nodeType string, nodeStatus string, nodeFlavor string) {
	// Cache frequently accessed values to avoid repeated PathSearch calls
	nodeID := utils.PathSearch("id", nodeRaw, nil)
	nodeName := utils.PathSearch("name", nodeRaw, nil)
	nodeTypeRaw := utils.PathSearch("type", nodeRaw, nil)
	nodeStatusRaw := utils.PathSearch("status", nodeRaw, nil)
	nodeAZ := utils.PathSearch("az_code", nodeRaw, nil)
	nodeFlavorRaw := utils.PathSearch("flavor_ref", nodeRaw, nil)

	node = map[string]interface{}{
		"id":                nodeID,
		"name":              nodeName,
		"type":              nodeTypeRaw,
		"status":            nodeStatusRaw,
		"availability_zone": nodeAZ,
	}

	// Set private_read_ip with safe type assertion
	if nodePrivateIpsRaw := utils.PathSearch("private_read_ips", nodeRaw, nil); nodePrivateIpsRaw != nil {
		if nodePrivateIps, ok := nodePrivateIpsRaw.([]interface{}); ok && len(nodePrivateIps) > 0 {
			node["private_read_ip"] = nodePrivateIps[0]
		}
	}

	// Extract type and status as strings for counting
	if nodeTypeRaw != nil {
		if t, ok := nodeTypeRaw.(string); ok {
			nodeType = t
		}
	}

	if nodeStatusRaw != nil {
		if s, ok := nodeStatusRaw.(string); ok {
			nodeStatus = s
		}
	}

	if nodeFlavorRaw != nil {
		if f, ok := nodeFlavorRaw.(string); ok {
			nodeFlavor = f
		}
	}

	return node, nodeType, nodeStatus, nodeFlavor
}
