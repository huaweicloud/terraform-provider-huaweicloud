package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCacheRefreshResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.ShowHistoryTaskDetailsRequest{
		EnterpriseProjectId: utils.StringIgnoreEmpty(state.Primary.Attributes["enterprise_project_id"]),
		HistoryTasksId:      state.Primary.ID,
	}

	resp, err := hcCdnClient.ShowHistoryTaskDetails(request)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CDN cache refresh: %s", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("error retrieving CDN cache refresh: Task is not found in API response")
	}
	return resp, nil
}

func TestAccCacheRefresh_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_cdn_cache_refresh.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCacheRefreshResourceFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNURL(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCacheRefresh_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "urls.0", acceptance.HW_CDN_DOMAIN_URL),
					resource.TestCheckResourceAttr(rName, "type", "directory"),
					resource.TestCheckResourceAttr(rName, "mode", "detect_modify_refresh"),
					resource.TestCheckResourceAttr(rName, "zh_url_encode", "true"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "processing"),
					resource.TestCheckResourceAttrSet(rName, "succeed"),
					resource.TestCheckResourceAttrSet(rName, "failed"),
					resource.TestCheckResourceAttrSet(rName, "total"),
				),
			},
		},
	})
}

func testCacheRefresh_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_cache_refresh" "test" {
  urls          = ["%s"]
  type          = "directory"
  mode          = "detect_modify_refresh"
  zh_url_encode = true
}
`, acceptance.HW_CDN_DOMAIN_URL)
}
