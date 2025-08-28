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

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetApplicationById(
		client,
		state.Primary.Attributes["environment_id"],
		state.Primary.ID,
		state.Primary.Attributes["enterprise_project_id"],
	)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		invalidName = "-tf-test-invalid-name"
		name        = acceptance.RandomAccResourceNameWithDash()

		obj          interface{}
		resourceName = "huaweicloud_cae_application.test.0"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getApplicationFunc)

		withNotDefaultEpsId   = "huaweicloud_cae_application.test.1"
		rcWithNotDefaultEpsId = acceptance.InitResourceCheck(withNotDefaultEpsId, &obj, getApplicationFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplication_basic(invalidName),
				ExpectError: regexp.MustCompile(`CAE.01500214`),
			},
			{
				Config: testAccApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf("^%s\\d$", name))),
					resource.TestCheckNoResourceAttr(resourceName, "enterprise_project_id"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcWithNotDefaultEpsId.CheckResourceExists(),
					resource.TestMatchResourceAttr(withNotDefaultEpsId, "name", regexp.MustCompile(fmt.Sprintf("^%s\\d$", name))),
					resource.TestCheckResourceAttrPair(withNotDefaultEpsId, "enterprise_project_id",
						"data.huaweicloud_cae_environments.test", "environments.0.annotations.enterprise_project_id"),
				),
			},
			{
				ResourceName:      "huaweicloud_cae_application.test[0]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportStateFunc(resourceName, false),
			},
			{
				ResourceName:      "huaweicloud_cae_application.test[1]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportStateFunc(withNotDefaultEpsId, true),
			},
		},
	})
}

func testAccApplicationImportStateFunc(rsName string, isNotDefaultEpsId bool) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		var (
			environmentId = rs.Primary.Attributes["environment_id"]
			applicationId = rs.Primary.ID
		)
		if environmentId == "" || applicationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<id>', but got '%s/%s'",
				environmentId, applicationId)
		}

		if isNotDefaultEpsId {
			epsId := rs.Primary.Attributes["enterprise_project_id"]
			if epsId == "" {
				return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<id>/<enterprise_project_id>', but got '%s/%s/%s'",
					environmentId, applicationId, epsId)
			}
			return fmt.Sprintf("%s/%s/%s", environmentId, applicationId, epsId), nil
		}

		return fmt.Sprintf("%s/%s", environmentId, applicationId), nil
	}
}

func testAccApplication_basic(name string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

# Query by environment ID under non-default enterprise project ID.
data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[1]
}

resource "huaweicloud_cae_application" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  name                  = "%[2]s${count.index}"
  enterprise_project_id = count.index == 1 ? try(data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id,
  null) : null
}
`, acceptance.HW_CAE_ENVIRONMENT_IDs, name)
}
