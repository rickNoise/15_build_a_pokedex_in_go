package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "     HELLO   WoRlD ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "three words here",
			expected: []string{"three", "words", "here"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Errorf("cleanInput output is not the right length")
        }
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("wrong word: expected %v, got %v", expectedWord, word)
			}
		}
	}
}
