package cae

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
)

func getResourceDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetDomainById(client, state.Primary.Attributes["environment_id"], state.Primary.ID)
}

func TestAccResourceDomain_basic(t *testing.T) {
	var (
		obj interface{}

		domainName = fmt.Sprintf("%s.com", acceptance.RandomAccResourceNameWithDash())

		rName = "huaweicloud_cae_domain.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getResourceDomainFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomain_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "environment_id", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr(rName, "name", domainName),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config:      testAccResourceDomain_existed(domainName),
				ExpectError: regexp.MustCompile(`the domain has already existed`),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDomainImportFunc(rName),
			},
		},
	})
}

func testAccResourceDomain_basic(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cae_domain" "test" {
  environment_id = "%[1]s"
  name           = "%[2]s"
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, domainName)
}

func testAccResourceDomain_existed(domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_domain" "test2" {
  environment_id = "%[2]s"
  name           = "%[3]s"
}
`, testAccResourceDomain_basic(domainName), acceptance.HW_CAE_ENVIRONMENT_ID, domainName)
}

func testAccDomainImportFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var (
			environmentId = rs.Primary.Attributes["environment_id"]
			domainName    = rs.Primary.Attributes["name"]
		)
		if environmentId == "" || domainName == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>', but got '%s/%s'",
				environmentId, domainName)
		}

		return fmt.Sprintf("%s/%s", environmentId, domainName), nil
	}
}
