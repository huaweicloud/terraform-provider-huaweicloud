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

func getV2HPAResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetV2HPA(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2HPA_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_hpa.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getV2HPAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2HPA_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "api_version"),
					resource.TestCheckResourceAttrSet(resourceName, "kind"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttr(resourceName, "min_replicas", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_replicas", "5"),
					resource.TestCheckResourceAttr(resourceName, "scale_target_ref.0.kind", "Deployment"),
					resource.TestCheckResourceAttr(resourceName, "scale_target_ref.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "scale_target_ref.0.api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "metrics.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2HPAImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2HPAImportStateFunc(name string) resource.ImportStateIdFunc {
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

func testAccV2HPA_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_hpa" "test" {
  name      = "%[2]s"
  namespace = huaweicloud_cciv2_namespace.test.name

  min_replicas = 1
  max_replicas = 5

  scale_target_ref {
    kind        = huaweicloud_cciv2_deployment.test.kind
    name        = huaweicloud_cciv2_deployment.test.name
    api_version = huaweicloud_cciv2_deployment.test.api_version
  }

  metrics {
    type = "Resource"

    resources {
      name = "memory"

      target {
        type                = "Utilization"
        average_utilization = 50
      }
    }
  }

  metrics {
    type = "Resource"

    resources {
      name = "cpu"

      target {
        type                = "Utilization"
        average_utilization = 50
      }
    }
  }
}
`, testAccV2Deployment_basic(rName), rName)
}
