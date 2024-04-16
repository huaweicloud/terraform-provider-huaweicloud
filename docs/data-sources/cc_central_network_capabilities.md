---
subcategory: "Cloud Connect (CC)"
---

# huaweicloud_cc_central_network_capabilities

Use this data source to get the list of CC central network capabilities.

## Example Usage

```hcl
variable "capability" {}

data "huaweicloud_cc_central_network_capabilities" "test" {
  capability = var.capability
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `capability` - (Optional, String) Specifies the capability of the central network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - Central network capability list.
  The [capabilities](#Capabilities) structure is documented below.

<a name="Capabilities"></a>
The `capabilities` block supports:

* `capability` - The capability of the central network.

* `domain_id` - The ID of the account that the central network belongs to.

* `specifications` - The specifications of the central network capability.
