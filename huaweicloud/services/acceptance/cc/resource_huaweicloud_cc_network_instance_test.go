package cc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getNetworkInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getNetworkInstance: Query the Network instance
	var (
		getNetworkInstanceHttpUrl = "v3/{domain_id}/ccaas/network-instances/{id}"
		getNetworkInstanceProduct = "cc"
	)
	getNetworkInstanceClient, err := config.NewServiceClient(getNetworkInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating NetworkInstance Client: %s", err)
	}

	getNetworkInstancePath := getNetworkInstanceClient.Endpoint + getNetworkInstanceHttpUrl
	getNetworkInstancePath = strings.Replace(getNetworkInstancePath, "{domain_id}", config.DomainID, -1)
	getNetworkInstancePath = strings.Replace(getNetworkInstancePath, "{id}", state.Primary.ID, -1)

	getNetworkInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getNetworkInstanceResp, err := getNetworkInstanceClient.Request("GET", getNetworkInstancePath, &getNetworkInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving NetworkInstance: %s", err)
	}
	return utils.FlattenResponse(getNetworkInstanceResp)
}

func TestAccNetworkInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_network_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getNetworkInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testNetworkInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "region_id", "huaweicloud_vpc.test", "region"),
					resource.TestCheckResourceAttrPair(rName, "cloud_connection_id",
						"huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testNetworkInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
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

func testNetworkInstanceRef(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "10.12.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "10.12.2.0/24"
  gateway_ip = "10.12.2.1"
}

resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "0"
  description           = "accDemo"
}
`, name, name, name)
}

func testNetworkInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test.id
  project_id          = "%s"
  region_id           = huaweicloud_vpc.test.region
  cidrs = [
    "10.12.2.0/24"
  ]
}
`, testNetworkInstanceRef(name), acceptance.HW_PROJECT_ID)
}

func testNetworkInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test.id
  project_id          = "%s"
  region_id           = huaweicloud_vpc.test.region
  description         = "demo_description"
  cidrs = [
    "10.12.2.0/24"
  ]
}
`, testNetworkInstanceRef(name), acceptance.HW_PROJECT_ID)
}
