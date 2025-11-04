package swrenterprise

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

func getResourceSwrEnterpriseInstanceRegistry(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<instance_id>/<registry_id>', but got '%s'", state.Primary.ID)
	}
	instanceId := parts[0]
	id := parts[1]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/registries/{registry_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{registry_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR registry: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccSwrEnterpriseInstanceRegistry_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_instance_registry.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseInstanceRegistry,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseInstanceRegistry_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "huawei-SWR"),
					resource.TestCheckResourceAttr(resourceName, "insecure", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc"),
				),
			},
			{
				Config: testAccSwrEnterpriseInstanceRegistry_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "type", "huawei-SWR"),
					resource.TestCheckResourceAttr(resourceName, "insecure", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc-update"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential.0.access_secret"},
			},
		},
	})
}

func testAccSwrEnterpriseInstanceRegistry_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

resource "huaweicloud_swr_temporary_login_command" "test" {}

locals {
  x_swr_docker_login = split(" ", huaweicloud_swr_temporary_login_command.test.x_swr_docker_login)
}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = "%[1]s"
  type        = "huawei-SWR"
  insecure    = true
  description = "desc"
  url         = "https://${huaweicloud_swr_temporary_login_command.test.auths.0.key}"
  
  credential {
    access_key    = local.x_swr_docker_login[3]
    access_secret = local.x_swr_docker_login[5]
	type          = "basic"
  }
}
`, rName)
}

func testAccSwrEnterpriseInstanceRegistry_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

resource "huaweicloud_swr_temporary_login_command" "test" {}

locals {
  x_swr_docker_login = split(" ", huaweicloud_swr_temporary_login_command.test.x_swr_docker_login)
}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = "%[1]s-update"
  type        = "huawei-SWR"
  insecure    = true
  description = "desc-update"
  url         = "https://${huaweicloud_swr_temporary_login_command.test.auths.0.key}"
  
  credential {
    access_key    = local.x_swr_docker_login[3]
    access_secret = local.x_swr_docker_login[5]
	type          = "basic"
  }
}
`, rName)
}

func TestAccSwrEnterpriseInstanceRegistry_enterprise(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_instance_registry.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseInstanceRegistry,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseInstanceRegistry_enterprise(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "swr-pro-internal"),
					resource.TestCheckResourceAttr(resourceName, "insecure", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc"),
				),
			},
			{
				Config: testAccSwrEnterpriseInstanceRegistry_enterprise_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "type", "swr-pro-internal"),
					resource.TestCheckResourceAttr(resourceName, "insecure", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc-update"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential.0.access_secret"},
			},
		},
	})
}

func testAccSwrEnterpriseInstanceRegistry_enterprise(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_instances" "beijing" {
  region = "cn-north-4"
}

data "huaweicloud_identity_projects" "beijing" {
  name = "cn-north-4"
}

resource "huaweicloud_swr_enterprise_temporary_credential" "test" {
  region      = "cn-north-4"
  instance_id = data.huaweicloud_swr_enterprise_instances.beijing.instances[0].id
}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = "%[1]s"
  type        = "swr-pro-internal"
  insecure    = true
  description = "desc"

  url                = "https://${data.huaweicloud_swr_enterprise_instances.beijing.instances[0].access_address}"
  region_id          = "cn-north-4"
  project_id         = data.huaweicloud_identity_projects.beijing.projects[0].id
  target_instance_id = data.huaweicloud_swr_enterprise_instances.beijing.instances[0].id
  
  credential {
    access_key    = huaweicloud_swr_enterprise_temporary_credential.test.user_id
    access_secret = huaweicloud_swr_enterprise_temporary_credential.test.auth_token
	type          = "basic"
  }
}
`, rName)
}

func testAccSwrEnterpriseInstanceRegistry_enterprise_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_instances" "beijing" {
  region = "cn-north-4"
}

data "huaweicloud_identity_projects" "beijing" {
  name = "cn-north-4"
}

resource "huaweicloud_swr_enterprise_temporary_credential" "test" {
  region      = "cn-north-4"
  instance_id = data.huaweicloud_swr_enterprise_instances.beijing.instances[0].id
}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = "%[1]s-update"
  type        = "swr-pro-internal"
  insecure    = true
  description = "desc-update"

  url                = "https://${data.huaweicloud_swr_enterprise_instances.beijing.instances[0].access_address}"
  region_id          = "cn-north-4"
  project_id         = data.huaweicloud_identity_projects.beijing.projects[0].id
  target_instance_id = data.huaweicloud_swr_enterprise_instances.beijing.instances[0].id
  
  credential {
    access_key    = huaweicloud_swr_enterprise_temporary_credential.test.user_id
    access_secret = huaweicloud_swr_enterprise_temporary_credential.test.auth_token
	type          = "basic"
  }
}
`, rName)
}
