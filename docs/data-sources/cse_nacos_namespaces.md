---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_nacos_namespaces"
description: |-
  Use this data source to query available Nacos namespaces within HuaweiCloud.
---

# huaweicloud_cse_nacos_namespaces

Use this data source to query available Nacos namespaces within HuaweiCloud.

## Example Usage

```hcl
variable "nacos_engine_id" {}

data "huaweicloud_cse_nacos_namespaces" "test" {
  engine_id = var.nacos_engine_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Nacos namespaces are located.  
  If omitted, the provider-level region will be used.

* `engine_id` - (Required, String) Specifies the ID of the Nacos microservice engine to which the namespaces belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namespaces` - All queried Nacos namespaces.  
  The [namespaces](#cse_nacos_namespaces) structure is documented below.

<a name="cse_nacos_namespaces"></a>
The `namespaces` block supports:

* `id` - The ID of the Nacos namespace.

  -> The reserved namespace (**public**) does not have the ID.

* `name` - The name of the Nacos namespace.
