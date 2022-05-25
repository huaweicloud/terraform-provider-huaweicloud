package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRecordCallbackResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcLiveV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Live v1 client: %s", err)
	}
	return client.ShowRecordCallbackConfig(&model.ShowRecordCallbackConfigRequest{Id: state.Primary.ID})
}

func TestAccRecordCallback_basic(t *testing.T) {
	var obj model.ShowRecordCallbackConfigResponse

	pushDomainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	rName := "huaweicloud_live_record_callback.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRecordCallbackResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCallBack_basic(pushDomainName, "http://mycallback.com.cn/record_notify"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", pushDomainName),
					resource.TestCheckResourceAttr(rName, "types.0", "RECORD_NEW_FILE_START"),
					resource.TestCheckResourceAttr(rName, "url", "http://mycallback.com.cn/record_notify"),
				),
			},
			{
				Config: testCallBack_basic(pushDomainName, "http://mycallback.com.cn/record_notify2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "domain_name", pushDomainName),
					resource.TestCheckResourceAttr(rName, "types.0", "RECORD_NEW_FILE_START"),
					resource.TestCheckResourceAttr(rName, "url", "http://mycallback.com.cn/record_notify2"),
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

func testCallBack_basic(pushDomainName, callbackUrl string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_record_callback" "test" {
  domain_name = huaweicloud_live_domain.ingestDomain.name
  types       = ["RECORD_NEW_FILE_START"]
  url         = "%s"
}
`, pushDomainName, callbackUrl)
}
