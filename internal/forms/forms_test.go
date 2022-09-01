package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields missing")
	}

	postedData := url.Values {}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	
	r, _ = http.NewRequest("POST", "whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required fields when it does")
	}	
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("Form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Shows form does not have when it should")
	}
}

func TestForm_MaxLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MaxLength("x", -4)
	if form.Valid() {
		t.Error("Form shows max length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedValues := url.Values {}
	postedValues.Add("some_field", "some value")
	form = New(postedValues)

	form.MaxLength("some_field", 4)
	if form.Valid() {
		t.Error("Shows max length of 4 met when data is longer")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abc")
	form = New(postedValues)

	form.MaxLength("another_field", 4)
	if !form.Valid() {
		t.Error("Shows max length of 4 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got valid for invalid email address")
	}
}