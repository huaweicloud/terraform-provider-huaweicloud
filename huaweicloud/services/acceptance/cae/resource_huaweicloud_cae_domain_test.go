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

	return cae.GetDomainById(
		client,
		state.Primary.Attributes["environment_id"],
		state.Primary.ID,
		state.Primary.Attributes["enterprise_project_id"],
	)
}

func TestAccResourceDomain_basic(t *testing.T) {
	var (
		domainName        = fmt.Sprintf("%s.com", acceptance.RandomAccResourceNameWithDash())
		domainNameWithEps = fmt.Sprintf("%s.com", acceptance.RandomAccResourceNameWithDash())

		domainObj interface{}
		rName     = "huaweicloud_cae_domain.test.0"
		rc        = acceptance.InitResourceCheck(rName, &domainObj, getResourceDomainFunc)

		nameWithSpecifiedEps = "huaweicloud_cae_domain.test.1"
		rcWithSpecifiedEps   = acceptance.InitResourceCheck(nameWithSpecifiedEps, &domainObj, getResourceDomainFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithSpecifiedEps.CheckResourceDestroy(),
		),

		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomain_basic(domainName, domainNameWithEps),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "environment_id"),
					resource.TestCheckResourceAttr(rName, "name", domainName),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckNoResourceAttr(rName, "enterprise_project_id"),
					rcWithSpecifiedEps.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(nameWithSpecifiedEps, "environment_id"),
					resource.TestCheckResourceAttr(nameWithSpecifiedEps, "name", domainNameWithEps),
					resource.TestCheckResourceAttrPair(nameWithSpecifiedEps, "enterprise_project_id",
						"data.huaweicloud_cae_environments.test", "environments.0.annotations.enterprise_project_id"),
					resource.TestMatchResourceAttr(nameWithSpecifiedEps, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config:      testAccResourceDomain_existed(domainName, domainNameWithEps),
				ExpectError: regexp.MustCompile(`the domain has already existed`),
			},
			{
				ResourceName:      "huaweicloud_cae_domain.test[0]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDomainImportFunc(rName, true),
			},
			{
				ResourceName:      "huaweicloud_cae_domain.test[1]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDomainImportFunc(nameWithSpecifiedEps, false),
			},
		},
	})
}

func testAccResourceDomain_basic(domainName, domainNameWithEps string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[1]
}

# Create two domains, the first one belonging to the default enterprise project and the second one belonging to the
# specified enterprise project.
resource "huaweicloud_cae_domain" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  name                  = count.index == 0 ? "%[2]s" : "%[3]s"
  enterprise_project_id = count.index == 1 ? try(data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id,
  null) : null
}
`, acceptance.HW_CAE_ENVIRONMENT_IDS, domainName, domainNameWithEps)
}

func testAccResourceDomain_existed(domainName, domainNameWithEps string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_domain" "test2" {
  environment_id = local.env_ids[0]
  name           = "%[2]s"
}
`, testAccResourceDomain_basic(domainName, domainNameWithEps), domainName)
}

func testAccDomainImportFunc(name string, isDefaultEps bool) resource.ImportStateIdFunc {
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
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>' or "+
				"'<environment_id>/<name>/<enterprise_project_id>', but got '%s/%s'",
				environmentId, domainName)
		}

		if isDefaultEps {
			return fmt.Sprintf("%s/%s", environmentId, domainName), nil
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if epsId == "" {
			return "", fmt.Errorf("enterprise_project_id is missing, want '<environment_id>/<name>/<enterprise_project_id>', "+
				"but got '%s/%s'", environmentId, domainName)
		}

		return fmt.Sprintf("%s/%s/%s", environmentId, domainName, epsId), nil
	}
}
