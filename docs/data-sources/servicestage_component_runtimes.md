---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_component_runtimes"
description: ""
---

# huaweicloud_servicestage_component_runtimes

Use this data source to query available runtimes within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestage_component_runtimes" "test" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the component runtimes.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the runtime name to use for filtering.
  For the runtime names corresponding to each type of component, please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-servicestage/servicestage_user_0411.html).

* `default_port` - (Optional, Int) Specifies the default container port to use for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `runtimes` - The list of runtime details.

The `runtimes` block contains:

* `name` - The runtime name.

* `default_port` - The default container port.

* `description` - The runtime description.
