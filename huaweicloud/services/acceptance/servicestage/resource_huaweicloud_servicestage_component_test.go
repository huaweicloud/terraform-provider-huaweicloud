package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/components"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getComponentFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage V2 client: %s", err)
	}
	return components.Get(c, state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccComponent_basic(t *testing.T) {
	var (
		component    components.Component
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_component.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&component,
		getComponentFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponent_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "MicroService"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "framework", "Mesher"),
				),
			},
			// Unable to update the component name because the update method do not work.
			{
				Config: testAccComponent_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "Common"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "framework", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccComponentImportStateIdFunc(),
			},
		},
	})
}

func TestAccComponent_web(t *testing.T) {
	var (
		component    components.Component
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_component.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&component,
		getComponentFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoTokenAuth(t)
			acceptance.TestAccPreCheckComponent(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponent_web(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "Webapp"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Nodejs14"),
					resource.TestCheckResourceAttr(resourceName, "framework", "Web"),
					resource.TestCheckResourceAttr(resourceName, "source.0.type", "GitHub"),
					resource.TestCheckResourceAttrPair(resourceName, "source.0.authorization",
						"huaweicloud_servicestage_repo_token_authorization.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "source.0.url", acceptance.HW_GITHUB_REPO_URL),
					resource.TestCheckResourceAttr(resourceName, "builder.0.organization", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttrPair(resourceName, "builder.0.cluster_id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "builder.0.node_label.owner", "terraform"),
				),
			},
			{
				Config: testAccComponent_webUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "Webapp"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Nodejs14"),
					resource.TestCheckResourceAttr(resourceName, "framework", "Web"),
					resource.TestCheckResourceAttr(resourceName, "source.0.type", "package"),
					resource.TestCheckResourceAttr(resourceName, "source.0.storage_type", "obs"),
					resource.TestCheckResourceAttr(resourceName, "source.0.url", acceptance.HW_OBS_STORAGE_URL),
					resource.TestCheckResourceAttr(resourceName, "builder.0.organization", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttrPair(resourceName, "builder.0.cluster_id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "builder.0.node_label.foo", "bar"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccComponentImportStateIdFunc(),
			},
		},
	})
}

func testAccComponentImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var appId, componentId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_servicestage_component" {
				appId = rs.Primary.Attributes["application_id"]
				componentId = rs.Primary.ID
			}
		}
		if appId == "" || componentId == "" {
			return "", fmt.Errorf("resource not found: %s/%s", appId, componentId)
		}
		return fmt.Sprintf("%s/%s", appId, componentId), nil
	}
}

func testAccComponent_basic_step1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestage_application" "test" {
  name = "%[1]s"
}

resource "huaweicloud_servicestage_component" "test" {
  application_id = huaweicloud_servicestage_application.test.id

  name = "%[1]s"

  type      = "MicroService"
  runtime   = "Docker"
  framework = "Mesher"
}`, rName)
}

func testAccComponent_basic_step2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestage_application" "test" {
  name = "%[1]s"
}

resource "huaweicloud_servicestage_component" "test" {
  application_id = huaweicloud_servicestage_application.test.id

  name = "%[1]s"

  type    = "Common"
  runtime = "Docker"
}`, rName)
}

func testAccComponent_buildConfig(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[1]s"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  flavor_id              = "cce.s2.medium"
  container_network_type = "vpc-router"
  cluster_version        = "v1.19"
  cluster_type           = "VirtualMachine"

  kube_proxy_mode = "iptables"

  dynamic "masters" {
    for_each = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

    content {
      availability_zone = masters.value
    }
  }
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%[1]s"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_kps_keypair.test.name

  root_volume {
    volumetype = "SSD"
    size       = 100
  }

  data_volumes {
    volumetype = "SSD"
    size       = 100
  }

  tags = {
    owner = "terraform"
    foo   = "bar"
  }
}

resource "huaweicloud_servicestage_repo_token_authorization" "test" {
  type  = "github"
  name  = "%[1]s"
  token = "%[2]s"
}
`, rName, acceptance.HW_GITHUB_PERSONAL_TOKEN)
}

func testAccComponent_web(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestage_application" "test" {
  name = "%[2]s"
}

resource "huaweicloud_servicestage_component" "test" {
  depends_on = [huaweicloud_cce_node.test]

  application_id = huaweicloud_servicestage_application.test.id
  type           = "Webapp"
  runtime        = "Nodejs14"
  framework      = "Web"

  name = "%[2]s"

  source {
    type          = "GitHub"
    authorization = huaweicloud_servicestage_repo_token_authorization.test.name
    url           = "%[3]s"
    repo_ref      = "master"
  }

  builder {
    organization = "%[4]s"
    cluster_id   = huaweicloud_cce_cluster.test.id

    node_label = {
      owner = "terraform"
    }
  }
}
`, testAccComponent_buildConfig(rName), rName, acceptance.HW_GITHUB_REPO_URL, acceptance.HW_DOMAIN_NAME)
}

func testAccComponent_webUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestage_application" "test" {
  name = "%[2]s"
}

resource "huaweicloud_servicestage_component" "test" {
  depends_on = [huaweicloud_cce_node.test]

  application_id = huaweicloud_servicestage_application.test.id
  type           = "Webapp"
  runtime        = "Nodejs14"
  framework      = "Web"

  name = "%[2]s"

  source {
    type         = "package"
    storage_type = "obs"
    url          = "%[3]s"
  }

  builder {
    organization = "%[4]s"
    cluster_id   = huaweicloud_cce_cluster.test.id

    node_label = {
      foo = "bar"
    }
  }
}
`, testAccComponent_buildConfig(rName), rName, acceptance.HW_OBS_STORAGE_URL, acceptance.HW_DOMAIN_NAME)
}
