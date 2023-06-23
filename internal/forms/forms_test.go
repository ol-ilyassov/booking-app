package forms

import (
	"net/url"
	"testing"
)

// Name convention: TestForm_Valid
// > Form - struct
// > Valid - receiver function

func TestForm_Valid(t *testing.T) {
	// r := httptest.NewRequest("POST", "/some-url", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	// #1
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	// #2
	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")
	form = New(postedData)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}

}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("field1")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("field1", "value1")
	form = New(postedData)

	has = form.Has("field1")
	if !has {
		t.Error("form shows form does not have a field when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	// #1
	postedValues := url.Values{}
	form := New(postedValues)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	// #2.1
	postedValues = url.Values{}
	postedValues.Add("some_field", "some_value")
	form = New(postedValues)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("form shows min length of 100 met when data is shorter")
	}

	// #2.2
	isError := form.Errors.Get("some_field")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// #3.1
	postedValues = url.Values{}
	postedValues.Add("another_field", "abc123")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("form shows min length of 1 is not met when it is")
	}

	// #3.2
	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_Email(t *testing.T) {
	// #1
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	// #2
	postedValues = url.Values{}
	postedValues.Add("email", "test.com")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}

	// #3
	postedValues = url.Values{}
	postedValues.Add("email", "test@gmail.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when it should have")
	}
}
