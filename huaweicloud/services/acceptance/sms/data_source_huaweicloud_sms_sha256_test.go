package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsSha256_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_sha256.test"
	key, _ := uuid.GenerateUUID()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsSha256_basic(key),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "value"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsSha256_basic(key string) string {
	return fmt.Sprintf(`
data "huaweicloud_sms_sha256" "test" {
  key = "%s"
}
`, key)
}
