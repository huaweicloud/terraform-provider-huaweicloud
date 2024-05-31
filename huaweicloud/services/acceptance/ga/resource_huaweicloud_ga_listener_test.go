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

func getListenerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getListener: Query the GA Listener detail
	var (
		getListenerHttpUrl = "v1/listeners/{id}"
		getListenerProduct = "ga"
	)
	getListenerClient, err := conf.NewServiceClient(getListenerProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Listener Client: %s", err)
	}

	getListenerPath := getListenerClient.Endpoint + getListenerHttpUrl
	getListenerPath = strings.ReplaceAll(getListenerPath, "{id}", state.Primary.ID)

	getListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getListenerResp, err := getListenerClient.Request("GET", getListenerPath, &getListenerOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Listener: %s", err)
	}
	return utils.FlattenResponse(getListenerResp)
}

func TestAccListener_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ga_listener.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testListener_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "description", "Terraform test"),
					resource.TestCheckResourceAttr(rName, "port_ranges.0.from_port", "4000"),
					resource.TestCheckResourceAttr(rName, "port_ranges.0.to_port", "4200"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "accelerator_id",
						"huaweicloud_ga_accelerator.test", "id"),
				),
			},
			{
				Config: testListener_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "description", "Terraform test update"),
					resource.TestCheckResourceAttr(rName, "port_ranges.0.from_port", "5000"),
					resource.TestCheckResourceAttr(rName, "port_ranges.0.to_port", "5200"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
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

func testListener_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_listener" "test" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = "%s"
  protocol       = "TCP"
  description    = "Terraform test"

  port_ranges {
    from_port = 4000
    to_port   = 4200
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccelerator_basic(name), name)
}

func testListener_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_listener" "test" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = "%s-update"
  protocol       = "TCP"
  description    = "Terraform test update"

  port_ranges {
    from_port = 5000
    to_port   = 5200
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccelerator_basic(name), name)
}
