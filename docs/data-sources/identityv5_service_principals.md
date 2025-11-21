---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_service_principals"
description: |-
  Use this data source to get the list of services principals in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_service_principals

## Example Usage

```hcl
data "huaweicloud_identityv5_service_principals" "test" {}
```

## Argument Reference

The following arguments are supported:

* `language` - (Optional, String) Specifies select the language of the information returned by the interface,
  which can be Chinese ("zh-cn") or English ("en-us"), with Chinese as the default.
  Default value: zh-cn

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `service_principals` - A list of service principals. Each element contains the following attributes:
  The [service_principals](#Identityv5_principals) structure is documented below.

<a name="Identityv5_principals"></a>
The `service_principals` block supports:

* `service_principal` - Indicates the ID of the service principal.

* `description` - Indicates the description of the service principal.

* `display_name` - Indicates the display name of the service principal.

* `service_catalog` - Indicates the cloud service name.
