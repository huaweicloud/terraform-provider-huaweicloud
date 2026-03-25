package css

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

func getLogstashPipelineResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	cssV1Client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getPipelinesHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listpipelines"
	getPipelinesPath := cssV1Client.Endpoint + getPipelinesHttpUrl
	getPipelinesPath = strings.ReplaceAll(getPipelinesPath, "{project_id}", cssV1Client.ProjectID)
	getPipelinesPath = strings.ReplaceAll(getPipelinesPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])

	getPipelineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPipelineResp, err := cssV1Client.Request("GET", getPipelinesPath, &getPipelineOpt)
	if err != nil {
		return nil, fmt.Errorf("error query CSS logstash cluster pipeline: %s", err)
	}
	getPipelineRespBody, err := utils.FlattenResponse(getPipelineResp)
	if err != nil {
		return nil, err
	}

	findStr := "pipelines|[?status!='stopped']"
	pipelines := utils.PathSearch(findStr, getPipelineRespBody, make([]interface{}, 0)).([]interface{})
	if len(pipelines) == 0 {
		return pipelines, golangsdk.ErrDefault404{}
	}

	return pipelines, nil
}

func TestAccLogstashPipeline_basic(t *testing.T) {

	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_css_logstash_pipeline.test"

		configContent = `input { \r\n    redis {\r\n        data_type => \"pattern_channel\"\r\n        key` +
			` => \"lgs-*\"\r\n        host => \"xxx.xxx.xxx.xxxx\"\r\n        port => 6379\r\n    }\r\n}\r\n` +
			`\r\nfilter {\r\n    mutate {\r\n        remove_field => [\"@timestamp\", \"@version\"] \r\n    }` +
			` \r\n} \r\n\r\noutput { \r\n    elasticsearch { \r\n        hosts => [\"http://xxx.xxx.xxx.xxx:9200\",` +
			` \"http://xxx.xxx.xxx.xxx:9200\", \"http://xxx.xxx.xxx.xxx:9200\"]\r\n` +
			`        index => \"xxxxxx\"\r\n    } \r\n}`
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogstashPipelineResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogstashPipeline_basic(name, configContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CSS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "names.#", "1"),
					resource.TestCheckResourceAttr(rName, "pipelines.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.name"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.keep_alive"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.status"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.updated_at"),
				),
			},
			{
				Config: testLogstashPipeline_update_hostStart(name, configContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "names.#", "3"),
					resource.TestCheckResourceAttr(rName, "pipelines.#", "3"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.1.name"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.1.keep_alive"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.1.status"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.1.updated_at"),
				),
			},
			{
				Config: testLogstashPipeline_update_hostStop(name, configContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "names.#", "1"),
					resource.TestCheckResourceAttr(rName, "pipelines.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.name"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.keep_alive"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.status"),
					resource.TestCheckResourceAttrSet(rName, "pipelines.0.updated_at"),
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

func testLogstashPipeline_basic(rName, configContent string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_pipeline" "test" {
  cluster_id = "%[2]s"
  names      = [huaweicloud_css_logstash_configuration.test_1.name]
}
`, logstashCluster_configurations(rName, configContent), acceptance.HW_CSS_CLUSTER_ID)
}

func testLogstashPipeline_update_hostStart(rName, configContent string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_pipeline" "test" {
  cluster_id = "%[2]s"
  names      = [
    huaweicloud_css_logstash_configuration.test_1.name,
    huaweicloud_css_logstash_configuration.test_2.name,
    huaweicloud_css_logstash_configuration.test_3.name
  ]
}
`, logstashCluster_configurations(rName, configContent), acceptance.HW_CSS_CLUSTER_ID)
}

func testLogstashPipeline_update_hostStop(rName, configContent string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_pipeline" "test" {
  cluster_id = "%[2]s"
  names      = [huaweicloud_css_logstash_configuration.test_3.name]
}
`, logstashCluster_configurations(rName, configContent), acceptance.HW_CSS_CLUSTER_ID)
}

func logstashCluster_configurations(rName, configContent string) string {
	return fmt.Sprintf(`
resource "huaweicloud_css_logstash_configuration" "test_1" {
  cluster_id   = "%[1]s"
  name         = "%[2]s_test1"
  conf_content = "%[3]s"

  setting {
    workers                  = 4
    batch_size               = 125
    batch_delay_ms           = 50
    queue_type               = "memory"
    queue_check_point_writes = 1024
    queue_max_bytes_mb       = 1024
  }
}

resource "huaweicloud_css_logstash_configuration" "test_2" {
  cluster_id   = "%[1]s"
  name         = "%[2]s_test2"
  conf_content = "%[3]s"

  setting {
    workers                  = 4
    batch_size               = 125
    batch_delay_ms           = 50
    queue_type               = "memory"
    queue_check_point_writes = 1024
    queue_max_bytes_mb       = 1024
  }
}

resource "huaweicloud_css_logstash_configuration" "test_3" {
  cluster_id   = "%[1]s"
  name         = "%[2]s_test3"
  conf_content = "%[3]s"

  setting {
    workers                  = 4
    batch_size               = 125
    batch_delay_ms           = 50
    queue_type               = "memory"
    queue_check_point_writes = 1024
    queue_max_bytes_mb       = 1024
  }
}
`, acceptance.HW_CSS_CLUSTER_ID, rName, configContent)
}
