package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cc"
)

func getBandwidthPackageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product = "cc"
		region  = acceptance.HW_REGION_NAME
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	return cc.GetBandwidthPackage(client, cfg.DomainID, state.Primary.ID)
}

func TestAccBandwidthPackage_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cc_bandwidth_package.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBandwidthPackageResourceFunc,
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
				Config: testBandwidthPackage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "5"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
				),
			},
			{
				Config: testBandwidthPackage_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "6"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform_test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_type", "cloud_connection"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
				),
			},
			{
				Config: testBandwidthPackage_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "5"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_type", "cloud_connection"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_cc_connection.test_another", "id"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testBandwidthPackage_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = "%[1]s"
  local_area_id         = "Chinese-Mainland"
  remote_area_id        = "Chinese-Mainland"
  charge_mode           = "bandwidth"
  billing_mode          = 3
  bandwidth             = 5
  description           = "This is an accaptance test"
  enterprise_project_id = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testBandwidthPackage_basic_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = "%[1]s_update"
  local_area_id         = "Chinese-Mainland"
  remote_area_id        = "Chinese-Mainland"
  charge_mode           = "bandwidth"
  billing_mode          = 3
  bandwidth             = 6
  description           = "This is an accaptance test update"
  resource_id           = huaweicloud_cc_connection.test.id
  resource_type         = "cloud_connection"
  enterprise_project_id = "%[2]s"

  tags = {
    foo   = "bar"
    owner = "terraform_test"
  }
}
`, name, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testBandwidthPackage_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test_another" {
  name = "%[1]s_another"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = "%[1]s"
  local_area_id         = "Chinese-Mainland"
  remote_area_id        = "Chinese-Mainland"
  charge_mode           = "bandwidth"
  billing_mode          = 3
  bandwidth             = 5
  description           = "This is an accaptance test"
  resource_id           = huaweicloud_cc_connection.test_another.id
  resource_type         = "cloud_connection"
  enterprise_project_id = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func TestAccBandwidthPackage_regionalInterflow(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_bandwidth_package.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBandwidthPackageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBandwidthPackage_regionalInterflow(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_area_id", "cn-north-4"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "cn-south-1"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "5"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "4"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "value"),
					resource.TestCheckResourceAttr(rName, "interflow_mode", "Region"),
					resource.TestCheckResourceAttr(rName, "spec_code", "Beijing4toGuangzhou"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_type", "cloud_connection"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_cc_connection.test", "id"),
				),
			},
		},
	})
}

func testBandwidthPackage_regionalInterflow(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = "%[1]s"
  local_area_id         = "cn-north-4"
  remote_area_id        = "cn-south-1"
  charge_mode           = "bandwidth"
  billing_mode          = 5
  bandwidth             = 4
  resource_id           = huaweicloud_cc_connection.test.id
  resource_type         = "cloud_connection"
  interflow_mode        = "Region"
  spec_code             = "Beijing4toGuangzhou"

  tags = {
    foo   = "bar"
    owner = "value"
  }
}
`, name)
}
