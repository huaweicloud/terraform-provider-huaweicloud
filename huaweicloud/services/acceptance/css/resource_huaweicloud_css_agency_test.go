package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
)

func getResourceAgencyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return css.GetAgency(client, state.Primary.ID)
}

func TestAccResourceAgency_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_css_agency.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceAgencyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAgency_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_id", acceptance.HW_DOMAIN_ID),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"domain_name",
					"type",
				},
			},
		},
	})
}

func testAccAgency_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_css_agency" "test" {
  domain_id   = "%[1]s"
  domain_name = "%[2]s"
  type        = "vpc"
}
`, acceptance.HW_DOMAIN_ID, acceptance.HW_DOMAIN_NAME)
}
