---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_authorize_statistic"
description: |-
  Use this data source to get the application authorize statistics under the specified dedicated instance within
  HuaweiCloud.
---

# huaweicloud_apig_application_authorize_statistic

Use this data source to get the application authorize statistics under the specified dedicated instance within
HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_application_authorize_statistic" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the applications belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authed_nums` - The number of authorized applications.

* `unauthed_nums` - The number of unauthorized applications.
