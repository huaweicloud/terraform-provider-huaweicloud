---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_templates"
description: |
  Use this data source to get the Identity Center application templates.
---

# huaweicloud_identitycenter_application_templates

Use this data source to get the Identity Center application templates.

## Example Usage

```hcl
var application_id {}

data "huaweicloud_identitycenter_application_templates" "test"{
  application_id = var.application_id
}
```

## Argument Reference

* `application_id` - (Required, String) Specifies the ID of the application.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `application_templates` - The list of the application templates.
  The [application_templates](#application_templates_struct) structure is documented below.

<a name="application_templates_struct"></a>
The `application_templates` block supports:

* `application_type` - The type of the application.

* `description` - The description of the application.

* `display_name` - The display name of the application.

* `sso_protocol` - The sso protocol supported by application.

* `template_id` - The ID of the application template.

* `template_version` - The version of the application template.

* `response_config` - The response configuration of the application.

* `response_schema_config` - The response schema configuration of the application.

* `security_config` - The security configuration of the application.
  The [security_config](#security_config_struct) structure is documented below.

* `service_provider_config` - The service provider configuration of the application.
  The [service_provider_config](#service_provider_config_struct) structure is documented below.

<a name="security_config_struct"></a>
The `security_config` block supports:

* `ttl` - The time to alive of the application instance certificate.

<a name="service_provider_config_struct"></a>
The `service_provider_config` block supports:

* `audience` - The audience of the application.

* `require_request_signature` - Whether the application instance requires request signature.

* `consumers` - The consumers of the application instance.
  The [consumers](#consumers_struct) structure is documented below.

* `start_url` - The start url of the application instance.

<a name="consumers_struct"></a>
The `consumers` block supports:

* `location` - The location url of the application instance.

* `binding` - The bind method of the application instance.

* `default_value` - Whether this consumer is a default one.
