package iam

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

func getServiceAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getAgencyHttpUrl := "v5/agencies/{agency_id}"
	getAgencyPath := client.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{agency_id}", state.Primary.ID)
	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAgencyResp, err := client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM service agency: %s", err)
	}
	return utils.FlattenResponse(getAgencyResp)
}

func TestAccServiceAgency_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_service_agency.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getServiceAgencyResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAgency_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "service.APIG"),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "trust_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccServiceAgency_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "duration", "7200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.new_key", "new_value"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceAgency_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_service_agency" "test" {
  name                   = "%[1]s"
  delegated_service_name = "service.APIG"
  policy_names           = ["NATReadOnlyPolicy"]
  description            = "Created by terraform script"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccServiceAgency_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_service_agency" "test" {
  name                   = "%[1]s"
  delegated_service_name = "service.APIG"
  policy_names           = ["NATReadOnlyPolicy", "RDSReadOnlyPolicy"]
  duration               = 7200
  description            = "Updated by terraform script"

  tags = {
    foo     = "baar"
    new_key = "new_value"
  }
}
`, name)
}
