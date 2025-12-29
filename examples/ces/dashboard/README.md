# Create a CES dashboard example

This example provides best practice code for using Terraform to create a dashboard in HuaweiCloud CES service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The name of the CES dashboard.
* `row_widget_num` - The monitoring view display mode.

#### Optional Variables

* `dashboard_id` - The copied dashboard ID.
* `enterprise_project_id` - The enterprise project ID of the dashboard.
* `is_favorite` - Whether the dashboard is favorite.
* `extend_info` - The information about the extension.
  + `filter` - The metric aggregation method.
  + `period` - The metric aggregation period.
  + `display_time` - The display time.
  + `refresh_time` - The refresh time.
  + `from` - The start time.
  + `to` - The end time.
  + `screen_color` - The monitoring screen background color.
  + `enable_screen_auto_play` - Whether the monitoring screen switches automatically.
  + `time_interval` - The automatic switching time interval of the monitoring screen.
  + `enable_legend` - Whether to enable the legend.
  + `full_screen_widget_num` - The number of large screen display views.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  name           = "tf_test_ces_dashboard_name"
  row_widget_num = "tf_test_ces_dashboard_row_widget_num"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.77.1 |
