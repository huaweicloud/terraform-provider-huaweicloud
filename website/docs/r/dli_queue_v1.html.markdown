---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_queue_v1"
sidebar_current: "docs-huaweicloud-resource-dli-queue-v1"
description: |-
  queue management
---

# huaweicloud\_dli\_queue\_v1

queue management

## Example Usage

### create a queue

```hcl
resource "huaweicloud_dli_queue_v1" "queue" {
  name = "terraform_dli_queue_v1_test"
  cu_count = 4
}
```

## Argument Reference

The following arguments are supported:

* `cu_count` -
  (Required)
  Minimum number of CUs that are bound to a queue. The value can be 4,
  16, or 64. Changing this parameter will create a new resource.

* `name` -
  (Required)
  Name of a queue. The name can contain only digits, letters, and
  underscores (_), but cannot contain only digits or start with an
  underscore (_). Changing this parameter will create a new resource.

* `description` -
  (Optional)
  Description of a queue. Changing this parameter will create a new resource.

* `management_subnet_cidr` -
  (Optional)
  CIDR of the management subnet. Changing this parameter will create a new resource.

* `subnet_cidr` -
  (Optional)
  Subnet CIDR. Changing this parameter will create a new resource.

* `vpc_cidr` -
  (Optional)
  VPC CIDR. Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `create_time` -
  Time when a queue is created.
