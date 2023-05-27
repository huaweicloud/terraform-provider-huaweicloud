package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getDcsResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DcsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID)
}

func TestAccDcsInstances_basic(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6389"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "10:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at"},
			},
		},
	})
}

func TestAccDcsInstances_withEpsId(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
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
				Config: testAccDcsV1Instance_epsId(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccDcsInstances_whitelists(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_whitelists(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "test-group1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "192.168.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "192.168.0.0/24"),
				),
			},
			{
				Config: testAccDcsV1Instance_whitelists_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "test-group2"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "172.16.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "172.16.0.0/24"),
				),
			},
		},
	})
}

func TestAccDcsInstances_tiny(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_tiny(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
				),
			},
		},
	})
}

func TestAccDcsInstances_single(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_single(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func TestAccDcsInstances_prePaid(t *testing.T) {
	var instance instances.DcsInstance
	var rName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsInstance_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccDcsInstance_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit"},
			},
		},
	})
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 0.125
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "02:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsV1Instance_updated(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 1
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6389
  capacity           = 1
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "10:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "01:00-02:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4]
    save_days   = 2
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  tags = {
    key   = "value_update"
    owner = "terraform_update"
  }
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsV1Instance_epsId(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }
  enterprise_project_id = "%s"
}`, common.TestVpc(instanceName), instanceName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDcsV1Instance_tiny(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsV1Instance_single(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "single"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsV1Instance_whitelists(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }

  whitelists {
    group_name = "test-group1"
    ip_address = ["192.168.10.100", "192.168.0.0/24"]
  }
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsV1Instance_whitelists_update(instanceName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }

  whitelists {
    group_name = "test-group2"
    ip_address = ["172.16.10.100", "172.16.0.0/24"]
  }
}`, common.TestVpc(instanceName), instanceName)
}

func testAccDcsInstance_prePaid(instanceName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  name               = "%s"
  engine             = "Redis"
  engine_version     = "5.0"
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  capacity           = 0.125
  password           = "Huawei_test"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%v"
}`, common.TestVpc(instanceName), instanceName, isAutoRenew)
}
