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
