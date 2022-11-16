// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package validator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/project-radius/radius/pkg/radlogger"
)

var (
	ErrSpecDocumentNotFound = errors.New("not found OpenAPI specification document")
)

// Loader is the OpenAPI spec loader implementation.
type Loader struct {
	validators        map[string]validator
	supportedVersions map[string][]string
	providerName      string
	rootScopePrefix   string
	rootScopeParam    string
	specFiles         fs.FS
}

// Name returns the name of loader.
func (l *Loader) Name() string {
	return l.providerName
}

// SupportedVersions returns supported api version for resource type
func (l *Loader) SupportedVersions(resourceType string) []string {
	if versions, ok := l.supportedVersions[resourceType]; ok {
		return versions
	}
	return []string{}
}

// GetValidator returns the cached validator.
func (l *Loader) GetValidator(resourceType, version string) (Validator, bool) {
	// ARM types are compared case-insensitively
	v, ok := l.validators[getValidatorKey(resourceType, version)]
	if ok {
		return &v, true
	}
	return nil, false
}

// LoadSpec loads the swagger files and caches the validator.
func LoadSpec(ctx context.Context, providerName string, specs fs.FS, rootScopePrefix string, rootScopeParam string) (*Loader, error) {
	log := logr.FromContextOrDiscard(ctx)
	l := &Loader{
		providerName:      providerName,
		validators:        map[string]validator{},
		supportedVersions: map[string][]string{},
		rootScopePrefix:   rootScopePrefix,
		rootScopeParam:    rootScopeParam,
		specFiles:         specs,
	}

	// Walk through embedded files to load OpenAPI spec document.
	err := fs.WalkDir(l.specFiles, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}

		// Skip the shared common-types
		if strings.HasPrefix(path, "specification/common-types") {
			return nil
		}

		// Check if specification file pathname is valid and skip global.json.
		parsed := parseSpecFilePath(path)
		if parsed == nil {
			log.V(radlogger.Warn).Info(fmt.Sprintf("failed to parse %s", path))
			return nil
		}

		if pn, ok := parsed["provider"]; !ok || !strings.EqualFold(pn, l.providerName) || parsed["resourcetype"] == "global" {
			return nil
		}

		// Load OpenAPI Spec
		specDoc, err := loads.Spec(
			path,
			loads.WithDocLoader(func(path string) (json.RawMessage, error) {
				data, err := fs.ReadFile(l.specFiles, path)
				return json.RawMessage(data), err
			}))
		if err != nil {
			return err
		}

		// Expand $ref external references.
		wDoc, err := specDoc.Expanded(&spec.ExpandOptions{
			RelativeBase: path,
			PathLoader: func(path string) (json.RawMessage, error) {
				// Trim before 'specification' to convert relative path.
				first := strings.Index(path, "specification")
				data, err := fs.ReadFile(l.specFiles, path[first:])
				return json.RawMessage(data), err
			},
		})
		if err != nil {
			return err
		}

		qualifiedType := parsed["provider"] + "/" + parsed["resourcetype"]
		key := getValidatorKey(qualifiedType, parsed["version"])
		l.validators[key] = validator{
			TypeName:        qualifiedType,
			APIVersion:      parsed["version"],
			specDoc:         wDoc,
			rootScopePrefix: l.rootScopePrefix,
			rootScopeParam:  l.rootScopeParam,
			paramCache:      make(map[string]map[string]spec.Parameter),
			paramCacheMu:    &sync.RWMutex{},
		}
		l.supportedVersions[qualifiedType] = append(l.supportedVersions[qualifiedType], parsed["version"])

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(l.validators) == 0 {
		return nil, ErrSpecDocumentNotFound
	}

	return l, nil
}

func getValidatorKey(resourceType, version string) string {
	return strings.ToLower(resourceType + "-" + version)
}

func parseSpecFilePath(path string) map[string]string {
	// OpenAPI specs are stored under swagger/ directory structure based on this spec - https://github.com/Azure/azure-rest-api-specs/wiki#directory-structure
	// This regex extracts the information from the filepath.
	re := regexp.MustCompile(`.*specification\/(?P<productname>.+)\/resource-manager\/(?P<provider>.+)\/(?P<state>.+)\/(?P<version>.+)\/(?P<resourcetype>.+)\.json$`)
	values := re.FindStringSubmatch(path)
	keys := re.SubexpNames()
	if len(keys) < 6 {
		return nil
	}

	d := map[string]string{}
	for i := 1; i < len(keys); i++ {
		d[keys[i]] = strings.ToLower(values[i])
	}
	return d
}
