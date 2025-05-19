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
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.memory", "2G"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.memory", "2G"),
					resource.TestCheckResourceAttr(resourceName, "image_pull_secrets.0.name", "imagepull-secret"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "containers"},
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
    image = "nginx:stable-alpine-perl"
    name  = "c1"

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
`, testAccV2Network_basic(rName), rName)
}
