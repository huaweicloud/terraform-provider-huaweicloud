package huaweicloud
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSmsTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMSTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSmsTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "name", "template_test_1"),
				),
			},
		},
	})
}

var testSmsTemplateConfig = `
resource "huaweicloud_sms_template" "test" {
	access = {
		AK =  "OYFWTNGHT3HG0V0MTXBZ"
        SK = "TYkhq0NqrCmOjAckdVC4fhnCW0AgGdJFAZ2q1oCQ"
	}

	name				= "template_test_1"
	is_template			= "true"
	region				= "ap-southeast-1"
	projectid			= "0e1b256e9680f3722fdec005b172cc1f"
	availability_zone	= "ap-southeast-1a"
}
`

func testAccCheckSMSTemplateDestroy(s *terraform.State) error {
	return nil
}