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

func getCentralNetworkPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCentralNetworkPolicy: Query the central network policy
	var (
		getCentralNetworkPolicyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies?id={id}"
		getCentralNetworkPolicyProduct = "cc"
	)
	getCentralNetworkPolicyClient, err := cfg.NewServiceClient(getCentralNetworkPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPolicyPath := getCentralNetworkPolicyClient.Endpoint + getCentralNetworkPolicyHttpUrl
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{central_network_id}",
		state.Primary.Attributes["central_network_id"])
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{id}", state.Primary.ID)

	getCentralNetworkPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkPolicyResp, err := getCentralNetworkPolicyClient.Request("GET", getCentralNetworkPolicyPath, &getCentralNetworkPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CentralNetworkPolicy: %s", err)
	}

	getCentralNetworkPolicyRespBody, err := utils.FlattenResponse(getCentralNetworkPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CentralNetworkPolicy: %s", err)
	}

	jsonPath := fmt.Sprintf("central_network_policies[?id =='%s']|[0]", state.Primary.ID)
	getCentralNetworkPolicyRespBody = utils.PathSearch(jsonPath, getCentralNetworkPolicyRespBody, nil)
	if getCentralNetworkPolicyRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getCentralNetworkPolicyRespBody, nil
}

func TestAccCentralNetworkPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_central_network_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCentralNetworkPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "er_instances.0.enterprise_router_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttr(rName, "is_applied", "false"),
					resource.TestCheckResourceAttrSet(rName, "document_template_version"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCentralNetworkPolicyImportState(rName),
			},
		},
	})
}

func testCentralNetworkPolicy_basic(name string) string {
	return fmt.Sprintf(`

data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name                           = "%[1]s"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

 resource "huaweicloud_cc_central_network" "test" {
   name        = "%[1]s"
   description = "This is an accaptance test"
 }

resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id

  planes {
    associate_er_tables {
      project_id                 = "%[2]s"
      region_id                  = "%[3]s"
      enterprise_router_id       = huaweicloud_er_instance.test.id
      enterprise_router_table_id = huaweicloud_er_instance.test.default_association_route_table_id
    }
  }

  er_instances {
    project_id           = "%[2]s"
    region_id            = "%[3]s"
    enterprise_router_id = huaweicloud_er_instance.test.id
  }
}
`, name, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testCentralNetworkPolicyImportState(name string) resource.ImportStateIdFunc {
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
