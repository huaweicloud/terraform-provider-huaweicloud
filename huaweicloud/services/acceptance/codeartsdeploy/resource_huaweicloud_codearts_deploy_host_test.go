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

func getDeployHostResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/resources/host-groups/{group_id}/hosts/{host_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating codearts client: %s", err)
	}

	getHostPath := client.Endpoint + httpUrl
	getHostPath = strings.ReplaceAll(getHostPath, "{group_id}", state.Primary.Attributes["group_id"])
	getHostPath = strings.ReplaceAll(getHostPath, "{host_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getHostPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy host: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDeployHost_withProxyMode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	proxyName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_host.test"
	proxyRName := "huaweicloud_codearts_deploy_host.test_proxy"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployHostResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployHost_withProxyMode(proxyName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id",
						"huaweicloud_codearts_deploy_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_address",
						"huaweicloud_compute_instance.test.1", "public_ip"),
					resource.TestCheckResourceAttrPair(rName, "proxy_host_id",
						"huaweicloud_codearts_deploy_host.test_proxy", "id"),
					resource.TestCheckResourceAttr(rName, "port", "22"),
					resource.TestCheckResourceAttr(rName, "username", "root"),
					resource.TestCheckResourceAttr(rName, "password", "Test@123"),
					resource.TestCheckResourceAttr(rName, "os_type", "linux"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "install_icagent", "true"),
					resource.TestCheckResourceAttr(rName, "sync", "true"),
					resource.TestCheckResourceAttr(rName, "as_proxy", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "lastest_connection_at"),
					resource.TestCheckResourceAttrSet(rName, "connection_status"),
					resource.TestCheckResourceAttrSet(rName, "permission.#"),

					resource.TestCheckResourceAttrPair(proxyRName, "ip_address",
						"huaweicloud_compute_instance.test.0", "public_ip"),
					resource.TestCheckResourceAttr(proxyRName, "as_proxy", "true"),
				),
			},
			{
				Config: testDeployHost_withProxyMode_update(proxyName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "ip_address",
						"huaweicloud_compute_instance.test.1", "access_ip_v4"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "install_icagent", "false"),
					resource.TestCheckResourceAttr(rName, "sync", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployHostImportState(rName),
				ImportStateVerifyIgnore: []string{
					"password",
					"private_key",
					"install_icagent",
					"sync",
					"connection_status",
					"lastest_connection_at",
				},
			},
		},
	})
}

func TestAccDeployHost_withoutProxyMode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_host.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployHostResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployHost_withoutProxyMode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id",
						"huaweicloud_codearts_deploy_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_address",
						"huaweicloud_compute_instance.test.1", "public_ip"),
					resource.TestCheckResourceAttr(rName, "port", "22"),
					resource.TestCheckResourceAttr(rName, "username", "root"),
					resource.TestCheckResourceAttr(rName, "password", "Test@123"),
					resource.TestCheckResourceAttr(rName, "os_type", "windows"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "install_icagent", "true"),
					resource.TestCheckResourceAttr(rName, "sync", "true"),
					resource.TestCheckResourceAttr(rName, "as_proxy", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "lastest_connection_at"),
					resource.TestCheckResourceAttrSet(rName, "connection_status"),
					resource.TestCheckResourceAttrSet(rName, "permission.#"),
				),
			},
			{
				Config: testDeployHost_withoutProxyMode_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "install_icagent", "false"),
					resource.TestCheckResourceAttr(rName, "sync", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployHostImportState(rName),
				ImportStateVerifyIgnore: []string{
					"password",
					"private_key",
					"install_icagent",
					"sync",
					"connection_status",
					"lastest_connection_at",
				},
			},
		},
	})
}

func testComputedInstance(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                        = "%[2]s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  admin_pass                  = "Test@123"
  delete_disks_on_termination = true
  
  system_disk_type = "SAS"
  system_disk_size = 40
  
  eip_type = "5_bgp"
  bandwidth {
    share_type  = "PER"
    size        = 1
    charge_mode = ""
  }

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      image_id, tags, name
    ]
  }
}
`, common.TestBaseComputeResources(name), name)
}

func testDeployHost_base_withProxyMode(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s"
  os_type     = "linux"
  description = "test description"
}
`, testComputedInstance(name), testProject_basic(name), name)
}

func testDeployHost_base_withoutProxyMode(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id    = huaweicloud_codearts_project.test.id
  name          = "%[3]s"
  os_type       = "windows"
  description   = "test description"
  is_proxy_mode = 0
}
`, testComputedInstance(name), testProject_basic(name), name)
}

func testDeployHost_withProxyMode(proxyName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_host" "test_proxy" {
  group_id   = huaweicloud_codearts_deploy_group.test.id
  ip_address = huaweicloud_compute_instance.test[0].public_ip
  port       = 22
  username   = "root"
  password   = "Test@123"
  os_type    = "linux"
  name       = "%[2]s"
  as_proxy   = true
}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[1].public_ip
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "linux"
  proxy_host_id   = huaweicloud_codearts_deploy_host.test_proxy.id
  name            = "%[3]s"
  install_icagent = true
  sync            = true
}
`, testDeployHost_base_withProxyMode(name), proxyName, name)
}

func testDeployHost_withProxyMode_update(proxyName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_host" "test_proxy" {
  group_id   = huaweicloud_codearts_deploy_group.test.id
  ip_address = huaweicloud_compute_instance.test[0].public_ip
  port       = 22
  username   = "root"
  password   = "Test@123"
  os_type    = "linux"
  name       = "%[2]s"
  as_proxy   = true
}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[1].access_ip_v4
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "linux"
  proxy_host_id   = huaweicloud_codearts_deploy_host.test_proxy.id
  name            = "%[3]s_update"
  install_icagent = false
  sync            = false
}
`, testDeployHost_base_withProxyMode(name), proxyName, name)
}

func testDeployHost_withoutProxyMode(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[1].public_ip
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "windows"
  name            = "%[2]s"
  install_icagent = true
  sync            = true
}
`, testDeployHost_base_withoutProxyMode(name), name)
}

func testDeployHost_withoutProxyMode_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id        = huaweicloud_codearts_deploy_group.test.id
  ip_address      = huaweicloud_compute_instance.test[1].public_ip
  port            = 22
  username        = "root"
  password        = "Test@123"
  os_type         = "windows"
  name            = "%[2]s_update"
  install_icagent = false
  sync            = false
}
`, testDeployHost_base_withoutProxyMode(name), name)
}

// testDeployHostImportState use to return an ID with format <group_id>/<id>
func testDeployHostImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return "", fmt.Errorf("attribute (group_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", groupID, rs.Primary.ID), nil
	}
}
