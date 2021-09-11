package dcs

import (
	"fmt"
	"testing"
	"time"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DcsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DCS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID)
}

//lintignore:R018
func TestAccResourceDcsRedisInstance_basic(t *testing.T) {
	var redisInstance instances.DcsInstance
	resourceName := "huaweicloud_dcs_redis_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&redisInstance,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDcsRedisInstance_basic(
					map[string]string{
						"name":     name,
						"version":  "5.0",
						"capacity": "0.5",
						"specCode": "redis.ha.xu1.tiny.r2.512",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.ha.xu1.tiny.r2.512"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "available_zones.0",
						"${data.huaweicloud_availability_zones.zones.names.0}"),
					resource.TestCheckResourceAttr(resourceName, "port", "6380"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "00:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0700"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "A"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.1.group_name", "group_2"),
				),
			},
			{
				Config: testAccResourceDcsRedisInstance_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.ha.xu1.tiny.r2.256"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "02:00-04:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0800"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "B"),
					resource.TestCheckResourceAttr(resourceName, "tags.test", "abc"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.1.group_name", "group_3"),
					func(s *terraform.State) error {
						// After changing the specification, we need to wait for the xx order to be closed.
						// Exist a unknown order, need to wait to be closed.
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands"},
			},
		},
	})
}

func TestAccResourceDcsRedisInstance_masterStandby(t *testing.T) {
	var redisInstance instances.DcsInstance
	resourceName := "huaweicloud_dcs_redis_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&redisInstance,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDcsRedisInstance_basic(
					map[string]string{
						"name":     name,
						"version":  "5.0",
						"capacity": "0.125",
						"specCode": "redis.ha.xu1.tiny.r2.128",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.ha.xu1.tiny.r2.128"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "available_zones.0",
						"${data.huaweicloud_availability_zones.zones.names.0}"),
					resource.TestCheckResourceAttr(resourceName, "port", "6380"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "00:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0700"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "A"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceDcsRedisInstance_cluster(t *testing.T) {
	var redisInstance instances.DcsInstance
	resourceName := "huaweicloud_dcs_redis_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&redisInstance,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDcsRedisInstance_basic(
					map[string]string{
						"name":     name,
						"version":  "5.0",
						"capacity": "4",
						"specCode": "redis.proxy.xu1.large.4",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.proxy.xu1.large.4"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "available_zones.0",
						"${data.huaweicloud_availability_zones.zones.names.0}"),
					resource.TestCheckResourceAttr(resourceName, "port", "6380"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "00:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0700"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "A"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceDcsRedisInstance_rwsplit(t *testing.T) {
	var redisInstance instances.DcsInstance
	resourceName := "huaweicloud_dcs_redis_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&redisInstance,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDcsRedisInstance_basic(
					map[string]string{
						"name":     name,
						"version":  "5.0",
						"capacity": "8",
						"specCode": "redis.ha.xu1.large.p2.8",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.ha.xu1.large.p2.8"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "available_zones.0",
						"${data.huaweicloud_availability_zones.zones.names.0}"),
					resource.TestCheckResourceAttr(resourceName, "port", "6380"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "00:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0700"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "A"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceDcsRedisInstance_postPaid(t *testing.T) {
	var redisInstance instances.DcsInstance
	resourceName := "huaweicloud_dcs_redis_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&redisInstance,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDcsRedisInstance_postPaid(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(resourceName, "resource_spec_code", "redis.ha.xu1.tiny.r2.128"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "available_zones.0",
						"${data.huaweicloud_availability_zones.zones.names.0}"),
					resource.TestCheckResourceAttr(resourceName, "port", "6389"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "02:00-04:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.timezone_offset", "+0800"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.level", "A"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
				),
			},
		},
	})
}

func dcsInstanceDependResource(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s_dcs"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = "%s_dcs"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s_dcs"
  description = "terraform security group acceptance test"
}

data "huaweicloud_availability_zones" "zones" {}
`, name, name, name)
}

func testAccResourceDcsRedisInstance_basic(params map[string]string) string {
	var (
		name     = params["name"]
		version  = params["version"]
		capacity = params["capacity"]
		specCode = params["specCode"]
	)

	return fmt.Sprintf(`
%s
resource "huaweicloud_dcs_redis_instance" "instance_1" {
  name               = "%s"
  engine_version     = "%s"
  capacity           = %s
  resource_spec_code = "%s"
  available_zones    = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  port               = 6380
  maintain_begin     = "22:00:00"
  maintain_end       = "02:00:00"
  password           = "Abc@123"

  charging_mode    = "prePaid"
  period_unit      = "month"
  auto_renew       = "false"
  period           = "1"
  description      = "tf demo"
  whitelist_enable = true

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
  whitelists {
    group_name = "group_2"
    ip_address = ["192.168.2.0/24"]
  }

  backup_policy {
    backup_type     = "auto"
    save_days       = 3
    period_type     = "weekly"
    backup_at       = [1, 3, 5]
    begin_at        = "00:00-02:00"
    timezone_offset = "+0700"
  }

  rename_commands = {
    "command": "cmd",
    "keys": "key",
    "flushdb": "flshdb",
    "flushall": "flsall",
    "hgetall": "getall"
  }

  tags = {
    "level": "A",
    "order": ""
  }
}
`, dcsInstanceDependResource(name), name, version, capacity, specCode)
}

func testAccResourceDcsRedisInstance_update(name string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_dcs_redis_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  capacity           = 0.25
  resource_spec_code = "redis.ha.xu1.tiny.r2.256"
  available_zones    = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  port               = 6388
  maintain_begin     = "02:00:00"
  maintain_end       = "06:00:00"
  password           = "Abc@321"

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "false"
  period        = "1"
  description   = "tf demo"
  whitelist_enable = false

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
  whitelists {
    group_name = "group_3"
    ip_address = ["192.168.3.0/24"]
  }

  backup_policy {
    backup_type     = "auto"
    save_days       = 3
    period_type     = "weekly"
    backup_at       = [1, 2, 3, 4, 5, 6, 7]
    begin_at        = "02:00-04:00"
    timezone_offset = "+0800"
  }

  rename_commands = {
    "command": "cmd",
    "keys": "key",
    "flushdb": "flshdb",
    "flushall": "flsall",
    "hgetall": "getall"
  }

  tags = {
    "level": "B",
    "order": "",
    "test": "abc"
  }
}
`, dcsInstanceDependResource(name), name)
}

func testAccResourceDcsRedisInstance_postPaid(name string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_dcs_redis_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  capacity           = 0.125
  resource_spec_code = "redis.ha.xu1.tiny.r2.128"
  available_zones    = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  port               = 6389
  maintain_begin     = "02:00:00"
  maintain_end       = "06:00:00"

  charging_mode = "postPaid"
  description   = "tf demo"
  whitelist_enable = true

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
  whitelists {
    group_name = "group_2"
    ip_address = ["192.168.3.0/24"]
  }

  backup_policy {
    backup_type     = "auto"
    save_days       = 3
    period_type     = "weekly"
    backup_at       = [1, 2, 3, 4, 5, 6, 7]
    begin_at        = "02:00-04:00"
    timezone_offset = "+0800"
  }

  rename_commands = {
    "command": "cmd",
    "keys": "key",
    "flushdb": "flshdb",
    "flushall": "flsall",
    "hgetall": "getall"
  }

  tags = {
    "level": "A",
    "order": ""
  }
}
`, dcsInstanceDependResource(name), name)
}
