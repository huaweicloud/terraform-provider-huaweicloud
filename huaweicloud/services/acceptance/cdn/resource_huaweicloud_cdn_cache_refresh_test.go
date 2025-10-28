package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getCacheRefreshResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	getRespBody, err := cdn.GetCacheDetailById(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccCacheRefresh_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cdn_cache_refresh.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCacheRefreshResourceFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNURL(t)
			// The value of the enterprise project ID must be consistent with the enterprise project to which the
			// domain belongs.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCacheRefresh_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "urls.0", acceptance.HW_CDN_DOMAIN_URL),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func testAccCacheRefresh_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_cache_refresh" "test" {
  urls                  = ["%[1]s"]
  enterprise_project_id = "%[2]s"
  type                  = "directory"
  mode                  = "detect_modify_refresh"
  zh_url_encode         = true
}
`, acceptance.HW_CDN_DOMAIN_URL, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
