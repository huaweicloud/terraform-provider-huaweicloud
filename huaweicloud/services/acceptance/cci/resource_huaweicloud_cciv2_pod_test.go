package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getV2PodResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetV2Pod(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Pod_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_pod.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getV2PodResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Pod_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "containers.0.image", "nginx:stable-alpine-perl"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.name", "c1"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.command.0", "nginx"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.command.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.command.2", "/etc/nginx/nginx.conf"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.args.0", "-g"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.args.1", "daemon off;"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.startup_probe.0.exec.0.command.0", "sh"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.startup_probe.0.exec.0.command.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.startup_probe.0.exec.0.command.2", "exit 0"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.memory", "2G"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.memory", "2G"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.0.name", "init-first"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.0.command.0", "sh"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.0.command.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.0.command.2", "exit 0"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.1.name", "init-second"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.1.command.0", "sh"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.1.command.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "init_containers.1.command.2", "exit 0"),
					resource.TestCheckResourceAttr(resourceName, "image_pull_secrets.0.name", "imagepull-secret"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "containers", "resource_version", "status"},
				ImportStateIdFunc:       testAccV2PodImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2PodImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["namespace"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("the namespace (%s) or name(%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccV2Pod_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "huaweicloud_vpcep_service_summary" "swr" {
  endpoint_service_name = "com.myhuaweicloud.cn-north-4.swr"
}

data "huaweicloud_vpcep_service_summary" "obs" {
  endpoint_service_name = "cn-north-4.com.myhuaweicloud.v4.obsv2.OBSCluster9"
}

resource "huaweicloud_vpcep_endpoint" "swr" {
  service_id = data.huaweicloud_vpcep_service_summary.swr.id
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  enable_dns = true
}

resource "huaweicloud_vpcep_endpoint" "obs" {
  service_id = data.huaweicloud_vpcep_service_summary.obs.id
  vpc_id     = huaweicloud_vpc.test.id
  enable_dns = true
}

resource "huaweicloud_cciv2_namespace" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cciv2_network" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[1]s"

  annotations = {
    "yangtse.io/project-id"                 = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/project-id"],
    "yangtse.io/domain-id"                  = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/domain-id"],
    "yangtse.io/warm-pool-size"             = "10",
    "yangtse.io/warm-pool-recycle-interval" = "2",
  }

  subnets {
    subnet_id = huaweicloud_vpc_subnet.test.subnet_id
  }

  security_group_ids = [huaweicloud_networking_secgroup.test.id]
}
`, rName)
}

func testAccV2Pod_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_pod" "test" {
  depends_on = [huaweicloud_cciv2_network.test]

  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  annotations = {
    "description"                    = "test",
    "resource.cci.io/pod-size-specs" = "2.00_2.0",
    "resource.cci.io/instance-type"  = "general-computing",
  }

  containers {
    image   = "nginx:stable-alpine-perl"
    name    = "c1"
    command = ["nginx", "-c", "/etc/nginx/nginx.conf"]
    args    = ["-g", "daemon off;"]

    startup_probe {
      exec {
        command = ["sh", "-c", "exit 0"]
      }

      termination_grace_period_seconds = 30
    }

    resources {
      limits = {
        cpu    = 2
        memory = "2G"
      }

      requests = {
        cpu    = 2
        memory = "2G"
      }
    }
  }

  init_containers {
    image   = "nginx:stable-alpine-perl"
    name    = "init-first"
    command = ["sh", "-c", "exit 0"]

    resources {
      limits = {
        cpu    = 2
        memory = "2G"
      }

      requests = {
        cpu    = 2
        memory = "2G"
      }
    }
  }

  init_containers {
    image   = "nginx:stable-alpine-perl"
    name    = "init-second"
    command = ["sh", "-c", "exit 0"]

    resources {
      limits = {
        cpu    = 2
        memory = "2G"
      }

      requests = {
        cpu    = 2
        memory = "2G"
      }
    }
  }

  image_pull_secrets {
    name = "imagepull-secret"
  }

  lifecycle {
    ignore_changes = [
      annotations, containers,
    ]
  }
}
`, testAccV2Pod_base(rName), rName)
}
