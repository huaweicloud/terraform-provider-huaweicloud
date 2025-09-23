package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dc"
)

func getResourceDcGlobalGatewayPeerLinkFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	return dc.GetPeerLinkDetail(client, state.Primary.Attributes["global_dc_gateway_id"], state.Primary.ID)
}

func TestAccResourceDcGlobalGatewayPeerLink_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dc_global_gateway_peer_link.test"
		randName     = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceDcGlobalGatewayPeerLinkFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			// Please configure a global gateway containing virtual interface, and there are no peer links below it.
			acceptance.TestAccPreCheckDcGlobalGatewayID(t)
			acceptance.TestAccPreCheckERInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceDcGlobalGatewayPeerLink_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "global_dc_gateway_id", acceptance.HW_DC_GLOBAL_GATEWAY_ID),
					resource.TestCheckResourceAttr(resourceName, "peer_site.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_site.0.gateway_id", acceptance.HW_ER_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "peer_site.0.project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "peer_site.0.region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.#", "1"),

					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.site_code"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "create_owner"),
				),
			},
			{
				Config: testResourceDcGlobalGatewayPeerLink_update1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.site_code"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "create_owner"),
				),
			},
			{
				Config: testResourceDcGlobalGatewayPeerLink_update2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.site_code"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_site.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "create_owner"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDcGlobalGatewayPeerLinkImportState(resourceName),
			},
		},
	})
}

func testResourceDcGlobalGatewayPeerLink_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway_peer_link" "test" {
  name                 = "%[1]s"
  global_dc_gateway_id = "%[2]s"
  description          = "test description"

  peer_site {
    gateway_id = "%[3]s"
    project_id = "%[4]s"
    region_id  = "%[5]s"
  }
}
`, name, acceptance.HW_DC_GLOBAL_GATEWAY_ID, acceptance.HW_ER_INSTANCE_ID, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testResourceDcGlobalGatewayPeerLink_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway_peer_link" "test" {
  name                 = "%[1]s_update"
  global_dc_gateway_id = "%[2]s"
  description          = "test description update"

  peer_site {
    gateway_id = "%[3]s"
    project_id = "%[4]s"
    region_id  = "%[5]s"
  }
}
`, name, acceptance.HW_DC_GLOBAL_GATEWAY_ID, acceptance.HW_ER_INSTANCE_ID, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testResourceDcGlobalGatewayPeerLink_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway_peer_link" "test" {
  name                 = "%[1]s_update"
  global_dc_gateway_id = "%[2]s"

  peer_site {
    gateway_id = "%[3]s"
    project_id = "%[4]s"
    region_id  = "%[5]s"
  }
}
`, name, acceptance.HW_DC_GLOBAL_GATEWAY_ID, acceptance.HW_ER_INSTANCE_ID, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testDcGlobalGatewayPeerLinkImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		gatewayID := rs.Primary.Attributes["global_dc_gateway_id"]
		if gatewayID == "" {
			return "", fmt.Errorf("attribute (global_dc_gateway_id) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", gatewayID, rs.Primary.ID), nil
	}
}
