package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccReadReplicaInstance_basic(t *testing.T) {
	var instance interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getResourceInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccReadReplicaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.replica", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test_secgroup.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "connect_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "14"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "09:00"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccReadReplicaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description_update"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.replica", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test_secgroup.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8889"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "div_precision_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "12"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "15:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "17:00"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"parameters",
				},
			},
		},
	})
}

func TestAccReadReplicaInstance_prePaid(t *testing.T) {
	var instance interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getResourceInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccReadReplicaInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
				),
			},
			{
				Config: testAccReadReplicaInstance_prePaid_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
				),
			},
		},
	})
}

func testAccReadReplicaInstance_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccReadReplicaInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}

resource "huaweicloud_networking_secgroup" "test_secgroup" {
  count = 2

  name                 = "%[2]s_${count.index}"
  delete_default_rules = true
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[2]s"
  description         = "test_description"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = huaweicloud_networking_secgroup.test_secgroup[0].id
  ssl_enable          = true
  maintain_begin      = "06:00"
  maintain_end        = "09:00"

  db {
    port = 8888
  }

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 10
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccReadReplicaInstance_base(name), name)
}

func testAccReadReplicaInstance_update(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 8
  vcpus         = 2
}

resource "huaweicloud_networking_secgroup" "test_secgroup" {
  count = 2

  name                 = "%[2]s_${count.index}"
  delete_default_rules = true
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[2]s"
  description         = "test_description_update"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = huaweicloud_networking_secgroup.test_secgroup[1].id
  ssl_enable          = false
  maintain_begin      = "15:00"
  maintain_end        = "17:00"

  db {
    port = 8889
  }

  volume {
    type              = "CLOUDSSD"
    size              = 60
    limit_size        = 500
    trigger_threshold = 15
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }

  tags = {
    key_update = "value_update"
    foo_update = "bar_update"
  }
}
`, testAccReadReplicaInstance_base(name), updateName)
}

func testAccReadReplicaInstance_prePaid(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[2]s"
  description         = "test_description"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  ssl_enable          = true
  maintain_begin      = "06:00"
  maintain_end        = "09:00"

  db {
    port = 8888
  }

  volume {
    type       = "CLOUDSSD"
    size       = 50
    limit_size = 0
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testAccReadReplicaInstance_base(name), name)
}

func testAccReadReplicaInstance_prePaid_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[2]s"
  description         = "test_description_update"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = data.huaweicloud_networking_secgroup.test.id

  db {
    port = 8889
  }

  volume {
    type       = "CLOUDSSD"
    size       = 60
    limit_size = 0
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, testAccReadReplicaInstance_base(name), name)
}
