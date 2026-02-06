---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_objectstorage_cloud_type"
description: |-
  Use this data source to query OMS supported cloud vendors.
---

# huaweicloud_oms_objectstorage_cloud_type

Use this data source to query OMS supported cloud vendors.

## Example Usage

```hcl
data "huaweicloud_oms_objectstorage_cloud_type" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the connection end type.
  The value can be **src** (source) or **dst** (destination).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The data source ID.

* `vendors` - The list of supported cloud vendors.
