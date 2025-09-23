package identitycenter

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

func getSystemIdentityPolicyAttachmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	listSystemPoliciesClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	listSystemPoliciesHttpUrl := "v1/instances/{instance_id}/permission-sets/{permission_set_id}/managed-policies"
	listSystemPoliciesPath := listSystemPoliciesClient.Endpoint + listSystemPoliciesHttpUrl
	listSystemPoliciesPath = strings.ReplaceAll(listSystemPoliciesPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	listSystemPoliciesPath = strings.ReplaceAll(listSystemPoliciesPath, "{permission_set_id}", state.Primary.Attributes["permission_set_id"])

	listSystemPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listSystemPoliciesResp, err := listSystemPoliciesClient.Request("GET", listSystemPoliciesPath, &listSystemPoliciesOpt)
	if err != nil {
		return nil, err
	}

	listSystemPoliciesRespBody, err := utils.FlattenResponse(listSystemPoliciesResp)
	if err != nil {
		return nil, err
	}

	attachPolicies := utils.PathSearch("attached_managed_policies", listSystemPoliciesRespBody, make([]interface{}, 0))
	if len(attachPolicies.([]interface{})) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return attachPolicies, nil
}

func TestAccSystemIdentityPolicyAttachment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_system_identity_policy_attachment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSystemIdentityPolicyAttachmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIdentityCenterIdentiyPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSystemIdentityPolicyAttachment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "policy_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "attached_policies.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
				),
			},
			{
				Config: testSystemIdentiyPolicyAttachment_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "policy_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "attached_policies.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testIdentityPolicyAttachmentImportState(rName),
				ImportStateVerify: true,
			},
		},
	})
}

func testIdentityPolicyAttachmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found: %s", name, rs)
		}

		return instanceID + "/" + rs.Primary.ID, nil
	}
}

func testSystemIdentityPolicyAttachment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_system_identity_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id

  policy_ids = [%[2]s]
}
`, testPermissionSet_basic(name), acceptance.HW_IDENTITY_CENTER_IDENTITY_POLICY_ID)
}

func testSystemIdentiyPolicyAttachment_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_system_identity_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id

  policy_ids = [%[2]s]
}
`, testPermissionSet_basic(name), acceptance.HW_IDENTITY_CENTER_IDENTITY_POLICY_ID_UPDATE)
}
