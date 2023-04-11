package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initalizes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MaLength checks for string maximum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("Minimum %d characters", length))
		return false
	}
	return true
}

// MaLength checks for string maximum length
func (f *Form) MaxLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) > length {
		f.Errors.Add(field, fmt.Sprintf("Maximum %d characters", length))
		return false
	}
	return true
}

func (f *Form) IsDuplicate(field string, field2 string, msg string) bool {
	x := f.Get(field)
	y := f.Get(field2)
	if x != y {
		f.Errors.Add(field, msg)
		f.Errors.Add(field2, msg)
		return false
	}
	return true
}

func (f *Form) AreDifferent(field string, field2 string, msg string) bool {
	x := f.Get(field)
	y := f.Get(field2)
	if x == y {
		f.Errors.Add(field, msg)
		f.Errors.Add(field2, msg)
		return false
	}
	return true
}

// IsEmail checks for valid e-mail address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid e-mail address")
	}
}

// ValidPassword checks for valid password
func (f *Form) ValidPassword(field string) {
	pw := f.Get(field)
	if len(pw) < 7 {
		f.Errors.Add(field, "1 uppercase, 1 lowercase, 1 number, minimum length 7 characters")
		return
	}
	if !govalidator.HasUpperCase(pw) && !govalidator.HasLowerCase(pw) {
		f.Errors.Add(field, "1 uppercase, 1 lowercase, 1 number, minimum length 7 characters")
		return
	}
	for _, char := range pw {
		if govalidator.IsNumeric(string(char)) {
			return
		}
	}
	f.Errors.Add(field, "1 uppercase, 1 lowercase, 1 number, minimum length 7 characters")
}

// Alphanumeric checks for alphanumeric and spaces
func (f *Form) AlphaNumeric(fields ...string) {
	regexp, _ := regexp.Compile(`^[a-zA-Z0-9õäöüÕÄÖÜ ]*$`)
	for _, field := range fields {
		value := f.Get(field)
		match := regexp.MatchString(value)
		if !match {
			f.Errors.Add(field, "Alphanumeric values only")
		}
	}
}

// IsUpper checks for uppercase letters and numbers
func (f *Form) IsUpper(field string) {
	match, _ := regexp.MatchString(`^[A-Z0-9 ]*$`, f.Get(field))
	if !match {
		f.Errors.Add(field, "Only uppercase letters and numbers allowed")
	}
}

func ValidExtension(fileName string, extensions ...string) bool {
	for _, extension := range extensions {
		regex := fmt.Sprintf("\\.%s", extension)
		re, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
		if re.MatchString(fileName) {
			return true
		}
	}
	return false
}
