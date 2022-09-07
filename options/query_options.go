package options

import (
	"fmt"
	"github.com/imgproxy/imgproxy/v3/ierrors"
	"net/http"
	"net/url"
	"strings"
)

func ParseQuery(path string, headers http.Header) (*ProcessingOptions, string, error) {
	if path == "" || path == "/" {
		return nil, "", ierrors.New(404, fmt.Sprintf("Invalid path: %s", path), "Invalid URL")
	}

	parsedUrl, err := url.ParseRequestURI(path)

	if queryStart := strings.IndexByte(path, '?'); queryStart >= 0 {
		path = path[:queryStart]
	}

	if err != nil {
		return nil, "", ierrors.New(404, fmt.Sprintf("Invalid Query: %s", path), "Invalid Query")
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	return parseQueryOptions(parts, parsedUrl.Query(), headers)
}

func parseQueryOptions(pathParts []string, ops url.Values, headers http.Header) (*ProcessingOptions, string, error) {

	po, err := defaultProcessingOptions(headers)
	if err != nil {
		return nil, "", err
	}

	options := parseQueryToURLOptions(ops)

	if err = applyURLOptions(po, options); err != nil {
		return nil, "", err
	}

	url, extension, err := DecodeURL(pathParts)
	if err != nil {
		return nil, "", err
	}

	if len(extension) > 0 {
		if err = applyFormatOption(po, []string{extension}); err != nil {
			return nil, "", err
		}
	}

	return po, url, nil
}

func parseQueryToURLOptions(opts url.Values) urlOptions {
	parsed := make(urlOptions, 0, len(opts))

	for i, opt := range opts {
		args := strings.Split(opt[len(opt)-1], ":")

		parsed = append(parsed, urlOption{Name: i, Args: args})
	}

	return parsed
}
