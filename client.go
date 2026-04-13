package iwmo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// Sender abstracts message transport. The default implementation uses HTTP(S),
// but callers may substitute any backend (e.g. VECOZO message bus, file-based
// exchange) by passing [WithSender] to [NewClient].
type Sender interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}

// Client sends and receives iWMO v3.2 messages.
//
// Construct one with [NewClient] and at least [WithBaseURL] or [WithSender].
type Client struct {
	httpClient   *http.Client
	baseURL      string
	agbCode      string
	gemeenteCode string
	logger       *slog.Logger
	sender       Sender
}

// Option is a functional option for [NewClient].
type Option func(*Client)

// WithHTTPClient replaces the default [http.DefaultClient] with c.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithBaseURL sets the base URL of the iWMO endpoint.
func WithBaseURL(url string) Option {
	return func(cl *Client) { cl.baseURL = url }
}

// WithAGBCode sets the AGB-Z (healthcare provider) code of the zorgaanbieder.
func WithAGBCode(code string) Option {
	return func(cl *Client) { cl.agbCode = code }
}

// WithGemeenteCode sets the CBS municipality code (gemeentecode) of the gemeente.
func WithGemeenteCode(code string) Option {
	return func(cl *Client) { cl.gemeenteCode = code }
}

// WithLogger sets the structured logger used for debug output.
func WithLogger(l *slog.Logger) Option {
	return func(cl *Client) { cl.logger = l }
}

// WithSender replaces the default HTTP sender with s.
// When a custom Sender is provided, [WithBaseURL] is not required.
func WithSender(s Sender) Option {
	return func(cl *Client) { cl.sender = s }
}

// NewClient creates a new Client, applies opts in order, and validates the
// resulting configuration. Returns an error if the configuration is invalid
// (e.g. neither baseURL nor a custom Sender was provided).
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: http.DefaultClient,
		logger:     slog.Default(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.sender == nil && c.baseURL == "" {
		return nil, fmt.Errorf("%w: WithBaseURL or WithSender is required", ErrConfiguration)
	}
	if c.sender == nil {
		c.sender = &httpSender{client: c.httpClient, baseURL: c.baseURL}
	}
	return c, nil
}

// SendToewijzing sends a [WMO301] (Toewijzing) message from gemeente to
// zorgaanbieder and returns the [WMO304] acknowledgement.
func (c *Client) SendToewijzing(ctx context.Context, msg *WMO301) (*WMO304, error) {
	return c.sendMessage(ctx, msg)
}

// SendVerzoekToewijzing sends a [WMO302] (Verzoek om Toewijzing) message from
// zorgaanbieder to gemeente and returns the [WMO304] acknowledgement.
func (c *Client) SendVerzoekToewijzing(ctx context.Context, msg *WMO302) (*WMO304, error) {
	return c.sendMessage(ctx, msg)
}

// SendDeclaratie sends a [WMO303] (Declaratie) message from zorgaanbieder to
// gemeente and returns the [WMO304] acknowledgement.
func (c *Client) SendDeclaratie(ctx context.Context, msg *WMO303) (*WMO304, error) {
	return c.sendMessage(ctx, msg)
}

// SendMutatie sends a [WMO305] (Mutatie) message from zorgaanbieder to
// gemeente and returns the [WMO304] acknowledgement.
func (c *Client) SendMutatie(ctx context.Context, msg *WMO305) (*WMO304, error) {
	return c.sendMessage(ctx, msg)
}

// SendStatusmelding sends a [WMO315] (Statusmelding) message from
// zorgaanbieder to gemeente and returns the [WMO304] acknowledgement.
func (c *Client) SendStatusmelding(ctx context.Context, msg *WMO315) (*WMO304, error) {
	return c.sendMessage(ctx, msg)
}

func (c *Client) sendMessage(ctx context.Context, msg Message) (*WMO304, error) {
	if err := msg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidMessage, err)
	}
	data, err := Encode(msg)
	if err != nil {
		return nil, err
	}
	c.logger.DebugContext(ctx, "sending iWMO message",
		"type", msg.MessageType(),
		"bytes", len(data),
	)
	resp, err := c.sender.Send(ctx, data)
	if err != nil {
		return nil, err
	}
	retour, err := DecodeAs[WMO304](resp)
	if err != nil {
		return nil, fmt.Errorf("%w: decoding WMO304 retour: %v", ErrTransport, err)
	}
	if verr := retour.Validate(); verr != nil {
		return nil, fmt.Errorf("%w: invalid WMO304 retour: %w", ErrTransport, verr)
	}
	c.logger.DebugContext(ctx, "received WMO304 retour",
		"codes", len(retour.RetourCodes),
	)
	return retour, nil
}

// httpSender is the default Sender implementation; it POSTs XML to the
// configured base URL and returns the response body.
type httpSender struct {
	client  *http.Client
	baseURL string
}

func (s *httpSender) Send(ctx context.Context, data []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: creating request: %v", ErrTransport, err)
	}
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransport, err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusUnauthorized,
		resp.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("%w: HTTP %d", ErrAuthentication, resp.StatusCode)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return nil, fmt.Errorf("%w: HTTP %d", ErrTransport, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, fmt.Errorf("%w: reading response body: %v", ErrTransport, err)
	}
	return body, nil
}
