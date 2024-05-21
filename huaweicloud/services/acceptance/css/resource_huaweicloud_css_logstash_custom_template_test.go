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

func getLogstashCustomTemplateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	cssV1Client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getCustomTemplateHttpUrl := "v1.0/{project_id}/lgsconf/template"
	getCustomTemplatePath := cssV1Client.Endpoint + getCustomTemplateHttpUrl
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{project_id}", cssV1Client.ProjectID)

	getCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCustomTemplateResp, err := cssV1Client.Request("GET", getCustomTemplatePath, &getCustomTemplateOpt)
	if err != nil {
		return nil, err
	}

	getCustomTemplateRespBody, err := utils.FlattenResponse(getCustomTemplateResp)
	if err != nil {
		return nil, err
	}

	getCustomTemplateExp := fmt.Sprintf("customTemplates[?name=='%s']|[0]", state.Primary.ID)
	customTemplate := utils.PathSearch(getCustomTemplateExp, getCustomTemplateRespBody, nil)
	if customTemplate == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return customTemplate, nil
}

func TestAccLogstashCustomTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_logstash_custom_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogstashCustomTemplateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogstashCustomTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					resource.TestCheckResourceAttrSet(rName, "conf_content"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "configuration_name"},
			},
		},
	})
}

func testLogstashCustomTemplate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_custom_template" "test" {
  cluster_id         = huaweicloud_css_logstash_cluster.test.id
  name               = "%[2]s_custom_template"
  configuration_name = huaweicloud_css_logstash_configuration.test.name
  description        = "custom template test"
}
`, logstashFavoriteConfiguration(name), name)
}

func logstashFavoriteConfiguration(rName string) string {
	confContent := `input { \r\n    redis {\r\n        data_type => \"pattern_channel\"\r\n        key` +
		` => \"lgs-*\"\r\n        host => \"xxx.xxx.xxx.xxxx\"\r\n        port => 6379\r\n    }\r\n}\r\n` +
		`\r\nfilter {\r\n    mutate {\r\n        remove_field => [\"@timestamp\", \"@version\"] \r\n    }` +
		` \r\n} \r\n\r\noutput { \r\n    elasticsearch { \r\n        hosts => [\"http://xxx.xxx.xxx.xxx:9200\",` +
		` \"http://xxx.xxx.xxx.xxx:9200\", \"http://xxx.xxx.xxx.xxx:9200\"]\r\n` +
		`        index => \"xxxxxx\"\r\n    } \r\n}`
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_configuration" "test" {
  cluster_id   = huaweicloud_css_logstash_cluster.test.id
  name         = "%[2]s_favorite"
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
`, testAccLogstashCluster_basic(rName, 1, "bar"), rName, confContent)
}
