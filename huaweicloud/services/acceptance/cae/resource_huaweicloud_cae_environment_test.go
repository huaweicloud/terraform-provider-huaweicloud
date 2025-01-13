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

func getEnvironmentFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetEnvironmentById(client, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, state.Primary.ID)
}

func TestAccEnvironment_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cae_environment.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getEnvironmentFunc)

		invalidName = "-tf-test-invalid-name"
		name        = acceptance.RandomAccResourceNameWithDash()
		baseConfig  = testAccEnvironment_base(name)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironment_basic(baseConfig, invalidName),
				ExpectError: regexp.MustCompile(`CAE.01500001`), // Invalid parameters input.
			},
			{
				Config: testAccEnvironment_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "annotations.type", "exclusive"),
					resource.TestCheckResourceAttrPair(resourceName, "annotations.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "annotations.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "annotations.security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "annotations.group_name",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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
				ImportStateVerifyIgnore: []string{
					"max_retries",
				},
			},
		},
	})
}

func testAccEnvironment_base(name string) string {
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
`, name)
}

func testAccEnvironment_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_environment" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"

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
`, baseConfig, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
