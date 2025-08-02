package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"

	"github.com/pachirode/iam_study/pkg/validation/field"
)

const (
	qnameCharFmt        string = "[A-Za-z0-9]"
	qnameExtCharFmt     string = "[-A-Za-z0-9_.]"
	qualifiedNameFmt           = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	labelValueFmt              = "(" + qualifiedNameFmt + ")"
	dns1123LabelFmt     string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	dns1123SubdomainFmt        = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"
	percentFmt          string = "[0-9]+%"
)

const (
	qualifiedNameErrMsg    string = "must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
	labelValueErrMsg       string = "a valid label must be empty or consist of alphanumeric characters, '-', '_' or '.', and must start and end with alphanumeric"
	dns1123LabelErrMsg     string = "a DNS-1123 label must consist of lower alphanumeric characters or '-', and must start and end with alphanumeric"
	dns1123SubdomainErrMsg string = "a DNS-1123 subdomain must consist of lower alphanumeric characters, '-', or '.', and must start and end with alphanumeric"
	percentErrMsg          string = "a valid percent must be a numeric string followed by an ending: '%'"

	qualifiedNameMaxLength    int = 63
	LabelValueMaxLength       int = 63
	DNS1123LabelMaxLength     int = 63
	DNS1123SubdomainMaxLength int = 253

	minPassLength = 8
	maxPassLength = 16
)

var (
	qualifiedNameRegexp    = regexp.MustCompile("^" + qualifiedNameFmt + "$")
	labelValueRegexp       = regexp.MustCompile("^" + labelValueFmt + "$")
	dns1123LabelRegexp     = regexp.MustCompile("^" + dns1123LabelFmt + "$")
	dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")
	percentRegexp          = regexp.MustCompile("^" + percentFmt + "$")
)

func IsvalidateValue(value string) []string {
	var errs []string

	if len(value) > LabelValueMaxLength {
		errs = append(errs, MaxLenError(LabelValueMaxLength))
	}
	if !labelValueRegexp.MatchString(value) {
		errs = append(errs, RegexError(labelValueErrMsg, labelValueFmt, "MyValue", "my_value", "12345"))
	}
	return errs
}

func IsDNS1123Label(value string) []string {
	var errs []string
	if len(value) > DNS1123LabelMaxLength {
		errs = append(errs, MaxLenError(DNS1123LabelMaxLength))
	}
	if !dns1123LabelRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123LabelErrMsg, dns1123LabelFmt, "my-name", "123-abc"))
	}
	return errs
}

func IsDNS1123Subdomain(value string) []string {
	var errs []string
	if len(value) > DNS1123SubdomainMaxLength {
		errs = append(errs, MaxLenError(DNS1123SubdomainMaxLength))
	}
	if !dns1123SubdomainRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123SubdomainErrMsg, dns1123SubdomainFmt, "example.com"))
	}
	return errs
}

func IsValidPortNum(port int) []string {
	if 1 <= port && port <= 65535 {
		return nil
	}
	return []string{InclusiveRangeError(1, 65535)}
}

func IsInRange(value int, min int, max int) []string {
	if value >= min && value <= max {
		return nil
	}
	return []string{InclusiveRangeError(min, max)}
}

func IsValidIP(value string) []string {
	if net.ParseIP(value) == nil {
		return []string{"must be a valid IP address, (e.g. 10.9.8.7)"}
	}
	return nil
}

func IsValidIPv4Address(fldPath *field.Path, value string) field.ErrorList {
	var allErrors field.ErrorList
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		allErrors = append(allErrors, field.Invalid(fldPath, value, "must be a valid IPv4 address"))
	}
	return allErrors
}

func IsValidIPv6Address(fldPath *field.Path, value string) field.ErrorList {
	var allErrors field.ErrorList
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil {
		allErrors = append(allErrors, field.Invalid(fldPath, value, "must be a valid IPv6 address"))
	}
	return allErrors
}

func IsValidPercent(percent string) []string {
	if !percentRegexp.MatchString(percent) {
		return []string{RegexError(percentErrMsg, percentFmt, "1%", "93%")}
	}
	return nil
}

func IsValidPassword(password string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
			passLen++
		case unicode.IsUpper(ch):
			hasUpper = true
			passLen++
		case unicode.IsLower(ch):
			hasLower = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !hasLower {
		appendError("lowercase letter missing")
	}
	if !hasUpper {
		appendError("uppercase letter missing")
	}
	if !hasNumber {
		appendError("at least one numeric character required")
	}
	if !hasSpecial {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(
			fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength),
		)
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}

	return nil
}

func IsQualifiedName(value string) []string {
	var errs []string
	parts := strings.Split(value, "/")
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	// nolint:gomnd // no need
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			errs = append(errs, "prefix part "+EmptyError())
		} else if msgs := IsDNS1123Subdomain(prefix); len(msgs) != 0 {
			errs = append(errs, prefixEach(msgs, "prefix part ")...)
		}
	default:
		return append(
			errs,
			"a qualified name "+RegexError(
				qualifiedNameErrMsg,
				qualifiedNameFmt,
				"MyName",
				"my.name",
				"123-abc",
			)+" with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')",
		)
	}

	if len(name) == 0 {
		errs = append(errs, "name part "+EmptyError())
	} else if len(name) > qualifiedNameMaxLength {
		errs = append(errs, "name part "+MaxLenError(qualifiedNameMaxLength))
	}
	if !qualifiedNameRegexp.MatchString(name) {
		errs = append(
			errs,
			"name part "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"),
		)
	}
	return errs
}
