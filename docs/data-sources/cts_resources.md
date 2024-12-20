---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_resources"
description: |-
  Use this data source to get the list of resources involved in the traces.
---

# huaweicloud_cts_resources

Use this data source to get the list of resources involved in the traces.

## Example Usage

```hcl
data "huaweicloud_cts_resources" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The resource list.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `service_type` - The cloud service type.

* `resource` - The resources corresponding to the cloud services.
