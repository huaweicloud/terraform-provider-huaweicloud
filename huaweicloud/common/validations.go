package common

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func StandardVerify(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^\\w*$"),
			"The value can only consist of letters, digits and underscores (_)."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphens(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\\w-]*$"),
			"The value can only consist of letters, digits, underscores (_) and hyphens (-)."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensAndDots(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\\w-.]*$"),
			"The value can only consist of letters, digits, underscores (_), hyphens (-) and dots (.)."),
		validation.StringLenBetween(min, max),
	)
}

func SimpleVerify(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9]*$"),
			"The value can only consist of letters and digits."),
		validation.StringLenBetween(min, max),
	)
}

func SimpleVerifyWithHyphens(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9-]*$"),
			"The value can only consist of letters, digits and hyphens (-)."),
		validation.StringLenBetween(min, max),
	)
}

func LowercaseVerify(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[a-z0-9]*$"),
			"The value can only consist of lowercase letters and digits."),
		validation.StringLenBetween(min, max),
	)
}

func LowercaseVerifyWithUnderscores(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[a-z0-9_]*$"),
			"The value can only consist of lowercase letters, digits and underscores (_)."),
		validation.StringLenBetween(min, max),
	)
}

func LowercaseVerifyWithHyphensAndStartEnd(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`),
			"The value can only consist of lowercase letters, digits and hyphens (-), "+
				"and it must start and end with a letter or digit."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensAndChineses(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-]*$"),
			"The value can only consist of letters, digits, underscores (_), hyphens (-) and chinese characters."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensDotsAndChineses(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-.]*$"), "The value can only consist of letters, "+
			"digits, underscores (_), hyphens (-), dots (.) and chinese characters."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithStart(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z]\\w*$"),
			"The value can only consist of letters, digits and underscores (_), and it must start with a letter."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensAndStart(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z][\\w-]*$"), "The value can only consist of letters, "+
			"digits, underscores (_) and hyphens (-), and it must start with a letter."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensDotsAndStart(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z][\\w-.]*?$"),
			"The value can only consist of letters, digits, underscores (_), hyphens (-) and dots (.),"+
				"and it must start with a letter."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithChinesesAndStart(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5\\w]*?$"),
			"The value can only consist of letters, digits, underscores (_) and chinese characters"+
				"and it must start with a letter or chinese character."),
		validation.StringLenBetween(min, max),
	)
}

func StandardVerifyWithHyphensChinesesAndStart(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5\\w-]*?$"),
			"The value can only consist of letters, digits, underscores (_), hyphens (-) and chinese characters"+
				"and it must start with a letter or chinese character."),
		validation.StringLenBetween(min, max),
	)
}

func CustomVerify(regex, message string, min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(regex), message),
		validation.StringLenBetween(min, max),
	)
}

func StringVerifyWithoutAngleBrackets(min, max int) schema.SchemaValidateFunc {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
			"The angle brackets (< and >) are not allowed for this value."),
		validation.StringLenBetween(min, max),
	)
}
