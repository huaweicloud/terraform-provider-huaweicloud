package ddm

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

func getDdmInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query DDM instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/instances/{instance_id}"
		getInstanceProduct = "ddm"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{instance_id}", fmt.Sprintf("%v", state.Primary.ID))

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmInstance: %s", err)
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", getInstanceRespBody, nil)
	if status == "DELETED" {
		return nil, fmt.Errorf("error get DDM instance")
	}
	return getInstanceRespBody, nil
}

func TestAccDdmInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ddm_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.huaweicloud_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.huaweicloud_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "bind_table"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "test_value"),
					resource.TestCheckResourceAttrPair(rName, "enterprise_project_id",
						"huaweicloud_enterprise_project.test.0", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "4"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.huaweicloud_ddm_flavors.test", "flavors.1.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.huaweicloud_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "concurrent_execution_level"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "RDS_INSTANCE"),
					resource.TestCheckResourceAttrPair(rName, "enterprise_project_id",
						"huaweicloud_enterprise_project.test.1", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update_reduce(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.huaweicloud_ddm_flavors.test", "flavors.1.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.huaweicloud_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password", "engine_id", "flavor_id", "parameters"},
			},
		},
	})
}

func TestAccDdmInstance_prepaid(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ddm_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmInstance_prepaid(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.huaweicloud_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.huaweicloud_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "bind_table"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "test_value"),
				),
			},
			{
				Config: testDdmInstance_prepaid_update(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.huaweicloud_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.huaweicloud_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "concurrent_execution_level"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "RDS_INSTANCE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"admin_password", "engine_id", "flavor_id",
					"charging_mode", "auto_renew", "period", "period_unit", "parameters"},
			},
		},
	})
}

func testDdmInstance_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}
`, common.TestBaseNetwork(name))
}

func testDdmInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_enterprise_project" "test" {
  count = 2

  name = "%[2]s_${count.index}"
}

resource "huaweicloud_ddm_instance" "test" {
  name                  = "%[2]s"
  flavor_id             = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num              = 2
  engine_id             = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  admin_user            = "test_user_1"
  admin_password        = "test_password_123"
  enterprise_project_id = huaweicloud_enterprise_project.test[0].id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  parameters {
    name  = "bind_table"
    value = "test_value"
  }
}
`, testDdmInstance_base(name), name)
}

func testDdmInstance_basic_update(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_enterprise_project" "test" {
  count = 2

  name = "%[2]s_${count.index}"
}

resource "huaweicloud_networking_secgroup" "test_update" {
  name = "%[2]s"
}

resource "huaweicloud_ddm_instance" "test" {
  name                  = "%[2]s"
  flavor_id             = data.huaweicloud_ddm_flavors.test.flavors[1].id
  node_num              = 4
  engine_id             = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test_update.id
  admin_user            = "test_user_1"
  admin_password        = "test_password_456"
  enterprise_project_id = huaweicloud_enterprise_project.test[1].id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  parameters {
    name  = "concurrent_execution_level"
    value = "RDS_INSTANCE"
  }
}
`, testDdmInstance_base(name), updateName)
}

func testDdmInstance_basic_update_reduce(name, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test_update" {
  name = "%[2]s"
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[1].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test_update.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  parameters {
    name  = "concurrent_execution_level"
    value = "RDS_INSTANCE"
  }
}
`, testDdmInstance_base(name), updateName)
}

func testDdmInstance_prepaid(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  parameters {
    name  = "bind_table"
    value = "test_value"
  }
}
`, testDdmInstance_base(name), name)
}

func testDdmInstance_prepaid_update(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test_update" {
  name = "%[2]s"
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test_update.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_456"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  parameters {
    name  = "concurrent_execution_level"
    value = "RDS_INSTANCE"
  }
}
`, testDdmInstance_base(name), updateName)
}
