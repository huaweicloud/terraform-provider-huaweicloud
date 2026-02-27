package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV3TrustAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getAgencyHttpUrl := "v5/agencies/{agency_id}"
	getAgencyPath := client.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{agency_id}", state.Primary.ID)
	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAgencyResp, err := client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM trust agency: %s", err)
	}
	return utils.FlattenResponse(getAgencyResp)
}

func TestAccV3TrustAgency_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_trust_agency.test"
		rName        = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &obj, getV3TrustAgencyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckIAMV5(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3TrustAgency_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "trust_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccV3TrustAgency_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "duration", "7200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
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

func testAccV3TrustAgency_basic_step1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%s"
  policy_names = ["NATReadOnlyPolicy"]
  description  = "Created by terraform script"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.APIG",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccV3TrustAgency_basic_step2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%s"
  policy_names = ["NATReadOnlyPolicy", "RDSReadOnlyPolicy"]
  duration     = 7200
  description  = "Updated by terraform script"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.APIG",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )

  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
}
`, rName)
}
