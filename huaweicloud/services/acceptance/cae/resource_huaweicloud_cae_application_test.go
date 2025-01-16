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

	envId := state.Primary.Attributes["environment_id"]
	return cae.GetApplicationById(client, envId, state.Primary.ID)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cae_application.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getApplicationFunc)

		invalidName = "-tf-test-invalid-name"
		name        = acceptance.RandomAccResourceNameWithDash()
		baseConfig  = testAccApplication_base(name)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplication_basic(baseConfig, invalidName),
				ExpectError: regexp.MustCompile(`CAE.01500214`),
			},
			{
				Config: testAccApplication_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportStateFunc(resourceName),
			},
		},
	})
}

func testAccApplicationImportStateFunc(rsName string) resource.ImportStateIdFunc {
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

		return fmt.Sprintf("%s/%s", environmentId, applicationId), nil
	}
}

func testAccApplication_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_swr_organization" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cae_environment" "test" {
  name = "%[1]s"

  annotations = {
    type              = "exclusive"
    vpc_id            = huaweicloud_vpc.test.id
    subnet_id         = huaweicloud_vpc_subnet.test.id
    security_group_id = huaweicloud_networking_secgroup.test.id
    group_name        = huaweicloud_swr_organization.test.name
  }

  // To avoid k8s container deploy failed.
  max_retries = 1
}
`, name)
}

func testAccApplication_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_application" "test" {
  environment_id = huaweicloud_cae_environment.test.id
  name           = "%[2]s"
}
`, baseConfig, name)
}
