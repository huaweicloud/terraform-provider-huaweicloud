package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTrafficMirrorFilterRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return "", fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorFilterRuleHttpUrl := "vpc/traffic-mirror-filter-rules/" + state.Primary.ID
	getTrafficMirrorFilterRulePath := client.ResourceBaseURL() + getTrafficMirrorFilterRuleHttpUrl
	getTrafficMirrorFilterRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorFilterRuleResp, err := client.Request("GET", getTrafficMirrorFilterRulePath, &getTrafficMirrorFilterRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving traffic mirror filter rule: %s", err)
	}

	return utils.FlattenResponse(getTrafficMirrorFilterRuleResp)
}

func TestAccTrafficMirrorFilterRule_basic(t *testing.T) {
	var (
		trafficMirrorFilterRule interface{}
		name                    = acceptance.RandomAccResourceNameWithDash()
		resourceName            = "huaweicloud_vpc_traffic_mirror_filter_rule.test"
		rc                      = acceptance.InitResourceCheck(
			resourceName,
			&trafficMirrorFilterRule,
			getTrafficMirrorFilterRuleResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficMirrorFilterRule_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_filter_id",
						"huaweicloud_vpc_traffic_mirror_filter.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					// the IP address is shown as all in console, but it is empty from API response
					resource.TestCheckResourceAttr(resourceName, "source_port_range", ""),
					resource.TestCheckResourceAttr(resourceName, "source_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "source_port_range", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "action", "accept"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "create VPC traffic mirror filter rule"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccTrafficMirrorFilterRule_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "source_cidr_block", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "source_port_range", "80-81"),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "destination_port_range", "1-65535"),
					resource.TestCheckResourceAttr(resourceName, "action", "reject"),
					resource.TestCheckResourceAttr(resourceName, "priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				Config: testAccTrafficMirrorFilterRule_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", "icmpv6"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv6"),
					resource.TestCheckResourceAttr(resourceName, "source_cidr_block", "2002:50::44/128"),
					resource.TestCheckResourceAttr(resourceName, "source_port_range", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "::/0"),
					resource.TestCheckResourceAttr(resourceName, "destination_port_range", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTrafficMirrorFilterRule_egress(t *testing.T) {
	var (
		trafficMirrorFilterRule interface{}
		name                    = acceptance.RandomAccResourceNameWithDash()
		resourceName            = "huaweicloud_vpc_traffic_mirror_filter_rule.test"
		rc                      = acceptance.InitResourceCheck(
			resourceName,
			&trafficMirrorFilterRule,
			getTrafficMirrorFilterRuleResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficMirrorFilterRule_egress(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_filter_id",
						"huaweicloud_vpc_traffic_mirror_filter.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "direction", "egress"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "all"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "source_cidr_block", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "action", "accept"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttr(resourceName, "source_port_range", ""),
					resource.TestCheckResourceAttr(resourceName, "source_port_range", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTrafficMirrorFilterRule_step1(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv4"
  direction                = "ingress"
  protocol                 = "tcp"
  action                   = "accept"
  priority                 = 1
  description              = "create VPC traffic mirror filter rule"
}
`, testAccTrafficMirrorFilter_base(name, ""))
}

func testAccTrafficMirrorFilterRule_step2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv4"
  direction                = "ingress"
  protocol                 = "udp"
  action                   = "reject"
  priority                 = 2
  source_cidr_block        = "192.168.0.0/24"
  source_port_range        = "80-81"
  destination_cidr_block   = "0.0.0.0/0"
  destination_port_range   = "1-65535"
  description              = ""
}
`, testAccTrafficMirrorFilter_base(name, ""))
}

func testAccTrafficMirrorFilterRule_step3(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv6"
  direction                = "ingress"
  protocol                 = "icmpv6"
  action                   = "reject"
  priority                 = 2
  source_cidr_block        = "2002:50::44/128"
  destination_cidr_block   = "::/0"
  description              = ""
}
`, testAccTrafficMirrorFilter_base(name, ""))
}

func testAccTrafficMirrorFilterRule_egress(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv4"
  direction                = "egress"
  protocol                 = "all"
  action                   = "accept"
  priority                 = 20
  source_cidr_block        = "192.168.1.0/24"
}
`, testAccTrafficMirrorFilter_base(name, ""))
}
