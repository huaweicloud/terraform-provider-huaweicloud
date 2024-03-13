package css

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCssConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getConfiguration: Query the CSS configuration.
	var (
		getConfigurationHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ymls/template"
		getConfigurationProduct = "css"
	)
	getConfigurationClient, err := cfg.NewServiceClient(getConfigurationProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS client: %s", err)
	}

	getConfigurationPath := getConfigurationClient.Endpoint + getConfigurationHttpUrl
	getConfigurationPath = strings.ReplaceAll(getConfigurationPath, "{project_id}", getConfigurationClient.ProjectID)
	getConfigurationPath = strings.ReplaceAll(getConfigurationPath, "{cluster_id}", state.Primary.ID)

	getConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getConfigurationResp, err := getConfigurationClient.Request("GET", getConfigurationPath, &getConfigurationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSS configuration: %s", err)
	}

	getConfigurationRespBody, err := utils.FlattenResponse(getConfigurationResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSS configuration: %s", err)
	}

	// If the value of one configuration is different with the default value, the resource is considered to exist.
	statusRaw := utils.PathSearch(`configurations.*.[value == defaultValue][]`, getConfigurationRespBody, []interface{}{})
	log.Printf("[DEBUG] CSS configuration status: %#v", statusRaw)
	if statusArray, ok := statusRaw.([]interface{}); ok {
		for _, v := range statusArray {
			if !v.(bool) {
				return getConfigurationRespBody, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccCssConfiguration_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCssConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCssConfiguration_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "thread_pool_force_merge_size", "3"),
					resource.TestCheckResourceAttr(rName, "http_cors_allow_credetials", "true"),
				),
			},
			{
				Config: testCssConfiguration_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "thread_pool_force_merge_size", "4"),
					resource.TestCheckResourceAttr(rName, "http_cors_allow_credetials", "true"),
					resource.TestCheckResourceAttr(rName, "http_cors_allow_headers", "X-Requested-With, Content-Type"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCssConfiguration_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_configuration" "test" {
  cluster_id                   = huaweicloud_css_cluster.test.id
  thread_pool_force_merge_size = "3"
  http_cors_allow_credetials   = true
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 7, "bar"))
}

func testCssConfiguration_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_configuration" "test" {
  cluster_id                   = huaweicloud_css_cluster.test.id
  thread_pool_force_merge_size = "4"
  http_cors_allow_credetials   = true
  http_cors_allow_headers      = "X-Requested-With, Content-Type"
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 7, "bar"))
}
