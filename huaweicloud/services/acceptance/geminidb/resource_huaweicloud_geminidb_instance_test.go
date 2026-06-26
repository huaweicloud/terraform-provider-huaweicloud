package geminidb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGeminiDbInstance(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	instance := utils.PathSearch("instances[0]", getRespBody, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return instance, nil
}

func TestAccGeminiDbInstance_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDbInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "redis"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "16"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "03:00-04:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "14"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "on"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),

					resource.TestCheckResourceAttrSet(resourceName, "datastore.0.patch_available"),
					resource.TestCheckResourceAttrSet(resourceName, "datastore.0.whole_version"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.0.used"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.spec_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.support_reduce"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
				),
			},
			{
				Config: testAccGeminiDbInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "5"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "24"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "07:00-08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "off"),
					resource.TestCheckResourceAttr(resourceName, "port", "9999"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				Config: testAccGeminiDbInstance_reduce_node(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_period(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDbInstance_period(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "redis"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "16"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "03:00-04:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "14"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "on"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccGeminiDbInstance_period_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "24"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "07:00-08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "off"),
					resource.TestCheckResourceAttr(resourceName, "port", "9999"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
					"auto_renew",
					"period",
					"period_unit",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_dynamodb(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDbInstance_dynamodb(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "dynamodb"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.availability_zones.0"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "200"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "03:00-04:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "14"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "on"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),

					resource.TestCheckResourceAttrSet(resourceName, "datastore.0.patch_available"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.volume.0.used"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.spec_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.support_reduce"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
				),
			},
			{
				Config: testAccGeminiDbInstance_dynamodb_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "5"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "400"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "07:00-08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "off"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_configuration(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBConfiguration_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "redis"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "3"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "16"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.storage", "ULTRAHIGH"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor.0.spec_code",
						"data.huaweicloud_gaussdb_nosql_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "switch_option", "on"),
					resource.TestCheckResourceAttr(resourceName, "second_level_monitoring_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "lb_ip_address", "192.168.0.153"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_start_time", "04:00"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_end_time", "08:00"),

					resource.TestCheckResourceAttr(resourceName, "access_control.0.type", "blackList"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.ip_groups.0.ip", "123.123.123.0/24"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.ip_groups.0.description", "test"),

					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.#"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.0.threshold"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.0.step"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.0.size"),
				),
			},
			{
				Config: testAccGeminiDBConfiguration_update_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "switch_option", "on"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.threshold", "85"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.step", "15"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.size", "20"),
					resource.TestCheckResourceAttr(resourceName, "second_level_monitoring_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "config_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "lb_ip_address", "192.168.0.118"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_start_time", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_end_time", "10:00"),

					resource.TestCheckResourceAttr(resourceName, "access_control.0.type", "blackList"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.ip_groups.0.ip", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.ip_groups.0.description", "test update"),
				),
			},
			{
				Config: testAccGeminiDBConfiguration_update_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "switch_option", "off"),
					resource.TestCheckResourceAttr(resourceName, "config_ips.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "access_control.0.type", "blackList"),
					resource.TestCheckResourceAttr(resourceName, "access_control.0.enabled", "false"),

					resource.TestCheckResourceAttrSet(resourceName, "maintenance_start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "maintenance_end_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_dataExport(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBDataExport_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "influxdb"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttr(resourceName, "mode", "EnhancedCluster"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.spec_code", "geminidb.influxdb.sqlstore.large.4"),
					resource.TestCheckResourceAttr(resourceName, "data_export_switch", "open"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket_name", "huaweicloud_obs_bucket.test1", "bucket"),
				),
			},
			{
				Config: testAccGeminiDBDataExport_update(name, "open"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "data_export_switch", "open"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket_name", "huaweicloud_obs_bucket.test2", "bucket"),
				),
			},
			{
				Config: testAccGeminiDBDataExport_update(name, "close"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "data_export_switch", "close"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket_name", "huaweicloud_obs_bucket.test2", "bucket"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
					"bucket_name",
					"data_export_switch",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_coldStorage(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBColdStorage_basic(name, 500),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "influxdb"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttr(resourceName, "mode", "InfluxdbSingle"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.num", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "flavor.0.spec_code", "geminidb.influxdb.single.xlarge.2"),
					resource.TestCheckResourceAttr(resourceName, "cold_storage_size", "500"),
				),
			},
			{
				Config: testAccGeminiDBColdStorage_basic(name, 505),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cold_storage_size", "505"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
					"cold_storage_size",
				},
			},
		},
	})
}

func TestAccGeminiDbInstance_nodeAutoExpansionPolicy(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBNodeAutoExpansion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "cassandra"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "3.11"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.storage_engine", "rocksDB"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "node_switch_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "overload_node_threshold", "60"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold", "60"),
					resource.TestCheckResourceAttr(resourceName, "mem_threshold", "60"),
					resource.TestCheckResourceAttr(resourceName, "step", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_limit", "4"),
				),
			},
			{
				Config: testAccGeminiDBNodeAutoExpansion_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "node_switch_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "overload_node_threshold", "90"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold", "90"),
					resource.TestCheckResourceAttr(resourceName, "mem_threshold", "90"),
					resource.TestCheckResourceAttr(resourceName, "step", "3"),
					resource.TestCheckResourceAttr(resourceName, "node_limit", "6"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"flavor.0.storage",
					"ssl_option",
					"delete_node_list",
					"cold_storage_size",
				},
			},
		},
	})
}

func testAccGeminiDbInstance_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  count = 2

  name                 = "%[1]s_${count.index}"
  delete_default_rules = true
}
`, rName)
}

func testAccGeminiDbInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[0].id
  password              = "test_1234"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  port                  = 8888
  ssl_option            = "on"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccGeminiDbInstance_base(rName), rName)
}

func testAccGeminiDbInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[1].id
  password              = "test_123456"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  port                  = 9999
  ssl_option            = "off"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "5"
    size      = "24"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  }

  backup_strategy {
    start_time = "07:00-08:00"
    keep_days  = 10
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, testAccGeminiDbInstance_base(rName), updateName)
}

func testAccGeminiDbInstance_reduce_node(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[1].id
  password              = "test_123456"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  port                  = 9999
  ssl_option            = "off"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "24"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  }

  backup_strategy {
    start_time = "07:00-08:00"
    keep_days  = 10
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, testAccGeminiDbInstance_base(rName), updateName)
}

func testAccGeminiDbInstance_period(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[0].id
  password              = "test_1234"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  port                  = 8888
  ssl_option            = "on"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "true"
  period        = 1
}
`, testAccGeminiDbInstance_base(rName), rName)
}

func testAccGeminiDbInstance_period_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[1].id
  password              = "test_123456"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  port                  = 9999
  ssl_option            = "off"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "24"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  }

  backup_strategy {
    start_time = "07:00-08:00"
    keep_days  = 10
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "true"
  period        = 1
}
`, testAccGeminiDbInstance_base(rName), updateName)
}

func testAccGeminiDbInstance_dynamodb(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus  = 2
  engine = "cassandra"
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].availability_zones[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[0].id
  password              = "Test_159357"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  ssl_option            = "on"

  datastore {
    type           = "dynamodb"
    version        = ""
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "200"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccGeminiDbInstance_base(rName), rName)
}

func testAccGeminiDbInstance_dynamodb_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus  = 4
  engine = "cassandra"
}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].availability_zones[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test[1].id
  password              = "Test_357159"
  mode                  = "Cluster"
  enterprise_project_id = "0"
  ssl_option            = "off"

  datastore {
    type           = "dynamodb"
    version        = ""
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "5"
    size      = "400"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  backup_strategy {
    start_time = "07:00-08:00"
    keep_days  = 10
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, testAccGeminiDbInstance_base(rName), updateName)
}

func testAccGeminiDBConfiguration_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  # setting auto enlarge policy
  switch_option = "on"

  # enable second level monitor
  second_level_monitoring_enabled = "true"

  # enable password-free configuration
  config_ips = ["192.168.1.15","192.168.1.23/24"]

  # update loadbalancer IP address
  lb_ip_address = "192.168.0.153"

  # setting maintenance window
  maintenance_start_time = "04:00"
  maintenance_end_time   = "08:00"

  # Configuring the Blacklist or Whitelist of Load Balancer IP Addresses
  access_control {
    type    = "blackList"
    enabled = "true"
    ip_groups { 
      ip          = "123.123.123.0/24"
      description = "test" 
    }
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccGeminiDBConfiguration_update_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  switch_option = "on"

  policy {
    threshold = 85
    step      = 15
    size      = 20
  }

  second_level_monitoring_enabled = "false"
  config_ips                      = ["192.168.1.15","192.168.1.23/24","192.168.1.38"]
  lb_ip_address                   = "192.168.0.118"
  maintenance_start_time          = "06:00"
  maintenance_end_time            = "10:00"

  access_control {
    type    = "blackList"
    enabled = "true"
    ip_groups { 
      ip          = "192.168.0.0/24"
      description = "test update" 
    }
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccGeminiDBConfiguration_update_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  switch_option                   = "off"
  second_level_monitoring_enabled = false
  config_ips                      = []
  lb_ip_address                   = "192.168.0.118"
  maintenance_start_time          = "06:00"
  maintenance_end_time            = "10:00"

  access_control {
    type    = "blackList"
    enabled = "false"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccGeminiDBDataExport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test1" {
  bucket        = "%[2]s-b1"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket" "test2" {
  bucket        = "%[2]s-b2"
  storage_class = "STANDARD"
  acl           = "private"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test@1234"
  mode              = "EnhancedCluster"

  datastore {
    type           = "influxdb"
    version        = "1.8"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "2"
    size      = "100"
    storage   = "ULTRAHIGH"
    spec_code = "geminidb.influxdb.sqlstore.large.4"
  }

  data_export_switch = "open"
  bucket_name        = huaweicloud_obs_bucket.test1.bucket
}
`, common.TestBaseNetwork(name), name)
}

func testAccGeminiDBDataExport_update(name, dataSwitch string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test1" {
  bucket        = "%[2]s-b1"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket" "test2" {
  bucket        = "%[2]s-b2"
  storage_class = "STANDARD"
  acl           = "private"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test@1234"
  mode              = "EnhancedCluster"

  datastore {
    type           = "influxdb"
    version        = "1.8"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "2"
    size      = "100"
    storage   = "ULTRAHIGH"
    spec_code = "geminidb.influxdb.sqlstore.large.4"
  }

  data_export_switch = "%[3]s"
  bucket_name        = huaweicloud_obs_bucket.test2.bucket
}
`, common.TestBaseNetwork(name), name, dataSwitch)
}

func testAccGeminiDBColdStorage_basic(name string, storageSize int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test@1234"
  mode              = "InfluxdbSingle"

  datastore {
    type           = "influxdb"
    version        = "1.8"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "1"
    size      = "100"
    storage   = "ULTRAHIGH"
    spec_code = "geminidb.influxdb.single.xlarge.2"
  }

  cold_storage_size = %[3]d
}
`, common.TestBaseNetwork(name), name, storageSize)
}

func testAccGeminiDBNodeAutoExpansion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus  = 4
  engine = "cassandra"
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = "cn-north-4a,cn-north-4b,cn-north-4c"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test@1234"
  mode              = "Cluster"

  datastore {
    type           = "cassandra"
    version        = "3.11"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "100"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  node_switch_option      = true
  overload_node_threshold = 60
  cpu_threshold           = 60
  mem_threshold           = 60
  step                    = 2
  node_limit              = 4
}
`, common.TestBaseNetwork(name), name)
}

func testAccGeminiDBNodeAutoExpansion_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus  = 4
  engine = "cassandra"
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = "cn-north-4a,cn-north-4b,cn-north-4c"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test@1234"
  mode              = "Cluster"

  datastore {
    type           = "cassandra"
    version        = "3.11"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "100"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  node_switch_option      = false
  overload_node_threshold = 90
  cpu_threshold           = 90
  mem_threshold           = 90
  step                    = 3
  node_limit              = 6
}
`, common.TestBaseNetwork(name), name)
}
