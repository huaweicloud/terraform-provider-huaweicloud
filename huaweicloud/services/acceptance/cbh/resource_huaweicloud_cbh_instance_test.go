package cbh

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbh"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCBHInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region             = acceptance.HW_REGION_NAME
		getInstanceHttpUrl = "v2/{project_id}/cbs/instance/list"
		getInstanceProduct = "cbh"
	)
	client, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBH client: %s", err)
	}

	getInstancePath := client.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", client.ProjectID)
	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceResp, err := client.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH instance list: %s", err)
	}
	getCbhInstancesRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", state.Primary.ID)
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return instance, nil
}

// Considering security vulnerabilities, the test cases do not include the binding and unbinding content of EIP.
func TestAccCBHInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cbh_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCBHInstanceResourceFunc,
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
				Config: testCBHInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "flavor_id", "cbh.basic.10"),
					resource.TestCheckResourceAttr(rName, "password", "test_123456"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "period_unit", "month"),
					resource.TestCheckResourceAttr(rName, "period", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(rName, "subnet_address", "192.168.0.154"),
					resource.TestCheckResourceAttr(rName, "public_ip_id", ""),
					resource.TestCheckResourceAttr(rName, "public_ip", ""),
					// The built-in disk for this flavor instance is **0.2TB**, increase disk size by **1TB** through
					// the `attach_disk_size` parameter.
					resource.TestCheckResourceAttr(rName, "data_disk_size", "1.2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),

					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),

					resource.TestCheckResourceAttrSet(rName, "security_group_id"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
			{
				Config: testCBHInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "password", "test_147258"),
					resource.TestCheckResourceAttr(rName, "flavor_id", "cbh.basic.20"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(rName, "data_disk_size", "3.2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),

					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode",
					"password",
					"period",
					"period_unit",
					"auto_renew",
					"ipv6_enable",
					"attach_disk_size",
				},
			},
		},
	})
}

func TestAccCBHInstance_WithPowerAction(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cbh_instance.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCBHInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBHInstance_doAction_config(name, cbh.Stop),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testCBHInstance_doAction_config(name, cbh.Start),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testCBHInstance_doAction_config(name, cbh.SoftReboot),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testCBHInstance_doAction_config(name, cbh.HardReboot),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode",
					"password",
					"period",
					"period_unit",
					"power_action",
				},
			},
		},
	})
}

func testCBHInstance_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test2" {
  name                 = "%[2]s_2"
  delete_default_rules = true
}

data "huaweicloud_availability_zones" "test" {}
`, common.TestBaseNetwork(name), name)
}

func testCBHInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id             = "cbh.basic.10"
  name                  = "%s"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  subnet_address        = "192.168.0.154"
  security_group_id     = join(",", [huaweicloud_networking_secgroup.test.id, huaweicloud_networking_secgroup.test2.id])
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  password              = "test_123456"
  charging_mode         = "prePaid"
  period_unit           = "month"
  auto_renew            = "false"
  period                = 1
  attach_disk_size      = 1
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testCBHInstance_base(name), name)
}

func testCBHInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id             = "cbh.basic.20"
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  subnet_address        = "192.168.0.154"
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  password              = "test_147258"
  charging_mode         = "prePaid"
  period_unit           = "month"
  auto_renew            = "true"
  period                = 1
  attach_disk_size      = 2
  enterprise_project_id = "%[3]s"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, testCBHInstance_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testCBHInstance_doAction_config(name, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id             = "cbh.basic.10"
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  password              = "test_123456"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  power_action          = "%[3]s"
}
`, testCBHInstance_base(name), name, action)
}

func TestAccCBHInstance_updateVpc(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cbh_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCBHInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		// 1.Because the CBH instance will automatically create an extended elastic network card and bind it to the
		// instance and security group. After switching to VPC, the original extended elastic network card will not be
		// deleted, so it cannot be completed through automated test cases. So use environment variables to inject
		// `vpc_id`, `subnet_id`, and `security_group_id`.
		// 2.After updating the 'subnet_address' parameter, the elastic network card corresponding to the original
		// subnet address will remain, and you need to manually delete it in the console.
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVpcId(t)
			acceptance.TestAccPreCheckSubnetId(t)
			acceptance.TestAccPreCheckSecurityGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBHInstance_updateVpc_stp1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "vpc_id", acceptance.HW_VPC_ID),
					resource.TestCheckResourceAttr(rName, "subnet_id", acceptance.HW_SUBNET_ID),
					resource.TestCheckResourceAttr(rName, "security_group_id", acceptance.HW_SECURITY_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "subnet_address", "192.168.0.154"),
				),
			},
			{
				Config: testCBHInstance_updateVpc_stp2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "subnet_address", "192.168.0.155"),

					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode",
					"password",
					"period",
					"period_unit",
				},
			},
		},
	})
}

func testCBHInstance_updateVpc_stp1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cbh_flavors" "test" {
  type = "basic"
}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = data.huaweicloud_cbh_flavors.test.flavors[0].id
  name              = "%[1]s"
  vpc_id            = "%[2]s"
  subnet_id         = "%[3]s"
  subnet_address    = "192.168.0.154"
  security_group_id = "%[4]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "test_123456"
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}
`, name, acceptance.HW_VPC_ID, acceptance.HW_SUBNET_ID, acceptance.HW_SECURITY_GROUP_ID)
}

func testCBHInstance_updateVpc_stp2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cbh_flavors" "test" {
  type = "basic"
}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = data.huaweicloud_cbh_flavors.test.flavors[0].id
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  subnet_address    = "192.168.0.155"
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "test_123456"
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}
`, common.TestBaseNetwork(name), name)
}
