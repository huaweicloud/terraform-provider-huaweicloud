# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  region     = "cn-south-1"
  access_key = "IT1BRDH8SNDJN1SUED8O"
  secret_key = "iZAsJnE1IoGmsi1sjEBkJG2JTzHPRZd1MIkiaaqH"
}

# Create a VPC
resource "huaweicloud_vpc" "example" {
  name = "zhangting_vpc"
  cidr = "192.168.0.0/16"
}
