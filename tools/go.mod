module github.com/terraform-linters/tflint/tools

go 1.13

require (
	github.com/hashicorp/hcl/v2 v2.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.2.0
	github.com/serenize/snaker v0.0.0-20171204205717-a683aaf2d516
	github.com/terraform-providers/terraform-provider-aws v2.35.0+incompatible
)

replace github.com/terraform-providers/terraform-provider-aws v2.35.0+incompatible => github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191107200508-571f1b4da739
