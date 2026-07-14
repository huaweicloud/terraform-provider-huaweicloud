---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publication_subscription_profiles"
description: |-
  Use this data source to query the list of RDS publication and subscription profiles within HuaweiCloud.
---

# huaweicloud_rds_publication_subscription_profiles

Use this data source to query the list of RDS publication and subscription profiles within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_publication_subscription_profiles" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the replication profiles.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `agent_type` - (Optional, String) Specifies the agent type to filter profiles.  
  The valid values are as follows:
  + **snapshot**: Snapshot agent.
  + **log_reader**: Log reader agent.
  + **distribution**: Distribution agent.
  + **merge**: Merge agent.
  + **queue_reader**: Queue reader agent.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `profiles` - The list of replication profiles.  
  The [profiles](#replication_profiles_profiles) structure is documented below.

<a name="replication_profiles_profiles"></a>
The `profiles` block supports:

* `profile_id` - The ID of the profile.

* `profile_name` - The name of the profile.

* `agent_type` - The agent type of the profile.  
  The valid values are as follows:
  + **snapshot**: Snapshot agent.
  + **log_reader**: Log reader agent.
  + **distribution**: Distribution agent.
  + **merge**: Merge agent.
  + **queue_reader**: Queue reader agent.

* `description` - The description of the profile.

* `is_def_profile` - Whether the profile is the default configuration.
