package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getThrottlingPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return throttles.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccThrottlingPolicy_basic(t *testing.T) {
	var (
		policy throttles.ThrottlingPolicy

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccApigThrottlingPolicy_base(name)

		rName = "huaweicloud_apig_throttling_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &policy, getThrottlingPolicyFunc)

		rNamePre = "huaweicloud_apig_throttling_policy.pre_test"
		rcPre    = acceptance.InitResourceCheck(rNamePre, &policy, getThrottlingPolicyFunc)
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
				// Check whether illegal type values ​​can be intercepted normally (create phase).
				Config:      testAccApigThrottlingPolicy_basic_step1(baseConfig, name),
				ExpectError: regexp.MustCompile("invalid throttling policy type: NON-Exist-Type"),
			},
			{
				// Check whether illegal application ID ​​can be intercepted normally (create phase).
				Config:      testAccApigThrottlingPolicy_basic_step2(baseConfig, name),
				ExpectError: regexp.MustCompile("error creating special application throttling policy"),
			},
			{
				// Check whether illegal call limit ​​can be intercepted normally (create phase).
				// The API does not check whether the value or format of the user ID is legal.
				Config:      testAccApigThrottlingPolicy_basic_step3(baseConfig, name),
				ExpectError: regexp.MustCompile("The parameter value is too small"),
			},
			{
				Config: testAccApigThrottlingPolicy_basic_step4(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcPre.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePre, "name", name+"_pre_test"),
					resource.TestCheckResourceAttr(rNamePre, "description", "Created by script"),
					resource.TestCheckResourceAttr(rNamePre, "type", "API-based"),
					resource.TestCheckResourceAttr(rNamePre, "period", "15000"),
					resource.TestCheckResourceAttr(rNamePre, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(rNamePre, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(rNamePre, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(rNamePre, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(rNamePre, "max_ip_requests", "60"),
					resource.TestCheckResourceAttr(rNamePre, "app_throttles.#", "0"),
					resource.TestCheckResourceAttr(rNamePre, "user_throttles.#", "0"),
					resource.TestMatchResourceAttr(rNamePre, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				// Check whether illegal type values ​​can be intercepted normally (update phase).
				Config:      testAccApigThrottlingPolicy_basic_step5(baseConfig, name),
				ExpectError: regexp.MustCompile("invalid throttling policy type: NON-Exist-Type"),
			},
			{
				// Check whether illegal application ID ​​can be intercepted normally (update phase).
				Config:      testAccApigThrottlingPolicy_basic_step6(baseConfig, name),
				ExpectError: regexp.MustCompile("error updating special app throttles"),
			},
			{
				// Check whether illegal call limit ​​can be intercepted normally (update phase).
				// The API does not check whether the value or format of the user ID is legal.
				Config:      testAccApigThrottlingPolicy_basic_step7(baseConfig, name),
				ExpectError: regexp.MustCompile("The parameter value is too small"),
			},
			{
				Config: testAccApigThrottlingPolicy_basic_step8(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "type", "API-based"),
					resource.TestCheckResourceAttr(rName, "period", "15000"),
					resource.TestCheckResourceAttr(rName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(rName, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(rName, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(rName, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(rName, "max_ip_requests", "60"),
					resource.TestCheckResourceAttr(rName, "app_throttles.#", "2"),
					resource.TestCheckResourceAttr(rName, "user_throttles.#", "2"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_basic_step9(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(rName, "type", "API-shared"),
					resource.TestCheckResourceAttr(rName, "period", "10"),
					resource.TestCheckResourceAttr(rName, "period_unit", "MINUTE"),
					resource.TestCheckResourceAttr(rName, "max_api_requests", "70"),
					resource.TestCheckResourceAttr(rName, "max_user_requests", "45"),
					resource.TestCheckResourceAttr(rName, "max_app_requests", "45"),
					resource.TestCheckResourceAttr(rName, "max_ip_requests", "45"),
					resource.TestCheckResourceAttr(rName, "app_throttles.#", "2"),
					resource.TestCheckResourceAttr(rName, "user_throttles.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccThrottlingPolicyImportStateFunc(rName),
			},
		},
	})
}

func testAccThrottlingPolicyImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApigThrottlingPolicy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identity_users" "test" {}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

# If you want to test a resource (huaweicloud_apig_throttling_policy) that does not contain app-specific throttling
# policy after step 2, then the application resources must be retained.
# If you delete this script (application resources definition), the application deletion operation and the special
# throttling policies update operation will be executed in parallel.
# If the APP is deleted before the throttling policies update operation complete, it will cause an error during the
# update API. The order of the two operations cannot be controlled.
resource "huaweicloud_apig_application" "test" {
  count = 3

  instance_id = local.instance_id
  name        = "%[2]s_${count.index}"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccApigThrottlingPolicy_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "invalid_type" {
  instance_id      = local.instance_id
  name             = "%[2]s_invalid_type"
  type             = "NON-Exist-Type"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step2(baseConfig, name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "invalid_app_id" {
  instance_id      = local.instance_id
  name             = "%[2]s_invalid_app_id"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100

  app_throttles {
    max_api_requests     = 30
    throttling_object_id = "%[3]s"
  }
}
`, baseConfig, name, randUUID)
}

func testAccApigThrottlingPolicy_basic_step3(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "invalid_user_id" {
  instance_id      = local.instance_id
  name             = "%[2]s_invalid_user_id"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100

  user_throttles {
    max_api_requests     = -1
    throttling_object_id = "INVALID_OBJECT_ID_CONTENT"
  }
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step4(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "pre_test" {
  instance_id       = local.instance_id
  name              = "%[2]s_pre_test"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step5(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "pre_test" {
  instance_id       = local.instance_id
  name              = "%[2]s_pre_test"
  type              = "NON-Exist-Type"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step6(baseConfig, name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "pre_test" {
  instance_id       = local.instance_id
  name              = "%[2]s_pre_test"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"

  app_throttles {
    max_api_requests     = 30
    throttling_object_id = "%[3]s"
  }
}
`, baseConfig, name, randUUID)
}

func testAccApigThrottlingPolicy_basic_step7(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "pre_test" {
  instance_id       = local.instance_id
  name              = "%[2]s_pre_test"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"

  user_throttles {
    max_api_requests     = -1
    throttling_object_id = "INVALID_OBJECT_ID_CONTENT"
  }
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step8(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = local.instance_id
  name              = "%[2]s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"

  dynamic "app_throttles" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 0, 2)

    content {
      max_api_requests     = 30
      throttling_object_id = app_throttles.value
    }
  }

  dynamic "user_throttles" {
    for_each = slice(data.huaweicloud_identity_users.test.users[*].id, 0, 2)

    content {
      max_api_requests     = 30
      throttling_object_id = user_throttles.value
    }
  }
}
`, baseConfig, name)
}

func testAccApigThrottlingPolicy_basic_step9(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = local.instance_id
  name              = "%[2]s"
  type              = "API-shared"
  period            = 10
  period_unit       = "MINUTE"
  max_api_requests  = 70
  max_user_requests = 45
  max_app_requests  = 45
  max_ip_requests   = 45
  description       = "Updated by script"

  dynamic "app_throttles" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 1, 3)

    content {
      max_api_requests     = 40
      throttling_object_id = app_throttles.value
    }
  }

  dynamic "user_throttles" {
    for_each = slice(data.huaweicloud_identity_users.test.users[*].id, 1, 3)

    content {
      max_api_requests     = 40
      throttling_object_id = user_throttles.value
    }
  }
}
`, baseConfig, name)
}
