package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEnvironmentVariableFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return environments.GetVariable(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccEnvironmentVariable_basic(t *testing.T) {
	var (
		variable environments.Variable
		rName    = "huaweicloud_apig_environment_variable.test"
		name     = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&variable,
		getEnvironmentVariableFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "value", "/stage/demo"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "env_id", "huaweicloud_apig_environment.test", "id"),
				),
			},
			{
				Config: testAccEnvironmentVariable_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "value", "/stage/terraform"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEnvironmentVariableImportStateFunc(),
			},
		},
	})
}

func testAccEnvironmentVariableImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_environment_variable.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		groupId := rs.Primary.Attributes["group_id"]
		variableName := rs.Primary.Attributes["name"]
		if instanceId == "" || groupId == "" || variableName == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<group_id>/<name>', but '%s/%s/%s'",
				instanceId, groupId, variableName)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, groupId, variableName), nil
	}
}

func testAccEnvironmentVariable_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_environment" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
  description = "Created by script"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccEnvironmentVariable_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_environment_variable" "test" {
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  env_id      = huaweicloud_apig_environment.test.id
  name        = "%[2]s"
  value       = "/stage/demo"
}
`, testAccEnvironmentVariable_base(name), name)
}

func testAccEnvironmentVariable_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_environment_variable" "test" {
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  env_id      = huaweicloud_apig_environment.test.id
  name        = "%[2]s"
  value       = "/stage/terraform"
}
`, testAccEnvironmentVariable_base(name), name)
}
