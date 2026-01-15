package cce

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

func getReleaseFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getReleaseHttpUrl = "cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		getReleaseProduct = "cce"
	)
	getReleaseClient, err := cfg.NewServiceClient(getReleaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", err)
	}

	getReleaseHttpPath := getReleaseClient.Endpoint + getReleaseHttpUrl
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{namespace}", state.Primary.Attributes["namespace"])
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{name}", state.Primary.ID)

	getReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getReleaseResp, err := getReleaseClient.Request("GET", getReleaseHttpPath, &getReleaseOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE release: %s", err)
	}

	return utils.FlattenResponse(getReleaseResp)
}

func TestAccRelease_basic(t *testing.T) {
	var (
		release      interface{}
		resourceName = "huaweicloud_cce_release.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&release,
			getReleaseFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRelease_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "chart_id", "huaweicloud_cce_chart.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "version", "4.9.0"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_name"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "status_description"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEReleaseImportStateIdFunc("default", name),
				ImportStateVerifyIgnore: []string{
					"version", "values_json", "chart_id", "description", "parameters",
				},
			},
		},
	})
}

func testAccCCEReleaseImportStateIdFunc(namespace, name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["huaweicloud_cce_cluster.test"]
		if !ok {
			return "", fmt.Errorf("cluster not found: %s", cluster)
		}
		if cluster.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", cluster.Primary.ID, namespace, name)
		}
		return fmt.Sprintf("%s/%s/%s", cluster.Primary.ID, namespace, name), nil
	}
}

func testAccRelease_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cce_release" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  chart_id   = huaweicloud_cce_chart.test.id
  name       = "%[3]s"
  namespace  = "default"
  version    = "4.9.0"

  values_json = jsonencode({
    "key1" : ["value1"]
    "key2" : "value2"
    "key3" : "value3"
    "key4" : {
      "key1" : "value1",
      "key2" : "value2",
      "key3" : {
         "sub_key1" : "sub_value1",
         "sub_key2" : "sub_value2"
      }
    }
  })

  description = "created by terraform"

  parameters {
    dry_run = false
  }
}
`, testAccCluster_basic(name), testAccChart_basic(), name)
}
