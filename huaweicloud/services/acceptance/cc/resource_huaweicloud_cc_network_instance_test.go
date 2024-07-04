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

func getNetworkInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getNetworkInstance: Query the Network instance
	var (
		getNetworkInstanceHttpUrl = "v3/{domain_id}/ccaas/network-instances/{id}"
		getNetworkInstanceProduct = "cc"
	)
	getNetworkInstanceClient, err := conf.NewServiceClient(getNetworkInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating NetworkInstance Client: %s", err)
	}

	getNetworkInstancePath := getNetworkInstanceClient.Endpoint + getNetworkInstanceHttpUrl
	getNetworkInstancePath = strings.ReplaceAll(getNetworkInstancePath, "{domain_id}", conf.DomainID)
	getNetworkInstancePath = strings.ReplaceAll(getNetworkInstancePath, "{id}", state.Primary.ID)

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
				Config: testNetworkInstance_basic(name, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "region_id", "huaweicloud_vpc.test.0", "region"),
					resource.TestCheckResourceAttrPair(rName, "cloud_connection_id", "huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
				),
			},
			{
				Config: testNetworkInstance_basic_update(name, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func TestAccNetworkInstance_multiple(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_network_instance.test2"

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
				Config: testNetworkInstance_multiple(name, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr("huaweicloud_cc_network_instance.test1", "type", "vpc"),
					resource.TestCheckResourceAttr("huaweicloud_cc_network_instance.test2", "type", "vpc"),
					resource.TestCheckResourceAttr("huaweicloud_cc_network_instance.test1", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("huaweicloud_cc_network_instance.test2", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("huaweicloud_cc_network_instance.test1", "domain_id"),
					resource.TestCheckResourceAttrSet("huaweicloud_cc_network_instance.test2", "domain_id"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test1", "instance_id",
						"huaweicloud_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test2", "instance_id",
						"huaweicloud_vpc.test.1", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test1", "region_id",
						"huaweicloud_vpc.test.0", "region"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test2", "region_id",
						"huaweicloud_vpc.test.1", "region"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test1", "cloud_connection_id",
						"huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_cc_network_instance.test2", "cloud_connection_id",
						"huaweicloud_cc_connection.test", "id"),
				),
			},
		},
	})
}

func testNetworkInstanceRef(name string, count int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count = %[1]d

  name = "%[2]s_${count.index}"
  cidr = cidrsubnet("10.12.0.0/16", 4, count.index)
}

resource "huaweicloud_vpc_subnet" "test" {
  count = %[1]d

  name = "%[2]s_${count.index}"
  vpc_id     = huaweicloud_vpc.test[count.index].id
  cidr       = cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 4, 1), 1)
}

resource "huaweicloud_cc_connection" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
  description           = "accDemo"
}
`, count, name)
}

func testNetworkInstance_basic(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = try(huaweicloud_vpc.test[0].id, "")
  project_id          = "%[2]s"
  region_id           = try(huaweicloud_vpc.test[0].region, "")
  description         = "demo_description"

  cidrs = [
    try(huaweicloud_vpc_subnet.test[0].cidr, ""),
  ]
}
`, testNetworkInstanceRef(name, count), acceptance.HW_PROJECT_ID)
}

func testNetworkInstance_basic_update(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = try(huaweicloud_vpc.test[0].id, "")
  project_id          = "%[2]s"
  region_id           = try(huaweicloud_vpc.test[0].region, "")

  cidrs = [
    try(huaweicloud_vpc_subnet.test[0].cidr, ""),
  ]
}
`, testNetworkInstanceRef(name, count), acceptance.HW_PROJECT_ID)
}

func testNetworkInstance_multiple(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_network_instance" "test1" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test[0].id
  project_id          = "%[2]s"
  region_id           = huaweicloud_vpc.test[0].region
  
  cidrs = [
    huaweicloud_vpc_subnet.test[0].cidr,
  ]
}
  
resource "huaweicloud_cc_network_instance" "test2" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test[1].id
  project_id          = "%[2]s"
  region_id           = huaweicloud_vpc.test[1].region
  
  cidrs = [
    huaweicloud_vpc_subnet.test[1].cidr,
  ]

  depends_on = [
    huaweicloud_cc_network_instance.test1,
  ]
}
`, testNetworkInstanceRef(name, count), acceptance.HW_PROJECT_ID)
}
