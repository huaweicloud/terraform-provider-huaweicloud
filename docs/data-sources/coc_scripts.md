---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_scripts"
description: |-
  Use this data source to get the list of COC scripts.
---

# huaweicloud_coc_scripts

Use this data source to get the list of COC scripts.

## Example Usage

```hcl
data "huaweicloud_coc_scripts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name_like` - (Optional, String) Specifies the fuzzy query script name.

* `creator` - (Optional, String) Specifies the creator.

* `risk_level` - (Optional, String) Specifies the risk level.
  Values can be as follows:
  + **LOW**: Low risk.
  + **MEDIUM**: Medium risk.
  + **HIGH**: High risk.

* `type` - (Optional, String) Specifies the script type.
  Values can be as follows:
  + **SHELL**: Shell script.
  + **PYTHON**: Python script.
  + **BAT**: Bat script.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of the scripts.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the script ID.

* `name` - Indicates the script name.

* `type` - Indicates the script type.
  Values can be as follows:
  + **SHELL**: Shell script.
  + **PYTHON**: Python script.
  + **BAT**: Bat script.

* `creator` - Indicates the creator.

* `creator_id` - Indicates the creator ID.

* `operator` - Indicates the modifier.

* `gmt_created` - Indicates the creation time.

* `gmt_modified` - Indicates the modification time.

* `status` - Indicates the script status.
  Values can be as follows:
  + **PENDING_APPROVE**: To be approved.
  + **APPROVED**: Normal (approved).
  + **REJECTED**: Rejected (indicates that the script is rejected by the approver).

* `script_uuid` - Indicates the script UUID.

* `usage_count` - Indicates the usage count.

* `properties` - Indicates the script label.

  The [properties](#data_properties_struct) structure is documented below.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `resource_tags` - Indicates the resource tags.

  The [resource_tags](#data_resource_tags_struct) structure is documented below.

<a name="data_properties_struct"></a>
The `properties` block supports:

* `risk_level` - Indicates the risk level.
  Values can be as follows:
  + **LOW**: Low risk.
  + **MEDIUM**: Medium risk.
  + **HIGH**: High risk.

* `reviewers` - Indicates the reviewers.

  The [reviewers](#properties_reviewers_struct) structure is documented below.

* `version` - Indicates the script version number.

* `protocol` - Indicates the approval message notification agreement, which is used to notify approvers.
  Values can be as follows:
  + **DEFAULT**: Default value.
  + **SMS**: Sms.
  + **EMAIL**: Email.
  + **DING_TALK**: DingTalk.
  + **WE_LINK**: WeLink.
  + **WECHAT**: WeChat.
  + **CALLNOTIFY**: Voice.
  + **NOT_TO_NOTIFY**: Not notify.

<a name="properties_reviewers_struct"></a>
The `reviewers` block supports:

* `reviewer_id` - Indicates the reviewer ID (IAM user ID).

* `reviewer_name` - Indicates the reviewer name (IAM username).

<a name="data_resource_tags_struct"></a>
The `resource_tags` block supports:

* `key` - Indicates the resource tag key.

* `value` - Indicates the resource tag value.
