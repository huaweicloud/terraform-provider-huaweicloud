package gaussdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDbInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
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

func TestAccGaussDbInstance_basic(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "port", "9000"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "eventual"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "20:00-21:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "6"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "dn:check_disconnect_query"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "off"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_tde_switch.0.kms_tde_key_id",
						"huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_switch.0.kms_project_name", "cn-north-4"),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_switch.0.kms_tde_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.name", "ilm"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.value", "on"),
					resource.TestCheckResourceAttr(resourceName, "wdr_snapshot_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "asp_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "error_log_switch_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "balance_status", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccGaussDbInstance_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "port", "8000"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "cn:auto_increment_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "1000"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_tde_switch.0.kms_tde_key_id",
						"huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_switch.0.kms_project_name", "cn-north-4"),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_switch.0.kms_tde_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.name", "ilm"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.value", "off"),
					resource.TestCheckResourceAttr(resourceName, "wdr_snapshot_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "asp_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "error_log_switch_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"advance_features",
					"auto_renew",
					"period",
					"period_unit",
					"replica_num",
					"availability_zone",
					"configuration_id",
					"parameters",
					"kms_tde_switch",
					"alias",
				},
			},
		},
	})
}

func TestAccGaussDbInstance_node_num(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_node_num(rName, password, 1, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
				),
			},
			{
				Config: testAccGaussDbInstance_node_num(rName, password, 2, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
				),
			},
			{
				Config: testAccGaussDbInstance_node_num(rName, password, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
				),
			},
		},
	})
}

func TestAccGaussDbInstance_node_num_prepaid(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_node_num_prepaid(rName, password, 1, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
				),
			},
			{
				Config: testAccGaussDbInstance_node_num_prepaid(rName, password, 2, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
				),
			},
			{
				Config: testAccGaussDbInstance_node_num_prepaid(rName, password, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
				),
			},
		},
	})
}

func TestAccGaussDbInstance_flavor(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_flavor(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_flavors.test", "flavors.0.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "mysql_compatibility_port", "12345"),
				),
			},
			{
				Config: testAccGaussDbInstance_flavor_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_flavors.test", "flavors.1.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "mysql_compatibility_port", "12346"),
				),
			},
		},
	})
}

func TestAccGaussDbInstance_haModeCentralized(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_haModeCentralized(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_flavors.test", "flavors.0.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "centralization_standard"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "eventual"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.instance_mode", "basic"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccGaussDbInstance_haModeCentralizedUpdate(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_flavors.test", "flavors.1.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
				),
			},
		},
	})
}

func TestAccGaussDbInstance_auto_scaling(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_auto_scaling(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.switch_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.limit_volume_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.trigger_available_percent", "25"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.step_percent", "50"),
				),
			},
			{
				Config: testAccGaussDbInstance_auto_scaling_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.switch_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.limit_volume_size", "600"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.trigger_available_percent", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling.0.step_size", "40"),
				),
			},
		},
	})
}

func TestAccGaussDbInstance_dn_node_deploy_mode(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getGaussDbInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstance_dn_node_deploy_mode(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dn_node_deploy_mode", "primary_standby_log"),
				),
			},
			{
				Config: testAccGaussDbInstance_dn_node_deploy_mode_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dn_node_deploy_mode", "primary_standby"),
				),
			},
		},
	})
}

func testAccGaussDbInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup" "test" {
  count = 2

  name = "%s"
}
`, common.TestVpc(rName), rName)
}

func testAccGaussDbInstance_basic(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[2]s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "test acc"
  enterprise_project_id = "%[4]s"
  pending_days          = "7"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test[0].id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  alias             = "%[2]s_alias"
  password          = "%[3]s"
  port              = "9000"
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  wdr_snapshot_status     = "OFF"
  asp_status              = "OFF"
  error_log_switch_status = "ON"
  enterprise_project_id   = "%[4]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  backup_strategy {
    start_time = "20:00-21:00"
    keep_days  = 6
  }

  parameters {
    name  = "dn:check_disconnect_query"
    value = "off"
  }

  advance_features {
    name  = "ilm"
    value = "on"
  }

  tags = {
    foo = "bar"
    key = "value"
  }

 kms_tde_switch {
    kms_tde_key_id   = huaweicloud_kms_key.test.id
    kms_project_name = "cn-north-4"
    kms_tde_status   = "on"
  }
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[2]s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "test acc"
  enterprise_project_id = "%[4]s"
  pending_days          = "7"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_gaussdb_parameter_template" "test" {
  name           = "%[2]s"
  engine_version = "8.218"
  instance_mode  = "independent"

  parameters {
    name  = "cma:log_max_count"
    value = "5000"
  }

  parameters {
    name  = "dn:a_format_date_timestamp"
    value = "on"
  }
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test[1].id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s-update"
  alias             = "%[2]s_alias_update"
  password          = "%[3]s"
  port              = "8000"
  sharding_num      = 1
  coordinator_num   = 2
  configuration_id  = huaweicloud_gaussdb_parameter_template.test.id
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  wdr_snapshot_status     = "ON"
  asp_status              = "ON"
  error_log_switch_status = "OFF"
  enterprise_project_id   = "%[4]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }

  parameters {
    name  = "cn:auto_increment_increment"
    value = "1000"
  }

  advance_features {
    name  = "ilm"
    value = "off"
  }

  tags = {
    foo_update = "bar"
    key        = "value_update"
  }

  kms_tde_switch {
    kms_tde_key_id   = huaweicloud_kms_key.test.id
    kms_project_name = "cn-north-4"
    kms_tde_status   = "on"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_node_num(rName, password string, shardingNum, coordinatorNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test[0].id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = %[4]d
  coordinator_num   = %[5]d
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[6]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccGaussDbInstance_base(rName), rName, password, shardingNum, coordinatorNum, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_node_num_prepaid(rName, password string, shardingNum, coordinatorNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test[0].id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = %[4]d
  coordinator_num   = %[5]d
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[6]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccGaussDbInstance_base(rName), rName, password, shardingNum, coordinatorNum, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_flavor(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  mysql_compatibility_port = "12345"
  enterprise_project_id    = "%[4]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_flavor_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[1].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  mysql_compatibility_port = "12346"
  enterprise_project_id    = "%[4]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_haModeCentralized(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[4]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_haModeCentralizedUpdate(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[1].spec_code
  name              = "%[2]s-update"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[4]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccGaussDbInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_auto_scaling(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[4]s"

  auto_scaling {
    switch_option             = true
    limit_volume_size         = 400
    trigger_available_percent = 25
    step_percent              = 50
  }

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_auto_scaling_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[4]s"

  auto_scaling {
    switch_option             = false
    limit_volume_size         = 600
    trigger_available_percent = 50
    step_size                 = 40
  }

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_dn_node_deploy_mode(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "enterprise"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 1
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  dn_node_deploy_mode = "primary_standby_log"

  enterprise_project_id = "%[4]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDbInstance_dn_node_deploy_mode_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "enterprise"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 1
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  dn_node_deploy_mode = "primary_standby"

  enterprise_project_id = "%[4]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
