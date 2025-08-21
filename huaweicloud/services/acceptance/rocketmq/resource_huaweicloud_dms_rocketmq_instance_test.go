package rocketmq

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

func getDmsRocketMQInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRocketmqInstance: Query DMS rocketmq instance
	var (
		getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstanceProduct = "dmsv2"
	)
	getRocketmqInstanceClient, err := cfg.NewServiceClient(getRocketmqInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	getRocketmqInstancePath := getRocketmqInstanceClient.Endpoint + getRocketmqInstanceHttpUrl
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", getRocketmqInstanceClient.ProjectID)
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", state.Primary.ID))

	getRocketmqInstanceOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getRocketmqInstanceResp, err := getRocketmqInstanceClient.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQInstance: %s", err)
	}
	return utils.FlattenResponse(getRocketmqInstanceResp)
}

func TestAccDmsRocketMQInstance_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "500"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_dms_rocketmq_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls_mode", "SSL"),
				),
			},
			{
				Config: testDmsRocketMQInstance_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key3", "value3_update"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "1200"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_dms_rocketmq_flavors.test", "flavors.1.id"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.name", "fileReservedTime"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.value", "72"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls_mode", "PLAINTEXT"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configs"},
			},
		},
	})
}

func TestAccDmsRocketMQInstance_prepaid_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_prepaid_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
				),
			},
			{
				Config: testDmsRocketMQInstance_prepaid_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key3", "value3_update"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "period", "period_unit"},
			},
		},
	})
}

func TestAccDmsRocketMQInstance_broker_publicip(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_broker_publicip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "300"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_dms_rocketmq_flavors.test", "flavors.0.id"),
				),
			},
			{
				Config: testDmsRocketMQInstance_broker_publicip_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "800"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
					resource.TestCheckResourceAttrPair(resourceName, "engine_version", "data.huaweicloud_dms_rocketmq_flavors.test", "versions.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_dms_rocketmq_flavors.test", "flavors.0.id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDmsRocketMQInstance_publicip(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_publicip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "false"),
				),
			},
			{
				Config: testDmsRocketMQInstance_publicip_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDmsRocketMQInstance_updateWithEpsId(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_withEpsId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testDmsRocketMQInstance_withEpsId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testDmsRocketMQInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  enable_acl        = true
  flavor_id         = local.flavor.id
  storage_space     = 500  
  broker_num        = 1
  tls_mode          = "SSL"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_update(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
  newFlavor     = data.huaweicloud_dms_rocketmq_flavors.test.flavors[1]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[2],
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  enable_acl        = false
  flavor_id         = local.newFlavor.id
  storage_space     = 1200 
  broker_num        = 2
  tls_mode          = "PLAINTEXT"

  configs {
    name  = "fileReservedTime"
    value = "72"
  }

  tags = {
    key1 = "value1_update"
    key2 = "value2_update"
    key3 = "value3_update"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_prepaid_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  storage_space     = 300
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 1
  enable_acl        = true

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_prepaid_update(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  storage_space     = 300
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[2],
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 1
  enable_acl        = false

  tags = {
    key1 = "value1_update"
    key2 = "value2_update"
    key3 = "value3_update"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_broker_publicip(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster.small"
}

resource "huaweicloud_vpc_eip" "test_eip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test_eip_${count.index}"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  count = 3
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
  publicip_id   = "${huaweicloud_vpc_eip.test_eip[0].id},${huaweicloud_vpc_eip.test_eip[1].id},${huaweicloud_vpc_eip.test_eip[2].id}"
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  enable_acl        = true
  flavor_id         = local.flavor.id
  storage_space     = 300  
  broker_num        = 1
  enable_publicip   = true
  publicip_id       = local.publicip_id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_broker_publicip_update(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster.small"
}

resource "huaweicloud_vpc_eip" "test_eip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test_eip_${count.index}"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  count = 6
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
  publicip_id   = format("%%s,%%s,%%s,%%s,%%s,%%s", huaweicloud_vpc_eip.test_eip[0].id,huaweicloud_vpc_eip.test_eip[1].id,
  huaweicloud_vpc_eip.test_eip[2].id,huaweicloud_vpc_eip.test_eip[3].id,huaweicloud_vpc_eip.test_eip[4].id,huaweicloud_vpc_eip.test_eip[5].id)
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  enable_acl        = false
  flavor_id         = local.flavor.id
  storage_space     = 800
  broker_num        = 2  
  enable_publicip   = true
  publicip_id       = local.publicip_id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_publicip(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "single.basic"
}
  
locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%[2]s"
  engine_version    = local.query_results.versions[length(local.query_results.versions)-1]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  flavor_id         = local.flavor.id
  storage_space     = 300
  enable_publicip   = false
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_publicip_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "single.basic"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%[2]s"
  engine_version    = local.query_results.versions[length(local.query_results.versions)-1]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  storage_spec_code = local.flavor.ios[0].storage_spec_code
  flavor_id         = local.flavor.id
  storage_space     = 300
  enable_publicip   = true
  publicip_id       = huaweicloud_vpc_eip.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name                  = "%s"
  engine_version        = "4.8.0"
  storage_space         = 300
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]
  
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = true

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}`, common.TestBaseNetwork(name), name, epsId)
}
