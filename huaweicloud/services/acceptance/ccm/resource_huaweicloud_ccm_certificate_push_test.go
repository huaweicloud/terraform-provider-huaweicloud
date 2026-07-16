package ccm

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	ccmcertificatepush "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ccm"
)

func getCertificatePushResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	targetsLength, err := strconv.Atoi(state.Primary.Attributes["targets.#"])
	if err != nil {
		return nil, fmt.Errorf("[ERROR] convert the string %q to int failed", state.Primary.Attributes["targets.#"])
	}

	rst := make([]interface{}, 0, len(state.Primary.Attributes["targets.#"]))
	for index := 0; index < targetsLength; index++ {
		project := state.Primary.Attributes[fmt.Sprintf("targets.%d.project_name", index)]
		certId := state.Primary.Attributes[fmt.Sprintf("targets.%d.cert_id", index)]

		switch state.Primary.Attributes["service"] {
		case "ELB":
			if elbRstMap := ccmcertificatepush.FlattenElbCertificateDetail(conf, project, certId); elbRstMap != nil {
				rst = append(rst, elbRstMap)
			}
		case "WAF":
			if wafRstMap := ccmcertificatepush.FlattenWafCertificateDetail(conf, project, certId); wafRstMap != nil {
				rst = append(rst, wafRstMap)
			}
		}
	}

	if len(rst) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return rst, nil
}

func TestAccCcmCertificatePush_elb(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_ccm_certificate_push.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificatePushResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMPushCertificateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      tesCcmCertificatePush_basic_step1(),
				ExpectError: regexp.MustCompile("all attempts to push the certificate to the services failed in creation"),
			},
			{
				Config:             tesCcmCertificatePush_basic_step2(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.project_name", "cn-north-4"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
				),
			},
			{
				Config:             tesCcmCertificatePush_basic_step3(),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("all attempts to push the certificate to the services failed in update operation"),
			},
			{
				Config: tesCcmCertificatePush_basic_step4(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
			{
				Config:             tesCcmCertificatePush_basic_step5(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
			{
				Config: tesCcmCertificatePush_basic_step6(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.project_name", "cn-north-4"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
				),
			},
		},
	})
}

// tesCcmCertificatePush_basic_step1 using to test the scenario of all certificate push failures in creation operation.
func tesCcmCertificatePush_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "unexpected-error-4"
  }

  targets {
    project_name = "unexpected-error-5"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// tesCcmCertificatePush_basic_step2 using to test the scenario of partial push success.
// This test step needs to ignore changes to `targets`.
func tesCcmCertificatePush_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "cn-north-4"
  }

  targets {
    project_name = "unexpected-error-5"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// tesCcmCertificatePush_basic_step3 using to test the scenario of all certificate push failures in update operation.
func tesCcmCertificatePush_basic_step3() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "unexpected-error-4"
  }

  targets {
    project_name = "unexpected-error-5"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// tesCcmCertificatePush_basic_step4 using to test that all pushes are successful.
func tesCcmCertificatePush_basic_step4() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "cn-north-4"
  }

  targets {
    project_name = "ap-southeast-3"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// tesCcmCertificatePush_basic_step5 using to test editing part push succeeds and the other fails.
func tesCcmCertificatePush_basic_step5() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "cn-north-4"
  }

  targets {
    project_name = "cn-east-3"
  }

  targets {
    project_name = "unexpected-error-5"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// tesCcmCertificatePush_basic_step6 using to test the successful push of a certificate.
func tesCcmCertificatePush_basic_step6() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "ELB"

  targets {
    project_name = "cn-north-4"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID)
}

// Before executing this test case, please confirm that the WAF instances has been purchased in the target project.
// Otherwise, the test case may fail to execute.
func TestAccCcmCertificatePush_waf(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_ccm_certificate_push.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificatePushResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMPushCertificateID(t)
			// Please ensure that there is an available WAF instance in the custom region.
			acceptance.TestAccPrecheckCustomRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:             tesCcmCertificatePush_waf_step1(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.project_name", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
				),
			},
			{
				Config: tesCcmCertificatePush_waf_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.project_name", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cert_name"),
				),
			},
		},
	})
}

// tesCcmCertificatePush_waf_step1 using to test the scenario where one push succeeds and the other fails.
func tesCcmCertificatePush_waf_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "WAF"

  targets {
    project_name = "%s"
  }

  targets {
    project_name = "unexpected-error-4"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID, acceptance.HW_CUSTOM_REGION_NAME)
}

// tesCcmCertificatePush_waf_step2 using to test the scenario where one push succeeds.
func tesCcmCertificatePush_waf_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = "%s"
  service        = "WAF"

  targets {
    project_name = "%s"
  }
}`, acceptance.HW_CERT_BATCH_PUSH_ID, acceptance.HW_CUSTOM_REGION_NAME)
}
