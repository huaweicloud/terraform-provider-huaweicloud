---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_service_principals"
description: |-
  Use this data source to get the list of services principals within HuaweiCloud.
---

# huaweicloud_identityv5_service_principals

Use this data source to get the list of services principals within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_service_principals" "test" {}
```

## Argument Reference

The following arguments are supported:

* `language` - (Optional, String) Specifies the language of the information returned by the interface.  
  The valid values are as follows:
  + **zh-cn**
  + **en-us**

  Defaults to **zh-cn**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_principals` - The list of service principals.  
  The [service_principals](#Identityv5_service_principals) structure is documented below.

<a name="Identityv5_service_principals"></a>
The `service_principals` block supports:

* `service_principal` - The ID of the service principal.

* `description` - The description of the service principal.

* `display_name` - The display name of the service principal.

* `service_catalog` - The cloud service name.
