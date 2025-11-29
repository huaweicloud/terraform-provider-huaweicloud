package cdn

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getStatisticSubscriptionTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.GetStatisticSubscriptionTaskById(client, state.Primary.ID)
}

func TestAccStatisticSubscriptionTask_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_cdn_statistic_subscription_task.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getStatisticSubscriptionTaskResourceFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStatisticSubscriptionTask_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "period_type", "0"),
					resource.TestCheckResourceAttrSet(rName, "emails"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
					resource.TestCheckResourceAttrSet(rName, "report_type"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccStatisticSubscriptionTask_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "period_type", "1"),
					resource.TestCheckResourceAttrSet(rName, "emails"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
					resource.TestCheckResourceAttrSet(rName, "report_type"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"emails",
				},
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testStatisticSubscriptionTaskImportStateWithName(rName),
				ImportStateVerifyIgnore: []string{
					"emails",
				},
			},
		},
	})
}

func testStatisticSubscriptionTaskImportStateWithName(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		taskName := rs.Primary.Attributes["name"]
		if taskName == "" {
			return "", errors.New("the subscription task name is missing, want '<name>'")
		}
		return taskName, nil
	}
}

func testAccStatisticSubscriptionTask_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_statistic_subscription_task" "test" {
  name        = "%[1]s"
  period_type = 0
  emails      = "test@example.com"
  domain_name = "%[2]s"
  report_type = "0,1,2"
}
`, name, acceptance.HW_CDN_DOMAIN_NAME)
}

func testAccStatisticSubscriptionTask_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_statistic_subscription_task" "test" {
  name        = "%[1]s"
  period_type = 1
  emails      = "test@example.com,test2@example.com"
  domain_name = "%[2]s"
  report_type = "0,1,2,3,4,5"
}
`, name, acceptance.HW_CDN_DOMAIN_NAME)
}
