package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
)

func getLogstashConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	return css.GetLogstashConfigDetails(client, state.Primary.Attributes["cluster_id"], state.Primary.Attributes["name"])
}

func TestAccLogstashConfiguration_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_logstash_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogstashConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogstashConfiguration_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "conf_content", "test-abc"),
					resource.TestCheckResourceAttr(rName, "setting.0.workers", "4"),
					resource.TestCheckResourceAttr(rName, "setting.0.batch_size", "125"),
					resource.TestCheckResourceAttr(rName, "setting.0.batch_delay_ms", "50"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_type", "memory"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_check_point_writes", "1024"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_max_bytes_mb", "1024"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testlogstashConfiguration_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "conf_content", "test_update"),
					resource.TestCheckResourceAttr(rName, "setting.0.workers", "8"),
					resource.TestCheckResourceAttr(rName, "setting.0.batch_size", "150"),
					resource.TestCheckResourceAttr(rName, "setting.0.batch_delay_ms", "60"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_type", "persisted"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_check_point_writes", "2048"),
					resource.TestCheckResourceAttr(rName, "setting.0.queue_max_bytes_mb", "2048"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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

func testLogstashConfiguration_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_configuration" "test" {
  cluster_id   = huaweicloud_css_logstash_cluster.test.id
  name         = "%[2]s"
  conf_content = "test-abc"
  setting {
    workers                  = 4
    batch_size               = 125
    batch_delay_ms           = 50
    queue_type               = "memory"
    queue_check_point_writes = 1024
    queue_max_bytes_mb       = 1024
  }
}
`, testAccLogstashCluster_basic(name, 1, "bar"), name)
}

func testlogstashConfiguration_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_configuration" "test" {
  cluster_id   = huaweicloud_css_logstash_cluster.test.id
  name         = "%[2]s"
  conf_content = "test_update"
  setting {
    workers                  = 8
    batch_size               = 150
    batch_delay_ms           = 60
    queue_type               = "persisted"
    queue_check_point_writes = 2048
    queue_max_bytes_mb       = 2048
  }
}
`, testAccLogstashCluster_basic(name, 1, "bar"), name)
}
