package identitycenter

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		listProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	listPath = strings.ReplaceAll(listPath, "{application_instance_id}", state.Primary.Attributes["application_instance_id"])

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	certificateId := state.Primary.ID

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center application certificate: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	certificate := utils.PathSearch(fmt.Sprintf("application_instance_certificates[?certificate_id =='%s']|[0]", certificateId), listRespBody, nil)
	if certificate == nil {
		return nil, fmt.Errorf("error get Identity Center application certificate")
	}
	return certificate, nil
}

func TestAccIdentityCenterApplicationCertificate_basic(t *testing.T) {
	var obj interface{}

	uuid, _ := uuid.GenerateUUID()
	rName := "huaweicloud_identitycenter_application_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getApplicationCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testApplicationCertificate_basic(uuid),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "application_instance_id",
						"huaweicloud_identitycenter_application_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "algorithm"),
					resource.TestCheckResourceAttrSet(rName, "certificate"),
					resource.TestCheckResourceAttrSet(rName, "expiry_date"),
					resource.TestCheckResourceAttr(rName, "status", "INACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "key_size"),
					resource.TestCheckResourceAttrSet(rName, "issue_date"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterApplicationCertificateImportState(rName),
			},
		},
	})
}

func testApplicationCertificate_basic(uuid string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_application_certificate" "test"{
  depends_on              = [huaweicloud_identitycenter_application_instance.test]
  application_instance_id = huaweicloud_identitycenter_application_instance.test.id
  instance_id             = data.huaweicloud_identitycenter_instance.test.id
}

`, testApplicationInstance_basic(uuid))
}

func testIdentityCenterApplicationCertificateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}

		applicationInstanceID := rs.Primary.Attributes["application_instance_id"]
		if applicationInstanceID == "" {
			return "", fmt.Errorf("attribute (application_instance_id) of Resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", instanceID, applicationInstanceID, rs.Primary.ID), nil
	}
}
