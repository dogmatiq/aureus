package loader

import (
	"fmt"
	"slices"
	"strings"
)

// ParseRequirements parses the comma-separated names in a requires attribute.
func ParseRequirements(value string) ([]string, error) {
	var names []string

	for part := range strings.SplitSeq(value, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			return nil, fmt.Errorf("fixture requirements must not contain an empty name")
		}

		names = append(names, name)
	}

	return names, nil
}

// resolveFixtures includes the pair's requirements in dependency order.
func resolveFixtures(
	definitions map[string]ContentEnvelope,
	groupFixtures []ContentEnvelope,
	input, output ContentEnvelope,
) ([]ContentEnvelope, error) {
	var resolved []ContentEnvelope
	included := make(map[string]bool)
	var path []string

	// Visit requirements before including the fixture that declares them.
	// The active path distinguishes a cycle from a shared dependency.
	var visit func(ContentEnvelope) error
	visit = func(source ContentEnvelope) error {
		for _, name := range source.Content.Requires {
			fixture, ok := definitions[name]
			if !ok {
				return fmt.Errorf("%s requires unknown fixture %q", location(source, true), name)
			}

			if group := groupName(fixture.Content.Group); group != "" {
				return fmt.Errorf(
					"%s requires fixture %q, which belongs to group %q",
					location(source, true),
					name,
					group,
				)
			}

			if included[name] {
				continue
			}

			if slices.Contains(path, name) {
				return fmt.Errorf(
					"fixture dependency cycle at %s: %s -> %s",
					location(source, true),
					strings.Join(path, " -> "),
					name,
				)
			}

			path = append(path, name)
			if err := visit(fixture); err != nil {
				return err
			}
			path = path[:len(path)-1]

			resolved = append(resolved, fixture)
			included[name] = true
		}

		return nil
	}

	// Group fixtures are included directly; requirements can only refer to
	// ungrouped fixtures, so no group fixture can have been included already.
	for _, fixture := range groupFixtures {
		if err := visit(fixture); err != nil {
			return nil, err
		}

		resolved = append(resolved, fixture)
	}

	for _, source := range []ContentEnvelope{input, output} {
		if err := visit(source); err != nil {
			return nil, err
		}
	}

	return resolved, nil
}
