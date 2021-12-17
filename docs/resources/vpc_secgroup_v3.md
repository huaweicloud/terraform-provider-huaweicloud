---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_secgroup_v3

Manages a Security Group resource within HuaweiCloud. 

## Example Usage

```hcl
resource "huaweicloud_vpc_secgroup_v3" "test" {
  region                = "cn-north-4"
  dry_run               = false
  name                  = "aaa"
  description           = "123"
  enterprise_project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the security group resource. If omitted, the
  provider-level region will be used. Changing this creates a new security group resource.

* `dry_run` - (Optional)Pre check this request only
Value range:True: when sending a check request, the security group will not be created. The check items include whether the required parameters, request format and business restrictions are filled in. If the check fails, the corresponding error is returned. If the check passes, the response code 202 is returned.
False (default): Send a normal request and directly create a security group.

* `name` - (Required, String) Specifies a unique name for the security group.

* `description` - (Optional, String) Specifies the description for the security group.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the security group.
  Changing this creates a new security group.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minute.

