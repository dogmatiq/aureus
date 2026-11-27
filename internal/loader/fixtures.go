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

	// Group fixtures are included directly. References may only select
	// ungrouped fixtures, even if a grouped fixture is already included.
	var include func(ContentEnvelope) error
	require := func(name string, source ContentEnvelope) error {
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

		if slices.Contains(path, name) {
			cycle := append(slices.Clone(path), name)
			return fmt.Errorf(
				"fixture dependency cycle at %s: %s",
				location(source, true),
				strings.Join(cycle, " -> "),
			)
		}

		return include(fixture)
	}

	include = func(fixture ContentEnvelope) error {
		name := fixture.Content.FixtureName
		if included[name] {
			return nil
		}

		path = append(path, name)

		for _, dependency := range fixture.Content.Requires {
			if err := require(dependency, fixture); err != nil {
				return err
			}
		}

		path = path[:len(path)-1]
		included[name] = true
		resolved = append(resolved, fixture)
		return nil
	}

	for _, fixture := range groupFixtures {
		if err := include(fixture); err != nil {
			return nil, err
		}
	}

	for _, source := range []ContentEnvelope{input, output} {
		for _, name := range source.Content.Requires {
			if err := require(name, source); err != nil {
				return nil, err
			}
		}
	}

	return resolved, nil
}
