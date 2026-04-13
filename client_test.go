package iwmo

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// errBodyTransport returns a 200 OK response whose body yields an error on
// Read. This exercises the io.ReadAll error path in httpSender.Send.
type errBodyTransport struct{}

func (t *errBodyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&errBodyReader{}),
		Header:     make(http.Header),
	}, nil
}

type errBodyReader struct{}

func (e *errBodyReader) Read([]byte) (int, error) {
	return 0, errors.New("intentional read error")
}

// mockSender is a Sender that delegates to a function, for unit testing.
type mockSender struct {
	fn func(ctx context.Context, data []byte) ([]byte, error)
}

func (m *mockSender) Send(ctx context.Context, data []byte) ([]byte, error) {
	return m.fn(ctx, data)
}

func TestNewClient_RequiresBaseURLOrSender(t *testing.T) {
	_, err := NewClient()
	if err == nil {
		t.Fatal("NewClient() error = nil, want error when neither baseURL nor Sender given")
	}
	if !errors.Is(err, ErrConfiguration) {
		t.Errorf("errors.Is(err, ErrConfiguration) = false, got: %v", err)
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	c, err := NewClient(WithBaseURL("https://example.nl/iwmo"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClient_WithSender(t *testing.T) {
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) { return nil, nil }}
	c, err := NewClient(WithSender(ms))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClient_OptionsApplied(t *testing.T) {
	c, err := NewClient(
		WithBaseURL("https://example.nl"),
		WithAGBCode("12345678"),
		WithGemeenteCode("0363"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.agbCode != "12345678" {
		t.Errorf("agbCode = %q, want 12345678", c.agbCode)
	}
	if c.gemeenteCode != "0363" {
		t.Errorf("gemeenteCode = %q, want 0363", c.gemeenteCode)
	}
}

func TestClient_SendVerzoekToewijzing_ValidatesFirst(t *testing.T) {
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		t.Error("Send should not be called when message is invalid")
		return nil, nil
	}}
	c, err := NewClient(WithSender(ms))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	bad := &WMO302{} // empty, will fail Validate()
	_, err = c.SendVerzoekToewijzing(context.Background(), bad)
	if err == nil {
		t.Fatal("SendVerzoekToewijzing() error = nil for invalid message")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestClient_SendVerzoekToewijzing_Success(t *testing.T) {
	retour, encErr := Encode(validWMO304())
	if encErr != nil {
		t.Fatalf("Encode(validWMO304()) error = %v", encErr)
	}

	ms := &mockSender{fn: func(_ context.Context, data []byte) ([]byte, error) {
		return retour, nil
	}}
	c, err := NewClient(WithSender(ms))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err != nil {
		t.Fatalf("SendVerzoekToewijzing() error = %v", err)
	}
	if resp == nil {
		t.Fatal("SendVerzoekToewijzing() returned nil response")
	}
	if resp.Header.BerichtCode != "304" {
		t.Errorf("response BerichtCode = %q, want 304", resp.Header.BerichtCode)
	}
}

func TestClient_SendDeclaratie_Success(t *testing.T) {
	retour, _ := Encode(validWMO304())
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return retour, nil
	}}
	c, _ := NewClient(WithSender(ms))
	resp, err := c.SendDeclaratie(context.Background(), validWMO303())
	if err != nil {
		t.Fatalf("SendDeclaratie() error = %v", err)
	}
	if resp == nil {
		t.Fatal("SendDeclaratie() returned nil response")
	}
}

func TestClient_SendMutatie_Success(t *testing.T) {
	retour, _ := Encode(validWMO304())
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return retour, nil
	}}
	c, _ := NewClient(WithSender(ms))
	resp, err := c.SendMutatie(context.Background(), validWMO305())
	if err != nil {
		t.Fatalf("SendMutatie() error = %v", err)
	}
	if resp == nil {
		t.Fatal("SendMutatie() returned nil response")
	}
}

func TestClient_SendStatusmelding_Success(t *testing.T) {
	retour, _ := Encode(validWMO304())
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return retour, nil
	}}
	c, _ := NewClient(WithSender(ms))
	resp, err := c.SendStatusmelding(context.Background(), validWMO315())
	if err != nil {
		t.Fatalf("SendStatusmelding() error = %v", err)
	}
	if resp == nil {
		t.Fatal("SendStatusmelding() returned nil response")
	}
}

func TestClient_SendToewijzing_Success(t *testing.T) {
	retour, _ := Encode(validWMO304())
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return retour, nil
	}}
	c, _ := NewClient(WithSender(ms))
	resp, err := c.SendToewijzing(context.Background(), validWMO301())
	if err != nil {
		t.Fatalf("SendToewijzing() error = %v", err)
	}
	if resp == nil {
		t.Fatal("SendToewijzing() returned nil response")
	}
}

func TestClient_TransportError(t *testing.T) {
	transportErr := errors.New("connection refused")
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, transportErr
	}}
	c, _ := NewClient(WithSender(ms))
	_, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("SendVerzoekToewijzing() error = nil, want transport error")
	}
}

func TestHTTPSender_Post(t *testing.T) {
	retour, _ := Encode(validWMO304())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("HTTP method = %q, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/xml; charset=utf-8" {
			t.Errorf("Content-Type = %q, want application/xml; charset=utf-8", ct)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write(retour)
	}))
	defer ts.Close()

	c, err := NewClient(WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	resp, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err != nil {
		t.Fatalf("SendVerzoekToewijzing() error = %v", err)
	}
	if resp.RetourCodes[0].Code != "0000" {
		t.Errorf("RetourCode = %q, want 0000", resp.RetourCodes[0].Code)
	}
}

func TestHTTPSender_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewClient(WithBaseURL(ts.URL))
	_, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("error = nil, want ErrAuthentication")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Errorf("errors.Is(err, ErrAuthentication) = false, got: %v", err)
	}
}

func TestHTTPSender_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewClient(WithBaseURL(ts.URL))
	_, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("error = nil, want ErrTransport")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c, err := NewClient(WithBaseURL("https://example.nl"), WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.httpClient != custom {
		t.Error("httpClient was not replaced by WithHTTPClient")
	}
}

func TestWithLogger(t *testing.T) {
	logger := slog.Default()
	c, err := NewClient(WithBaseURL("https://example.nl"), WithLogger(logger))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.logger != logger {
		t.Error("logger was not replaced by WithLogger")
	}
}

func TestClient_InvalidRetour_ReturnsErrTransport(t *testing.T) {
	// Server returns a WMO304 that fails validation (no RetourCodes).
	invalidRetour := &WMO304{
		Header: WMO304Header{Header: validHeader("304")},
	}
	data, err := Encode(invalidRetour)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return data, nil
	}}
	c, _ := NewClient(WithSender(ms))
	_, err = c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("error = nil, want ErrTransport for invalid WMO304 retour")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestClient_GarbageRetour_ReturnsErrTransport(t *testing.T) {
	// Server returns bytes that cannot be decoded as WMO304 XML.
	ms := &mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		return []byte("not xml at all"), nil
	}}
	c, _ := NewClient(WithSender(ms))
	_, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("error = nil, want ErrTransport for non-XML retour")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestClient_SendMessage_EncodeError(t *testing.T) {
	// badMsg is defined in codec_test.go: Validate() passes but MarshalXML returns error.
	c, _ := NewClient(WithSender(&mockSender{fn: func(_ context.Context, _ []byte) ([]byte, error) {
		t.Error("Send should not be called when Encode fails")
		return nil, nil
	}}))
	_, err := c.sendMessage(context.Background(), &badMsg{})
	if err == nil {
		t.Fatal("sendMessage() error = nil, want error from Encode")
	}
}

func TestHTTPSender_InvalidURL(t *testing.T) {
	// A null byte in the URL causes url.Parse (called by NewRequestWithContext) to fail.
	s := &httpSender{client: http.DefaultClient, baseURL: "http://host\x00/"}
	_, err := s.Send(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("Send() error = nil, want ErrTransport")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestHTTPSender_DoError(t *testing.T) {
	// Close the server before making the request to force a connection-refused error.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := ts.URL
	ts.Close()

	s := &httpSender{client: &http.Client{}, baseURL: url}
	_, err := s.Send(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("Send() error = nil, want ErrTransport")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestHTTPSender_BodyReadError(t *testing.T) {
	// errBodyTransport returns a 200 OK with a body that errors on Read.
	s := &httpSender{
		client:  &http.Client{Transport: &errBodyTransport{}},
		baseURL: "http://example.com",
	}
	_, err := s.Send(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("Send() error = nil, want ErrTransport")
	}
	if !errors.Is(err, ErrTransport) {
		t.Errorf("errors.Is(err, ErrTransport) = false, got: %v", err)
	}
}

func TestHTTPSender_Forbidden(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewClient(WithBaseURL(ts.URL))
	_, err := c.SendVerzoekToewijzing(context.Background(), validWMO302())
	if err == nil {
		t.Fatal("error = nil, want ErrAuthentication")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Errorf("errors.Is(err, ErrAuthentication) = false, got: %v", err)
	}
}
