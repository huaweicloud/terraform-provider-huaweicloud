---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_appcodes"
description: |-
  Use this data source to query the APPCODEs of the specified APIG application within HuaweiCloud.
---

# huaweicloud_apig_appcodes

Use this data source to query the APPCODEs of the specified APIG application within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_id" {}

data "huaweicloud_apig_appcodes" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the APPCODEs are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the application belongs.

* `application_id` - (Required, String) Specifies the ID of the application to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `appcodes` - The list of the APPCODEs of the specified application.  
  The [appcodes](#apig_appcodes_attr) structure is documented below.

<a name="apig_appcodes_attr"></a>
The `appcodes` block supports:

* `id` - The ID of the APPCODE.

* `value` - The APPCODE value (content).

* `application_id` - The ID of the application.

* `created_at` - The creation time of the APPCODE, in RFC3339 format.
