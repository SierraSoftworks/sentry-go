package sentry

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type httpTransport struct {
	transport *http.Client
}

func newHTTPTransport() Transport {
	return &httpTransport{
		transport: http.DefaultClient,
	}
}

func (t *httpTransport) Send(dsn string, packet Packet) error {
	if dsn == "" {
		return nil
	}

	url, authHeader := t.parseDSN(dsn)

	body, contentType, err := t.serializePacket(packet)
	if err != nil {
		return errors.Wrap(err, "failed to serialize packet")
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return errors.Wrap(err, "failed to create new request")
	}

	req.Header.Set("X-Sentry-Auth", authHeader)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", fmt.Sprintf("sentry-go %s (Sierra Softworks; github.com/SierraSoftworks/sentry-go)", version))

	res, err := t.transport.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to submit request")
	}

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("got http status %d, expected 200", res.StatusCode)
	}

	return nil
}

func (t *httpTransport) parseDSN(dsn string) (url, authHeader string) {
	d, err := newDSN(dsn)
	if err != nil {
		// TODO: Indicate that this is an invalid DSN to the user
		return
	}

	return d.URL, d.AuthHeader()
}

func (t *httpTransport) serializePacket(packet Packet) (io.Reader, string, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(packet); err != nil {
		return nil, "", errors.Wrap(err, "failed to encode JSON payload data")
	}

	if buf.Len() < 1000 {
		return buf, "application/json; charset=utf8", nil
	}

	cbuf := bytes.NewBuffer([]byte{})
	b64 := base64.NewEncoder(base64.StdEncoding, cbuf)
	deflate, err := zlib.NewWriterLevel(b64, zlib.BestCompression)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to configure zlib deflate")
	}

	if _, err := io.Copy(deflate, buf); err != nil {
		return nil, "", errors.Wrap(err, "failed to deflate message")
	}

	deflate.Close()
	b64.Close()

	return cbuf, "application/octet-stream", nil
}
