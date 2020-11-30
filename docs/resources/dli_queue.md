---
subcategory: "Data Lake Insight (DLI)"
---

# huaweicloud\_dli\_queue

DLI Queue management
This is an alternative to `huaweicloud_dli_queue`

## Example Usage

### create a queue

```hcl
resource "huaweicloud_dli_queue" "queue" {
  name     = "terraform_dli_queue_test"
  cu_count = 4
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DLI queue resource. If omitted, the provider-level region will be used. Changing this creates a new DLI Queue resource.

* `cu_count` - (Required, Int, ForceNew) Minimum number of CUs that are bound to a queue. The value can be 4,
  16, or 64. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Name of a queue. The name can contain only digits, letters, and
  underscores (_), but cannot contain only digits or start with an
  underscore (_). Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Description of a queue. Changing this parameter will create a new resource.

* `management_subnet_cidr` - (Optional, String, ForceNew) CIDR of the management subnet. Changing this parameter will create a new resource.

* `subnet_cidr` - (Optional, String, ForceNew) Subnet CIDR. Changing this parameter will create a new resource.

* `vpc_cidr` - (Optional, String, ForceNew) VPC CIDR. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `create_time` -  Time when a queue is created.
