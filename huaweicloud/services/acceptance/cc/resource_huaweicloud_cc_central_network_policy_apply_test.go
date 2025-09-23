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

func getCentralNetworkPolicyApplyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getCentralNetworkPolicyApplyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies"
		getCentralNetworkPolicyApplyProduct = "cc"
	)
	getCentralNetworkPolicyApplyClient, err := cfg.NewServiceClient(getCentralNetworkPolicyApplyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPolicyApplyPath := getCentralNetworkPolicyApplyClient.Endpoint + getCentralNetworkPolicyApplyHttpUrl
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{central_network_id}",
		state.Primary.Attributes["central_network_id"])

	getCentralNetworkPolicyApplyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkPolicyApplyResp, err := getCentralNetworkPolicyApplyClient.Request("GET",
		getCentralNetworkPolicyApplyPath, &getCentralNetworkPolicyApplyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network policy: %s", err)
	}

	getCentralNetworkPolicyApplyRespBody, err := utils.FlattenResponse(getCentralNetworkPolicyApplyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network policy: %s", err)
	}

	jsonPath := fmt.Sprintf("central_network_policies[?id =='%s' && is_applied]|[0]", state.Primary.Attributes["policy_id"])
	getCentralNetworkPolicyApplyRespBody = utils.PathSearch(jsonPath, getCentralNetworkPolicyApplyRespBody, nil)
	if getCentralNetworkPolicyApplyRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getCentralNetworkPolicyApplyRespBody, nil
}

func TestAccCentralNetworkPolicyApply_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_central_network_policy_apply.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkPolicyApplyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCentralNetworkPolicyApply_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_cc_central_network_policy.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCentralNetworkPolicyApplyImportState(rName),
			},
		},
	})
}

func testCentralNetworkPolicyApply_basic(name string) string {
	policyConfig := testCentralNetworkPolicy_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test.id
}
`, policyConfig)
}

func testCentralNetworkPolicyApplyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		if rs.Primary.Attributes["policy_id"] == "" {
			return "", fmt.Errorf("attribute (policy_id) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.ID + "/" + rs.Primary.Attributes["policy_id"], nil
	}
}
