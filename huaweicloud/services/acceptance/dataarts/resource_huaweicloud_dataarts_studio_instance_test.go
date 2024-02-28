package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dayu/v1/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DataArtsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio v1 client, err=%s", err)
	}

	resp, err := instances.List(client, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range resp {
		if item.ID == state.Primary.ID {
			return &item, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestAccResourceInstance_basic(t *testing.T) {
	var dayuInstance instances.Instance
	resourceName := "huaweicloud_dataarts_studio_instance.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&dayuInstance,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)        // enterprise_project_id is required for this case
			acceptance.TestAccPreCheckChargingMode(t) // the resource only supports pre-paid charging mode
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "dayu.starter"),
					resource.TestCheckResourceAttr(resourceName, "status", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "order_id"),
					resource.TestCheckResourceAttrSet(resourceName, "expire_days"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags", "period_unit", "period", "auto_renew", "expire_days"},
			},
		},
	})
}

func testAccInstance_basic(rName string) string {
	baseNetwork := common.TestBaseNetwork(rName)

	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dataarts_studio_instance" "test" {
  name                  = "%s"
  version               = "dayu.starter"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = "%s"
  period_unit           = "month"
  period                = 1

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, baseNetwork, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccResourceInstance_updateWithEpsId(t *testing.T) {
	var dayuInstance instances.Instance
	resourceName := "huaweicloud_dataarts_studio_instance.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&dayuInstance,
		getInstanceResourceFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstance_withEpsId(name, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccInstance_withEpsId(name, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccInstance_withEpsId(rName, epsId string) string {
	baseNetwork := common.TestBaseNetwork(rName)

	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dataarts_studio_instance" "test" {
  name                  = "%s"
  version               = "dayu.starter"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = "%s"
  period_unit           = "month"
  period                = 1

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, baseNetwork, rName, epsId)
}
