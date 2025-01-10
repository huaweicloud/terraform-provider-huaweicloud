package rms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/rms/v1/policyassignments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getOrganizationalPolicyAssignmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	// getOrganizationalPolicyAssignment: Query the RMS organizational policy assignment
	var (
		getOrgPolicyAssignmentHttpUrl = "v1/resource-manager/organizations/{organization_id}/policy-assignments/{id}"
		getOrgPolicyAssignmentProduct = "rms"
	)
	getOrgPolicyAssignmentClient, err := cfg.NewServiceClient(getOrgPolicyAssignmentProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Config client: %s", err)
	}

	getOrgPolicyAssignmentPath := getOrgPolicyAssignmentClient.Endpoint + getOrgPolicyAssignmentHttpUrl
	getOrgPolicyAssignmentPath = strings.ReplaceAll(getOrgPolicyAssignmentPath, "{organization_id}",
		state.Primary.Attributes["organization_id"])
	getOrgPolicyAssignmentPath = strings.ReplaceAll(getOrgPolicyAssignmentPath, "{id}", state.Primary.ID)

	getOrgPolicyAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrgPolicyAssignmentResp, err := getOrgPolicyAssignmentClient.Request("GET", getOrgPolicyAssignmentPath,
		&getOrgPolicyAssignmentOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS organizational policy assignment: %s", err)
	}

	return utils.FlattenResponse(getOrgPolicyAssignmentResp)
}

// Test the builtin policy assignment.
func TestAccOrganizationalPolicyAssignment_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rms_organizational_policy_assignment.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationalPolicyAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalPolicyAssignment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization_id",
						"data.huaweicloud_organizations_organization.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "The maximum number of days without rotation. Default 90."),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "policy_definition_id",
						"data.huaweicloud_rms_policy_definitions.test", "definitions.0.id"),
					resource.TestCheckResourceAttr(rName, "period", "TwentyFour_Hours"),
					resource.TestCheckResourceAttr(rName, "parameters.maxAccessKeyAge", "\"90\""),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccOrganizationalPolicyAssignment_basicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization_id",
						"data.huaweicloud_organizations_organization.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "The maximum number of days without rotation. Default 60."),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "policy_definition_id",
						"data.huaweicloud_rms_policy_definitions.test", "definitions.0.id"),
					resource.TestCheckResourceAttr(rName, "period", "Twelve_Hours"),
					resource.TestCheckResourceAttr(rName, "parameters.maxAccessKeyAge", "\"60\""),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOrganizationalPolicyAssignmentImportState(rName),
			},
		},
	})
}

func testAccOrganizationalPolicyAssignment_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "access-keys-rotated"
}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id      = data.huaweicloud_organizations_organization.test.id
  name                 = "%[1]s"
  description          = "The maximum number of days without rotation. Default 90."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  period               = "TwentyFour_Hours"

  parameters = {
    maxAccessKeyAge = "\"90\""
  }
}
`, name)
}

func testAccOrganizationalPolicyAssignment_basicUpdate(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "access-keys-rotated"
}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id      = data.huaweicloud_organizations_organization.test.id
  name                 = "%[1]s"
  description          = "The maximum number of days without rotation. Default 60."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  period               = "Twelve_Hours"

  parameters = {
    maxAccessKeyAge = "\"60\""
  }
}
`, name)
}

// Test the custom policy assignment.
func TestAccOrganizationalPolicyAssignment_custom(t *testing.T) {
	var (
		obj policyassignments.Assignment

		rName        = "huaweicloud_rms_organizational_policy_assignment.test"
		name         = acceptance.RandomAccResourceNameWithDash()
		customConfig = testAccOrganizationalPolicyAssignment_customConfig(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationalPolicyAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckRMSExcludedAccounts(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalPolicyAssignment_custom(customConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization_id",
						"data.huaweicloud_organizations_organization.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"huaweicloud_fgs_function.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "This is a custom policy assignment."),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "parameters.string_test", "\"string_value\""),
					resource.TestCheckResourceAttr(rName, "parameters.array_test", "[\"array_element\"]"),
					resource.TestCheckResourceAttr(rName, "parameters.object_test", "{\"terraform_version\":\"1.xx.x\"}"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccOrganizationalPolicyAssignment_customUpdate(customConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization_id",
						"data.huaweicloud_organizations_organization.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"huaweicloud_fgs_function.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a custom policy assignment update."),
					resource.TestCheckResourceAttr(rName, "parameters.string_test", "\"update_string_value\""),
					resource.TestCheckResourceAttr(rName, "parameters.update_array_test", "[\"array_element\"]"),
					resource.TestCheckResourceAttr(rName, "parameters.object_test", "{\"update_terraform_version\":\"1.xx.xx\"}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOrganizationalPolicyAssignmentImportState(rName),
			},
		},
	})
}

func testAccOrganizationalPolicyAssignment_customConfig(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name                  = "%[1]s"
  code_type             = "inline"
  handler               = "index.handler"
  runtime               = "Node.js10.16"
  functiongraph_version = "v2"
  app                   = "default"
  enterprise_project_id = "0"
  memory_size           = 128
  timeout               = 3
}
`, name)
}

func testAccOrganizationalPolicyAssignment_custom(customConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = "%[2]s"
  description     = "This is a custom policy assignment."
  function_urn    = "${huaweicloud_fgs_function.test.urn}:${huaweicloud_fgs_function.test.version}"
  period          = "TwentyFour_Hours"

  excluded_accounts = [
    "%[3]s",
    "%[4]s",
  ]

  parameters = {
    string_test = "\"string_value\""
    array_test  = "[\"array_element\"]"
    object_test = jsonencode({"terraform_version": "1.xx.x"})
  }
}
`, customConfig, name, acceptance.HW_RMS_EXCLUDED_ACCOUNT_1, acceptance.HW_RMS_EXCLUDED_ACCOUNT_2)
}

func testAccOrganizationalPolicyAssignment_customUpdate(customConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = "%[2]s"
  description     = "This is a custom policy assignment update."
  function_urn    = "${huaweicloud_fgs_function.test.urn}:${huaweicloud_fgs_function.test.version}"
  period          = "Twelve_Hours"

  excluded_accounts = [
    "%[3]s",
    "%[4]s",
  ]

  parameters = {
    string_test       = "\"update_string_value\""
    update_array_test = "[\"array_element\"]"
    object_test       = jsonencode({"update_terraform_version": "1.xx.xx"})
  }
}
`, customConfig, name, acceptance.HW_RMS_EXCLUDED_ACCOUNT_1, acceptance.HW_RMS_EXCLUDED_ACCOUNT_2)
}

func testOrganizationalPolicyAssignmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		dataStandard, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, dataStandard)
		}

		var organizationID string
		if organizationID = dataStandard.Primary.Attributes["organization_id"]; organizationID == "" {
			return "", fmt.Errorf("attribute (organization_id) of Resource (%s) not found: %s", name, dataStandard)
		}
		return fmt.Sprintf("%s/%s", organizationID, dataStandard.Primary.ID), nil
	}
}
