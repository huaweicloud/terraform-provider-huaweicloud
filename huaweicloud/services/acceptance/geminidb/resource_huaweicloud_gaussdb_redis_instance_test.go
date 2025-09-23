package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussRedisInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("Acc%s@123", acctest.RandString(5))
	newPassword := fmt.Sprintf("Acc%sUpdate@123", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_redis_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGaussRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussRedisInstanceConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "16"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
				),
			},
			{
				Config: testAccGaussRedisInstanceConfig_update(rName, newPassword, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "5"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "24"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "ssl", "false"),
				),
			},
			{
				Config: testAccGaussRedisInstanceConfig_update(rName, newPassword, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "24"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"availability_zone", "password", "ssl"},
			},
		},
	})
}

func TestAccGaussRedisInstance_withReplication(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_redis_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGaussRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussRedisInstanceConfig_withReplication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "node_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Replication"),
				),
			},
		},
	})
}

func TestAccGaussRedisInstance_updateWithEpsId(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("Acc%s@123", acctest.RandString(5))
	newPassword := fmt.Sprintf("Acc%sUpdate@123", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_redis_instance.test"
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGaussRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussRedisInstanceConfig_withEpsId(rName, password, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccGaussRedisInstanceConfig_withEpsId(rName, newPassword, destEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussRedisInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccCheckGaussRedisInstanceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.GeminiDBV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating GaussRedis client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_redis_instance" {
			continue
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
		if found.Id != "" {
			return fmt.Errorf("instance <%s> still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGaussRedisInstanceExists(n string, instance *instances.GeminiDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.GeminiDBV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating GuassRedis client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("instance <%s> not found", rs.Primary.ID)
			}
			return err
		}
		if found.Id == "" {
			return fmt.Errorf("instance <%s> not found", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccGaussRedisInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_redis_instance" "test" {
  name        = "%[2]s"
  password    = "%[3]s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 16
  vpc_id      = data.huaweicloud_vpc.test.id
  subnet_id   = data.huaweicloud_vpc_subnet.test.id
  node_num    = 3
  port        = 8888
  ssl         = true

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestSecGroup(rName), rName, password)
}

func testAccGaussRedisInstanceConfig_update(rName, password string, nodeNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_redis_instance" "test" {
  name        = "%[2]s-update"
  password    = "%[3]s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 24
  vpc_id      = data.huaweicloud_vpc.test.id
  subnet_id   = data.huaweicloud_vpc_subnet.test.id
  node_num    = %[4]d
  port        = 8888
  ssl         = false

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestSecGroup(rName+"_update"), rName, password, nodeNum)
}

func testAccGaussRedisInstanceConfig_withReplication(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_gaussdb_redis_flavors" "test" {}

resource "huaweicloud_gaussdb_redis_instance" "test" {
  name              = "%[2]s"
  password          = "Huangwei!120521"
  flavor            = data.huaweicloud_gaussdb_redis_flavors.test.flavors[0].spec_code
  volume_size       = 2
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  node_num          = 2
  port              = 8888
  ssl               = true
  mode              = "Replication"
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = keys(data.huaweicloud_gaussdb_redis_flavors.test.flavors[0].az_status)[1]

  availability_zone_detail{
    primary_availability_zone   = split(",", keys(data.huaweicloud_gaussdb_redis_flavors.test.flavors[0].az_status)[1])[0]
    secondary_availability_zone = split(",", keys(data.huaweicloud_gaussdb_redis_flavors.test.flavors[0].az_status)[1])[1]
  }
}
`, common.TestSecGroup(rName), rName)
}

func testAccGaussRedisInstanceConfig_withEpsId(rName, password, epsId string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_redis_instance" "test" {
  name                  = "%[2]s"
  password              = "%[3]s"
  flavor                = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size           = 50
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  node_num              = 3
  port                  = 8888
  ssl                   = true
  enterprise_project_id = "%s"

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestSecGroup(rName), rName, password, epsId)
}
