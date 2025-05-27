---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_agent_configs"
description: |-
  Use this data source to get SMS agent configs.
---

# huaweicloud_sms_agent_configs

Use this data source to get SMS agent configs.

## Example Usage

```hcl
data "huaweicloud_sms_agent_configs" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `config` - Indicates the main region, obs domain, disk type and information to be added later.

* `regions` - Indicates the region list.
