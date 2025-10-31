---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_instance"
description: |
  Manages an Identity Center application instance resource within HuaweiCloud.
---

# huaweicloud_identitycenter_application_instance

Manages an Identity Center application instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "template_id" {}
variable "instance_id" {}

resource "huaweicloud_identitycenter_application_instance" "test"{
  name        = var.name
  template_id = var.template_id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM Identity Center instance.

* `name` - (Required, String, NonUpdatable) Specifies the name of the application instance.

* `template_id` - (Required, String, NonUpdatable) Specifies the ID of the application instance template.

* `description` - (Optional, String) Specifies the description of the application instance.

* `display_name` - (Optional, String) Specifies the display name of the application instance.

* `metadata` - (Optional, String) Specifies the metadata of the application instance.

* `status` - (Optional, String) Specifies the status of the application instance.

* `response_config` - (Optional, String) Specifies the response configuration of the application instance.

* `response_schema_config` - (Optional, String) Specifies the response schema configuration of the application instance.

* `security_config` - (Optional, List) Specifies the security configuration of the application instance.
  The [security_config](#security_config_struct) structure is documented below.

* `service_provider_config` - (Optional, List) Specifies the service provider configuration of the application instance.
  The [service_provider_config](#service_provider_config_struct) structure is documented below.

<a name="security_config_struct"></a>
The `security_config` block supports:

* `ttl` - (Optional, String) Specifies the time to alive of the application instance certificate.

<a name="service_provider_config_struct"></a>
The `service_provider_config` block supports:

* `audience` - (Optional, String) Specifies the audience of the application instance.

* `require_request_signature` - (Optional, Bool) Whether the application instance requires request signature.

* `consumers` - (Optional, List) Specifies the consumers of the application instance.
  The [consumers](#consumers_struct) structure is documented below.

* `start_url` - (Optional, String) Specifies the start url of the application instance.

<a name="consumers_struct"></a>
The `consumers` block supports:

* `location` - (Optional, String) Specifies the location url of the application instance.

* `binding` - (Optional, String) Specifies the bind method of the application instance.

* `default_value` - (Optional, Bool) Whether this consumer is a default one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `identity_provider_config` - The identity provider configuration of the application instance.
  The [identity_provider_config](#identity_provider_config_struct) structure is documented below.

* `active_certificate` - The active certificate of the application instance.
  The [active_certificate](#active_certificate_struct) structure is documented below.

* `visible` - Whether the application instance is visible.

* `client_id` - The client ID of the application instance.

* `end_user_visible` - Whether the application instance is visible for end user.

* `managed_account` - The managed account of the application instance.

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

## Import

The IdentityCenter application instance can be imported using the `instance_id` and `application_instance_id` separated
by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_application_instance.test <instance_id>/<application_instance_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template_id`, `metadata`. It is generally
recommended running `terraform plan` after importing an IdentityCenter application instance. You can then decide if
changes should be applied to the IdentityCenter application instance, or the resource definition should be updated to
align with the application instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_application_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      template_id,
      metadata,
    ]
  }
}
```
