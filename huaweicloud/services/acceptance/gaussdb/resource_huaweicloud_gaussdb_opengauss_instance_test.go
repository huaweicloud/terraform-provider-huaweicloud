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

func getOpenGaussInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccOpenGaussInstance_basic(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
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
				Config: testAccOpenGaussInstance_basic(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_opengauss_flavors.test", "flavors.0.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "eventual"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "20:00-21:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "6"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "dn:check_disconnect_query"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "off"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.name", "ilm"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.value", "on"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "balance_status", "true"),
					resource.TestCheckResourceAttr(resourceName, "error_log_switch_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccOpenGaussInstance_update(rName, newPassword, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "cn:auto_increment_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.name", "ilm"),
					resource.TestCheckResourceAttr(resourceName, "advance_features.0.value", "off"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_flavor(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
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
				Config: testAccOpenGaussInstance_flavor(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_opengauss_flavors.test", "flavors.0.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "mysql_compatibility_port", "12345"),
				),
			},
			{
				Config: testAccOpenGaussInstance_flavor_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_opengauss_flavors.test", "flavors.1.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "mysql_compatibility_port", "12346"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_haModeCentralized(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
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
				Config: testAccOpenGaussInstance_haModeCentralized(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_opengauss_flavors.test", "flavors.0.spec_code"),
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
				Config: testAccOpenGaussInstance_haModeCentralizedUpdate(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_opengauss_flavors.test", "flavors.1.spec_code"),
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

func testAccOpenGaussInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

// opengauss requires more sg ports open
resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, common.TestBaseNetwork(rName))
}

func testAccOpenGaussInstance_basic(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "enterprise"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 2
  replica_num       = %[4]d
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[5]s"

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
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstance_update(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "enterprise"
}

resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name           = "%[2]s"
  engine_version = "8.210"
  instance_mode  = "independent"

  parameters {
    name  = "cn:auto_explain_log_min_duration"
    value = "1000"
  }

  parameters {
    name  = "dn:a_format_date_timestamp"
    value = "on"
  }
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name              = "%[2]s-update"
  password          = "%[3]s"
  sharding_num      = 2
  coordinator_num   = 3
  replica_num       = %[4]d
  configuration_id  = huaweicloud_gaussdb_opengauss_parameter_template.test.id
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[5]s"

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
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstance_flavor(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
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
`, testAccOpenGaussInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstance_flavor_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[1].spec_code
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
`, testAccOpenGaussInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstance_haModeCentralized(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
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
`, testAccOpenGaussInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstance_haModeCentralizedUpdate(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[1].spec_code
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
`, testAccOpenGaussInstance_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
