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

func getCentralNetworkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCentralNetwork: Query the central network
	var (
		getCentralNetworkHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
		getCentralNetworkProduct = "cc"
	)
	getCentralNetworkClient, err := cfg.NewServiceClient(getCentralNetworkProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPath := getCentralNetworkClient.Endpoint + getCentralNetworkHttpUrl
	getCentralNetworkPath = strings.ReplaceAll(getCentralNetworkPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPath = strings.ReplaceAll(getCentralNetworkPath, "{id}", state.Primary.ID)

	getCentralNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkResp, err := getCentralNetworkClient.Request("GET", getCentralNetworkPath, &getCentralNetworkOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network: %s", err)
	}

	getCentralNetworkRespBody, err := utils.FlattenResponse(getCentralNetworkResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network: %s", err)
	}

	return getCentralNetworkRespBody, nil
}

func TestAccCentralNetwork_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_central_network.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCentralNetwork_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
				),
			},
			{
				Config: testCentralNetwork_basic_update(name + "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accaptance test update"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo2", "bar2"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
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

func testCentralNetwork_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_central_network" "test" {
  name        = "%s"
  description = "This is an accaptance test"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testCentralNetwork_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_central_network" "test" {
  name        = "%s"
  description = "This is an accaptance test update"

  tags = {
    foo2 = "bar2"
    key = "value"
  }
}
`, name)
}
