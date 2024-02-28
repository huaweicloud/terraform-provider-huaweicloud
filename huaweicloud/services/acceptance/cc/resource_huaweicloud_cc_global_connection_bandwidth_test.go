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

func getGCBResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}
	getGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	getGCBPath := client.Endpoint + getGCBHttpUrl
	getGCBPath = strings.ReplaceAll(getGCBPath, "{domain_id}", cfg.DomainID)
	getGCBPath = strings.ReplaceAll(getGCBPath, "{id}", state.Primary.ID)

	getGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGCBResp, err := client.Request("GET", getGCBPath, &getGCBOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global connection bandwidth: %s", err)
	}
	return utils.FlattenResponse(getGCBResp)
}

func TestAccGCB_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_global_connection_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGCBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGCB_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "Region"),
					resource.TestCheckResourceAttr(rName, "bordercross", "false"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bwd"),
					resource.TestCheckResourceAttr(rName, "size", "5"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "sla_level", "Ag"),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "frozen"),
				),
			},
			{
				Config: testAccGCB_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "95"),
					resource.TestCheckResourceAttr(rName, "size", "100"),
					resource.TestCheckResourceAttr(rName, "description", "test-update"),
					resource.TestCheckResourceAttr(rName, "sla_level", "Au"),
					resource.TestCheckResourceAttr(rName, "binding_service", "GEIP"),
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

func testAccGCB_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name                  = "%[1]s"
  type                  = "Region"  
  bordercross           = false
  charge_mode           = "bwd"
  size                  = 5
  enterprise_project_id = "%[2]s"
  description           = "test"
  sla_level             = "Ag"

  tags = {
    foo = "bar"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGCB_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name                  = "%[1]s-update"
  type                  = "Region"  
  bordercross           = false
  charge_mode           = "95"
  size                  = 100
  enterprise_project_id = "%[2]s"
  description           = "test-update"
  sla_level             = "Au"
  binding_service       = "GEIP"

  tags = {
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
