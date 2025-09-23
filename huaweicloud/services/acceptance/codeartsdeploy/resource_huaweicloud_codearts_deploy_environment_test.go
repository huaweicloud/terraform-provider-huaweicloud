package codeartsdeploy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDeployEnvironmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_deploy", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{application_id}/environments/{environment_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{application_id}", state.Primary.Attributes["application_id"])
	getPath = strings.ReplaceAll(getPath, "{environment_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy environment: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	environment := utils.PathSearch("result", getRespBody, nil)
	if environment == nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy environment: result is not found in API response")
	}

	return environment, nil
}

func TestAccDeployEnvironment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_environment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployEnvironmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployEnvironment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "application_id",
						"huaweicloud_codearts_deploy_application.test", "id"),
					resource.TestCheckResourceAttr(rName, "deploy_type", "0"),
					resource.TestCheckResourceAttr(rName, "os_type", "linux"),
					resource.TestCheckResourceAttr(rName, "description", "demo"),
					resource.TestCheckResourceAttr(rName, "hosts.#", "1"),
					resource.TestCheckResourceAttr(rName, "proxies.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "proxies.0.host_id",
						"huaweicloud_codearts_deploy_host.test_proxy", "id"),
					resource.TestCheckResourceAttrSet(rName, "instance_count"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "created_by.#"),
					resource.TestCheckResourceAttrSet(rName, "created_by.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "created_by.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "permission.#"),
					resource.TestCheckResourceAttrSet(rName, "permission.0.can_view"),
					resource.TestCheckResourceAttrSet(rName, "permission.0.can_edit"),
					resource.TestCheckResourceAttrSet(rName, "permission.0.can_delete"),
					resource.TestCheckResourceAttrSet(rName, "permission.0.can_manage"),
					resource.TestCheckResourceAttrSet(rName, "permission.0.can_deploy"),
					resource.TestCheckResourceAttrSet(rName, "permission_matrix.#"),
					resource.TestCheckResourceAttrSet(rName, "permission_matrix.0.role_id"),
					resource.TestCheckResourceAttrSet(rName, "permission_matrix.0.role_name"),
				),
			},
			{
				Config: testDeployEnvironment_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "hosts.#", "1"),
					resource.TestCheckResourceAttr(rName, "proxies.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployEnvironmentImportState(rName),
			},
		},
	})
}

func testDeployEnvironment_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 3

  name                        = "%[2]s-${count.index}"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  admin_pass                  = "Test@123"
  delete_disks_on_termination = true
  
  eip_type = "5_bgp"
  bandwidth {
    share_type  = "PER"
    size        = 1
    charge_mode = ""
  }

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

%[3]s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[2]s"
  os_type     = "linux"
  description = "test description"
}

resource "huaweicloud_codearts_deploy_host" "test_proxy" {
  group_id   = huaweicloud_codearts_deploy_group.test.id
  ip_address = huaweicloud_compute_instance.test[0].public_ip
  port       = 22
  username   = "root"
  password   = "Test@123"
  os_type    = "linux"
  name       = "%[2]s-proxy"
  as_proxy   = true
}

resource "huaweicloud_codearts_deploy_host" "test1" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[1].public_ip
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "linux"
  proxy_host_id   = huaweicloud_codearts_deploy_host.test_proxy.id
  name            = "%[2]s-1"
  install_icagent = true
}

resource "huaweicloud_codearts_deploy_host" "test2" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[2].public_ip
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "linux"
  proxy_host_id   = huaweicloud_codearts_deploy_host.test_proxy.id
  name            = "%[2]s-2"
  install_icagent = true
}
`, common.TestBaseComputeResources(name), name, testDeployApplication_basic(name))
}

func testDeployEnvironment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_environment" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  application_id = huaweicloud_codearts_deploy_application.test.id
  name           = "%[2]s"
  deploy_type    = 0
  os_type        = "linux"
  description    = "demo"

  hosts {
    group_id = huaweicloud_codearts_deploy_group.test.id
    host_id  = huaweicloud_codearts_deploy_host.test1.id
  }
}
`, testDeployEnvironment_base(name), name)
}

func testDeployEnvironment_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_environment" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  application_id = huaweicloud_codearts_deploy_application.test.id
  name           = "%[2]s-update"
  deploy_type    = 0
  os_type        = "linux"

  hosts {
    group_id = huaweicloud_codearts_deploy_group.test.id
    host_id  = huaweicloud_codearts_deploy_host.test2.id
  }
}
`, testDeployEnvironment_base(name), name)
}

// testDeployEnvironmentImportState use to return an ID with format <project_id>/<application_id>/<id>
func testDeployEnvironmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		projectId := rs.Primary.Attributes["project_id"]
		if projectId == "" {
			return "", fmt.Errorf("attribute (project_id) of resource (%s) not found: %s", name, rs)
		}

		applicationId := rs.Primary.Attributes["application_id"]
		if applicationId == "" {
			return "", fmt.Errorf("attribute (application_id) of resource (%s) not found: %s", name, rs)
		}
		return projectId + "/" + applicationId + "/" + rs.Primary.ID, nil
	}
}
