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

func getAcceleratorResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/accelerators/{id}"
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
		return nil, fmt.Errorf("error retrieving GA accelerator: %s", err)
	}
	return utils.FlattenResponse(resp)
}

func TestAccAccelerator_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_ga_accelerator.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAcceleratorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccelerator_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "ip_sets.0.area", "CM"),
					resource.TestCheckResourceAttr(rName, "ip_sets.0.ip_type", "IPV4"),
					resource.TestCheckResourceAttr(rName, "ip_sets.1.area", "CM"),
					resource.TestCheckResourceAttr(rName, "ip_sets.1.ip_type", "IPV6"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "ip_sets.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "ip_sets.1.ip_address"),
				),
			},
			{
				Config: testAccelerator_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
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

func testAccelerator_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%s"
  description = "terraform test"

  ip_sets {
    ip_type = "IPV4"
    area    = "CM"
  }

  ip_sets {
    ip_type = "IPV6"
    area    = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccelerator_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%s-update"
  description = "terraform test update"

  ip_sets {
    ip_type = "IPV4"
    area    = "CM"
  }

  ip_sets {
    ip_type = "IPV6"
    area    = "CM"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, name)
}
