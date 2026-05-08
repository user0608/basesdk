package answer

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"basesdk/errs"
)

type fakeTarget struct {
	code            int
	body            any
	jsonCalled      bool
	noContentCode   int
	noContentCalled bool
}

func (f *fakeTarget) JSON(code int, i any) error {
	f.code = code
	f.body = i
	f.jsonCalled = true
	return nil
}

func (f *fakeTarget) NoContent(code int) error {
	f.noContentCode = code
	f.noContentCalled = true
	return nil
}

func assertResponse(t *testing.T, c *fakeTarget, code int, message string, data any) {
	t.Helper()

	if !c.jsonCalled {
		t.Fatal("expected JSON to be called")
	}

	if c.code != code {
		t.Fatalf("expected status code %d, got %d", code, c.code)
	}

	response, ok := c.body.(*Response)
	if !ok {
		t.Fatalf("expected body type *Response, got %T", c.body)
	}

	if response.Message != message {
		t.Fatalf("expected message %q, got %q", message, response.Message)
	}

	if !reflect.DeepEqual(response.Data, data) {
		t.Fatalf("expected data %#v, got %#v", data, response.Data)
	}
}

func TestOk(t *testing.T) {
	c := &fakeTarget{}
	payload := map[string]string{"id": "123"}

	err := Ok(c, payload)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusOK, "", payload)
}

func TestCreated(t *testing.T) {
	c := &fakeTarget{}
	payload := map[string]string{"id": "123"}

	err := Created(c, payload)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusCreated, "", payload)
}

func TestAccepted(t *testing.T) {
	c := &fakeTarget{}
	payload := map[string]string{"status": "queued"}

	err := Accepted(c, payload)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusAccepted, "", payload)
}

func TestMessage(t *testing.T) {
	c := &fakeTarget{}

	err := Message(c, "custom message")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusOK, "custom message", nil)
}

func TestCreatedMessage(t *testing.T) {
	c := &fakeTarget{}

	err := CreatedMessage(c, "created successfully")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusCreated, "created successfully", nil)
}

func TestAcceptedMessage(t *testing.T) {
	c := &fakeTarget{}

	err := AcceptedMessage(c, "accepted successfully")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusAccepted, "accepted successfully", nil)
}

func TestSuccess(t *testing.T) {
	c := &fakeTarget{}

	err := Success(c)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusOK, "Operación completada exitosamente", nil)
}

func TestCreatedSuccess(t *testing.T) {
	c := &fakeTarget{}

	err := CreatedSuccess(c)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusCreated, "Recurso creado exitosamente", nil)
}

func TestAcceptedSuccess(t *testing.T) {
	c := &fakeTarget{}

	err := AcceptedSuccess(c)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusAccepted, "Operación aceptada exitosamente", nil)
}

func TestNoContent(t *testing.T) {
	c := &fakeTarget{}

	err := NoContent(c)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !c.noContentCalled {
		t.Fatal("expected NoContent to be called")
	}

	if c.noContentCode != http.StatusNoContent {
		t.Fatalf("expected status code %d, got %d", http.StatusNoContent, c.noContentCode)
	}

	if c.jsonCalled {
		t.Fatal("expected JSON not to be called")
	}
}

func TestUnwrapErrWithDomainError(t *testing.T) {
	err := errs.BadRequestError(nil, "Bad request: %s", "missing field")

	code, message := UnwrapErr(err)

	if code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, code)
	}

	if message != "Bad request: missing field" {
		t.Fatalf("expected message %q, got %q", "Bad request: missing field", message)
	}
}

func TestUnwrapErrWithWrappedDomainError(t *testing.T) {
	baseErr := fmt.Errorf("database failed")
	err := errs.WrapError(baseErr, "Unexpected failure", http.StatusInternalServerError)

	code, message := UnwrapErr(err)

	if code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, code)
	}

	if message != "Unexpected failure" {
		t.Fatalf("expected message %q, got %q", "Unexpected failure", message)
	}
}

func TestUnwrapErrWithPrefixedBadRequest(t *testing.T) {
	code, message := UnwrapErr(errors.New(":invalid request"))

	if code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, code)
	}

	if message != "invalid request" {
		t.Fatalf("expected message %q, got %q", "invalid request", message)
	}
}

func TestUnwrapErrWithInternalError(t *testing.T) {
	code, message := UnwrapErr(errors.New("database failed"))

	if code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, code)
	}

	expected := "Ocurrió un problema. Se produjo un error inesperado."
	if message != expected {
		t.Fatalf("expected message %q, got %q", expected, message)
	}
}

func TestErrWithDomainError(t *testing.T) {
	c := &fakeTarget{}
	err := errs.BadRequestError(nil, "Invalid value: %s", "email")

	got := Err(c, err)

	if got != nil {
		t.Fatalf("expected nil error, got %v", got)
	}

	assertResponse(t, c, http.StatusBadRequest, "Invalid value: email", nil)
}

func TestErrWithPrefixedBadRequest(t *testing.T) {
	c := &fakeTarget{}

	err := Err(c, errors.New(":invalid request"))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusBadRequest, "invalid request", nil)
}

func TestAutoWithNilError(t *testing.T) {
	c := &fakeTarget{}

	err := Auto(c, nil)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusOK, "Operación completada exitosamente", nil)
}

func TestAutoWithDomainError(t *testing.T) {
	c := &fakeTarget{}
	domainErr := errs.BadRequestError(nil, "Invalid payload")

	err := Auto(c, domainErr)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusBadRequest, "Invalid payload", nil)
}

func TestAutoWithPrefixedBadRequest(t *testing.T) {
	c := &fakeTarget{}

	err := Auto(c, errors.New(":invalid request"))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertResponse(t, c, http.StatusBadRequest, "invalid request", nil)
}
