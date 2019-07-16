---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_cluster_v1"
sidebar_current: "docs-huaweicloud-resource-cdm-cluster-v1"
description: |-
  cdm cluster management
---

# huaweicloud\_cdm\_cluster\_v1

cdm cluster management

## Example Usage

### create a cdm cluster

```hcl
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cdm_cluster_v1" "cluster" {
  availability_zone = "{{ availability_zone }}"
  flavor_id         = "{{ flavor_id }}"
  name              = "terraform_test_cdm_cluster"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  subnet_id         = "{{ network_id }}"
  vpc_id            = "{{ vpc_id }}"
  version           = "{{ version }}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` -
  (Required)
  Available zone.  Changing this parameter will create a new resource.

* `flavor_id` -
  (Required)
  Flavor id.  Changing this parameter will create a new resource.

* `name` -
  (Required)
  Cluster name.  Changing this parameter will create a new resource.

* `security_group_id` -
  (Required)
  Security group ID.  Changing this parameter will create a new resource.

* `subnet_id` -
  (Required)
  Subnet ID.  Changing this parameter will create a new resource.

* `version` -
  (Required)
  Cluster version.  Changing this parameter will create a new resource.

* `vpc_id` -
  (Required)
  VPC ID.  Changing this parameter will create a new resource.

- - -

* `email` -
  (Optional)
  Notification email addresses. The max number is 5.  Changing this parameter will create a new resource.

* `enterprise_project_id` -
  (Optional)
  The enterprise project id.  Changing this parameter will create a new resource.

* `is_auto_off` -
  (Optional)
  Whether to automatically shut down.  Changing this parameter will create a new resource.

* `phone_num` -
  (Optional)
  Notification phone numbers. The max number is 5.  Changing this parameter will create a new resource.

* `schedule_boot_time` -
  (Optional)
  Timed boot time.  Changing this parameter will create a new resource.

* `schedule_off_time` -
  (Optional)
  Timed shutdown time.  Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` -
  Create time.

* `instances` -
  Instance list. Structure is documented below.

* `publid_ip` -
  Public ip.

The `instances` block contains:

* `id` -
  Instance ID.

* `name` -
  Instance name.

* `public_ip` -
  Public IP.

* `role` -
  Role.

* `traffic_ip` -
  Traffic IP.

* `type` -
  Instance type.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
