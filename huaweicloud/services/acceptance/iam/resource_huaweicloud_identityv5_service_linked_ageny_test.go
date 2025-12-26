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

func getIdentityServiceLinkedAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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
		return nil, fmt.Errorf("error retrieving IAM service linked agency: %s", err)
	}
	return utils.FlattenResponse(getAgencyResp)
}

func TestAccIdentityV5ServiceLinkedService_basic(t *testing.T) {
	var object interface{}
	var servicePrincipal = "service.CBH"
	var description = "test"
	resourceName := "huaweicloud_identityv5_service_linked_agency.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentityServiceLinkedAgencyResourceFunc,
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
				Config: testAccIdentityServiceLinkedServiceV5_basic(servicePrincipal, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "trust_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "agency_id"),
					resource.TestCheckResourceAttrSet(resourceName, "agency_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "path", "service-linked-agency/"+servicePrincipal+"/"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "max_session_duration", "3600"),
				),
			},
		},
	})
}

func testAccIdentityServiceLinkedServiceV5_basic(servicePrincipal, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_service_linked_agency" "test" {
  service_principal = "%s"
  description       = "%s"
}
`, servicePrincipal, description)
}
