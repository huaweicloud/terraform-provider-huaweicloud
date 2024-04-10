package rms

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

func getOrgAssignmentPackageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getOrgAssignmentPackage: Query the RMS organizational assignment package
	var (
		getOrgAssignmentPackageHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}"
		getOrgAssignmentPackageProduct = "rms"
	)
	getOrgAssignmentPackageClient, err := cfg.NewServiceClient(getOrgAssignmentPackageProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Config client: %s", err)
	}

	getOrgAssignmentPackagePath := getOrgAssignmentPackageClient.Endpoint + getOrgAssignmentPackageHttpUrl
	getOrgAssignmentPackagePath = strings.ReplaceAll(getOrgAssignmentPackagePath, "{organization_id}",
		state.Primary.Attributes["organization_id"])
	getOrgAssignmentPackagePath = strings.ReplaceAll(getOrgAssignmentPackagePath, "{conformance_pack_id}", state.Primary.ID)

	getOrgAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrgAssignmentPackageResp, err := getOrgAssignmentPackageClient.Request("GET", getOrgAssignmentPackagePath,
		&getOrgAssignmentPackageOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS organizational assignment package: %s", err)
	}

	getOrgAssignmentPackageRespBody, err := utils.FlattenResponse(getOrgAssignmentPackageResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS organizational assignment package: %s", err)
	}

	return getOrgAssignmentPackageRespBody, nil
}

func TestAccOrgAssignmentPackage_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rms_organizational_assignment_package.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrgAssignmentPackageResourceFunc,
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
				Config: testOrgAssignmentPackage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization_id",
						"data.huaweicloud_organizations_organization.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "vars_structure",
						"data.huaweicloud_rms_assignment_package_templates.test", "templates.0.parameters"),
					resource.TestCheckResourceAttrSet(rName, "owner_id"),
					resource.TestCheckResourceAttrSet(rName, "org_conformance_pack_urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testOrgAssignmentPackageImportState(rName),
				ImportStateVerifyIgnore: []string{"template_key", "template_body", "template_uri"},
			},
		},
	})
}

func testOrgAssignmentPackage_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_assignment_package_templates" "test" {}

resource "huaweicloud_rms_organizational_assignment_package" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = "%s"
  template_key    = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters
    content {
      var_key   = vars_structure.value["name"]
      var_value = jsondecode(vars_structure.value["default_value"])
    }
  }
}
`, name)
}

func testOrgAssignmentPackageImportState(name string) resource.ImportStateIdFunc {
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
