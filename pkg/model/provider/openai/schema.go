package openai

import (
	"maps"
	"slices"

	"github.com/openai/openai-go/v3/shared"

	"github.com/docker/cagent/pkg/tools"
)

// ConvertParametersToSchema converts parameters to OpenAI Schema format
func ConvertParametersToSchema(params any) (shared.FunctionParameters, error) {
	p, err := tools.SchemaToMap(params)
	if err != nil {
		return nil, err
	}

	return normalizeUnionTypes(fixSchemaArrayItems(makeAllRequired(p))), nil
}

// makeAllRequired make all the parameters "required"
// because that's what the Response API wants, now.
func makeAllRequired(schema shared.FunctionParameters) shared.FunctionParameters {
	if schema == nil {
		return makeAllRequired(map[string]any{"type": "object", "properties": map[string]any{}})
	}

	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return schema
	}

	reallyRequired := map[string]bool{}
	if required, ok := schema["required"].([]any); ok {
		for _, name := range required {
			reallyRequired[name.(string)] = true
		}
	}

	// We can't use a nil 'required' attribute
	newRequired := []any{}

	// Sort property names for deterministic output
	propNames := slices.Sorted(maps.Keys(properties))

	for _, propName := range propNames {
		newRequired = append(newRequired, propName)
		if reallyRequired[propName] {
			continue
		}

		// Make its type nullable
		if propMap, ok := properties[propName].(map[string]any); ok {
			if typeValue, ok := propMap["type"].(string); ok {
				propMap["type"] = []string{typeValue, "null"}
			}
		}
	}

	schema["required"] = newRequired
	schema["additionalProperties"] = false
	return schema
}

// In Docker Desktop 4.52, the MCP Gateway produces an invalid tools shema for `mcp-config-set`.
func fixSchemaArrayItems(schema shared.FunctionParameters) shared.FunctionParameters {
	propertiesValue, ok := schema["properties"]
	if !ok {
		return schema
	}

	properties, ok := propertiesValue.(map[string]any)
	if !ok {
		return schema
	}

	for _, propValue := range properties {
		prop, ok := propValue.(map[string]any)
		if !ok {
			continue
		}

		checkForMissingItems := false
		switch t := prop["type"].(type) {
		case string:
			checkForMissingItems = t == "array"
		case []string:
			checkForMissingItems = slices.Contains(t, "array")
		}
		if !checkForMissingItems {
			continue
		}

		if _, ok := prop["items"]; !ok {
			prop["items"] = map[string]any{"type": "object"}
		}
	}

	return schema
}

// normalizeUnionTypes converts union types like ["array", "null"] back to simple types
// for compatibility with AI gateways that don't support JSON Schema union types.
// This is needed for Cloudflare AI Gateway and similar proxies.
func normalizeUnionTypes(schema shared.FunctionParameters) shared.FunctionParameters {
	if schema == nil {
		return schema
	}

	// Convert union types at the current level
	if typeArray, ok := schema["type"].([]any); ok {
		if len(typeArray) == 2 {
			// Find the non-null type
			for _, t := range typeArray {
				if tStr, ok := t.(string); ok && tStr != "null" {
					schema["type"] = tStr
					break
				}
			}
		}
	} else if typeArray, ok := schema["type"].([]string); ok {
		if len(typeArray) == 2 {
			// Find the non-null type
			for _, t := range typeArray {
				if t != "null" {
					schema["type"] = t
					break
				}
			}
		}
	}

	// Convert anyOf patterns like {"anyOf": [{"type":"string"},{"type":"null"}]} to {"type":"string"}
	// This is needed for Gemini via Cloudflare which doesn't support anyOf in tool parameters.
	if anyOf, ok := schema["anyOf"].([]any); ok {
		for _, item := range anyOf {
			if itemMap, ok := item.(map[string]any); ok {
				if typStr, ok := itemMap["type"].(string); ok && typStr != "null" {
					schema["type"] = typStr
					delete(schema, "anyOf")
					break
				}
			}
		}
	}

	// Recursively handle properties
	if propertiesValue, ok := schema["properties"]; ok {
		if properties, ok := propertiesValue.(map[string]any); ok {
			for _, propValue := range properties {
				if prop, ok := propValue.(map[string]any); ok {
					normalizeUnionTypes(prop)
				}
			}
		}
	}

	// Recursively handle items (for arrays)
	if items, ok := schema["items"].(map[string]any); ok {
		normalizeUnionTypes(items)
	}

	// Recursively handle additionalProperties
	if addProps, ok := schema["additionalProperties"].(map[string]any); ok {
		normalizeUnionTypes(addProps)
	}

	return schema
}
