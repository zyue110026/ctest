package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// GetFieldValuesFromFixtures searches fixtures (map[string]json.RawMessage) for values
// matching field. If types are provided, only those top-level keys are searched.
// Returns slice of json.RawMessage (each element is the JSON value found) and an error.
// If nothing found, returned slice is empty and error describes what was missing.
func GetFieldValuesFromFixtures(fixtures map[string]json.RawMessage, field string, types ...string) ([]json.RawMessage, error) {
	if field == "" {
		return nil, errors.New("field must not be empty")
	}

	// determine which top-level keys to search
	var toSearch []string
	if len(types) > 0 {
		// ensure the requested types exist in fixtures
		var missing []string
		for _, t := range types {
			if _, ok := fixtures[t]; !ok {
				missing = append(missing, t)
			}
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("requested fixture types not found: %v", missing)
		}
		toSearch = types
	} else {
		// search all keys present in fixtures
		for k := range fixtures {
			toSearch = append(toSearch, k)
		}
	}

	pathParts := strings.Split(field, ".")
	usePath := len(pathParts) > 1

	var results []json.RawMessage

	for _, topKey := range toSearch {
		raw := fixtures[topKey]
		if len(raw) == 0 || string(raw) == "null" {
			// skip nulls (per your requirement fixtures with null values are omitted by loader anyway)
			continue
		}

		var v interface{}
		if err := json.Unmarshal(raw, &v); err != nil {
			// if one top-level entry can't be parsed, return error (you can change to skip if desired)
			return nil, fmt.Errorf("failed to unmarshal fixture %q: %w", topKey, err)
		}

		if usePath {
			// strict path search
			found := findByPath(v, pathParts)
			for _, f := range found {
				b, err := json.Marshal(f)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal found value for %q: %w", field, err)
				}
				results = append(results, json.RawMessage(b))
			}
		} else {
			// recursive key-name search
			found := findByKeyRecursive(v, field)
			for _, f := range found {
				b, err := json.Marshal(f)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal found value for %q: %w", field, err)
				}
				results = append(results, json.RawMessage(b))
			}
		}
	}

	// If nothing found, return an error describing that
	if len(results) == 0 {
		return nil, fmt.Errorf("no values found for field %q in requested fixtures", field)
	}

	return results, nil
}

// findByPath walks obj following the exact path parts and returns the terminal values found.
// obj can be an array or object. The function collects all terminal nodes reached by walking the path.
func findByPath(obj interface{}, parts []string) []interface{} {
	if len(parts) == 0 {
		return nil
	}
	current := []interface{}{obj}
	for _, p := range parts {
		var next []interface{}
		for _, cur := range current {
			switch t := cur.(type) {
			case []interface{}:
				// iterate array, try to step into each element's map
				for _, elem := range t {
					if m, ok := elem.(map[string]interface{}); ok {
						if val, exists := m[p]; exists {
							next = append(next, val)
						}
					}
				}
			case map[string]interface{}:
				if val, exists := t[p]; exists {
					next = append(next, val)
				}
			default:
				// can't step further
			}
		}
		if len(next) == 0 {
			// no match for this part
			return nil
		}
		current = next
	}
	// current contains terminal nodes
	return current
}

// findByKeyRecursive searches obj recursively for any map key equal to keyName and returns values found.
func findByKeyRecursive(obj interface{}, keyName string) []interface{} {
	var out []interface{}
	switch t := obj.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if k == keyName {
				out = append(out, v)
			}
			// regardless, recurse into v
			out = append(out, findByKeyRecursive(v, keyName)...)
		}
	case []interface{}:
		for _, elem := range t {
			out = append(out, findByKeyRecursive(elem, keyName)...)
		}
	}
	// for other scalar types, nothing to do
	return out
}

// PrintJSONRawMessages pretty-prints a slice of json.RawMessage, each separated by a header.
func PrintJSONRawMessages(msgs []json.RawMessage) {
	for i, m := range msgs {
		var pretty json.RawMessage
		// Unmarshal then marshal with indent to ensure pretty format
		var tmp interface{}
		if err := json.Unmarshal(m, &tmp); err != nil {
			// fallback: print raw bytes
			fmt.Printf("Result #%d (raw): %s\n", i+1, string(m))
			continue
		}
		b, err := json.MarshalIndent(tmp, "", "  ")
		if err != nil {
			fmt.Printf("Result #%d (raw): %s\n", i+1, string(m))
			continue
		}
		pretty = json.RawMessage(b)
		fmt.Printf("Result #%d:\n%s\n\n", i+1, string(pretty))
	}
}
