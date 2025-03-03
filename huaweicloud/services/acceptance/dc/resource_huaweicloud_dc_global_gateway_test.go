package dc

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

func getResourceDcGlobalGatewayFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", state.Primary.ID)
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DC global gateway: %s", err)
	}
	return utils.FlattenResponse(resp)
}

func TestAccResourceDcGlobalGateway_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dc_global_gateway.test"
		randName     = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceDcGlobalGatewayFunc,
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
				Config: testResourceDcGlobalGateway_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "current_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "available_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "all_tags.%"),
				),
			},
			{
				Config: testResourceDcGlobalGateway_basic_update1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "address_family", "dual"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "current_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "available_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "all_tags.%"),
				),
			},
			{
				Config: testResourceDcGlobalGateway_basic_update2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "current_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "available_peer_link_count"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttrSet(resourceName, "all_tags.%"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The map type data is filtered before setting to local state.
				// The purpose of ignoring tags here is just to pass the test case.
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testResourceDcGlobalGateway_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway" "test" {
  name                  = "%[1]s"
  description           = "test description"
  bgp_asn               = 10
  enterprise_project_id = "%[2]s"
  address_family        = "ipv4"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testResourceDcGlobalGateway_basic_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway" "test" {
  name                  = "%[1]s_update"
  description           = "test description update"
  bgp_asn               = 10
  enterprise_project_id = "%[2]s"
  address_family        = "dual"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testResourceDcGlobalGateway_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway" "test" {
  name                  = "%[1]s_update"
  bgp_asn               = 10
  enterprise_project_id = "%[2]s"
  address_family        = "ipv4"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
