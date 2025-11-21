---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_instances"
description: |-
  Use this data source to get the Identity Center application instances.
---

# huaweicloud_identitycenter_application_instances

Use this data source to get the Identity Center application instances.

## Example Usage

```hcl
var instance_id {}

data "huaweicloud_identitycenter_application_instances" "test" { 
  instance_id = var.instance_id
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the Identity Center instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `application_instances` - The list of the application instances.
  The [application_instances](#application_instances_struct) structure is documented below.

<a name="application_instances_struct"></a>
The `application_instances` block supports:

* `description` - The description of the application instance.

* `display_name` - The display name of the application instance.

* `name` - The name of the application instance.

* `status` - The status of the application instance.

* `response_config` - The response configuration of the application instance.

* `response_schema_config` - The response schema configuration of the application instance.

* `security_config` - The security configuration of the application instance.
  The [security_config](#security_config_struct) structure is documented below.

* `service_provider_config` - The service provider configuration of the application instance.
  The [service_provider_config](#service_provider_config_struct) structure is documented below.

* `identity_provider_config` - The identity provider configuration of the application instance.
  The [identity_provider_config](#identity_provider_config_struct) structure is documented below.

* `active_certificate` - The active certificate of the application instance.
  The [active_certificate](#active_certificate_struct) structure is documented below.

* `visible` - Whether the application instance is visible.

* `client_id` - The client ID of the application instance.

* `end_user_visible` - Whether the application instance is visible for end user.

* `managed_account` - The managed account of the application instance.

<a name="security_config_struct"></a>
The `security_config` block supports:

* `ttl` - The time to alive of the application instance certificate.

<a name="service_provider_config_struct"></a>
The `service_provider_config` block supports:

* `audience` - The audience of the application instance.

* `require_request_signature` - Whether the application instance requires request signature.

* `consumers` - The consumers of the application instance.
  The [consumers](#consumers_struct) structure is documented below.

* `start_url` - The start url of the application instance.

<a name="consumers_struct"></a>
The `consumers` block supports:

* `location` - The location url of the application instance.

* `binding` - The bind method of the application instance.

* `default_value` - Whether this consumer is a default one.

<a name="identity_provider_config_struct"></a>
The `identity_provider_config` block supports:

* `issuer_url` - The issuer url of the application instance.

* `metadata_url` - The metadata url of the application instance.

* `remote_login_url` - The remote login url of the application instance.

* `remote_logout_url` - The remote logout url of the application instance.

<a name="active_certificate_struct"></a>
The `active_certificate` block supports:

* `algorithm` - The algorithm of the application instance certificate.

* `certificate` - The certificate of the application instance.

* `certificate_id` - The ID of the application instance certificate.

* `expiry_date` - The expiry date of the application instance certificate.

* `status` - The status of the application instance certificate.

* `key_size` - The key size of the application instance certificate.

* `issue_date` - The issue date of the application instance certificate.
