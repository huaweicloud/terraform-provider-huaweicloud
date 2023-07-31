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

func getAccountAssignmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAccountAssignment: Query Identity Center account assignment
	var (
		getAccountAssignmentHttpUrl = "v1/instances/{instance_id}/account-assignments"
		getAccountAssignmentProduct = "identitycenter"
	)
	getAccountAssignmentClient, err := cfg.NewServiceClient(getAccountAssignmentProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	getAccountAssignmentBasePath := getAccountAssignmentClient.Endpoint + getAccountAssignmentHttpUrl
	getAccountAssignmentBasePath = strings.ReplaceAll(getAccountAssignmentBasePath, "{instance_id}",
		state.Primary.Attributes["instance_id"])

	getAccountAssignmentPath := getAccountAssignmentBasePath + buildGetAccountAssignmentQueryParams(state, "")

	getAccountAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	principalId := state.Primary.Attributes["principal_id"]
	for {
		getAccountAssignmentResp, err := getAccountAssignmentClient.Request("GET", getAccountAssignmentPath,
			&getAccountAssignmentOpt)

		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center account assignment: %s", err)
		}

		getAccountAssignmentRespBody, err := utils.FlattenResponse(getAccountAssignmentResp)
		if err != nil {
			return nil, err
		}

		accountPermissions := utils.PathSearch("account_assignments", getAccountAssignmentRespBody, nil)

		if accountPermissions == nil {
			return nil, fmt.Errorf("error get Identity Center account assignment")
		}

		for _, v := range accountPermissions.([]interface{}) {
			if principalId == utils.PathSearch("principal_id", v, "") {
				return v, nil
			}
		}
		marker := utils.PathSearch("page_info.next_marker", getAccountAssignmentRespBody, nil)
		if marker == nil {
			break
		}

		getAccountAssignmentPath = getAccountAssignmentBasePath + buildGetAccountAssignmentQueryParams(state,
			marker.(string))
	}
	return nil, fmt.Errorf("error get Identity Center account permission")
}

func buildGetAccountAssignmentQueryParams(state *terraform.ResourceState, marker string) string {
	res := "?limit=100"
	if v, ok := state.Primary.Attributes["target_id"]; ok {
		res = fmt.Sprintf("%s&account_id=%v", res, v)
	}

	if v, ok := state.Primary.Attributes["permission_set_id"]; ok {
		res = fmt.Sprintf("%s&permission_set_id=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func TestAccIdentityCenterAccountAssignment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_account_assignment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccountAssignment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "principal_id",
						"huaweicloud_identitycenter_user.test", "id"),
					resource.TestCheckResourceAttr(rName, "principal_type", "USER"),
					resource.TestCheckResourceAttr(rName, "target_id", acceptance.HW_IDENTITY_CENTER_ACCOUNT_ID),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterAccountAssignmentImportState(rName),
			},
		},
	})
}

func testAccountAssignment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_permission_set" "test" {
  instance_id      = data.huaweicloud_identitycenter_instance.test.id
  name             = "%[2]s"
  session_duration = "PT8H"
}

resource "huaweicloud_identitycenter_account_assignment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.test.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  principal_id      = huaweicloud_identitycenter_user.test.id
  principal_type    = "USER"
  target_id         = "%[3]s"
  target_type       = "ACCOUNT"
}
`, testIdentityCenterUser_basic(name), name, acceptance.HW_IDENTITY_CENTER_ACCOUNT_ID)
}

// testIdentityCenterUserImportState use to return an id with format <identity_store_id>/<user_id>
func testIdentityCenterAccountAssignmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}
		permissionSetID := rs.Primary.Attributes["permission_set_id"]
		if permissionSetID == "" {
			return "", fmt.Errorf("attribute (permission_set_id) of Resource (%s) not found: %s", name, rs)
		}
		targetID := rs.Primary.Attributes["target_id"]
		if targetID == "" {
			return "", fmt.Errorf("attribute (target_id) of Resource (%s) not found: %s", name, rs)
		}
		principalID := rs.Primary.Attributes["principal_id"]
		if principalID == "" {
			return "", fmt.Errorf("attribute (principal_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s/%s/%s", instanceID, permissionSetID, targetID, principalID), nil
	}
}
