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

func TestAccCodeArtsDeployHostsCopy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCodeArtsDeployHostsCopy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					removeTheCopyHost("huaweicloud_codearts_deploy_group.test.1"),
				),
			},
		},
	})
}

func testCodeArtsDeployHostsCopy_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_compute_instance" "test" {
  name                        = "%[3]s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids          = [huaweicloud_networking_secgroup.test.id]
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  admin_pass                  = "Test@123"
  delete_disks_on_termination = true

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  eip_type = "5_bgp"

  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }
}

resource "huaweicloud_codearts_deploy_group" "test" {
  count = 2

  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s-${count.index}"
  os_type     = "linux"
  description = "test"
}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id   = huaweicloud_codearts_deploy_group.test[0].id
  ip_address = huaweicloud_compute_instance.test.public_ip
  port       = 22
  username   = "root"
  password   = "Test@123"
  os_type    = "linux"
  name       = "%[3]s"
  as_proxy   = true
}
`, common.TestBaseComputeResources(name), testProject_basic(name), name)
}

func testCodeArtsDeployHostsCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_hosts_copy" "test" {
  source_group_id = huaweicloud_codearts_deploy_group.test[0].id
  host_uuids      = [huaweicloud_codearts_deploy_host.test.id]
  target_group_id = huaweicloud_codearts_deploy_group.test[1].id
}
`, testCodeArtsDeployHostsCopy_base(name))
}

// Can not delete group when group having a host, delete host before destroy
func removeTheCopyHost(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// get group ID
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceName)
		}
		groupId := rs.Primary.ID
		if groupId == "" {
			return fmt.Errorf("attribute ID of Resource (%s) not found: %s", resourceName, rs)
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.NewServiceClient("codearts_deploy", acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CodeArts Deploy client, err: %s", err)
		}

		// get host ID from host list
		gethttpUrl := "v1/resources/host-groups/{group_id}/hosts"
		getPath := client.Endpoint + gethttpUrl
		getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return fmt.Errorf("error retrieving CodeArts deploy hosts: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return err
		}

		hostId := utils.PathSearch("result[0].uuid", getRespBody, "").(string)
		if hostId == "" {
			return fmt.Errorf("unable to find host ID from API response")
		}

		// delete host
		deletehttpUrl := "v1/resources/host-groups/{group_id}/hosts/{host_id}"
		deletePath := client.Endpoint + deletehttpUrl
		deletePath = strings.ReplaceAll(deletePath, "{group_id}", groupId)
		deletePath = strings.ReplaceAll(deletePath, "{host_id}", hostId)
		deleteOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
		}

		_, err = client.Request("DELETE", deletePath, &deleteOpt)
		if err != nil {
			return fmt.Errorf("error deleting CodeArts deploy host: %s", err)
		}

		return nil
	}
}
