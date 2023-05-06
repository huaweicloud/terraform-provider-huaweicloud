package live

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

func getLiveStreamDisableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getLiveStreamDisable: Query Live stream disable
	var (
		getLiveStreamDisableHttpUrl = "v1/{project_id}/stream/blocks"
		getLiveStreamDisableProduct = "live"
	)
	getLiveStreamDisableClient, err := cfg.NewServiceClient(getLiveStreamDisableProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <domain_name>/<app_name>/<stream_name>")
	}
	domainName := parts[0]
	appName := parts[1]
	streamName := parts[2]

	getLiveStreamDisablePath := getLiveStreamDisableClient.Endpoint + getLiveStreamDisableHttpUrl
	getLiveStreamDisablePath = strings.ReplaceAll(getLiveStreamDisablePath, "{project_id}",
		getLiveStreamDisableClient.ProjectID)

	getLiveStreamDisableQueryParams := buildGetLiveStreamDisableQueryParams(domainName, appName, streamName)
	getLiveStreamDisablePath += getLiveStreamDisableQueryParams

	getLiveStreamDisableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getLiveStreamDisableResp, err := getLiveStreamDisableClient.Request("GET", getLiveStreamDisablePath,
		&getLiveStreamDisableOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live stream disable: %s", err)
	}

	getLiveStreamDisableRespBody, err := utils.FlattenResponse(getLiveStreamDisableResp)
	if err != nil {
		return nil, err
	}

	blocks := utils.PathSearch("blocks", getLiveStreamDisableRespBody, make([]interface{}, 0)).([]interface{})
	if len(blocks) == 0 {
		return nil, fmt.Errorf("err get live stream disable")
	}

	return blocks[0], nil
}

func buildGetLiveStreamDisableQueryParams(domainName, appName, streamName string) string {
	return fmt.Sprintf("?domain=%s&app_name=%s&stream_name=%s", domainName, appName, streamName)
}

func TestAccLiveStreamDisable_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	pushDomainName := fmt.Sprintf("%s.huaweicloud.com", name)
	rName := "huaweicloud_live_stream_disable.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveStreamDisableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveStreamDisable_basic(pushDomainName, name, "2023-05-10T00:00:00Z"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name",
						"huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "app_name", "elb"),
					resource.TestCheckResourceAttr(rName, "stream_name", name),
					resource.TestCheckResourceAttr(rName, "resume_time", "2023-05-10T00:00:00Z"),
				),
			},
			{
				Config: testLiveStreamDisable_basic(pushDomainName, name, "2023-05-15T00:00:00Z"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name",
						"huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "app_name", "elb"),
					resource.TestCheckResourceAttr(rName, "stream_name", name),
					resource.TestCheckResourceAttr(rName, "resume_time", "2023-05-15T00:00:00Z"),
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

func testLiveStreamDisable_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "push"
}
`, name)
}

func testLiveStreamDisable_basic(pushDomainName, name, resumeTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_stream_disable" "test" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = "elb"
  stream_name = "%[2]s"
  resume_time = "%[3]s"
}
`, testLiveStreamDisable_base(pushDomainName), name, resumeTime)
}
