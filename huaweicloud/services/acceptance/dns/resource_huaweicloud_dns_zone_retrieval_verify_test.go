package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDNSRetrievalVerify_basic(t *testing.T) {
	resourceName := "huaweicloud_dns_zone_retrieval_verify.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDnsZoneRetrievalName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSRetrievalVerify_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "message"),
				),
			},
		},
	})
}

func testAccDNSRetrievalVerify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone_retrieval" "test" {
  zone_name = "%[1]s"
}

resource "huaweicloud_dns_zone_retrieval_verify" "test" {
  retrieval_id = huaweicloud_dns_zone_retrieval.test.retrieval_id
}
`, acceptance.HW_DNS_ZONE_RETRIEVAL_NAME)
}
