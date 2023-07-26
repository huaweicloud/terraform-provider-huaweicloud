package lts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getHostGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getHostGroup: Query the LTS HostGroup detail
	var (
		getHostGroupHttpUrl = "v3/{project_id}/lts/host-group-list"
		getHostGroupProduct = "lts"
	)
	getHostGroupClient, err := cfg.NewServiceClient(getHostGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS Client: %s", err)
	}

	getHostGroupPath := getHostGroupClient.Endpoint + getHostGroupHttpUrl
	getHostGroupPath = strings.ReplaceAll(getHostGroupPath, "{project_id}", getHostGroupClient.ProjectID)

	getHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getHostGroupOpt.JSONBody = utils.RemoveNil(lts.BuildGetOrDeleteHostGroupBodyParams(state.Primary.ID))
	getHostGroupResp, err := getHostGroupClient.Request("POST", getHostGroupPath, &getHostGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	getHostGroupRespBody, err := utils.FlattenResponse(getHostGroupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?host_group_id=='%s']|[0]", state.Primary.ID)
	getHostGroupRespBody = utils.PathSearch(jsonPath, getHostGroupRespBody, nil)
	if getHostGroupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getHostGroupRespBody, nil
}

func TestAccHostGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_host_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckOutput("is_host_id_different", "false"),
				),
			},
			{
				Config:            testHostGroup_import(name),
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testHostGroup_basic_update_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value"),
					resource.TestCheckOutput("is_host_id_different", "false"),
				),
			},
			{
				Config: testHostGroup_basic_update_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "0"),
				),
			},
		},
	})
}

func testHostGroup_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                = "%s-${count.index}"
  description         = "terraform test"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  # install IC agent
  user_data = <<EOF
#! /bin/bash
set +o history; curl http://icagent-cn-north-4.obs.cn-north-4.myhuaweicloud.com/ICAgent_linux/apm_agent_install.sh > \
apm_agent_install.sh && REGION=cn-north-4 bash apm_agent_install.sh -ak %s \
-sk %s -region cn-north-4 -projectid 0970dd7a1300f5672ff2c003c60ae115 \
-accessip 100.125.12.150 -obsdomain obs.cn-north-4.myhuaweicloud.com \
-accessdomain lts-access.cn-north-4.myhuaweicloud.com
  EOF
}

# wait 2 minutes for the lts service to discover the server
resource "null_resource" "test" {
  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = "sleep 120;"
  }

  depends_on = [huaweicloud_compute_instance.test]
}
`, name, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testHostGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_host_group" "test" {
  name     = "%s"
  type     = "linux"
  host_ids = [
    huaweicloud_compute_instance.test[0].id
  ]

  tags = {
    foo = "bar"
    key = "value"
  }

  depends_on = [null_resource.test]
}

output "is_host_id_different" {
  value = length(setsubtract(huaweicloud_lts_host_group.test.host_ids,
    tolist([huaweicloud_compute_instance.test[0].id]))) != 0
}
`, testHostGroup_base(name), name)
}

func testHostGroup_import(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_host_group" "test" {
  name     = "%s"
  type     = "linux"
  host_ids = [
    huaweicloud_compute_instance.test[0].id
  ]

  tags = {
    foo = "bar"
    key = "value"
  }

  depends_on = [null_resource.test]
}
`, testHostGroup_base(name), name)
}

func testHostGroup_basic_update_1(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_host_group" "test" {
  name     = "%s-update"
  type     = "linux"
  host_ids = [
    huaweicloud_compute_instance.test[0].id,
    huaweicloud_compute_instance.test[1].id
  ]

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }

  depends_on = [null_resource.test]
}

output "is_host_id_different" {
  value = length(setsubtract(huaweicloud_lts_host_group.test.host_ids,
    huaweicloud_compute_instance.test[*].id)) != 0
}
`, testHostGroup_base(name), name)
}

func testHostGroup_basic_update_2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_host_group" "test" {
  name = "%s-update"
  type = "linux"

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }

  depends_on = [null_resource.test]
}
`, testHostGroup_base(name), name)
}
