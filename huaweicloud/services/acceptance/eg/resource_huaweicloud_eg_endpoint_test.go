package eg

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEndpointResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getEndpoint: Query the EG Endpoint detail
	var (
		getEndpointHttpUrl = "v1/{project_id}/endpoints"
		getEndpointProduct = "eg"
	)
	getEndpointClient, err := cfg.NewServiceClient(getEndpointProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	getEndpointPath := getEndpointClient.Endpoint + getEndpointHttpUrl
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{project_id}", getEndpointClient.ProjectID)

	getEndpointqueryParams := eg.BuildGetEndpointQueryParams(state.Primary.Attributes["name"])
	getEndpointPath += getEndpointqueryParams

	getEndpointResp, err := pagination.ListAllItems(
		getEndpointClient,
		"offset",
		getEndpointPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Endpoint: %s", err)
	}

	getEndpointRespJson, err := json.Marshal(getEndpointResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Endpoint: %s", err)
	}
	var getEndpointRespBody interface{}
	err = json.Unmarshal(getEndpointRespJson, &getEndpointRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Endpoint: %s", err)
	}

	jsonPath := fmt.Sprintf("items[?id =='%s']|[0]", state.Primary.ID)
	getEndpointRespBody = utils.PathSearch(jsonPath, getEndpointRespBody, nil)
	if getEndpointRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getEndpointRespBody, nil
}

func TestAccEndpoint_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_eg_endpoint.test"

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
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testEndpoint_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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

func testEndpoint_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eg_endpoint" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  description = "created by terraform"
}
`, common.TestVpc(name), name)
}

func testEndpoint_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eg_endpoint" "test" {
  name      = "%s"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}
`, common.TestVpc(name), name)
}
