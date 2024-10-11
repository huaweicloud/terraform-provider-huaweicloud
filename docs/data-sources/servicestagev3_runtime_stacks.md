---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_runtime_stacks"
description: |-
  Use this data source to query the list of runtime stacks within HuaweiCloud.
---

# huaweicloud_servicestagev3_runtime_stacks

Use this data source to query the list of runtime stacks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestagev3_runtime_stacks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the runtime stacks are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `runtime_stacks` - All runtime stack details.  
  The [runtime_stacks](#servicestage_v3_runtime_stacks) structure is documented below.

<a name="servicestage_v3_runtime_stacks"></a>
The `runtime_stacks` block supports:

* `name` - The name of the runtime stack.

* `deploy_mode` - The deploy mode.
  + **container**
  + **virtualmachine**

* `version` - The version number.

* `url` - The URL of the runtime stack.

* `type` - The type of the runtime stack.
  + **Nodejs**
  + **Java**
  + **Tomcat**
  + **Python**
  + **Php**
  + **Docker**

* `status` - The status of the runtime stack.
  + **Supported**
  + **Deprecated**
