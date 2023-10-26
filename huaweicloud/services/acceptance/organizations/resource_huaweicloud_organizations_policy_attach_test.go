package organizations

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPolicyAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPolicyAttach: Query Organizations policy attach
	var (
		getPolicyAttachHttpUrl = "v1/organizations/policies/{policy_id}/attached-entities"
		getPolicyAttachProduct = "organizations"
	)
	getPolicyAttachClient, err := cfg.NewServiceClient(getPolicyAttachProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	// Split policy_id and entity_id from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <policy_id>/<entity_id>")
	}
	policyId := parts[0]
	entityId := parts[1]

	getPolicyAttachPath := getPolicyAttachClient.Endpoint + getPolicyAttachHttpUrl
	getPolicyAttachPath = strings.ReplaceAll(getPolicyAttachPath, "{policy_id}", policyId)

	getPolicyAttachResp, err := pagination.ListAllItems(
		getPolicyAttachClient,
		"marker",
		getPolicyAttachPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy attach: %s", err)
	}

	getPolicyAttachRespJson, err := json.Marshal(getPolicyAttachResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy attach: %s", err)
	}
	var getPolicyAttachRespBody interface{}
	err = json.Unmarshal(getPolicyAttachRespJson, &getPolicyAttachRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy attach: %s", err)
	}

	attachedEntity := utils.PathSearch(fmt.Sprintf("attached_entities[?id=='%s']|[0]", entityId),
		getPolicyAttachRespBody, nil)

	if attachedEntity == nil {
		return nil, fmt.Errorf("error retrieving Organizations policy attach: %s", err)
	}

	return getPolicyAttachRespBody, nil
}

func TestAccPolicyAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_policy_attach.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicyAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_organizations_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "entity_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "entity_name"),
					resource.TestCheckResourceAttrSet(rName, "entity_type"),
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

func testPolicyAttach_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = huaweicloud_organizations_organizational_unit.test.id
}
`, testPolicy_basic(name), testOrganizationalUnit_basic(name))
}
