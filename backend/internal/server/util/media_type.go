package util

import (
	"fmt"
	"mime"
	"sort"
	"strings"
)

func parseMediaType(s string) (*mediaType, error) {
	parsedType, params, err := mime.ParseMediaType(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse media type: %w", err)
	}

	parsedParts := strings.SplitN(parsedType, "/", 2)
	if len(parsedParts) != 2 {
		return nil, fmt.Errorf("invalid media type: %s", parsedType)
	}

	mt := mediaType{
		Type:    strings.TrimSpace(parsedParts[0]),
		Subtype: strings.TrimSpace(parsedParts[1]),
		Quality: 1.0,
		Params:  params,
	}

	if q, ok := params["q"]; ok {
		if count, err := fmt.Sscanf(q, "q=%f", &mt.Quality); count != 1 || err != nil {
			return nil, fmt.Errorf("failed to parse quality parameter: %w", err)
		}
	}

	return &mt, nil
}

type mediaType struct {
	Type    string
	Subtype string
	Quality float64
	Params  map[string]string
}

func (m mediaType) Accepts(contentType string) bool {
	parsedType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	parsedParts := strings.SplitN(parsedType, "/", 2)
	if len(parsedParts) != 2 {
		return false
	}

	if m.Type == "*" || (m.Type == parsedParts[0] && (m.Subtype == "*" || m.Subtype == parsedParts[1])) || (m.Type == "*" && m.Subtype == "*") {
		return true
	}

	return false
}

func parseAcceptHeader(header string) []mediaType {
	var mediaTypes []mediaType
	for _, part := range strings.Split(header, ",") {
		mt, err := parseMediaType(strings.TrimSpace(part))
		if err != nil {
			continue
		}
		mediaTypes = append(mediaTypes, *mt)
	}

	sort.Slice(mediaTypes, func(i, j int) bool {
		return mediaTypes[i].Quality > mediaTypes[j].Quality
	})

	return mediaTypes
}
