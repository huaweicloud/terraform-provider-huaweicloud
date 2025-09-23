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

func getRecordCallbackFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/record/callbacks/{id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live record callback configuration: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func TestAccRecordCallback_basic(t *testing.T) {
	var (
		callbackObj interface{}
		rName       = "huaweicloud_live_record_callback.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&callbackObj,
		getRecordCallbackFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveIngestDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCallBack_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "types.0", "RECORD_NEW_FILE_START"),
					resource.TestCheckResourceAttr(rName, "url", "https://mycallback.com.cn/record_notify"),
					resource.TestCheckResourceAttr(rName, "sign_type", "MD5"),
					resource.TestCheckResourceAttr(rName, "key", "d4a3345fd2f5b76fdb91f6ae4fe924cb"),
				),
			},
			{
				Config: testCallBack_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "types.0", "RECORD_NEW_FILE_START"),
					resource.TestCheckResourceAttr(rName, "types.1", "RECORD_FILE_COMPLETE"),
					resource.TestCheckResourceAttr(rName, "url", "https://mycallback.com.cn/record_notify2"),
					resource.TestCheckResourceAttr(rName, "sign_type", "HMACSHA256"),
					resource.TestCheckResourceAttr(rName, "key", "74b4464a8a00e72d5d0d66387543b6b7"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
		},
	})
}

func testCallBack_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_record_callback" "test" {
  domain_name = "%[1]s"
  types       = ["RECORD_NEW_FILE_START"]
  url         = "https://mycallback.com.cn/record_notify"
  sign_type   = "MD5"
  key         = "d4a3345fd2f5b76fdb91f6ae4fe924cb"
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testCallBack_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_record_callback" "test" {
  domain_name = "%[1]s"
  types       = ["RECORD_NEW_FILE_START", "RECORD_FILE_COMPLETE"]
  url         = "https://mycallback.com.cn/record_notify2"
  sign_type   = "HMACSHA256"
  key         = "74b4464a8a00e72d5d0d66387543b6b7"
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}
