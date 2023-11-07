package ga

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

func getEndpointResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getEndpoint: Query the GA Endpoint detail
	var (
		getEndpointHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{id}"
		getEndpointProduct = "ga"
	)
	getEndpointClient, err := config.NewServiceClient(getEndpointProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Endpoint Client: %s", err)
	}

	getEndpointPath := getEndpointClient.Endpoint + getEndpointHttpUrl
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{endpoint_group_id}", state.Primary.Attributes["endpoint_group_id"])
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{id}", state.Primary.ID)

	getEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getEndpointResp, err := getEndpointClient.Request("GET", getEndpointPath, &getEndpointOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Endpoint: %s", err)
	}
	return utils.FlattenResponse(getEndpointResp)
}

func TestAccEndpoint_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ga_endpoint.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEndpoint_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "endpoint_group_id", "huaweicloud_ga_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_address", "huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttr(rName, "weight", "1"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testEndpoint_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "weight", "10"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testEndpointImportState(rName),
			},
		},
	})
}

func testEndpoint_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  region = "cn-south-1"
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_ga_endpoint" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  resource_id       = huaweicloud_vpc_eip.test.id
  ip_address        = huaweicloud_vpc_eip.test.address
}
`, testEndpointGroup_basic(name), name)
}

func testEndpoint_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  region = "cn-south-1"
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_ga_endpoint" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  resource_id       = huaweicloud_vpc_eip.test.id
  ip_address        = huaweicloud_vpc_eip.test.address
  weight            = 10
}
`, testEndpointGroup_basic(name), name)
}

func testEndpointImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["endpoint_group_id"] == "" {
			return "", fmt.Errorf("Attribute (endpoint_group_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["endpoint_group_id"] + "/" +
			rs.Primary.ID, nil
	}
}
