package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/taurusdb"
)

func getHtapStarrocksInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}
	taurusdbInstanceId := state.Primary.Attributes["instance_id"]
	htapInstanceId := state.Primary.ID
	htapInstanceDetail, err := taurusdb.GetHtapInstanceDetail(client, taurusdbInstanceId, htapInstanceId)
	if err != nil {
		return nil, err
	}
	if htapInstanceDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return htapInstanceDetail, nil
}

func TestAccTaurusDBHtapStarrocksInstance_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_htap_starrocks_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHtapStarrocksInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "fe_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "be_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "az_mode", "single"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db_port", "3306"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_vip"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_taurusdb_instance.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_taurusdb_instance.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_taurusdb_instance.test", "security_group_id"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "users_sync_switch_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "can_enable_public_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "update_at"),
					resource.TestCheckResourceAttr(resourceName, "fe_node_volume_code", "gaussdb.sr.evs.ultrahighio"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_node_volume_size"),
					resource.TestCheckResourceAttr(resourceName, "be_node_volume_code", "gaussdb.sr.evs.ultrahighio"),
					resource.TestCheckResourceAttrSet(resourceName, "be_node_volume_size"),
					resource.TestCheckResourceAttr(resourceName, "support_data_replication", "true"),
					resource.TestCheckResourceAttr(resourceName, "new_version_available", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_store_type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version_id"),
					resource.TestCheckResourceAttr(resourceName, "db_user", "root"),
					resource.TestCheckResourceAttr(resourceName, "cluster_mode", "Single"),
					resource.TestCheckResourceAttr(resourceName, "is_frozen", "0"),
					resource.TestCheckResourceAttr(resourceName, "enable_users_sync", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.status", "normal"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.0.type", "SSD"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.mem"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.0.type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.version"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.frozen_flag", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.pay_model", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_description"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_type"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.region_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.update_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.iass_flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.max_connections"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.sg_id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.need_restart", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.nodes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ops_window.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.end_time"),
					resource.TestCheckResourceAttr(resourceName, "port_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_info.0.mysql_port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.sys_tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.sys_tags.0.key", "_sys_enterprise_project_id"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.sys_tags.0.value", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "be_configurations.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "be_configurations.0.configuration_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(resourceName, "be_configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "be_parameters.0.description"),
					resource.TestCheckResourceAttr(resourceName, "fe_configurations.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_configurations.0.configuration_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_parameters.0.description"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksInstance_basicUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "fe_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_group_id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.need_restart", "true"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.need_restart", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_users_sync", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBHtapStarrocksInstanceImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"db_root_pwd", "be_parameter_values", "fe_parameter_values", "period", "period_unit", "auto_renew",
				},
			},
		},
	})
}
func TestAccTaurusDBHtapStarrocksInstance_cluster(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_htap_starrocks_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHtapStarrocksInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksInstance_cluster(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "fe_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "be_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "az_mode", "single"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db_port", "3306"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_vip"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_taurusdb_instance.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_taurusdb_instance.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_taurusdb_instance.test", "security_group_id"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "users_sync_switch_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "can_enable_public_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "update_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_node_volume_code"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_node_volume_size"),
					resource.TestCheckResourceAttrSet(resourceName, "be_node_volume_code"),
					resource.TestCheckResourceAttrSet(resourceName, "be_node_volume_size"),
					resource.TestCheckResourceAttr(resourceName, "support_data_replication", "true"),
					resource.TestCheckResourceAttr(resourceName, "new_version_available", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_store_type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version_id"),
					resource.TestCheckResourceAttr(resourceName, "db_user", "root"),
					resource.TestCheckResourceAttr(resourceName, "cluster_mode", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "is_frozen", "0"),
					resource.TestCheckResourceAttr(resourceName, "enable_users_sync", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.status", "normal"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.0.type", "SSD"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.mem"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.0.type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.version"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.frozen_flag", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.pay_model", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_description"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_type"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.region_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.update_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.iass_flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.max_connections"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.sg_id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.need_restart", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.nodes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ops_window.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.end_time"),
					resource.TestCheckResourceAttr(resourceName, "port_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_info.0.mysql_port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.sys_tags.#", "1"),
				),
			},
		},
	})
}
func TestAccTaurusDBHtapStarrocksInstance_prePaid(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_htap_starrocks_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHtapStarrocksInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "fe_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "be_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "az_mode", "single"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "db_port", "3306"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_vip"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_taurusdb_instance.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_taurusdb_instance.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_taurusdb_instance.test", "security_group_id"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+07:00"),
					resource.TestCheckResourceAttr(resourceName, "ssl_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "users_sync_switch_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "can_enable_public_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "update_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_node_volume_code"),
					resource.TestCheckResourceAttrSet(resourceName, "fe_node_volume_size"),
					resource.TestCheckResourceAttrSet(resourceName, "be_node_volume_code"),
					resource.TestCheckResourceAttrSet(resourceName, "be_node_volume_size"),
					resource.TestCheckResourceAttr(resourceName, "support_data_replication", "true"),
					resource.TestCheckResourceAttr(resourceName, "new_version_available", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_store_type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version"),
					resource.TestCheckResourceAttrSet(resourceName, "data_store_version_id"),
					resource.TestCheckResourceAttr(resourceName, "db_user", "root"),
					resource.TestCheckResourceAttr(resourceName, "cluster_mode", "Single"),
					resource.TestCheckResourceAttr(resourceName, "is_frozen", "0"),
					resource.TestCheckResourceAttr(resourceName, "enable_users_sync", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.status", "normal"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.type"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.volume.0.type", "SSD"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.mem"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.datastore.0.type", "star-rocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.datastore.0.version"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.actions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.frozen_flag", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.pay_model", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_description"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.az_type"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.region_code"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.create_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.update_at"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.iass_flavor_ref"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.max_connections"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.0.nodes.0.sg_id"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.nodes.0.need_restart", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.name"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.group_type_name", "starrocks"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.1.group_node_type"),
					resource.TestCheckResourceAttr(resourceName, "groups.1.nodes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ops_window.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "ops_window.0.end_time"),
					resource.TestCheckResourceAttr(resourceName, "port_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_info.0.mysql_port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_info.0.sys_tags.#", "1"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksInstance_prePaid_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "fe_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "be_flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_group_id"),
					resource.TestCheckResourceAttr(resourceName, "users_sync_switch_on", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBHtapStarrocksInstanceImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"db_root_pwd", "be_parameter_values", "fe_parameter_values", "period", "period_unit", "auto_renew",
				},
			},
		},
	})
}

func testAccHtapInstanceConfig_base(rName, enterpriseProjectId string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_taurusdb_flavors" "test" {
 engine                 = "gaussdb-mysql"
 version                = "8.0"
 availability_zone_mode = "multi"
}

resource "huaweicloud_networking_secgroup" "test" {
 count                = 2
 name                 = "%[2]s_${count.index}"
 delete_default_rules = true
}

resource "huaweicloud_taurusdb_instance" "test" {
 name                     = "%[2]s"
 password                 = "Test@12345678"
 flavor                   = data.huaweicloud_taurusdb_flavors.test.flavors[0].name
 vpc_id                   = huaweicloud_vpc.test.id
 subnet_id                = huaweicloud_vpc_subnet.test.id
 security_group_id        = huaweicloud_networking_secgroup.test[0].id
 enterprise_project_id    = "%[3]s"
 master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
 availability_zone_mode   = "multi"
 read_replicas            = 2
}

data "huaweicloud_taurusdb_htap_flavors" "test" {
 engine_name            = "star-rocks"
 availability_zone_mode = "single"
}

data "huaweicloud_taurusdb_htap_datastores" "test" {
 engine_name = "star-rocks"
}

locals {
 all_az_normal = [for k, v in data.huaweicloud_taurusdb_htap_flavors.test.flavors[0].az_status : k if v == "normal"]
 az_code       = local.all_az_normal[0]

 be_flavors = [for f in data.huaweicloud_taurusdb_htap_flavors.test.flavors : f if length(regexall("sr-be", f.spec_code)) > 0 &&
 contains(keys(f.az_status), local.az_code) && f.az_status[local.az_code] == "normal"]

 fe_flavors = [for f in data.huaweicloud_taurusdb_htap_flavors.test.flavors : f if length(regexall("sr-fe", f.spec_code)) > 0 &&
 contains(keys(f.az_status), local.az_code) && f.az_status[local.az_code] == "normal"]

 be_base_flavor_id   = local.be_flavors[0].id
 be_resize_flavor_id = local.be_flavors[1].id
 fe_base_flavor_id   = local.fe_flavors[0].id
 fe_resize_flavor_id = local.fe_flavors[1].id
 security_group_id   = huaweicloud_networking_secgroup.test[1].id
 engine_version      = data.huaweicloud_taurusdb_htap_datastores.test.datastores[0].kernel_version
}
`, common.TestVpc(rName, enterpriseProjectId), rName, enterpriseProjectId)
}

func testAccTaurusDBHtapStarrocksInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = "%[2]s"
  fe_flavor_id      = local.fe_base_flavor_id
  be_flavor_id      = local.be_base_flavor_id
  db_root_pwd       = "Test@123456!"
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  az_code           = local.az_code
  enable_users_sync = "true"

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = "%[3]s"
    }
  }
  
  be_parameter_values = {
    "alter_tablet_worker_count"            = "1"
    "base_compaction_num_threads_per_disk" = "1"
  }
  
  fe_parameter_values = {
    "alter_table_timeout_second"     = "21600"
    "bdbje_heartbeat_timeout_second" = "10"
  }
}
`, testAccHtapInstanceConfig_base(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccTaurusDBHtapStarrocksInstance_basicUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = "%[2]s"
  fe_flavor_id      = local.fe_resize_flavor_id
  be_flavor_id      = local.be_resize_flavor_id
  db_root_pwd       = "Test@1234567!"
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  az_code           = local.az_code
  security_group_id = local.security_group_id
  enable_users_sync = "false"

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = "%[3]s"
    }
  }
  
  be_parameter_values = {
    "alter_tablet_worker_count"            = "20"
    "base_compaction_num_threads_per_disk" = "5"
  }
  
  fe_parameter_values = {
    "alter_table_timeout_second"     = "259200"
    "bdbje_heartbeat_timeout_second" = "100"
  }
}
`, testAccHtapInstanceConfig_base(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccTaurusDBHtapStarrocksInstanceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		htapInstanceId := rs.Primary.ID
		if instanceId == "" || htapInstanceId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, htapInstanceId)
		}
		return fmt.Sprintf("%s/%s", instanceId, htapInstanceId), nil
	}
}

func testAccTaurusDBHtapStarrocksInstance_cluster(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = "%[2]s"
  fe_flavor_id      = local.fe_base_flavor_id
  be_flavor_id      = local.be_base_flavor_id
  db_root_pwd       = "Test@123456!"
  fe_count          = 3
  be_count          = 3
  az_mode           = "single"
  az_code           = local.az_code
  enable_users_sync = "true"

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = "Cluster"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = "%[3]s"
    }
  }
}
`, testAccHtapInstanceConfig_base(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccTaurusDBHtapStarrocksInstance_prePaid(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = "%[2]s"
  fe_flavor_id      = local.fe_base_flavor_id
  be_flavor_id      = local.be_base_flavor_id
  db_root_pwd       = "Test@123456!"
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  az_code           = local.az_code
  time_zone         = "UTC+07:00"
  enable_users_sync = "true"
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
  auto_renew        = "false"

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = "%[3]s"
    }
  }
}
`, testAccHtapInstanceConfig_base(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccTaurusDBHtapStarrocksInstance_prePaid_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = "%[2]s"
  fe_flavor_id      = local.fe_resize_flavor_id
  be_flavor_id      = local.be_resize_flavor_id
  db_root_pwd       = "Test@1234567!"
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  az_code           = local.az_code
  time_zone         = "UTC+07:00"
  security_group_id = local.security_group_id
  enable_users_sync = "false"
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
  auto_renew        = "true"

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = "%[3]s"
    }
  }
}
`, testAccHtapInstanceConfig_base(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
