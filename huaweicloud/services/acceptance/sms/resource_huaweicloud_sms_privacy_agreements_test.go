package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPrivacyAgreementsResourceFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	_, err := conf.SmsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return nil, nil
}

func TestAccPrivacyAgreements_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_sms_privacy_agreements.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivacyAgreementsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivacyAgreements_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "flag", "true"),
				),
			},
		},
	})
}

func testAccPrivacyAgreements_basic() string {
	return `
resource "huaweicloud_sms_privacy_agreements" "test" {}
`
}
