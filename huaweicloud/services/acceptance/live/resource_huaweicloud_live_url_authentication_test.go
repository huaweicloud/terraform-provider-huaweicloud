package live

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUrlAuthentication_basic(t *testing.T) {
	var (
		rName1           = "huaweicloud_live_url_authentication.test1"
		rName2           = "huaweicloud_live_url_authentication.test2"
		ingestDomainName = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		streamDomainName = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		startTime        = time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
	)
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// The action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccUrlAuthentication_ingest(ingestDomainName, startTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName1, "key_chain.#", "1"),
				),
			},
			{
				// The action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccUrlAuthentication_stream(streamDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName2, "key_chain.#", "3"),
				),
			},
		},
	})
}

func testAccUrlAuthentication_ingest(name, startTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test1" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_live_url_validation" "test1" {
  domain_name = huaweicloud_live_domain.test1.name
  key         = "IbBIzklRGCyMEd18oPV9sxAuuwNIzT81"
  auth_type   = "c_aes"
  timeout     = 1000
}

resource "huaweicloud_live_url_authentication" "test1" {
  domain_name = huaweicloud_live_domain.test1.name
  type        = "push"
  app_name    = "live1"
  stream_name = "tf-test"
  check_level = 3
  start_time  = "%[2]s"

  depends_on = [huaweicloud_live_url_validation.test1]
}
`, name, startTime)
}

func testAccUrlAuthentication_stream(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test2" {
  name = "%s"
  type = "pull"
}

resource "huaweicloud_live_url_validation" "test2" {
  domain_name = huaweicloud_live_domain.test2.name
  key         = "IbBIzklRGCyMEd18oPV9sxAuuwNIzT81"
  auth_type   = "d_sha256"
  timeout     = 800
}

resource "huaweicloud_live_url_authentication" "test2" {
  domain_name = huaweicloud_live_domain.test2.name
  type        = "pull"
  app_name    = "live2"
  stream_name = "tf-try"

  depends_on = [huaweicloud_live_url_validation.test2]
}
`, name)
}
