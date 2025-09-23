package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
)

func getResourceDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	return live.GetDomain(client, state.Primary.ID)
}

func TestAccResourceDomain_basic(t *testing.T) {
	var (
		domainInfo         interface{}
		ingestDomainName   = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		ingestResourceName = "huaweicloud_live_domain.ingestDomain"
	)

	rc := acceptance.InitResourceCheck(
		ingestResourceName,
		&domainInfo,
		getResourceDomainFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIngestDomain_basic(ingestDomainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(ingestResourceName, "name", ingestDomainName),
					resource.TestCheckResourceAttr(ingestResourceName, "type", "push"),
					resource.TestCheckResourceAttr(ingestResourceName, "status", "on"),
					resource.TestCheckResourceAttr(ingestResourceName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttr(ingestResourceName, "is_ipv6", "true"),
					resource.TestCheckResourceAttr(ingestResourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccIngestDomain_update(ingestDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ingestResourceName, "is_ipv6", "false"),
				),
			},
			{
				ResourceName:      ingestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_streamDomain(t *testing.T) {
	var (
		domainInfo            interface{}
		streamingDomainName   = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		streamingResourceName = "huaweicloud_live_domain.streamingDomain"
	)

	rc := acceptance.InitResourceCheck(
		streamingResourceName,
		&domainInfo,
		getResourceDomainFunc,
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
				Config: testAccStreamDomain_basic(streamingDomainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(streamingResourceName, "name", streamingDomainName),
					resource.TestCheckResourceAttr(streamingResourceName, "type", "pull"),
					resource.TestCheckResourceAttr(streamingResourceName, "status", "on"),
					resource.TestCheckResourceAttr(streamingResourceName, "is_ipv6", "true"),
					resource.TestCheckResourceAttr(streamingResourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(streamingResourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(streamingResourceName, "ingest_domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
				),
			},
			{
				Config: testAccStreamDomain_update_1(streamingDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamingResourceName, "is_ipv6", "false"),
					resource.TestCheckResourceAttr(streamingResourceName, "ingest_domain_name", ""),
				),
			},
			{
				Config: testAccStreamDomain_update_2(streamingDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamingResourceName, "status", "off"),
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

func testAccIngestDomain_basic(pushDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name                  = "%s"
  type                  = "push"
  service_area          = "mainland_china"
  enterprise_project_id = "0"
  is_ipv6               = true
}
`, pushDomain)
}

func testAccIngestDomain_update(pushDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name                  = "%s"
  type                  = "push"
  service_area          = "mainland_china"
  enterprise_project_id = "0"
  is_ipv6               = false
}
`, pushDomain)
}

func testAccStreamDomain_basic(streamDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "streamingDomain" {
  name                  = "%[1]s"
  type                  = "pull"
  ingest_domain_name    = "%[2]s"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"
  is_ipv6               = true
}
`, streamDomain, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testAccStreamDomain_update_1(streamDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "streamingDomain" {
  name                  = "%s"
  type                  = "pull"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"
  is_ipv6               = false
}
`, streamDomain)
}

func testAccStreamDomain_update_2(streamDomain string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "streamingDomain" {
  name                  = "%s"
  type                  = "pull"
  status                = "off"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"
  is_ipv6               = false
}
`, streamDomain)
}
