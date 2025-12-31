---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_reports_emails"
description: |-
  Use this data source to query the recipient email status.
---

# huaweicloud_secmaster_reports_emails

Use this data source to query the recipient email status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "email_address" {}

data "huaweicloud_secmaster_reports_emails" "test" {
  workspace_id  = var.workspace_id
  email_address = var.email_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `email_address` - (Required, String) Specifies the recipient email.
  Support multiple email addresses, separated by semicolon (;). e.g. `test1@example.com;test2@example.com`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `emails` - The email status information.

  The [emails](#email_status_struct) structure is documented below.

<a name="email_status_struct"></a>
The `emails` block supports:

* `report_address` - The email address.

* `email_status` - The email status.
  The valid values are as follows:
  + **true**: Indicates it can be sent directly.
  + **false**: Indicates a confirmation email needs to be sent first.
