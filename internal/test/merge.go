package test

// Merge recursively combines tests with the same name into a single test.
func Merge(tests ...Test) []Test {
	var groups [][]Test

next:
	for _, t := range tests {
		for i, g := range groups {
			if g[0].Name == t.Name {
				groups[i] = append(g, t)
				continue next
			}
		}

		groups = append(groups, []Test{t})
	}

	var result []Test
	for _, g := range groups {
		result = append(result, mergeGroup(g))
	}

	return result
}

// mergeGroup merges a group of tests with the same name into a single test.
func mergeGroup(group []Test) Test {
	var skipped, unskipped []Test

	for _, t := range group {
		if t.Skip {
			skipped = append(skipped, t)
		} else {
			unskipped = append(unskipped, t)
		}
	}

	if len(skipped) == 0 {
		return mergeNaive(unskipped)
	} else if len(unskipped) == 0 {
		return mergeNaive(skipped)
	}

	// If we have both skipped and unskipped tests we can't flatten them, so
	// instead we return a new test that contains all of the unskipped tests
	// from the group, and a single skipped test which in turn contains the
	// skipped tests from the group.

	s := mergeNaive(skipped)
	s.Name += " (skipped)"

	u := mergeNaive(unskipped)
	u.SubTests = append(u.SubTests, s)

	return u
}

// mergeNaive returns a single test containing all of the sub-tests and
// assertions from the given tests.
//
// It uses the name and meta-data of the first test.
func mergeNaive(tests []Test) Test {
	if len(tests) == 0 {
		panic("no tests to merge")
	}

	if len(tests) == 1 {
		return tests[0]
	}

	var (
		subTests   []Test
		assertions []Assertion
	)

	for _, t := range tests {
		subTests = append(subTests, t.SubTests...)
		assertions = append(assertions, t.Assertions...)
	}

	return Test{
		Name:       tests[0].Name,
		Skip:       tests[0].Skip,
		SubTests:   Merge(subTests...),
		Assertions: assertions,
	}
}
