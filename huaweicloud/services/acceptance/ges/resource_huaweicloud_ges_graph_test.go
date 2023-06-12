package ges

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

func getGesGraphResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGraph: Query the GES graph.
	var (
		getGraphHttpUrl = "v2/{project_id}/graphs/{id}"
		getGraphProduct = "ges"
	)
	getGraphClient, err := cfg.NewServiceClient(getGraphProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GES Client: %s", err)
	}

	getGraphPath := getGraphClient.Endpoint + getGraphHttpUrl
	getGraphPath = strings.ReplaceAll(getGraphPath, "{project_id}", getGraphClient.ProjectID)
	getGraphPath = strings.ReplaceAll(getGraphPath, "{id}", state.Primary.ID)

	getGraphOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getGraphResp, err := getGraphClient.Request("GET", getGraphPath, &getGraphOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GesGraph: %s", err)
	}

	getGraphRespBody, err := utils.FlattenResponse(getGraphResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GesGraph: %s", err)
	}

	return getGraphRespBody, nil
}

func TestAccGesGraph_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ges_graph.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGesGraphResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGesGraph_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "graph_size_type_index", "1"),
					resource.TestCheckResourceAttr(rName, "cpu_arch", "x86_64"),
					resource.TestCheckResourceAttr(rName, "crypt_algorithm", "generalCipher"),
					resource.TestCheckResourceAttr(rName, "enable_https", "false"),
					resource.TestCheckResourceAttrSet(rName, "az_code"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "traffic_ip_list.#"),
				),
			},
			{
				Config: testGesGraph_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "graph_size_type_index", "2"),
					resource.TestCheckResourceAttr(rName, "cpu_arch", "x86_64"),
					resource.TestCheckResourceAttr(rName, "crypt_algorithm", "generalCipher"),
					resource.TestCheckResourceAttr(rName, "enable_https", "false"),
					resource.TestCheckResourceAttrSet(rName, "az_code"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "traffic_ip_list.#"),
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

func testGesGraph_basic(name string) string {
	baseNetwork := common.TestBaseNetwork(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_ges_graph" "test" {
  name                  = "%s"
  graph_size_type_index = "1"
  cpu_arch              = "x86_64"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = "generalCipher"
  enable_https          = false

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, baseNetwork, name)
}

func testGesGraph_basic_update(name string) string {
	baseNetwork := common.TestBaseNetwork(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_ges_graph" "test" {
  name                  = "%s"
  graph_size_type_index = "2"
  cpu_arch              = "x86_64"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = "generalCipher"
  enable_https          = false

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, baseNetwork, name)
}
