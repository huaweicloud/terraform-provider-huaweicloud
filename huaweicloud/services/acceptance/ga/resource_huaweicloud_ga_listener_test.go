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
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/listeners/{id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GA listener: %s", err)
	}
	return utils.FlattenResponse(resp)
}

func TestAccListener_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_ga_listener.test"
	)

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
