# Create an IoTDA device access environment

This example provides best practice code for using Terraform to create an IoTDA (Internet of Things Device Access)
device access environment in HuaweiCloud. It includes creating a resource space, a product model, and a device.

## Prerequisites

* A HuaweiCloud account with IoTDA service activated
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IoTDA resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `space_name` - The IoTDA resource space name
* `product_name` - The IoTDA product name
* `device_node_id` - The node ID of the device
* `device_name` - The device name

#### Optional Variables

* `iotda_endpoint` - The IoTDA service endpoint for standard/enterprise edition instances
* `product_device_type` - The device type of the product (default: "Thermometer")
* `product_protocol` - The protocol used by the product (default: "MQTT")
* `product_data_type` - The data type of the product (default: "json")
* `product_service_id` - The service ID of the product (default: "service_1")
* `product_service_type` - The service type of the product (default: "serv_type")
* `device_secret` - The secret of the device for authentication (default: "1234567890")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  space_name         = "your_space_name"
  product_name       = "your_product_name"
  product_service_id = "service_1"
  device_node_id     = "your_device_node_id"
  device_name        = "your_device_name"
  device_secret      = "1234567890"
  ```

* If you are using an IoTDA standard or enterprise edition instance, specify the endpoint:

  ```hcl
  iotda_endpoint = "your_instance_id.iotda-app.cn-north-4.myhuaweicloud.com"
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
* The device access workflow requires the resources to be created in order: Space -> Product -> Device
* For standard/enterprise edition IoTDA instances, you must specify the `iotda_endpoint` variable

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.38.0 |
