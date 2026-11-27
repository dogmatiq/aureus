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
	input, output ContentEnvelope,
) ([]ContentEnvelope, error) {
	var resolved []ContentEnvelope
	included := make(map[string]bool)
	var path []string

	// Completed fixtures are included once. The active path distinguishes a
	// cycle from a dependency shared by several fixtures.
	var visit func(string, ContentEnvelope) error
	visit = func(name string, source ContentEnvelope) error {
		fixture, ok := definitions[name]
		if !ok {
			return fmt.Errorf("%s requires unknown fixture %q", location(source, true), name)
		}

		if included[name] {
			return nil
		}

		if slices.Contains(path, name) {
			cycle := append(slices.Clone(path), name)
			return fmt.Errorf(
				"fixture dependency cycle at %s: %s",
				location(source, true),
				strings.Join(cycle, " -> "),
			)
		}

		path = append(path, name)

		for _, dependency := range fixture.Content.Requires {
			if err := visit(dependency, fixture); err != nil {
				return err
			}
		}

		path = path[:len(path)-1]
		included[name] = true
		resolved = append(resolved, fixture)
		return nil
	}

	for _, source := range []ContentEnvelope{input, output} {
		for _, name := range source.Content.Requires {
			if err := visit(name, source); err != nil {
				return nil, err
			}
		}
	}

	return resolved, nil
}
