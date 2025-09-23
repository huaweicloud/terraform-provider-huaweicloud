package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getAssociatedDomainInfoFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.GetGroupAssociatedDomainByUrl(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["group_id"], state.Primary.Attributes["url_domain"])
}

func TestAccGroupDomainAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceNameWithDash()

		resourceName = "huaweicloud_apig_group_domain_associate.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAssociatedDomainInfoFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDomainAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "url_domain", "huaweicloud_dns_zone.test", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "min_ssl_version"),
					resource.TestCheckResourceAttr(resourceName, "ingress_http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "ingress_https_port", "443"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			{
				Config: testAccGroupDomainAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "url_domain", "huaweicloud_dns_zone.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "min_ssl_version", "TLSv1.2"),
					resource.TestCheckResourceAttr(resourceName, "ingress_http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "ingress_https_port", "-1"),
					resource.TestCheckResourceAttr(resourceName, "is_http_redirect_to_https", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			{
				Config: testAccGroupDomainAssociate_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "url_domain", "huaweicloud_dns_zone.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttr(resourceName, "ingress_http_port", "-1"),
					resource.TestCheckResourceAttr(resourceName, "ingress_https_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "is_http_redirect_to_https", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			{
				Config: testAccGroupDomainAssociate_basic_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "url_domain", "huaweicloud_dns_zone.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttr(resourceName, "ingress_http_port", "1024"),
					resource.TestCheckResourceAttr(resourceName, "ingress_https_port", "1025"),
					resource.TestCheckResourceAttr(resourceName, "is_http_redirect_to_https", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"is_http_redirect_to_https",
				},
			},
		},
	})
}

func testAccGroupDomainAssociate_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name        = "script-test-%[1]s.com."
  description = "a zone"
  ttl         = 300
  status      = "DISABLE"

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}

resource "huaweicloud_apig_group" "test" {
  instance_id = "%[2]s"
  name        = "%[1]s"

  force_destroy = true

  lifecycle {
    ignore_changes = [
      url_domains,
    ]
  }
}
`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccGroupDomainAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id = huaweicloud_apig_group.test.instance_id
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = huaweicloud_dns_zone.test.name

  depends_on = [
    huaweicloud_dns_zone.test,
  ]
}
`, testAccGroupDomainAssociate_base(name))
}

func testAccGroupDomainAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id = huaweicloud_apig_group.test.instance_id
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = huaweicloud_dns_zone.test.name

  min_ssl_version           = "TLSv1.2"
  ingress_http_port         = 80
  ingress_https_port        = -1
  is_http_redirect_to_https = true

  depends_on = [
    huaweicloud_dns_zone.test,
  ]
}
`, testAccGroupDomainAssociate_base(name))
}

func testAccGroupDomainAssociate_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id = huaweicloud_apig_group.test.instance_id
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = huaweicloud_dns_zone.test.name

  min_ssl_version           = "TLSv1.1"
  ingress_http_port         = -1
  ingress_https_port        = 443
  is_http_redirect_to_https = false

  depends_on = [
    huaweicloud_dns_zone.test,
  ]
}
`, testAccGroupDomainAssociate_base(name))
}

func testAccGroupDomainAssociate_basic_step4(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id = huaweicloud_apig_group.test.instance_id
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = huaweicloud_dns_zone.test.name

  ingress_http_port         = 1024 # Make sure this custom inbound access port is opened.
  ingress_https_port        = 1025 # Make sure this custom inbound access port is opened.
  is_http_redirect_to_https = false

  depends_on = [
    huaweicloud_dns_zone.test,
  ]
}
`, testAccGroupDomainAssociate_base(name))
}
