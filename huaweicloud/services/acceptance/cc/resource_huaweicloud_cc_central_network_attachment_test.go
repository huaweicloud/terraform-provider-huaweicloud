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

func getAttachmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getAttachmentHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{id}"
		getAttachmentProduct = "cc"
	)
	getAttachmentClient, err := cfg.NewServiceClient(getAttachmentProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	getAttachmentPath := getAttachmentClient.Endpoint + getAttachmentHttpUrl
	getAttachmentPath = strings.ReplaceAll(getAttachmentPath, "{domain_id}", cfg.DomainID)
	getAttachmentPath = strings.ReplaceAll(getAttachmentPath, "{central_network_id}",
		state.Primary.Attributes["central_network_id"])
	getAttachmentPath = strings.ReplaceAll(getAttachmentPath, "{id}", state.Primary.ID)

	getAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getAttachmentResp, err := getAttachmentClient.Request("GET", getAttachmentPath, &getAttachmentOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network attachment: %s", err)
	}

	getAttachmentRespBody, err := utils.FlattenResponse(getAttachmentResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network attachment: %s", err)
	}

	return getAttachmentRespBody, nil
}

func TestAccCentralNetworkAttachment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_central_network_attachment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAttachmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckCCGlobalGateway(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCentralNetworkAttachment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttrPair(rName, "enterprise_router_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_router_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "enterprise_router_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_id", acceptance.HW_CC_GLOBAL_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
				),
			},
			{
				Config: testCentralNetworkAttachment_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo update"),
					resource.TestCheckResourceAttrPair(rName, "enterprise_router_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_router_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "enterprise_router_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_id", acceptance.HW_CC_GLOBAL_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "global_dc_gateway_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCentralNetworkAttachmentImportState(rName),
			},
		},
	})
}

func testCentralNetworkAttachment_basic(name string) string {
	policy := testCentralNetworkPolicy_basic(name)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_central_network_attachment" "test" {
  name                         = "%[2]s"
  description                  = "This is a demo"
  central_network_id           = huaweicloud_cc_central_network.test.id
  enterprise_router_id         = huaweicloud_er_instance.test.id
  enterprise_router_project_id = "%[3]s"
  enterprise_router_region_id  = "%[4]s"
  global_dc_gateway_id         = "%[5]s"
  global_dc_gateway_project_id = "%[3]s"
  global_dc_gateway_region_id  = "%[4]s"
}
`, policy, name, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME, acceptance.HW_CC_GLOBAL_GATEWAY_ID)
}

func testCentralNetworkAttachment_basic_update(name string) string {
	policy := testCentralNetworkPolicy_basic(name)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_central_network_attachment" "test" {
  name                         = "%[2]s_update"
  description                  = "This is a demo update"
  central_network_id           = huaweicloud_cc_central_network.test.id
  enterprise_router_id         = huaweicloud_er_instance.test.id
  enterprise_router_project_id = "%[3]s"
  enterprise_router_region_id  = "%[4]s"
  global_dc_gateway_id         = "%[5]s"
  global_dc_gateway_project_id = "%[3]s"
  global_dc_gateway_region_id  = "%[4]s"
}
`, policy, name, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME, acceptance.HW_CC_GLOBAL_GATEWAY_ID)
}

func testCentralNetworkAttachmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["central_network_id"] == "" {
			return "", fmt.Errorf("attribute (central_network_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["central_network_id"] + "/" +
			rs.Primary.ID, nil
	}
}
