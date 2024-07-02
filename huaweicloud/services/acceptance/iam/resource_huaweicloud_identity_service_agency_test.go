package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityServiceAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccIdentityServiceAgency_basic(t *testing.T) {
	var object agency.Agency
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_service_agency.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentityServiceAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckIAMV5(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityServiceAgency_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "service.APIG"),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "trust_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccIdentityServiceAgency_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "duration", "7200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform update"),
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

func testAccIdentityServiceAgency_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_service_agency" "test" {
  name                   = "%s"
  delegated_service_name = "service.APIG"
  policy_names           = ["NATReadOnlyPolicy"]
  description            = "test for terraform"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccIdentityServiceAgency_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_service_agency" "test" {
  name                   = "%s"
  delegated_service_name = "service.APIG"
  policy_names           = ["NATReadOnlyPolicy", "RDSReadOnlyPolicy"]
  duration               = 7200
  description            = "test for terraform update"

  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
}
`, rName)
}
