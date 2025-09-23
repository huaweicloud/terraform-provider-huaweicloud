package identitycenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getApplicationInstanceHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}"
		getApplicationInstanceProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(getApplicationInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IdentityCenter client: %s", err)
	}

	getApplicationInstancePath := client.Endpoint + getApplicationInstanceHttpUrl
	getApplicationInstancePath = strings.ReplaceAll(getApplicationInstancePath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getApplicationInstancePath = strings.ReplaceAll(getApplicationInstancePath, "{application_instance_id}", state.Primary.ID)

	getApplicationInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getApplicationInstanceResp, err := client.Request("GET",
		getApplicationInstancePath, &getApplicationInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center application instance: %s", err)
	}

	return utils.FlattenResponse(getApplicationInstanceResp)
}

func TestAccApplicationInstance_basic(t *testing.T) {
	var obj interface{}

	name, _ := uuid.GenerateUUID()
	rName := "huaweicloud_identitycenter_application_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getApplicationInstanceResourceFunc,
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
				Config: testApplicationInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "display_name", "create"),
					resource.TestCheckResourceAttr(rName, "description", "create"),
					resource.TestCheckResourceAttrSet(rName, "response_config"),
					resource.TestCheckResourceAttrSet(rName, "response_schema_config"),
					resource.TestCheckResourceAttr(rName, "service_provider_config.0.audience", "https://create.com"),
					resource.TestCheckResourceAttr(rName, "service_provider_config.0.consumers.0.location", "https://create.com"),
					resource.TestCheckResourceAttr(rName, "security_config.0.ttl", "P9M"),
					resource.TestCheckResourceAttr(rName, "status", "CREATED"),
				),
			},
			{
				Config: testApplicationInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "display_name", "update"),
					resource.TestCheckResourceAttr(rName, "description", "update"),
					resource.TestCheckResourceAttrSet(rName, "response_config"),
					resource.TestCheckResourceAttrSet(rName, "response_schema_config"),
					resource.TestCheckResourceAttr(rName, "service_provider_config.0.audience", "https://update.com"),
					resource.TestCheckResourceAttr(rName, "service_provider_config.0.consumers.0.location", "https://update.com"),
					resource.TestCheckResourceAttr(rName, "service_provider_config.0.start_url", "https://update.com"),
					resource.TestCheckResourceAttr(rName, "security_config.0.ttl", "P6M"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLED"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateIdFunc:       testApplicationInstanceImportState(rName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_id"},
			},
		},
	})
}

func testApplicationInstanceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found: %s", name, rs)
		}

		return instanceID + "/" + rs.Primary.ID, nil
	}
}

func testApplicationInstance_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_catalog_applications" "test"{}

data "huaweicloud_identitycenter_application_templates" "test"{
  application_id = data.huaweicloud_identitycenter_catalog_applications.test.applications[0].application_id
}

resource "huaweicloud_identitycenter_application_instance" "test"{
  name         = "%s"
  template_id  = data.huaweicloud_identitycenter_application_templates.test.application_templates[0].template_id
  instance_id  = data.huaweicloud_identitycenter_instance.test.id
  display_name = "create"
  description  = "create"

  response_schema_config = jsonencode(
    {
      "properties" : {
        "key1" : {
          "attr_name_format" : "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include" : "YES"
        },
        "key2" : {
          "attr_name_format" : "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
          "include" : "YES"
        },
        "key5" : {
          "attr_name_format" : "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include" : "YES"
        },
        "key3" : {
          "attr_name_format" : "urn:oasis:names:tc:SAML:2.0:attrname-format:uri",
          "include" : "YES"
        },
        "key4" : {
          "attr_name_format" : "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include" : "YES"
        }
      },
      "subject" : {
        "name_id_format" : "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
        "include" : "REQUIRED"
      },
      "supported_name_id_formats" : []
    }
  )

  response_config = jsonencode(
    {
      "properties" : {
        "key1" : {
          "source" : [
            "$${user:email}"
          ]
        },
        "key2" : {
          "source" : [
            "$${user:familyName}"
          ]
        },
        "key5" : {
          "source" : [
            "$${user:preferredUsername}"
          ]
        },
        "key3" : {
          "source" : [
            "$${user:givenName}"
          ]
        },
        "key4" : {
          "source" : [
            "$${user:familyName}"
          ]
        }
      },
      "subject" : {
        "source" : [
          "$${user:name}"
        ]
      },
      "relay_state" : null,
      "ttl" : "PT1H"
    }
  )

  security_config {
    ttl = "P9M"
  }

  service_provider_config {
    audience                  = "https://create.com"
    require_request_signature = false
    consumers {
      binding       = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
      default_value = true
      location      = "https://create.com"
    }
  }
}
`, name)
}

func testApplicationInstance_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_catalog_applications" "test"{}

data "huaweicloud_identitycenter_application_templates" "test"{
  application_id = data.huaweicloud_identitycenter_catalog_applications.test.applications[0].application_id
}

resource "huaweicloud_identitycenter_application_instance" "test"{
  name         = "%s"
  template_id  = data.huaweicloud_identitycenter_application_templates.test.application_templates[0].template_id
  instance_id  = data.huaweicloud_identitycenter_instance.test.id
  display_name = "update"
  description  = "update"

    response_schema_config = jsonencode(
    {
      "properties": {
        "key1": {
          "attr_name_format": "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include": "YES"
        },
        "key2": {
          "attr_name_format": "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
          "include": "YES"
        },
        "key5": {
          "attr_name_format": "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include": "YES"
        },
        "key3": {
          "attr_name_format": "urn:oasis:names:tc:SAML:2.0:attrname-format:uri",
          "include": "YES"
        },
        "key4": {
          "attr_name_format": "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
          "include": "YES"
        }
      },
      "subject": {
        "name_id_format": "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
        "include": "REQUIRED"
      },
      "supported_name_id_formats": []
    }
  )

  response_config = jsonencode(
    {
      "properties": {
        "key1": {
          "source": [
            "$${user:email}"
          ]
        },
        "key2": {
          "source": [
            "$${user:familyName}"
          ]
        },
        "key5": {
          "source": [
            "$${user:preferredUsername}"
          ]
        },
        "key3": {
          "source": [
            "$${user:givenName}"
          ]
        },
        "key4": {
          "source": [
            "$${user:familyName}"
          ]
        }
      },
      "subject": {
        "source": [
          "$${user:name}"
        ]
      },
      "relay_state": null,
      "ttl": "PT1H"
    }
  )

  security_config {
    ttl = "P6M"
  }

  service_provider_config {
    audience                  = "https://update.com"
    require_request_signature = false
    consumers{
        binding       = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
        default_value = true
        location      = "https://update.com"
    } 
    start_url = "https://update.com"
  }
  status = "ENABLED"
}
`, name)
}
