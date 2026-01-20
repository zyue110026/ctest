package ctest

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"log"
	// "path/filepath"
	"reflect"
	// "runtime"
	"strings"

	k8sjson "k8s.io/apimachinery/pkg/util/json"
	// "log"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	fixtures "k8s.io/kubernetes/test/ctest/fixtures"
	utils "k8s.io/kubernetes/test/ctest/utils"
)

// Mode controls how external fixtures and hardcoded defaults are combined.
type Mode int

const (
	ExtendOnly Mode = iota
	OverrideOnly
	Union
)

// // GenerateEffectiveConfig takes a single entry (one element) from
// // ctestglobals.HardcodedConfig (pass it as interface{}), plus the Mode.
// // It returns:
// //   - effectiveObj: the resulting k8s-typed object (typed Go value, e.g. []v1.Container or v1.ConfigMap) after conversion
// //   - rawJSON: the JSON bytes produced from the HardcodedConfig
// //   - error

// // NOTE: the three Mode-* functions are called but left as stubs (no implementation) per your request.
// func GenerateEffectiveConfig(entry interface{}, mode Mode) (effectiveObjs interface{}, objType reflect.Type, rawJSON []byte, err error) {
// 	// Sanity check: entry must be a struct or pointer-to-struct
// 	v := reflect.ValueOf(entry)
// 	if !v.IsValid() {
// 		return nil, nil, nil, errors.New("entry is nil or invalid")
// 	}
// 	// If pointer, dereference
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	if v.Kind() != reflect.Struct {
// 		return nil, nil, nil, fmt.Errorf("entry must be a struct or pointer to struct; got %T", entry)
// 	}

// 	// Find HardcodedConfig field
// 	fieldVal := v.FieldByName("HardcodedConfig")
// 	if !fieldVal.IsValid() {
// 		return nil, nil, nil, errors.New("entry does not have HardcodedConfig field")
// 	}
// 	if !fieldVal.CanInterface() {
// 		// sometimes unexported - try to make a copy via reflect
// 		// but as a strict check, return error
// 		return nil, nil, nil, errors.New("cannot access HardcodedConfig field (unexported?)")
// 	}

// 	hardcoded := fieldVal.Interface()
// 	if hardcoded == nil {
// 		return nil, nil, nil, errors.New("HardcodedConfig is nil")
// 	}

// 	// 1) Convert hardcoded -> JSON using k8s built-in json util.
// 	// We marshal with the std encoding/json first (to preserve normal Go marshaling),
// 	// then validate/unmarshal with k8s util if you wish. Using k8sutil solely for unmarshalling
// 	// below is sufficient, but we show both for clarity.
// 	rawJSON, err = stdjson.Marshal(hardcoded)
// 	if err != nil {
// 		return nil, nil, nil, fmt.Errorf("failed to marshal HardcodedConfig to JSON: %w", err)
// 	}

// 	// fmt.Print(rawJSON)
// 	// fmt.Println(hcType)

// 	// At this point reconstructed.Interface() is a typed k8s object (e.g. []v1.Container)
// 	// We'll call the Mode-specific function to combine fixture + hardcoded.
// 	// For now these stubs are placeholders; they should take the reconstructed typed
// 	// object and the fixture (if loaded) and return the effective object.
// 	// 1) Load all non-null fixtures

// 	hardcodedConfigField := v.FieldByName("Field")

// 	k8sObjects := v.FieldByName("K8sObjects")

// 	// Convert slice to comma-separated string
// 	var objectsList []string
// 	if k8sObjects.IsValid() && k8sObjects.Kind() == reflect.Slice {
// 		for i := 0; i < k8sObjects.Len(); i++ {
// 			item := k8sObjects.Index(i)
// 			if item.Kind() == reflect.String {
// 				objectsList = append(objectsList, item.String())
// 			}
// 		}
// 	}

// 	k8sObjectsCSVString := strings.Join(objectsList, ",")

// 	fixtures, err := fixtures.LoadFixturesAsJSON("./fixtures/"+ctestglobals.TestExternalFixtureFile, k8sObjectsCSVString)

// 	if err != nil {
// 		log.Fatalf("load all fixtures failed: %v", err)
// 	}
// 	// pretty, err := stdjson.MarshalIndent(fixtures, "", "  ")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(string(pretty))

// 	externalFieldValues, err := utils.GetFieldValuesFromFixtures(fixtures, hardcodedConfigField.String())
// 	if err != nil {
// 		fmt.Println("err:", err)
// 	}
// 	// utils.PrintJSONRawMessages(externalFieldValues)
// 	// Variable to store the effective objects that will be marshaled to JSON
// 	var objectsToMarshal interface{}
// 	// If no fixtures found, use hardcoded config as is
// 	if len(fixtures) == 0 {
// 		objectsToMarshal = hardcoded
// 	} else {
// 		// Apply mode-based combination
// 		switch mode {
// 		case ExtendOnly:
// 			fmt.Printf("Calling ExtendOnly with %d external values\n", len(externalFieldValues))
// 			objectsToMarshal, err = extendOnly(rawJSON, externalFieldValues)

// 		case OverrideOnly:
// 			objectsToMarshal, err = overrideOnly(rawJSON, externalFieldValues, KeepMissingOriginal)
// 		case Union:
// 			objectsToMarshal, err = union(rawJSON, externalFieldValues)
// 		default:
// 			return nil, nil, rawJSON, fmt.Errorf("unknown Mode: %v", mode)
// 		}

// 		if err != nil {
// 			return nil, nil, rawJSON, fmt.Errorf("mode-combination failed: %w", err)
// 		}
// 	}

// 	// Get the type of hardcoded config for return
// 	hardcodedType := reflect.TypeOf(hardcoded)

// 	// Process based on the type of objectsToMarshal
// 	var effectiveJSON []byte
// 	var finalEffectiveObjs interface{}

// 	// Check if objectsToMarshal is already [][]byte (like what your mode functions return)
// 	if objs, ok := objectsToMarshal.([][]byte); ok {
// 		log.Printf("Generated %d JSON objects", len(objs))

// 		// Create a slice to hold the unmarshaled objects
// 		var unmarshaledObjs []interface{}

// 		for i, jsonData := range objs {
// 			log.Printf("\n=== Converting Result %d/%d ===", i+1, len(objs))
// 			fmt.Println(string(jsonData))

// 			// Unmarshal each JSON object
// 			hcType := reflect.TypeOf(hardcoded)
// 			var targetPtr reflect.Value
// 			if hcType.Kind() == reflect.Ptr {
// 				targetPtr = reflect.New(hcType.Elem())
// 			} else {
// 				targetPtr = reflect.New(hcType)
// 			}

// 			if err := k8sjson.Unmarshal(jsonData, targetPtr.Interface()); err != nil {
// 				return nil, nil, jsonData, fmt.Errorf("k8s json unmarshal into %s failed: %w", hcType.String(), err)
// 			}

// 			var reconstructed reflect.Value
// 			if targetPtr.Kind() == reflect.Ptr {
// 				reconstructed = targetPtr.Elem()
// 			} else {
// 				reconstructed = targetPtr
// 			}

// 			unmarshaledObj := reconstructed.Interface()
// 			unmarshaledObjs = append(unmarshaledObjs, unmarshaledObj)
// 			fmt.Println(unmarshaledObj)
// 		}

// 		// Marshal the slice of unmarshaled objects to JSON for return
// 		effectiveJSON, err = stdjson.Marshal(unmarshaledObjs)
// 		if err != nil {
// 			return nil, nil, nil, fmt.Errorf("failed to marshal effective objects to JSON: %w", err)
// 		}

// 		finalEffectiveObjs = unmarshaledObjs

// 	} else if objectsToMarshal != nil {
// 		// If it's not [][]byte, marshal whatever we got
// 		effectiveJSON, err = stdjson.Marshal(objectsToMarshal)
// 		if err != nil {
// 			return nil, nil, nil, fmt.Errorf("failed to marshal effective objects to JSON: %w", err)
// 		}
// 		finalEffectiveObjs = objectsToMarshal
// 	} else {
// 		// If mode functions returned nil (stubs), use hardcoded config
// 		effectiveJSON, err = stdjson.Marshal(hardcoded)
// 		if err != nil {
// 			return nil, nil, nil, fmt.Errorf("failed to marshal hardcoded to JSON: %w", err)
// 		}
// 		finalEffectiveObjs = hardcoded
// 	}

// 	return finalEffectiveObjs, hardcodedType, effectiveJSON, nil
// }

// GenerateEffectiveConfig processes a configuration entry by combining hardcoded defaults with
// external fixture data based on the specified merge mode. It returns typed Kubernetes objects,
// their runtime type information, and the combined configuration as JSON.
//
// This generic function handles Kubernetes API objects of any type T. The hardcoded configuration
// is merged with external fixture values according to the selected Mode, producing one or more
// effective configurations.
//
// Parameters:
//   - entry: A struct or pointer-to-struct containing at least:
//   - HardcodedConfig field: The default Kubernetes object configuration
//   - Field field: The field name to look up in external fixtures
//   - K8sObjects field: A []string of Kubernetes object types to load from fixtures
//   - mode: The merge strategy to use when combining configurations:
//   - ExtendOnly: Adds missing fields from external fixtures without overriding existing values
//   - OverrideOnly: Overrides existing fields with external values, keeping missing fields unchanged
//   - Union: Performs both override and extend operations (override first, then extend)
//
// Returns:
//   - effectiveObjs: A slice of typed Kubernetes objects ([]T) resulting from the merge operation.
//     Each element corresponds to one effective configuration from the external fixtures.
//     If no external fixtures are found, returns a single element with the hardcoded configuration.
//   - objType: The reflect.Type of the hardcoded configuration, useful for runtime type inspection.
//   - rawJSON: JSON representation of ALL effective configurations as a JSON array.
//     For example: [{"config1": "value1"}, {"config2": "value2"}]
//     If no merge was performed, returns the JSON of the original hardcoded configuration.
//   - err: Any error encountered during processing, such as:
//   - Invalid entry structure or missing required fields
//   - JSON marshaling/unmarshaling failures
//   - Fixture loading errors
//   - Mode-specific merge failures
//
// Type Parameters:
//   - T: The Kubernetes API object type (e.g., v1.Container, v1.ConfigMap, []v1.Container).
//     Must match the type stored in the HardcodedConfig field.
//
// Example usage:
//
//	// For processing container configurations
//	containers, objType, rawJSON, err := GenerateEffectiveConfig[[]v1.Container](
//	    entry,
//	    ctest.Union,
//	)
//
//	// For processing config map configurations
//	configMaps, objType, rawJSON, err := GenerateEffectiveConfig[v1.ConfigMap](
//	    entry,
//	    ctest.ExtendOnly,
//	)
//
// Processing flow:
//  1. Extract hardcoded configuration from entry.HardcodedConfig
//  2. Load external fixtures based on entry.Field and entry.K8sObjects
//  3. Apply the specified merge mode (ExtendOnly/OverrideOnly/Union)
//  4. Convert JSON results to typed Kubernetes objects of type T
//  5. Return typed objects, type information, and combined JSON
//
// Notes:
//   - External fixtures are loaded from "./fixtures/{TestExternalFixtureFile}"
//   - The function uses k8s.io/apimachinery/pkg/util/json for Kubernetes-compatible JSON handling
//   - All merge operations preserve Kubernetes object semantics and type safety
func GenerateEffectiveConfigReturnType[T any](entry interface{}, mode Mode) (effectiveObjs []T, effectiveObjsJson []byte, err error) {
	fmt.Println("=== GENERATE EFFECTIVE CONFIG START ===")
	v := reflect.ValueOf(entry)
	if !v.IsValid() {
		fmt.Println(ctestglobals.DebugPrefix(), "entry is nil or invalid")
		return nil, nil, errors.New("entry is nil or invalid")
	}
	// If pointer, dereference
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		fmt.Println(ctestglobals.DebugPrefix(), "entry must be a struct or pointer to struct")
		return nil, nil, fmt.Errorf("entry must be a struct or pointer to struct; got %T", entry)
	}

	// Find HardcodedConfig field
	fieldVal := v.FieldByName("HardcodedConfig")
	if !fieldVal.IsValid() {
		fmt.Println(ctestglobals.DebugPrefix(), "entry does not have HardcodedConfig field")
		return nil, nil, errors.New("entry does not have HardcodedConfig field")
	}
	if !fieldVal.CanInterface() {
		fmt.Println(ctestglobals.DebugPrefix(), "cannot access HardcodedConfig field (unexported?)")
		return nil, nil, errors.New("cannot access HardcodedConfig field (unexported?)")
	}

	hardcoded := fieldVal.Interface()
	if hardcoded == nil {
		fmt.Println(ctestglobals.DebugPrefix(), "HardcodedConfig is nil")
		return nil, nil, errors.New("HardcodedConfig is nil")
	}

	// 1) Convert hardcoded -> JSON using k8s built-in json util.
	originalRawJSON, err := stdjson.Marshal(hardcoded)
	if err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "failed to marshal HardcodedConfig to JSON")
		return nil, nil, fmt.Errorf("failed to marshal HardcodedConfig to JSON: %w", err)
	}

	// Get the field name and k8s objects
	hardcodedConfigField := v.FieldByName("Field")
	k8sObjects := v.FieldByName("K8sObjects")

	fmt.Println(ctestglobals.DebugPrefix(), "K8sObjects: ")
	fmt.Println(k8sObjects)

	// Convert slice to comma-separated string
	var objectsList []string
	if k8sObjects.IsValid() && k8sObjects.Kind() == reflect.Slice {
		if k8sObjects.IsNil() {
			fmt.Println(ctestglobals.DebugPrefix(), "[DEBUG] K8sObjects is nil, using empty list")
		} else {
			hasEmptyStrings := false
			for i := 0; i < k8sObjects.Len(); i++ {
				item := k8sObjects.Index(i)
				if item.Kind() == reflect.String {
					str := item.String()
					if str == "" {
						hasEmptyStrings = true
						fmt.Println(ctestglobals.DebugPrefix(), "[DEBUG] Found empty string at index %d in K8sObjects\n", i)
					} else {
						objectsList = append(objectsList, str)
					}
				}
			}

			if hasEmptyStrings {
				fmt.Println(ctestglobals.DebugPrefix(), "[WARNING] K8sObjects contains empty strings which were filtered out")
			}

			if k8sObjects.Len() > 0 && len(objectsList) == 0 {
				fmt.Println(ctestglobals.DebugPrefix(), "[WARNING] All strings in K8sObjects were empty after filtering")
			}
		}
	} else {
		fmt.Printf(ctestglobals.DebugPrefix(), "[DEBUG] K8sObjects is not a slice or invalid (Kind: %v, IsValid: %v)\n",
			k8sObjects.Kind(), k8sObjects.IsValid())
	}

	fmt.Printf(ctestglobals.DebugPrefix(), "[DEBUG] Loading fixtures for types: %v (count: %d)\n", objectsList, len(objectsList))

	fixtures, err := fixtures.LoadFixturesAsJSON(
		ctestglobals.TestExternalFixtureFile,
		objectsList...,
	)

	if err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "load all fixtures failed")
		log.Fatalf("load all fixtures failed: %v", err)
	}
	externalFieldValues, err := utils.GetFieldValuesFromFixtures(fixtures, hardcodedConfigField.String())
	if err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "err:", err)
	}

	// Process the results based on mode
	var jsonResults [][]byte
	if len(fixtures) != 0 {
		switch mode {
		case ExtendOnly:
			// fmt.Printf(ctestglobals.DebugPrefix(), "Calling ExtendOnly with %d external values\n", len(externalFieldValues))
			jsonResults, err = extendOnly(originalRawJSON, externalFieldValues)
		case OverrideOnly:
			jsonResults, err = overrideOnly(originalRawJSON, externalFieldValues, KeepMissingOriginal)
		case Union:
			jsonResults, err = union(originalRawJSON, externalFieldValues)
		default:
			return nil, nil, fmt.Errorf("unknown Mode: %v", mode)
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("mode-combination failed: %w", err)
	}

	// If no fixtures were processed, use the original JSON
	if jsonResults == nil {
		jsonResults = [][]byte{originalRawJSON}
	}

	// Convert each JSON result to type T and filter out duplicates
	effectiveObjs = make([]T, 0, len(jsonResults))
	var effectiveObjectsJSON []stdjson.RawMessage

	// Convert original hardcoded to T for comparison
	var originalObj T
	if err := k8sjson.Unmarshal(originalRawJSON, &originalObj); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal original config to type %T: %w", originalObj, err)
	}

	// fmt.Printf(ctestglobals.DebugPrefix(), "Original hardcoded config (as %T): %+v\n", originalObj, originalObj)

	// Compare JSON strings by normalizing them (remove whitespace)
	normalizedOriginalJSON := normalizeJSON(originalRawJSON)
	fmt.Printf(ctestglobals.DebugPrefix(), "Normalized original JSON: %s\n", normalizedOriginalJSON)

	for i, jsonData := range jsonResults {
		// log.Printf("\n=== Processing Result %d/%d ===", i+1, len(jsonResults))

		// Create a new zero value of type T
		var target T

		// Unmarshal JSON into T using k8s json util
		if err := k8sjson.Unmarshal(jsonData, &target); err != nil {
			return nil, nil, fmt.Errorf("k8s json unmarshal into %T failed: %w", target, err)
		}

		// Check if this result is identical to original hardcoded config
		normalizedResultJSON := normalizeJSON(jsonData)
		isIdenticalToOriginal := normalizedResultJSON == normalizedOriginalJSON

		// Also compare the typed objects using DeepEqual for extra safety
		isDeepEqual := reflect.DeepEqual(target, originalObj)

		if isIdenticalToOriginal || isDeepEqual {
			// fmt.Printf(ctestglobals.DebugPrefix(), "⚠️  Result %d is identical to original hardcoded config, filtering out\n", i+1)
			// fmt.Printf(ctestglobals.DebugPrefix(), "   JSON identical: %v, DeepEqual: %v\n", isIdenticalToOriginal, isDeepEqual)
			continue // Skip this result
		}

		// Add to results if not identical
		effectiveObjs = append(effectiveObjs, target)

		// Convert the typed object back to JSON
		objJSON, err := stdjson.Marshal(target)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal effective object %d to JSON: %w", i, err)
		}
		effectiveObjectsJSON = append(effectiveObjectsJSON, objJSON)

		fmt.Println(ctestglobals.DebugPrefix(), "✅ Added Result %d as unique effective object\n", i+1)
		log.Printf(ctestglobals.DebugPrefix(), "Successfully converted to type %T", target)
		fmt.Println(ctestglobals.DebugPrefix(), "Result value: %+v\n", target)
	}

	// Check if we have any unique results after filtering
	if len(effectiveObjs) == 0 {
		fmt.Println(ctestglobals.DebugPrefix(), "⚠️  All results were identical to original hardcoded config, returning nil")
		return nil, nil, nil
	}

	fmt.Println(ctestglobals.DebugPrefix(), "✅ Generated %d unique effective object(s) after filtering\n", len(effectiveObjs))

	// Marshal ALL effective objects as JSON array for effectiveObjsJson
	if len(effectiveObjectsJSON) > 0 {
		effectiveObjsJson, err = stdjson.Marshal(effectiveObjectsJSON)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal effective objects to JSON array: %w", err)
		}
	} else {
		// This shouldn't happen since we checked len(effectiveObjs) > 0
		effectiveObjsJson = originalRawJSON
	}

	fmt.Println("=== GENERATE EFFECTIVE CONFIG COMPLETE ===")
	return effectiveObjs, effectiveObjsJson, nil
}

// Helper function to normalize JSON by removing whitespace
func normalizeJSON(data []byte) string {
	var v interface{}
	if err := stdjson.Unmarshal(data, &v); err != nil {
		// If can't unmarshal, return as-is (trimmed)
		return strings.TrimSpace(string(data))
	}
	normalized, err := stdjson.Marshal(v)
	if err != nil {
		return strings.TrimSpace(string(data))
	}
	return string(normalized)
}

// Alternative: Helper function to compare JSON equality
func jsonEqual(a, b []byte) bool {
	var objA, objB interface{}
	if err := stdjson.Unmarshal(a, &objA); err != nil {
		return false
	}
	if err := stdjson.Unmarshal(b, &objB); err != nil {
		return false
	}
	return reflect.DeepEqual(objA, objB)
}

/*
 * Mode-specific stub functions
 *
 * You asked to NOT implement their internals now. Implement them later with your
 * actual merge/combine logic (and a fixture-loading/parsing step).
 *
 * Signatures:
 *  - hardcodedObj: the typed object reconstructed from HardcodedConfig (e.g. []v1.Container)
 *  - fixtureObj:    (nil for now) when you add fixture-loading, pass parsed fixture object here
 *
 * Return:
 *  - combined effective object (typed)
 *  - error
 */

func union(baseJSON []byte, externalFieldValues []stdjson.RawMessage) ([][]byte, error) {
	log.Println("=== UNION FUNCTION START (OVERRIDE + EXTEND) ===")
	log.Printf("Base JSON size: %d bytes", len(baseJSON))
	log.Printf("Number of external values: %d", len(externalFieldValues))

	// Parse base JSON as generic interface
	var baseData interface{}
	if err := stdjson.Unmarshal(baseJSON, &baseData); err != nil {
		log.Printf("ERROR: Failed to parse base JSON: %v", err)
		return nil, fmt.Errorf("failed to unmarshal base JSON: %w", err)
	}

	//Pretty print base for debugging
	prettyBase, _ := stdjson.MarshalIndent(baseData, "", "  ")
	log.Printf("BASE DATA (type: %T):\n%s", baseData, prettyBase)

	results := make([][]byte, len(externalFieldValues))

	for i, externalRaw := range externalFieldValues {
		// log.Printf("\n--- Processing external %d/%d ---", i+1, len(externalFieldValues))

		// Parse external value
		var externalData interface{}
		if err := stdjson.Unmarshal(externalRaw, &externalData); err != nil {
			log.Printf("ERROR: Failed to parse external %d: %v", i, err)
			return nil, fmt.Errorf("external %d: %w", i, err)
		}

		// Pretty print external for debugging
		// prettyExt, _ := stdjson.MarshalIndent(externalData, "", "  ")
		// log.Printf("EXTERNAL %d (type: %T):\n%s", i+1, externalData, prettyExt)

		// Perform UNION: First override, then extend
		resultData := unionRecursive(baseData, externalData, "")

		// Marshal result
		resultJSON, err := stdjson.MarshalIndent(resultData, "", "  ")
		if err != nil {
			log.Printf("ERROR: Failed to marshal result %d: %v", i, err)
			return nil, err
		}

		results[i] = resultJSON
		// log.Printf("✅ Result %d size: %d bytes", i+1, len(resultJSON))
		// log.Printf("Result %d:\n%s", i+1, string(resultJSON))
	}

	log.Printf("\n=== UNION COMPLETE ===")
	log.Printf("Generated %d result(s)", len(results))
	return results, nil
}

// unionRecursive performs both override and extend operations
// 1. First, override existing fields (replace base values with external values)
// 2. Then, add any new fields from external that don't exist in base
func unionRecursive(base, external interface{}, path string) interface{} {
	// Handle maps/objects
	if baseMap, ok := base.(map[string]interface{}); ok {
		if extMap, ok := external.(map[string]interface{}); ok {
			result := make(map[string]interface{})

			// PHASE 1: Copy all base fields first
			for key, baseValue := range baseMap {
				currentPath := path
				if currentPath != "" {
					currentPath += "."
				}
				currentPath += key

				// Check if this key exists in external
				if extValue, exists := extMap[key]; exists {
					// Key exists in external - OVERRIDE recursively
					result[key] = unionRecursive(baseValue, extValue, currentPath)
					log.Printf("  [UNION OVERRIDE] %s: overridden", currentPath)
				} else {
					// Key doesn't exist in external - KEEP original
					result[key] = baseValue
					log.Printf("  [UNION KEEP] %s: kept original", currentPath)
				}
			}

			// PHASE 2: Add any new fields from external that don't exist in base
			for key, extValue := range extMap {
				currentPath := path
				if currentPath != "" {
					currentPath += "."
				}
				currentPath += key

				if _, exists := result[key]; !exists {
					// This is a new field from external - EXTEND
					result[key] = extValue
					log.Printf("  [UNION EXTEND] %s: added new field", currentPath)
				}
			}

			return result
		}

		// External is not a map, return external (override entire structure)
		log.Printf("  [UNION REPLACE] %s: entire structure replaced", path)
		return external
	}

	// Handle arrays
	if baseArr, ok := base.([]interface{}); ok {
		if extArr, ok := external.([]interface{}); ok {
			// For arrays, we need to handle both override and extend
			maxLen := len(baseArr)
			if len(extArr) > maxLen {
				maxLen = len(extArr)
			}

			result := make([]interface{}, maxLen)

			// PHASE 1: Override existing indices
			for i := 0; i < len(baseArr) && i < len(extArr); i++ {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				result[i] = unionRecursive(baseArr[i], extArr[i], arrayPath)
				log.Printf("  [UNION ARRAY] %s: overridden", arrayPath)
			}

			// PHASE 2: Keep base values where external doesn't exist
			for i := len(extArr); i < len(baseArr); i++ {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				result[i] = baseArr[i]
				log.Printf("  [UNION ARRAY] %s: kept original", arrayPath)
			}

			// PHASE 3: Extend with external values beyond base length
			for i := len(baseArr); i < len(extArr); i++ {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				result[i] = extArr[i]
				log.Printf("  [UNION ARRAY] %s: extended with external", arrayPath)
			}

			return result
		}

		// External is not an array, return external (override entire array)
		log.Printf("  [UNION REPLACE] %s: entire array replaced", path)
		return external
	}

	// For primitive values
	// Always use external value (override)
	if isCompatibleType(base, external) {
		log.Printf("  [UNION PRIMITIVE] %s: %v → %v", path, base, external)
		return external
	}

	// Types incompatible, keep external (it's an override)
	log.Printf("  [UNION TYPE MISMATCH] %s: using external value", path)
	return external
}

// overrideOnly merges external fixture values into the base hardcoded JSON,

// OverrideMode defines how to handle fields that are missing in external values
type OverrideMode int

const (
	// SetMissingToNil - fields missing in external become nil/null
	SetMissingToNil OverrideMode = iota
	// KeepMissingOriginal - fields missing in external keep their original values
	KeepMissingOriginal
)

func overrideOnly(baseJSON []byte, externalFieldValues []stdjson.RawMessage, mode OverrideMode) ([][]byte, error) {
	log.Println(ctestglobals.DebugPrefix(), "=== OVERRIDE ONLY FUNCTION START ===")
	log.Printf("Mode: %v", mode)
	log.Printf("Base JSON size: %d bytes", len(baseJSON))
	log.Printf("Number of external values: %d", len(externalFieldValues))

	// Parse base JSON as generic interface
	var baseData interface{}
	if err := stdjson.Unmarshal(baseJSON, &baseData); err != nil {
		log.Printf(ctestglobals.DebugPrefix(), "ERROR: Failed to parse base JSON: %v", err)
		return nil, fmt.Errorf("failed to unmarshal base JSON: %w", err)
	}

	// Pretty print base for debugging
	prettyBase, _ := stdjson.MarshalIndent(baseData, "", "  ")
	log.Printf(ctestglobals.DebugPrefix(), "BASE DATA (type: %T):\n%s", baseData, prettyBase)

	results := make([][]byte, 0, len(externalFieldValues))

	for i, externalRaw := range externalFieldValues {
		// log.Printf("\n--- Processing external %d/%d ---", i+1, len(externalFieldValues))

		// Parse external value
		var externalData interface{}
		if err := stdjson.Unmarshal(externalRaw, &externalData); err != nil {
			log.Printf(ctestglobals.DebugPrefix(), "ERROR: Failed to parse external %d: %v", i, err)
			return nil, fmt.Errorf("external %d: %w", i, err)
		}

		// Pretty print external for debugging
		prettyExt, _ := stdjson.MarshalIndent(externalData, "", "  ")
		log.Printf(ctestglobals.DebugPrefix(), "EXTERNAL %d (type: %T):\n%s", i+1, externalData, prettyExt)

		// Perform override-only merge based on mode
		var resultData interface{}
		switch mode {
		case SetMissingToNil:
			resultData = overrideSetMissingToNil(baseData, externalData, "")
		case KeepMissingOriginal:
			resultData = overrideKeepMissingOriginal(baseData, externalData, "")
		default:
			return nil, fmt.Errorf("unknown OverrideMode: %v", mode)
		}

		// Check if all values became nil (only for SetMissingToNil mode)
		if mode == SetMissingToNil && isAllNil(resultData) {
			log.Printf(ctestglobals.DebugPrefix(), "WARNING: Result %d is all nil, skipping", i+1)
			continue
		}

		// Marshal result
		resultJSON, err := stdjson.MarshalIndent(resultData, "", "  ")
		if err != nil {
			log.Printf(ctestglobals.DebugPrefix(), "ERROR: Failed to marshal result %d: %v", i, err)
			return nil, err
		}

		results = append(results, resultJSON)
		// log.Printf(ctestglobals.DebugPrefix(), "✅ Result %d size: %d bytes", i+1, len(resultJSON))
		// log.Printf(ctestglobals.DebugPrefix(), "Result %d:\n%s", i+1, string(resultJSON))
	}

	log.Printf(ctestglobals.DebugPrefix(), "\n=== OVERRIDE ONLY COMPLETE ===")
	log.Printf(ctestglobals.DebugPrefix(), "Generated %d valid result(s)", len(results))

	if len(results) == 0 && mode == SetMissingToNil {
		log.Println(ctestglobals.DebugPrefix(), "WARNING: No valid results generated (all became nil)")
		return nil, nil
	}

	return results, nil
}

// overrideSetMissingToNil - fields missing in external become nil
func overrideSetMissingToNil(base, external interface{}, path string) interface{} {
	// Handle maps/objects
	if baseMap, ok := base.(map[string]interface{}); ok {
		result := make(map[string]interface{})

		// For each key in base
		for key, baseValue := range baseMap {
			currentPath := path
			if currentPath != "" {
				currentPath += "."
			}
			currentPath += key

			// Check if this key exists in external
			if extMap, ok := external.(map[string]interface{}); ok {
				if extValue, exists := extMap[key]; exists {
					// Key exists in external, recursively override
					result[key] = overrideSetMissingToNil(baseValue, extValue, currentPath)
				} else {
					// Key doesn't exist in external, set to nil
					result[key] = nil
					log.Printf("  [OVERRIDE] %s → nil (missing in external)", currentPath)
				}
			} else {
				// External is not a map, can't override specific keys
				// Set entire base value to external (if external is primitive)
				return external
			}
		}

		return result
	}

	// Handle arrays
	if baseArr, ok := base.([]interface{}); ok {
		// If external is also an array, override element by element
		if extArr, ok := external.([]interface{}); ok {
			result := make([]interface{}, len(baseArr))

			for i := range baseArr {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)

				if i < len(extArr) {
					// External has element at this index, override
					result[i] = overrideSetMissingToNil(baseArr[i], extArr[i], arrayPath)
				} else {
					// External doesn't have element at this index, set to nil
					result[i] = nil
					log.Printf("  [OVERRIDE] %s → nil (missing in external)", arrayPath)
				}
			}

			return result
		}

		// External is not an array, can't override array elements
		// Set entire array to external (if external is primitive)
		return external
	}

	// For primitive values (string, number, bool, null)
	// Always replace base value with external value
	// If external is nil, return nil
	if external == nil {
		log.Printf("  [OVERRIDE] %s → nil", path)
		return nil
	}

	// Type check: if types are compatible, override
	if isCompatibleType(base, external) {
		log.Printf("  [OVERRIDE] %s: %v → %v", path, base, external)
		return external
	}

	// Types incompatible, keep external (or nil if external is wrong type)
	log.Printf("  [TYPE MISMATCH] %s: %T → %T", path, base, external)
	return external
}

// overrideKeepMissingOriginal - fields missing in external keep original values
func overrideKeepMissingOriginal(base, external interface{}, path string) interface{} {
	// Handle maps/objects
	if baseMap, ok := base.(map[string]interface{}); ok {
		result := make(map[string]interface{})

		// For each key in base
		for key, baseValue := range baseMap {
			currentPath := path
			if currentPath != "" {
				currentPath += "."
			}
			currentPath += key

			// Check if this key exists in external
			if extMap, ok := external.(map[string]interface{}); ok {
				if extValue, exists := extMap[key]; exists {
					// Key exists in external, recursively override
					result[key] = overrideKeepMissingOriginal(baseValue, extValue, currentPath)
					// log.Printf("  [OVERRIDE] %s: %v → %v", currentPath, baseValue, extValue)
				} else {
					// Key doesn't exist in external, keep original value
					result[key] = baseValue
					log.Printf("  [KEEP] %s: %v (missing in external)", currentPath, baseValue)
				}
			} else {
				// External is not a map, replace entire value
				log.Printf("  [REPLACE ALL] %s: entire structure replaced", path)
				return external
			}
		}

		return result
	}

	// Handle arrays
	if baseArr, ok := base.([]interface{}); ok {
		// If external is also an array, override element by element
		if extArr, ok := external.([]interface{}); ok {
			result := make([]interface{}, len(baseArr))

			for i := range baseArr {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)

				if i < len(extArr) {
					// External has element at this index, override
					result[i] = overrideKeepMissingOriginal(baseArr[i], extArr[i], arrayPath)
				} else {
					// External doesn't have element at this index, keep original
					result[i] = baseArr[i]
					log.Printf("  [KEEP] %s: %v (missing in external)", arrayPath, baseArr[i])
				}
			}

			return result
		}

		// External is not an array, replace entire value
		log.Printf("  [REPLACE ALL] %s: entire array replaced", path)
		return external
	}

	// For primitive values, always replace with external value
	if isCompatibleType(base, external) {
		log.Printf("  [OVERRIDE] %s: %v → %v", path, base, external)
		return external
	}

	// Types incompatible, keep base value
	log.Printf("  [TYPE MISMATCH] %s: keeping base value %v", path, base)
	return base
}

// Helper functions (same as before)
func isAllNil(data interface{}) bool {
	if data == nil {
		return true
	}

	switch v := data.(type) {
	case map[string]interface{}:
		for _, val := range v {
			if !isAllNil(val) {
				return false
			}
		}
		return len(v) > 0

	case []interface{}:
		for _, elem := range v {
			if !isAllNil(elem) {
				return false
			}
		}
		return len(v) > 0

	default:
		return false
	}
}

func isCompatibleType(base, external interface{}) bool {
	if base == nil || external == nil {
		return true
	}

	baseType := reflect.TypeOf(base)
	extType := reflect.TypeOf(external)

	if isNumber(base) && isNumber(external) {
		return true
	}

	return baseType == extType
}

func isNumber(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

// ExtendOnly merges external fixture values into the base hardcoded JSON,

func extendOnly(baseJSON []byte, externalFieldValues []stdjson.RawMessage) ([][]byte, error) {
	log.Println("=== EXTEND ONLY (RECURSIVE MERGE) ===")

	// Parse base as generic interface to preserve structure
	var baseData interface{}
	if err := stdjson.Unmarshal(baseJSON, &baseData); err != nil {
		log.Printf("Failed to parse base JSON: %v", err)
		return nil, fmt.Errorf("failed to unmarshal base JSON: %w", err)
	}

	// log.Printf("Base data type: %T\n", baseData)

	results := make([][]byte, len(externalFieldValues))

	for i, externalRaw := range externalFieldValues {
		// log.Printf("\n🔄 Processing external %d/%d", i+1, len(externalFieldValues))

		// Parse external as generic interface
		var externalData interface{}
		if err := stdjson.Unmarshal(externalRaw, &externalData); err != nil {
			log.Printf("Failed to parse external %d: %v", i, err)
			return nil, fmt.Errorf("external %d: %w", i, err)
		}

		// log.Printf("External data type: %T", externalData)

		// Deep merge: add missing fields at any level
		resultData := deepMergeAddMissing(baseData, externalData)

		// Marshal result
		resultJSON, err := stdjson.MarshalIndent(resultData, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal result %d: %v", i, err)
			return nil, err
		}

		results[i] = resultJSON
		// log.Printf("✅ Result %d generated\n%s", i+1, string(resultJSON))
	}

	log.Printf("\n=== COMPLETE: Generated %d results ===", len(results))
	return results, nil
}

// deepMergeAddMissing recursively adds missing fields from external to base
func deepMergeAddMissing(base, external interface{}) interface{} {
	// If base is a map
	if baseMap, ok := base.(map[string]interface{}); ok {
		// If external is also a map, merge them
		if extMap, ok := external.(map[string]interface{}); ok {
			result := make(map[string]interface{})

			// First, copy all base fields
			for key, baseValue := range baseMap {
				result[key] = baseValue
			}

			// Then, for each external field
			for key, extValue := range extMap {
				// If key doesn't exist in base, add it
				if baseValue, exists := result[key]; !exists {
					result[key] = extValue
				} else {
					// If key exists in both, recursively merge if both are objects/arrays
					result[key] = deepMergeAddMissing(baseValue, extValue)
				}
			}

			return result
		}
		// If external is not a map, return base unchanged
		return baseMap
	}

	// If base is an array
	if baseArr, ok := base.([]interface{}); ok {
		// If external is also an array, merge element by element
		if extArr, ok := external.([]interface{}); ok {
			result := make([]interface{}, len(baseArr))

			// Copy base array
			copy(result, baseArr)
			// for i := range baseArr {
			// 	result[i] = baseArr[i]
			// }

			// For each position where external has an element
			for i, extValue := range extArr {
				// If base has element at this position, merge them
				if i < len(result) {
					result[i] = deepMergeAddMissing(result[i], extValue)
				} else {
					// If base doesn't have element at this position, add it
					result = append(result, extValue)
				}
			}

			return result
		}
		// If external is not an array, return base unchanged
		return baseArr
	}

	// For primitive values (string, number, bool, null)
	// Return base value (don't override with external)
	return base
}
