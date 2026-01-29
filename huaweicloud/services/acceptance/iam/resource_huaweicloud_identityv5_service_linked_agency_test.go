package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5ServiceLinkedAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5ServiceLinkedAgencyById(client, state.Primary.ID)
}

func TestAccV5ServiceLinkedService_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_identityv5_service_linked_agency.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5ServiceLinkedAgencyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckServiceLinkedAgencyPrincipal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource has no delete logic.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV5ServiceLinkedService_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "service_principal", acceptance.HW_IAM_SERVICE_LINKED_AGENCY_PRINCIPAL),
					resource.TestCheckResourceAttr(rName, "description", "Create by terraform script"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "trust_policy"),
					resource.TestCheckResourceAttrSet(rName, "agency_id"),
					resource.TestCheckResourceAttrSet(rName, "agency_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttr(rName, "path",
						"service-linked-agency/"+acceptance.HW_IAM_SERVICE_LINKED_AGENCY_PRINCIPAL+"/"),
					resource.TestCheckResourceAttr(rName, "max_session_duration", "3600"),
				),
			},
		},
	})
}

func testAccV5ServiceLinkedService_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_service_linked_agency" "test" {
  service_principal = "%s"
  description       = "Create by terraform script"
}
`, acceptance.HW_IAM_SERVICE_LINKED_AGENCY_PRINCIPAL)
}
