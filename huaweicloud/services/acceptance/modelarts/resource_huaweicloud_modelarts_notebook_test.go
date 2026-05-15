package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getNotebookResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetNotebookById(client, state.Primary.ID)
}

func TestAccNotebook_basic(t *testing.T) {
	var (
		obj interface{}

		managed     = "huaweicloud_modelarts_notebook.managed"
		rcManaged   = acceptance.InitResourceCheck(managed, &obj, getNotebookResourceFunc)
		dedicated   = "huaweicloud_modelarts_notebook.dedicated"
		rcDedicated = acceptance.InitResourceCheck(dedicated, &obj, getNotebookResourceFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()

		baseConfig = testAccNotebook_basic_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRunnerPublicIPs(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcManaged.CheckResourceDestroy(),
			rcDedicated.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebook_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					// Managed notebook
					rcManaged.CheckResourceExists(),
					resource.TestCheckResourceAttr(managed, "name", name+"-managed"),
					resource.TestCheckResourceAttrSet(managed, "flavor_id"),
					resource.TestCheckResourceAttrSet(managed, "image_id"),
					resource.TestCheckResourceAttr(managed, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(managed, "allowed_access_ips.#", "0"),
					resource.TestCheckResourceAttr(managed, "volume.0.type", "EVS"),
					resource.TestCheckResourceAttr(managed, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(managed, "volume.0.size", "10"),
					resource.TestCheckResourceAttr(managed, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(managed, "image_type", "BUILD_IN"),
					resource.TestCheckResourceAttrSet(managed, "image_name"),
					resource.TestCheckResourceAttrSet(managed, "image_swr_path"),
					resource.TestCheckResourceAttrSet(managed, "created_at"),
					resource.TestCheckResourceAttrSet(managed, "updated_at"),
					resource.TestCheckResourceAttrSet(managed, "url"),
					resource.TestCheckResourceAttr(managed, "tags.%", "2"),
					resource.TestCheckResourceAttr(managed, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(managed, "tags.key", "value"),
					// Dedicated notebook
					rcDedicated.CheckResourceExists(),
					resource.TestCheckResourceAttr(dedicated, "name", name+"-dedicated"),
					resource.TestCheckResourceAttrSet(dedicated, "flavor_id"),
					resource.TestCheckResourceAttrSet(dedicated, "image_id"),
					resource.TestCheckResourceAttr(dedicated, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrPair(dedicated, "key_pair", "huaweicloud_kps_keypair.test", "name"),
					resource.TestCheckResourceAttrPair(dedicated, "pool_id", "huaweicloud_modelarts_resource_pool.test", "id"),
					resource.TestCheckResourceAttr(managed, "allowed_access_ips.#", "0"),
					resource.TestCheckResourceAttr(dedicated, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(dedicated, "volume.0.ownership", "DEDICATED"),
					resource.TestCheckResourceAttrSet(dedicated, "volume.0.uri"),
					resource.TestCheckResourceAttrPair(dedicated, "volume.0.id", "huaweicloud_sfs_turbo.test", "id"),
					resource.TestCheckResourceAttr(dedicated, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrSet(dedicated, "pool_name"),
					resource.TestCheckResourceAttr(dedicated, "image_type", "BUILD_IN"),
					resource.TestCheckResourceAttrSet(dedicated, "image_name"),
					resource.TestCheckResourceAttrSet(dedicated, "image_swr_path"),
					resource.TestCheckResourceAttrSet(dedicated, "created_at"),
					resource.TestCheckResourceAttrSet(dedicated, "updated_at"),
					resource.TestCheckResourceAttrSet(dedicated, "url"),
					resource.TestCheckResourceAttrSet(dedicated, "ssh_uri"),
					resource.TestCheckResourceAttr(dedicated, "tags.%", "0"),
				),
			},
			{
				Config: testAccNotebook_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					// Managed notebook
					rcManaged.CheckResourceExists(),
					resource.TestCheckResourceAttr(managed, "name", updateName+"-managed"),
					resource.TestCheckResourceAttrSet(managed, "flavor_id"),
					resource.TestCheckResourceAttrSet(managed, "image_id"),
					resource.TestCheckResourceAttr(managed, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(managed, "key_pair", ""),
					resource.TestCheckResourceAttr(managed, "allowed_access_ips.#", "0"),
					resource.TestCheckResourceAttr(managed, "volume.0.type", "EVS"),
					resource.TestCheckResourceAttr(managed, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(managed, "volume.0.size", "20"),
					resource.TestCheckResourceAttr(managed, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(managed, "image_type", "BUILD_IN"),
					resource.TestCheckResourceAttrSet(managed, "image_name"),
					resource.TestCheckResourceAttrSet(managed, "image_swr_path"),
					resource.TestCheckResourceAttrSet(managed, "created_at"),
					resource.TestCheckResourceAttrSet(managed, "updated_at"),
					resource.TestCheckResourceAttrSet(managed, "url"),
					resource.TestCheckResourceAttr(managed, "tags.%", "2"),
					resource.TestCheckResourceAttr(managed, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(managed, "tags.owner", "terraform"),
					// Dedicated notebook
					rcDedicated.CheckResourceExists(),
					resource.TestCheckResourceAttr(dedicated, "name", updateName+"-dedicated"),
					resource.TestCheckResourceAttrSet(dedicated, "flavor_id"),
					resource.TestCheckResourceAttrSet(dedicated, "image_id"),
					resource.TestCheckResourceAttr(dedicated, "description", ""),
					resource.TestCheckResourceAttrPair(dedicated, "key_pair", "huaweicloud_kps_keypair.test", "name"),
					resource.TestCheckResourceAttr(dedicated, "allowed_access_ips.#", "2"),
					resource.TestCheckResourceAttr(dedicated, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(dedicated, "volume.0.ownership", "DEDICATED"),
					resource.TestCheckResourceAttrSet(dedicated, "volume.0.uri"),
					resource.TestCheckResourceAttrPair(dedicated, "volume.0.id", "huaweicloud_sfs_turbo.test", "id"),
					resource.TestCheckResourceAttr(dedicated, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrPair(dedicated, "pool_id", "huaweicloud_modelarts_resource_pool.test", "id"),
					resource.TestCheckResourceAttrSet(dedicated, "pool_name"),
					resource.TestCheckResourceAttr(dedicated, "image_type", "BUILD_IN"),
					resource.TestCheckResourceAttrSet(dedicated, "image_name"),
					resource.TestCheckResourceAttrSet(dedicated, "image_swr_path"),
					resource.TestCheckResourceAttrSet(dedicated, "created_at"),
					resource.TestCheckResourceAttrSet(dedicated, "updated_at"),
					resource.TestCheckResourceAttrSet(dedicated, "url"),
					resource.TestCheckResourceAttrSet(dedicated, "ssh_uri"),
					resource.TestCheckResourceAttr(dedicated, "tags.%", "2"),
					resource.TestCheckResourceAttr(dedicated, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(dedicated, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      managed,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      dedicated,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebook_basic_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

variable "runner_public_ips" {
  type    = string
  default = "%[2]s"
}

locals {
  runner_public_ips = split(",", var.runner_public_ips)
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name                  = "%[3]s"
  cidr                  = "192.168.0.0/16"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[3]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)              # 192.168.0.0/24
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) # 192.168.0.1
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[3]s"
  delete_default_rules = true
}

# Make sure open the full ingress access for 111, 2048, 2049, 2051, 2052 and 20048 ports and about TCP and UDP protocols.
resource "huaweicloud_networking_secgroup_rule" "tcp_ingress_access" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "111,2048,2049,2051,2052,20048"
}

resource "huaweicloud_networking_secgroup_rule" "udp_ingress_access" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  ports             = "111,2048,2049,2051,2052,20048"
}

resource "huaweicloud_sfs_turbo" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.tcp_ingress_access,
    huaweicloud_networking_secgroup_rule.udp_ingress_access,
  ]

  name                  = "%[3]s"
  size                  = 1228
  share_proto           = "NFS"
  share_type            = "HPC"
  hpc_bandwidth         = "40M"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_modelarts_network" "test" {
  name = "%[3]s"
  cidr = "10.168.0.0/16" # The recommended connecting CIDR about SFS Turbo.

  sfs_turbos {
    name = huaweicloud_sfs_turbo.test.name
    id   = huaweicloud_sfs_turbo.test.id
  }
}

data "huaweicloud_modelarts_resource_flavors" "test" {
  type = "Dedicate"
}

locals {
  available_resource_flavors_with_onsale_and_less_memory = [
    for o in data.huaweicloud_modelarts_resource_flavors.test.flavors: o if
      lookup(o.az_status, data.huaweicloud_availability_zones.test.names[0], "soldout") == "normal" &&
      length(regexall("modelarts.vm.cpu", o.id)) > 0 && o.cpu <= 16
  ]
}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%[3]s"
  scope       = ["Notebook", "Train", "Infer"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = try(local.available_resource_flavors_with_onsale_and_less_memory[0].id, null)
    count     = 1
  }

  lifecycle {
    ignore_changes = [
      resources,
    ]
  }
}

data "huaweicloud_modelarts_notebook_flavors" "test" {
  type     = "MANAGED"
  category = "CPU"
}

locals {
  available_notebook_flavors_with_onsale = [
    for o in data.huaweicloud_modelarts_notebook_flavors.test.flavors: o if
      !o.sold_out && !strcontains(o.id, "free")
  ]
}

data "huaweicloud_modelarts_notebook_images" "test" {
  type     = "BUILD_IN"
  cpu_arch = try(data.huaweicloud_modelarts_notebook_flavors.test.flavors[0].arch, "x86_64")
}

locals {
  available_notebook_images_with_onsale = [
    for o in data.huaweicloud_modelarts_notebook_images.test.images: o if
      contains(o.resource_categories, "CPU") && o.status == "ACTIVE" && contains(o.dev_services, "NOTEBOOK")
  ]
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[3]s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_RUNNER_PUBLIC_IPS, name)
}

func testAccNotebook_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_notebook" "managed" {
  name        = "%[2]s-managed"
  flavor_id   = try(local.available_notebook_flavors_with_onsale[0].id, null)
  image_id    = try(data.huaweicloud_modelarts_notebook_images.test.images[0].id, null)
  description = "Created by terraform script"

  volume {
    type      = "EVS"
    ownership = "MANAGED"
    size      = 10
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_modelarts_notebook" "dedicated" {
  name        = "%[2]s-dedicated"
  flavor_id   = try(local.available_notebook_flavors_with_onsale[0].id, null)
  image_id    = try(local.available_notebook_images_with_onsale[0].id, null)
  description = "Created by terraform script"
  key_pair    = huaweicloud_kps_keypair.test.name
  pool_id     = huaweicloud_modelarts_resource_pool.test.id

  volume {
    type      = "EFS"
    ownership = "DEDICATED"
    uri       = format("%%s:/", huaweicloud_modelarts_network.test.sfs_turbos[0].uri)
    id        = huaweicloud_sfs_turbo.test.id
  }
}
`, baseConfig, name)
}

func testAccNotebook_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_notebook" "managed" {
  name        = "%[2]s-managed"
  flavor_id   = try(local.available_notebook_flavors_with_onsale[1].id, null)
  image_id    = try(local.available_notebook_images_with_onsale[1].id, null)
  description = "Updated by terraform script"

  volume {
    type      = "EVS"
    ownership = "MANAGED"
    size      = 20
  }

  tags = {
    foo   = "baar"
    owner = "terraform"
  }
}

resource "huaweicloud_modelarts_notebook" "dedicated" {
  name               = "%[2]s-dedicated"
  # The exclusive upgradeable flavors are not included in the basic list.
  flavor_id          = try(local.available_notebook_flavors_with_onsale[0].id, null)
  image_id           = try(local.available_notebook_images_with_onsale[1].id, null)
  key_pair           = huaweicloud_kps_keypair.test.name
  allowed_access_ips = try(slice(local.runner_public_ips, 0, 2), null)
  pool_id            = huaweicloud_modelarts_resource_pool.test.id

  volume {
    type      = "EFS"
    ownership = "DEDICATED"
    uri       = format("%%s:/", huaweicloud_modelarts_network.test.sfs_turbos[0].uri)
    id        = huaweicloud_sfs_turbo.test.id
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, baseConfig, name)
}
