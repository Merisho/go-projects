package problem_template

import (
	"io/ioutil"
)

func ProblemTemplateWithMultipleTests(name string) error {
	fname := name + ".cpp"

	tmp := `#include <bits/stdc++.h>
using namespace std;
using ll = long long;

int main() {
	int T;
	cin >> T;

	for (int test_case = 1; test_case <= T; ++test_case) {
		
	}
	
	return 0;
}
`
	err := ioutil.WriteFile(fname, []byte(tmp), 0777)
	return err
}

func ProblemTemplate(name string) error {
	fname := name + ".cpp"

	tmp := `#include <bits/stdc++.h>
using namespace std;
using ll = long long;

int main() {
	
	
	return 0;
}
`
	err := ioutil.WriteFile(fname, []byte(tmp), 0777)
	return err
}
