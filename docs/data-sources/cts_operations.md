---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_operations"
description: |-
  Use this data source to list all operations on a cloud service.
---

# huaweicloud_cts_operations

Use this data source to list all operations on a cloud service.

## Example Usage

```hcl
data "huaweicloud_cts_operations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `service_type` - (Optional, String) Specifies the type of the cloud service on which operations are performed.

* `resource_type` - (Optional, String) Specifies the type of the resource on which operations are performed.
  If this parameter is used, `service_type` is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `operations` - All operations on the cloud service.

  The [operations](#operations_struct) structure is documented below.

<a name="operations_struct"></a>
The `operations` block supports:

* `service_type` - The type of the cloud service on which operations are performed.

* `resource_type` - The type of the resource on which operations are performed.

* `operation_list` - The array of operation names.
