---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_datastores"
description: |-
  Use this data source to get the list of GaussDB OpenGauss engines.
---

# huaweicloud_gaussdb_opengauss_datastores

Use this data source to get the list of GaussDB OpenGauss engines.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_datastores" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastores` - Specifies the  DB engines.

  The [datastores](#datastores_struct) structure is documented below.

<a name="datastores_struct"></a>
The `datastores` block supports:

* `supported_versions` - Specifies the engine versions supported by the deployment model.

* `instance_mode` - Specifies the deployment model.
  The value can be:
  + **ha**: primary/standby
  + **independent**: independent
