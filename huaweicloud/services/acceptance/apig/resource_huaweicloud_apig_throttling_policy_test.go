package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccApigThrottlingPolicyV2_basic(t *testing.T) {
	var (
		// The dedicated instance name only allow letters, digits and underscores (_).
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_throttling_policy.test"
		policy       throttles.ThrottlingPolicy
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigThrottlingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigThrottlingPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "type", "API-based"),
					resource.TestCheckResourceAttr(resourceName, "period", "15000"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(resourceName, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_ip_requests", "60"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourceName, "type", "API-shared"),
					resource.TestCheckResourceAttr(resourceName, "period", "10"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "MINUTE"),
					resource.TestCheckResourceAttr(resourceName, "max_api_requests", "70"),
					resource.TestCheckResourceAttr(resourceName, "max_user_requests", "45"),
					resource.TestCheckResourceAttr(resourceName, "max_app_requests", "45"),
					resource.TestCheckResourceAttr(resourceName, "max_ip_requests", "45"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccApigThrottlingPolicyV2_spec(t *testing.T) {
	var (
		// The dedicated instance name only allow letters, digits and underscores (_).
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_throttling_policy.test"
		policy       throttles.ThrottlingPolicy
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigThrottlingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigThrottlingPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "type", "API-based"),
					resource.TestCheckResourceAttr(resourceName, "period", "15000"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(resourceName, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_ip_requests", "60"),
					resource.TestCheckResourceAttr(resourceName, "app_throttles.#", "0"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_spec(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourceName, "app_throttles.#", "1"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_specUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourceName, "app_throttles.#", "1"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigThrottlingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_throttles.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigThrottlingPolicyDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_throttling_policy" {
			continue
		}
		_, err := throttles.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 throttling policy (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigThrottlingPolicyExists(n string, app *throttles.ThrottlingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no throttling policy id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := throttles.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

func testAccApigThrottlingPolicy_base(rName, code string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_application" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id

  app_codes = ["%s"]
}
`, testAccApigApplication_base(rName), rName, utils.EncodeBase64String(code))
}

func testAccApigThrottlingPolicy_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigThrottlingPolicy_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%s_update"
  type              = "API-shared"
  period            = 10
  period_unit       = "MINUTE"
  max_api_requests  = 70
  max_user_requests = 45
  max_app_requests  = 45
  max_ip_requests   = 45
  description       = "Updated by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigThrottlingPolicy_spec(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Updated by script"
  
  app_throttles {
    max_api_requests     = 30
    throttling_object_id = huaweicloud_apig_application.test.id
  }
}
`, testAccApigThrottlingPolicy_base(rName, acctest.RandString(64)), rName)
}

func testAccApigThrottlingPolicy_specUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Updated by script"

  app_throttles {
    max_api_requests     = 45
    throttling_object_id = huaweicloud_apig_application.test.id
  }
}
`, testAccApigThrottlingPolicy_base(rName, acctest.RandString(64)), rName)
}
