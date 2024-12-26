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

func getDomainResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcLiveV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Live v1 client: %s", err)
	}

	return client.ShowDomain(&model.ShowDomainRequest{Domain: &state.Primary.ID})
}

func TestAccDomain_basic(t *testing.T) {
	var obj model.CreateDomainMappingResponse

	ingestDomainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	streamingDomainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	ingestResourceName := "huaweicloud_live_domain.ingestDomain"
	streamingResourceName := "huaweicloud_live_domain.streamingDomain"

	rc := acceptance.InitResourceCheck(
		ingestResourceName,
		&obj,
		getDomainResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDomain_basic(ingestDomainName, streamingDomainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(ingestResourceName, "name", ingestDomainName),
					resource.TestCheckResourceAttr(ingestResourceName, "type", "push"),
					resource.TestCheckResourceAttr(ingestResourceName, "status", "on"),
					resource.TestCheckResourceAttr(ingestResourceName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttr(ingestResourceName, "is_ipv6", "true"),
					resource.TestCheckResourceAttr(ingestResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),

					resource.TestCheckResourceAttr(streamingResourceName, "name", streamingDomainName),
					resource.TestCheckResourceAttr(streamingResourceName, "type", "pull"),
					resource.TestCheckResourceAttr(streamingResourceName, "status", "on"),
					resource.TestCheckResourceAttr(streamingResourceName, "ingest_domain_name", ingestDomainName),
					resource.TestCheckResourceAttr(streamingResourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(streamingResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testDomain_basic_update(ingestDomainName, streamingDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ingestResourceName, "name", ingestDomainName),
					resource.TestCheckResourceAttr(ingestResourceName, "type", "push"),
					resource.TestCheckResourceAttr(ingestResourceName, "status", "on"),
					resource.TestCheckResourceAttr(ingestResourceName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttr(ingestResourceName, "is_ipv6", "false"),
					resource.TestCheckResourceAttr(ingestResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),

					resource.TestCheckResourceAttr(streamingResourceName, "name", streamingDomainName),
					resource.TestCheckResourceAttr(streamingResourceName, "type", "pull"),
					resource.TestCheckResourceAttr(streamingResourceName, "status", "off"),
					resource.TestCheckResourceAttr(streamingResourceName, "ingest_domain_name", ""),
					resource.TestCheckResourceAttr(streamingResourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(streamingResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      streamingResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDomain_basic(pushDomain, pullDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name                  = "%[1]s"
  type                  = "push"
  service_area          = "mainland_china"
  is_ipv6               = true
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_live_domain" "streamingDomain" {
  name                  = "%[2]s"
  type                  = "pull"
  ingest_domain_name    = huaweicloud_live_domain.ingestDomain.name
  service_area          = "outside_mainland_china"
  enterprise_project_id = "%[3]s"
}
`, pushDomain, pullDomain, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDomain_basic_update(pushDomain, pullDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name                  = "%[1]s"
  type                  = "push"
  service_area          = "mainland_china"
  is_ipv6               = false
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_live_domain" "streamingDomain" {
  name                  = "%[2]s"
  type                  = "pull"
  status                = "off"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "%[3]s"
}
`, pushDomain, pullDomain, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
