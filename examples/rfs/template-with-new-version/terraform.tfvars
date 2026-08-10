template_name = "tf-test-template"
template_body = <<EOF
resource "huaweicloud_vpc" "test" {
  name = "tf-test-vpc"
  cidr = "172.16.0.0/16"
}
EOF

template_version_body = <<EOF
resource "huaweicloud_vpc" "test" {
  name = "tf-test-vpc"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "tf-test-subnet"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "172.16.1.0/24"
  gateway_ip = "172.16.1.1"
}
EOF
