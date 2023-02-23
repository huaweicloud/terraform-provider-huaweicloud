---
subcategory: "Distributed Database Middleware (DDM)"
---

# huaweicloud_ddm_flavors

Use this data source to get the list of DDM flavors.

## Example Usage

```hcl
data "huaweicloud_ddm_engines" test {
  version    = "3.0.8.5"
}

data "huaweicloud_ddm_flavor" test {
  engine_id  = data.huaweicloud_ddm_engine.test.id
  group_type = "X86"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_id` - (Required, String) Specifies the ID of a engine.

* `group_type` - (Required, String) Specifies the compute resource architecture type. The options are **X86** and **ARM**.

* `type_code` - (Optional, String) Specifies the resource type code.

* `code` - (Optional, String) Specifies the VM flavor types recorded in DDM.

* `iaas_code` - (Optional, String) Specifies the VM flavor types recorded by the IaaS layer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `flavors` - Indicates the list of DDM compute flavors.
  The [Flavor](#DdmFlavors_Flavor) structure is documented below.

<a name="DdmFlavors_Flavor"></a>
The `Flavor` block supports:

* `id` - Indicates the compute resource architecture type.

* `type_code` - Indicates the resource type code.

* `code` - Indicates the VM flavor types recorded in DDM.

* `iaas_code` - Indicates the VM flavor types recorded by the IaaS layer.
