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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getModelartsResourcePoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getModelartsResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		getModelartsResourcePoolProduct = "modelarts"
	)
	getModelartsResourcePoolClient, err := cfg.NewServiceClient(getModelartsResourcePoolProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsResourcePoolPath := getModelartsResourcePoolClient.Endpoint + getModelartsResourcePoolHttpUrl
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{project_id}", getModelartsResourcePoolClient.ProjectID)
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{id}", state.Primary.ID)

	getModelartsResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsResourcePoolResp, err := getModelartsResourcePoolClient.Request("GET", getModelartsResourcePoolPath, &getModelartsResourcePoolOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts resource pool: %s", err)
	}

	getModelartsResourcePoolRespBody, err := utils.FlattenResponse(getModelartsResourcePoolResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts resource pool: %s", err)
	}

	return getModelartsResourcePoolRespBody, nil
}

func TestAccModelartsResourcePool_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_modelarts_resource_pool.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsResourcePoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsResourcePool_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "scope.#", "3"),
					resource.TestCheckResourceAttrPair(rName, "network_id",
						"huaweicloud_modelarts_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
				),
			},
			{
				Config: testModelartsResourcePool_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo update"),
					resource.TestCheckResourceAttr(rName, "scope.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "network_id",
						"huaweicloud_modelarts_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
					resource.TestCheckResourceAttr(rName, "resources.1.flavor_id", "modelarts.vm.cpu.16u64g.d"),
					resource.TestCheckResourceAttr(rName, "resources.1.count", "1"),
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

func testModelartsResourcePool_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%s"
  description = "This is a demo"
  scope       = ["Train", "Infer", "Notebook"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = "modelarts.vm.cpu.8ud"
    count     = 1
  }
}
`, testModelartsResourcePool_base(name), name)
}

func testModelartsResourcePool_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%s"
  description = "This is a demo update"
  scope       = ["Train", "Infer"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = "modelarts.vm.cpu.8ud"
    count     = 1
  }

  resources {
    flavor_id = "modelarts.vm.cpu.16u64g.d"
    count     = 1
  }
}
`, testModelartsResourcePool_base(name), name)
}

func testModelartsResourcePool_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_network" "test" {
  name = "%s"
  cidr = "172.16.0.0/12"
}`, name)
}
