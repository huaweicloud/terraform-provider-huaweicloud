package modelarts

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

func getModelartsNetworkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getModelartsNetwork: Query the Modelarts network.
	var (
		getModelartsNetworkHttpUrl = "v1/{project_id}/networks/{id}"
		getModelartsNetworkProduct = "modelarts"
	)
	getModelartsNetworkClient, err := cfg.NewServiceClient(getModelartsNetworkProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsNetworkPath := getModelartsNetworkClient.Endpoint + getModelartsNetworkHttpUrl
	getModelartsNetworkPath = strings.ReplaceAll(getModelartsNetworkPath, "{project_id}", getModelartsNetworkClient.ProjectID)
	getModelartsNetworkPath = strings.ReplaceAll(getModelartsNetworkPath, "{id}", state.Primary.ID)

	getModelartsNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsNetworkResp, err := getModelartsNetworkClient.Request("GET", getModelartsNetworkPath, &getModelartsNetworkOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts network: %s", err)
	}

	getModelartsNetworkRespBody, err := utils.FlattenResponse(getModelartsNetworkResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts network: %s", err)
	}

	return getModelartsNetworkRespBody, nil
}

func TestAccModelartsNetwork_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_modelarts_network.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsNetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsNetwork_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "0"),
				),
			},
			{
				Config: testModelartsNetwork_basic_update(name), // add a connection
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testModelartsNetwork_basic(name), // remove a connection
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "0"),
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

func testModelartsNetwork_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_network" "test" {
  name = "%s"
  cidr = "192.168.20.0/24"

  depends_on = [huaweicloud_vpc.test, huaweicloud_vpc_subnet.test]
}
`, common.TestVpc(name), name)
}

func testModelartsNetwork_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_network" "test" {
  name = "%s"
  cidr = "192.168.20.0/24"

  peer_connections {
    vpc_id    = huaweicloud_vpc.test.id
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  depends_on = [huaweicloud_vpc.test, huaweicloud_vpc_subnet.test]
}
`, common.TestVpc(name), name)
}
