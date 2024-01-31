package cc

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

func getBandwidthPackageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBandwidthPackage: Query the bandwidth package
	var (
		getBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
		getBandwidthPackageProduct = "cc"
	)
	getBandwidthPackageClient, err := cfg.NewServiceClient(getBandwidthPackageProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC Client: %s", err)
	}

	getBandwidthPackagePath := getBandwidthPackageClient.Endpoint + getBandwidthPackageHttpUrl
	getBandwidthPackagePath = strings.ReplaceAll(getBandwidthPackagePath, "{domain_id}", cfg.DomainID)
	getBandwidthPackagePath = strings.ReplaceAll(getBandwidthPackagePath, "{id}", state.Primary.ID)

	getBandwidthPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getBandwidthPackageResp, err := getBandwidthPackageClient.Request("GET", getBandwidthPackagePath, &getBandwidthPackageOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving bandwidth package: %s", err)
	}

	getBandwidthPackageRespBody, err := utils.FlattenResponse(getBandwidthPackageResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving bandwidth package: %s", err)
	}

	return getBandwidthPackageRespBody, nil
}

func TestAccBandwidthPackage_basic(t *testing.T) {
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
				Config: testBandwidthPackage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "5"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testBandwidthPackage_basic_update(name + "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "6"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform_test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_type", "cloud_connection"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_cc_connection.test", "id")),
			},
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
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_id", ""),
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
  name           = "%s"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is an accaptance test"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testBandwidthPackage_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = "%[1]s"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 6
  description    = "This is an accaptance test update"
  resource_id    = huaweicloud_cc_connection.test.id
  resource_type  = "cloud_connection"

  tags = {
    foo   = "bar"
    owner = "terraform_test"
  }
}
`, name)
}

func TestAccBandwidthPackage_withEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_bandwidth_package.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBandwidthPackageResourceFunc,
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
				Config: testBandwidthPackage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "5"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testBandwidthPackage_updateWithEpsId(name + "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
					resource.TestCheckResourceAttr(rName, "local_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "remote_area_id", "Chinese-Mainland"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "billing_mode", "3"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "6"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform_test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resource_type", "cloud_connection"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_cc_connection.test", "id")),
			},
		},
	})
}

func testBandwidthPackage_updateWithEpsId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = "%[1]s"
  local_area_id         = "Chinese-Mainland"
  remote_area_id        = "Chinese-Mainland"
  charge_mode           = "bandwidth"
  billing_mode          = 3
  bandwidth             = 6
  description           = "This is an accaptance test update"
  enterprise_project_id = "%[2]s"
  resource_id           = huaweicloud_cc_connection.test.id
  resource_type         = "cloud_connection"

  tags = {
    foo   = "bar"
    owner = "terraform_test"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
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
