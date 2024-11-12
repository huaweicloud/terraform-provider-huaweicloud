---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_subscription"
description: ""
---

# huaweicloud_eg_event_subscription

Using this resource to manage an EG event subscription within Huaweicloud.

## Example Usage

### Event subscription created based on official channel (OBS to FGS)

```hcl
variable "official_channel_id" {}
variable "subscription_name" {}
variable "custom_id_for_obs_source" {} # Custom ID that used to prevent ID change, in UUID format
variable "source_obs_bucket_name" {}
variable "custom_id_for_fgs_target" {} # Custom ID that used to prevent ID change, in UUID format
variable "target_function_urn" {}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id = var.official_channel_id
  name       = var.subscription_name

  sources {
    id            = var.custom_id_for_obs_source
    name          = "HC.OBS.DWR"
    provider_type = "OFFICIAL"
    detail_name   = "detail"
    detail        = jsonencode({
      "bucket": var.source_obs_bucket_name,
      "objectKeyEncode": true
    })
    filter_rule  = jsonencode({
      "source": [{
        "op": "StringIn",
        "values": ["HC.OBS.DWR"]
      }],
      "subject": {
        "and": [
          {
            "op": "StringStartsWith",
            "values": [
              "test"
            ]
          },
          {
            "op": "StringEndsWith",
            "values": [
              "jpg"
            ]
          }
        ]
      },
      "type": [{
        "op": "StringIn",
        "values": [
          "OBS:DWR:ObjectCreated:PUT",
          "OBS:DWR:ObjectCreated:POST",
          "OBS:DWR:ObjectCreated:COPY",
          "OBS:DWR:ObjectCreated:CompleteMultipartUpload",
          "OBS:DWR:ObjectRemoved:Delete",
          "OBS:DWR:ObjectRemoved:DeleteMarkerCreated"
        ]
      }],
      "data": {
        "obs": {
          "bucket": {
            "name": [{
              "op": "StringIn",
              "values": [var.source_obs_bucket_name]
            }]
          }
        }
      }
    })
  }

  targets {
    id            = var.custom_id_for_fgs_target
    name          = "HC.FunctionGraph"
    provider_type = "OFFICIAL"
    detail_name   = "detail"
    detail        = jsonencode({
      "urn": var.target_function_urn,
      "invoke_type": "ASYNC",
      "agency_name": "EG_TARGET_AGENCY"
    })
    transform     = jsonencode({
      "type": "VARIABLE",
      "value": "{\"name\":\"$.data.name\"}",
      "template": "My name is $${name}."
    })
  }
}
```

### Event subscription created based on custom channel (custom source to EG and SMN)

```hcl
variable "subscription_name" {}
variable "source_channel_id" {}
variable "target_channel_id" {}
variable "custom_id_for_custom_source" {} # Custom ID that used to prevent ID change, in UUID format
variable "custom_event_source_name" {}
variable "custom_id_for_eg_target" {}     # Custom ID that used to prevent ID change, in UUID format
variable "custom_id_for_smn_target" {}    # Custom ID that used to prevent ID change, in UUID format
variable "project_id" {}
variable "region_name" {}
variable "smn_topic_urn" {}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id = var.source_channel_id
  name       = var.subscription_name

  sources {
    id            = var.custom_id_for_custom_source
    provider_type = "CUSTOM"
    name          = var.custom_event_source_name
    filter_rule   = jsonencode({
      "source": [{
        "op":"StringIn",
        "values":["${var.custom_event_source_name}"]
      }]
    })
  }

  targets {
    id            = var.custom_id_for_eg_target
    provider_type = "OFFICIAL"
    name          = "HC.EG"
    detail_name   = "eg_detail"
    detail        = jsonencode({
      "agency_name": "EG_TARGET_AGENCY",
      "cross_account": false,
      "cross_region": false,
      "target_channel_id": "${var.target_channel_id}",
      "target_project_id": "${var.project_id}",
      "target_region": "${var.region_name}"
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }

  targets {
    id            = var.custom_id_for_smn_target
    provider_type = "OFFICIAL"
    name          = "HC.SMN"
    detail_name   = "smn_detail"
    detail        = jsonencode({
      "subject_transform": {
        "type": "CONSTANT",
        "value": "TEST_CONDTANT"
      },
      "urn": "${var.smn_topic_urn}",
      "agency_name": "EG_TARGET_AGENCY",
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event subscription is located.  
  If omitted, the provider-level region will be used.
  Changing this will create a new resource.

* `channel_id` - (Required, String, ForceNew) Specifies the channel ID to which the event subscription belongs.
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the event subscription.  
  The valid length is limited from `1` to `128`, only letters, digits, hyphens (-), underscores (_) and dots (.) are
  allowed. The name must start with a letter or digit. Changing this will create a new resource.

* `sources` - (Required, List) Specifies the list of the event sources.
  The [sources](#subscription_sources) structure is documented below.

* `targets` - (Required, List) Specifies the list of the event targets.
  The [targets](#subscription_targets) structure is documented below.

* `description` - (Optional, String) Specifies the description of the event subscription.

<a name="subscription_sources"></a>
The `sources` block supports:

* `id` - (Required, String) Specifies the custom ID of the event source, in UUID format.

  -> This `id` field is only used for internal management of event subscription resource and has no association with the
     parameters or attributes of other resources.

* `provider_type` - (Required, String) Specifies the provider type of the event source.
  The valid values are as follows:
  + **CUSTOM**
  + **OFFICIAL**

* `name` - (Required, String) Specifies the name of the event source.
  The valid length is limited from `1` to `128`.

* `detail_name` - (Optional, String) Specifies the name (key) of the target detail configuration.
  The valid values are as follows:
  + **detail**: Used by some official event sources, such as OBS.

* `detail` - (Optional, String) Specifies the configuration detail of the event source, in JSON format.
  The valid length is limited from `1` to `1,024`.

* `filter_rule` - (Required, String) Specifies the filter rule of the event source, in JSON format.
  The valid length is limited from `1` to `2,048`.

-> The JSON result (`detail` and `filter_rule`) returned by the EG service may be different from the definition of the
   resource created (especially the official event source). Changes can be handled by `lifecycle.ignore_changes` or
   manual synchronization.

<a name="subscription_targets"></a>
The `targets` block supports:

* `id` - (Required, String) Specifies the custom ID of the event target, in UUID format.

* `provider_type` - (Required, String) Specifies the provider type of the event target.
  The valid values are as follows:
  + **CUSTOM**
  + **OFFICIAL**

* `name` - (Required, String) Specifies the name of the event target.
  The valid length is limited from `1` to `128`.

* `detail_name` - (Required, String) Specifies the name (key) of the target detail configuration.
  The valid values are as follows:
  + **detail**: Custom event targets and FunctionGraph event targets are used.
  + **smn_detail**: SMN event targets are used.
  + **kafka_detail**: DMS kafka event targets are used.
  + **eg_detail**: EG event targets are used.

* `detail` - (Required, String) Specifies the configuration detail of the event target, in JSON format.
  The valid length is limited from `1` to `1,024`.

* `transform` - (Required, String) Specifies the transform configuration of the event target, in JSON format.

* `connection_id` - (Optional, String) Specifies the connection ID of the EG event target.

* `dead_letter_queue` - (Optional, String) Specifies the specified queue to which failure events sent, in JSON format.

-> The JSON result (`detail`, `transform`and `dead_letter_queue`) returned by the EG service may be different from the
   definition of the resource created (especially the official event target). Changes can be handled by
   `lifecycle.ignore_changes` or manual synchronization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - The status of the event subscription.

* `sources` - The list of the event sources.
  The [sources](#subscription_sources_attr) structure is documented below.

* `targets` - The list of the event targets.
  The [targets](#subscription_targets_attr) structure is documented below.

* `created_at` - The (UTC) creation time of the event subscription, in RFC3339 format.

* `updated_at` - The (UTC) update time of the event subscription, in RFC3339 format.

<a name="subscription_sources_attr"></a>
The `sources` block supports:

* `created_at` - The (UTC) creation time of the event source, in RFC3339 format.

* `updated_at` - The (UTC) update time of the event source, in RFC3339 format.

<a name="subscription_targets_attr"></a>
The `targets` block supports:

* `created_at` - The (UTC) creation time of the event target, in RFC3339 format.

* `updated_at` - The (UTC) update time of the event target, in RFC3339 format.

## Import

Subscriptions can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_eg_event_subscription.test <id>
```
