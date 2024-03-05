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

func getCloudConnectionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCloudConnection: Query the Cloud Connection
	var (
		getCloudConnectionHttpUrl = "v3/{domain_id}/ccaas/cloud-connections/{id}"
		getCloudConnectionProduct = "cc"
	)
	getCloudConnectionClient, err := conf.NewServiceClient(getCloudConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CloudConnection Client: %s", err)
	}

	getCloudConnectionPath := getCloudConnectionClient.Endpoint + getCloudConnectionHttpUrl
	getCloudConnectionPath = strings.ReplaceAll(getCloudConnectionPath, "{domain_id}", conf.DomainID)
	getCloudConnectionPath = strings.ReplaceAll(getCloudConnectionPath, "{id}", state.Primary.ID)

	getCloudConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCloudConnectionResp, err := getCloudConnectionClient.Request("GET", getCloudConnectionPath, &getCloudConnectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CloudConnection: %s", err)
	}
	return utils.FlattenResponse(getCloudConnectionResp)
}

func TestAccCloudConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCloudConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCloudConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "used_scene", "vpc"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.sam", "tool"),
					resource.TestCheckResourceAttrSet(rName, "network_instance_number"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_package_number"),
					resource.TestCheckResourceAttrSet(rName, "inter_region_bandwidth_number"),
				),
			},
			{
				Config: testCloudConnection_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
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

func testCloudConnection_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    sam = "tool"
  }
}
`, name)
}

func testCloudConnection_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "0"
  description           = "demo_description"

  tags = {
    key = "value"
  }
}
`, name)
}

func TestAccCloudConnection_withEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCloudConnectionResourceFunc,
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
				Config: testCloudConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "used_scene", "vpc"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(rName, "network_instance_number"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_package_number"),
					resource.TestCheckResourceAttrSet(rName, "inter_region_bandwidth_number"),
				),
			},
			{
				Config: testCloudConnection_updateWithEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
				),
			},
		},
	})
}

func testCloudConnection_updateWithEpsId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "%s"
  description           = "demo_description"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
