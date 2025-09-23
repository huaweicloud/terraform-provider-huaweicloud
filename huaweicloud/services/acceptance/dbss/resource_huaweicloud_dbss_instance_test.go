package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"
)

func getInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "dbss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DBSS client: %s", err)
	}

	return dbss.QueryTargetDBSSInstance(client, state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var (
		obj        interface{}
		name       = acceptance.RandomAccResourceName()
		updataName = acceptance.RandomAccResourceName()
		rName      = "huaweicloud_dbss_instance.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dbss_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", "dbss.bypassaudit.low"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"data.huaweicloud_networking_secgroup.test", "id"),
				),
			},
			{
				Config: testInstance_update_1(updataName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updataName),
					resource.TestCheckResourceAttr(rName, "status", "SHUTOFF"),
					resource.TestCheckResourceAttr(rName, "description", "test desc"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "test"),
					resource.TestCheckResourceAttr(rName, "tags.acc", "value"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(rName, "action", "stop"),
				),
			},
			{
				Config: testInstance_update_2(updataName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "action", "start"),
				),
			},
			{
				Config: testInstance_update_3(updataName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "action", "reboot"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode", "enterprise_project_id", "flavor", "period", "period_unit",
					"product_spec_desc", "tags", "action",
				},
			},
		},
	})
}

func testInstance_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name   = "subnet-default"
  vpc_id = data.huaweicloud_vpc.test.id
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dbss_flavors" "test" {} 

locals {
  vpc_name    = data.huaweicloud_vpc.test.name
  subnet_name = data.huaweicloud_vpc_subnet.test.name

  # The splicing specification of this field lacks documentation. Please refer to the following format when writing scripts.
  product_spec_desc = jsonencode(
    {
      "specDesc" : {
        "zh-cn" : {
          "主机名称" : "%[1]s",
          "虚拟私有云" : local.vpc_name,
          "子网" : local.subnet_name
        },
        "en-us" : {
          "Instance Name" : "%[1]s",
          "VPC" : local.vpc_name,
          "Subnet" : local.subnet_name
        }
      }
    }
  )
}
`, name)
}

func testInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dbss_instance" "test" {
  name               = "%s"
  description        = "terraform test"
  flavor             = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code = "dbss.bypassaudit.low"
  product_spec_desc  = local.product_spec_desc
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = data.huaweicloud_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testInstance_base(name), name)
}

func testInstance_update_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[2]s"
  description = "security group acceptance test"
}

resource "huaweicloud_dbss_instance" "test" {
  name               = "%[2]s"
  description        = "test desc"
  flavor             = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code = "dbss.bypassaudit.low"
  product_spec_desc  = local.product_spec_desc
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
  action             = "stop"

  tags = {
    foo = "test"
    acc = "value"
  }
}
`, testInstance_base(name), name)
}

func testInstance_update_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[2]s"
  description = "security group acceptance test"
}

resource "huaweicloud_dbss_instance" "test" {
  name               = "%[2]s"
  description        = "test desc"
  flavor             = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code = "dbss.bypassaudit.low"
  product_spec_desc  = local.product_spec_desc
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
  action             = "start"

  tags = {
    foo = "test"
    acc = "value"
  }
}
`, testInstance_base(name), name)
}

func testInstance_update_3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[2]s"
  description = "security group acceptance test"
}

resource "huaweicloud_dbss_instance" "test" {
  name               = "%[2]s"
  description        = "test desc"
  flavor             = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code = "dbss.bypassaudit.low"
  product_spec_desc  = local.product_spec_desc
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
  action             = "reboot"

  tags = {
    foo = "test"
    acc = "value"
  }
}
`, testInstance_base(name), name)
}

func TestAccInstance_updateWithEpsId(t *testing.T) {
	var (
		obj     interface{}
		name    = acceptance.RandomAccResourceName()
		rName   = "huaweicloud_dbss_instance.test"
		srcEPS  = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		destEPS = acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_withEpsId(name, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testInstance_withEpsId(name, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testInstance_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dbss_instance" "test" {
  name                   = "%s"
  description            = "terraform test"
  flavor                 = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code     = "dbss.bypassaudit.low"
  product_spec_desc      = local.product_spec_desc
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                 = data.huaweicloud_vpc.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  charging_mode          = "prePaid"
  period_unit            = "month"
  period                 = 1
  enterprise_project_id  = "%s"
}
`, testInstance_base(name), name, epsId)
}
