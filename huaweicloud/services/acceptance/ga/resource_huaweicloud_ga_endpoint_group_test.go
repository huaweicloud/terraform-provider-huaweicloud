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

func getEndpointGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getEndpointGroup: Query the GA Endpoint Group detail
	var (
		getEndpointGroupHttpUrl = "v1/endpoint-groups/{id}"
		getEndpointGroupProduct = "ga"
	)
	getEndpointGroupClient, err := conf.NewServiceClient(getEndpointGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EndpointGroup Client: %s", err)
	}

	getEndpointGroupPath := getEndpointGroupClient.Endpoint + getEndpointGroupHttpUrl
	getEndpointGroupPath = strings.ReplaceAll(getEndpointGroupPath, "{id}", state.Primary.ID)

	getEndpointGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getEndpointGroupResp, err := getEndpointGroupClient.Request("GET", getEndpointGroupPath, &getEndpointGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving EndpointGroup: %s", err)
	}
	return utils.FlattenResponse(getEndpointGroupResp)
}

func TestAccEndpointGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ga_endpoint_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEndpointGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEndpointGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "region_id", "cn-south-1"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "listeners.0.id", "huaweicloud_ga_listener.test", "id"),
				),
			},
			{
				Config: testEndpointGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
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

func testEndpointGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_endpoint_group" "test" {
  name        = "%s"
  description = "terraform test"
  region_id   = "cn-south-1"

  listeners {
    id = huaweicloud_ga_listener.test.id
  }
}
`, testListener_basic(name), name)
}

func testEndpointGroup_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_endpoint_group" "test" {
  name        = "%s-update"
  description = "terraform test update"
  region_id   = "cn-south-1"

  listeners {
    id = huaweicloud_ga_listener.test.id
  }
}
`, testListener_basic(name), name)
}
