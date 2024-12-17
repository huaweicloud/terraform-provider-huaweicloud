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

	pushDomainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	pullDomainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	pushResourceName := "huaweicloud_live_domain.ingestDomain"
	pullResourceName := "huaweicloud_live_domain.streamingDomain"

	rc := acceptance.InitResourceCheck(
		pushResourceName,
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
				Config: testDomain_basic(pushDomainName, pullDomainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pushResourceName, "name", pushDomainName),
					resource.TestCheckResourceAttr(pushResourceName, "type", "push"),
					resource.TestCheckResourceAttr(pushResourceName, "status", "on"),
					resource.TestCheckResourceAttr(pushResourceName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttr(pushResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),

					resource.TestCheckResourceAttr(pullResourceName, "name", pullDomainName),
					resource.TestCheckResourceAttr(pullResourceName, "type", "pull"),
					resource.TestCheckResourceAttr(pullResourceName, "status", "on"),
					resource.TestCheckResourceAttr(pullResourceName, "ingest_domain_name", pushDomainName),
					resource.TestCheckResourceAttr(pullResourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(pullResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testDomain_basic_update(pushDomainName, pullDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pushResourceName, "name", pushDomainName),
					resource.TestCheckResourceAttr(pushResourceName, "type", "push"),
					resource.TestCheckResourceAttr(pushResourceName, "status", "on"),
					resource.TestCheckResourceAttr(pushResourceName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttr(pushResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),

					resource.TestCheckResourceAttr(pullResourceName, "name", pullDomainName),
					resource.TestCheckResourceAttr(pullResourceName, "type", "pull"),
					resource.TestCheckResourceAttr(pullResourceName, "status", "off"),
					resource.TestCheckResourceAttr(pullResourceName, "ingest_domain_name", ""),
					resource.TestCheckResourceAttr(pullResourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(pullResourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      pullResourceName,
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
