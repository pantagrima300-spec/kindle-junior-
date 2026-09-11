package round2

type LanguageCode string

const (
	LangC      LanguageCode = "C"
	LangPython LanguageCode = "Python"
)

type TestCase struct {
	Input          string
	ExpectedOutput string
	IsHidden       bool
}

type LanguageSpecificContent struct {
	InitialCode string
}

type QuestionType string

const (
	TypeCaseBased QuestionType = "CASE STUDY"
	TypeDebugging QuestionType = "DEBUGGING"
)

type Question struct {
	ID          int
	Type        QuestionType
	Title       string
	Description string
	Languages   map[LanguageCode]LanguageSpecificContent
	TestCases   []TestCase
}

// ──────────────────────────────────────────────────────────────────────────────
// Debugging Starter Code (Cleaned - Hints Removed)
// ──────────────────────────────────────────────────────────────────────────────

// Q3 Grade & Attendance Eligibility Checker
var q3StarterC = `#include <stdio.h>

int main() {
    int attendance, marks;
    scanf("%d %d", &attendance, &marks);
    
    if (attendance >= 75 || marks >= 40) {
        printf("Pass");
    } else if (attendance < 75 && marks >= 80) {
        printf("Special Permission");
    } else {
        printf("Fail");
    }
    return 0;
}`

var q3StarterPy = `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    attendance, marks = map(int, line.split())

    if attendance >= 75 or marks >= 40:
        print("Pass")
    elif attendance < 75 and marks >= 80:
        print("Special Permission")
    else:
        print("Fail")
`

// Q4 Income Tax Slab Calculator
var q4StarterC = `#include <stdio.h>

int main() {
    double income, tax = 0;
    if (scanf("%lf", &income) == 1) {
        if (income > 500000) {
            tax = income * 0.20;
        } else if (income > 250000) {
            tax = income * 0.05;
        } else {
            tax = 0;
        }
        printf("%.2f\n", tax);
    }
    return 0;
}`

var q4StarterPy = `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    income = float(line)
    tax = 0.0

    if income > 500000:
        tax = income * 0.20
    elif income > 250000:
        tax = income * 0.05
    else:
        tax = 0.0

    print(f"{tax:.2f}")
`

// Q5 Max of Three Numbers
var q5StarterC = `#include <stdio.h>

void solve(int a, int b, int c) {
    int max = a;

    if (b > max) {
        max = b;
    }

    if (c < max) {
        max = c;
    }

    printf("%d\n", max);
}

int main() {
    int a, b, c;

    while (scanf("%d %d %d", &a, &b, &c) == 3) {
        solve(a, b, c);
    }

    return 0;
}`

var q5StarterPy = `import sys

def solve(a, b, c):
    max_val = a

    if b > max_val:
        max_val = b

    if c < max_val:
        max_val = c

    print(max_val)

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        parts = line.split()

        if len(parts) == 3:
            a, b, c = map(int, parts)
            solve(a, b, c)
`

// Q6 Electricity Bill Calculator
var q6StarterC = `#include <stdio.h>

void solve(int units) {
    float bill = 0;

    if (units <= 100) {
        bill = units * 5;
    } else if (units <= 200) {
        bill = (100 * 5) + ((units - 100) * 7);
    } else {
        bill = (100 * 5) + (100 * 7) + (units * 10);
    }

    printf("%.2f\n", bill);
}

int main() {
    int units;

    while (scanf("%d", &units) == 1) {
        solve(units);
    }

    return 0;
}`

var q6StarterPy = `import sys

def solve(units):
    if units <= 100:
        bill = units * 5
    elif units <= 200:
        bill = (100 * 5) + ((units - 100) * 7)
    else:
        bill = (100 * 5) + (100 * 7) + (units * 10)

    print(f"{bill:.2f}")

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        units = int(line)
        solve(units)
`

// Q7 Valid Triangle & Type
var q7StarterC = `#include <stdio.h>

void solve(int a, int b, int c) {
    if (a + b > c && a + c > b && b + c > a) {

        if (a == b && b == c) {
            printf("Equilateral\n");

        } else if (a == b || b == c) {
            printf("Isosceles\n");

        } else {
            printf("Scalene\n");
        }

    } else {
        printf("Invalid\n");
    }
}

int main() {
    int a, b, c;

    while (scanf("%d %d %d", &a, &b, &c) == 3) {
        solve(a, b, c);
    }

    return 0;
}`

var q7StarterPy = `import sys

def solve(a, b, c):
    if (a + b > c) and (a + c > b) and (b + c > a):

        if a == b == c:
            print("Equilateral")

        elif a == b or b == c:
            print("Isosceles")

        else:
            print("Scalene")

    else:
        print("Invalid")

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        parts = line.split()

        if len(parts) == 3:
            a, b, c = map(int, parts)
            solve(a, b, c)
`

// ──────────────────────────────────────────────────────────────────────────────
// Question Bank (IDs 1 through 7)
// ──────────────────────────────────────────────────────────────────────────────

var QuestionBank = []Question{
	// ID 1: Utility Water Tariff Calculator (Case Study)
	{
		ID:          1,
		Type:        TypeCaseBased,
		Title:       "Utility Water Tariff Calculator",
		Description: "Calculate the water tariff based on total water consumption. 0-100L is free. 101-300L costs ₹2/L for units over 100. 301-500L costs ₹400 plus ₹5/L for units over 300. Above 500L costs ₹1400 plus ₹8/L for units over 500, along with a ₹100 surcharge.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC: {
				InitialCode: `#include <stdio.h>

int main() {
    int litres;
    scanf("%d", &litres);

    // Write your solution here

    return 0;
}`,
			},
			LangPython: {
				InitialCode: `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    litres = int(line)

    # Write your solution here
`,
			},
		},
		TestCases: []TestCase{
			{Input: "80", ExpectedOutput: "0", IsHidden: false},
			{Input: "200", ExpectedOutput: "200", IsHidden: false},
			{Input: "400", ExpectedOutput: "900", IsHidden: false},
			{Input: "600", ExpectedOutput: "2300", IsHidden: false},
		},
	},

	// ID 2: Flight Baggage Allowance Auditor (Case Study)
	{
		ID:          2,
		Type:        TypeCaseBased,
		Title:       "Flight Baggage Allowance Auditor",
		Description: "Calculate the excess baggage charge based on baggage weight and flight type. Domestic flights allow up to 15kg free, charge ₹500/kg for 16-25kg over the allowance, and charge ₹5000 plus ₹1000/kg over 25kg above that. International flights allow up to 25kg free and charge ₹1200/kg for every kilogram over 25kg.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC: {
				InitialCode: `#include <stdio.h>
#include <string.h>

int main() {
    int weight;
    char type[20];

    scanf("%d %s", &weight, type);

    // Write your solution here

    return 0;
}`,
			},
			LangPython: {
				InitialCode: `weight, flight_type = input().split()
weight = int(weight)

# Write your solution here
`,
			},
		},
		TestCases: []TestCase{
			{Input: "12 DOMESTIC", ExpectedOutput: "0", IsHidden: false},
			{Input: "20 DOMESTIC", ExpectedOutput: "2500", IsHidden: false},
			{Input: "28 DOMESTIC", ExpectedOutput: "8000", IsHidden: false},
			{Input: "30 INTERNATIONAL", ExpectedOutput: "6000", IsHidden: false},
		},
	},

	// ID 3: Grade & Attendance Eligibility Checker (Debugging)
	{
		ID:          3,
		Type:        TypeDebugging,
		Title:       "Grade & Attendance Eligibility Checker",
		Description: "Attendance >= 75% AND Marks >= 40 results in 'Pass'. Attendance < 75% BUT Marks >= 80 results in 'Special Permission'. Otherwise, 'Fail'. Fix the logical bug in the condition.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC:      {InitialCode: q3StarterC},
			LangPython: {InitialCode: q3StarterPy},
		},
		TestCases: []TestCase{
			{Input: "80 50", ExpectedOutput: "Pass", IsHidden: false},
			{Input: "60 85", ExpectedOutput: "Special Permission", IsHidden: false},
			{Input: "60 50", ExpectedOutput: "Fail", IsHidden: false},
			{Input: "75 40", ExpectedOutput: "Pass", IsHidden: false},
			{Input: "80 35", ExpectedOutput: "Fail", IsHidden: true},
			{Input: "70 90", ExpectedOutput: "Special Permission", IsHidden: true},
		},
	},

	// ID 4: Income Tax Slab Calculator (Debugging)
	{
		ID:          4,
		Type:        TypeDebugging,
		Title:       "Income Tax Slab Calculator",
		Description: "Calculate tax on income: 0-250000 = 0%, 250001-500000 = 5% on amount over 250000, Above 500000 = 12500 + 20% on amount over 500000. Fix the flat-rate calculation bug.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC:      {InitialCode: q4StarterC},
			LangPython: {InitialCode: q4StarterPy},
		},
		TestCases: []TestCase{
			{Input: "200000", ExpectedOutput: "0.00", IsHidden: false},
			{Input: "400000", ExpectedOutput: "7500.00", IsHidden: false},
			{Input: "600000", ExpectedOutput: "32500.00", IsHidden: false},
			{Input: "500000", ExpectedOutput: "12500.00", IsHidden: false},
			{Input: "1000000", ExpectedOutput: "112500.00", IsHidden: true},
			{Input: "300000", ExpectedOutput: "2500.00", IsHidden: true},
		},
	},

	// ID 5: Max of Three Numbers (Debugging)
	{
		ID:          5,
		Type:        TypeDebugging,
		Title:       "Max of Three Numbers",
		Description: "Find the maximum of three integers. Fix the comparison bug in the code.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC:      {InitialCode: q5StarterC},
			LangPython: {InitialCode: q5StarterPy},
		},
		TestCases: []TestCase{
			{Input: "3 7 5", ExpectedOutput: "7", IsHidden: false},
			{Input: "10 10 10", ExpectedOutput: "10", IsHidden: false},
			{Input: "1 2 3", ExpectedOutput: "3", IsHidden: false},
			{Input: "-1 -2 -3", ExpectedOutput: "-1", IsHidden: false},
			{Input: "100 50 75", ExpectedOutput: "100", IsHidden: true},
			{Input: "0 0 1", ExpectedOutput: "1", IsHidden: true},
		},
	},

	// ID 6: Electricity Bill Calculator (Debugging)
	{
		ID:          6,
		Type:        TypeDebugging,
		Title:       "Electricity Bill Calculator",
		Description: "Calculate the electricity bill based on units consumed. Fix the calculation bug for high consumption.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC:      {InitialCode: q6StarterC},
			LangPython: {InitialCode: q6StarterPy},
		},
		TestCases: []TestCase{
			{Input: "50", ExpectedOutput: "250.00", IsHidden: false},
			{Input: "150", ExpectedOutput: "850.00", IsHidden: false},
			{Input: "250", ExpectedOutput: "2200.00", IsHidden: false},
			{Input: "100", ExpectedOutput: "500.00", IsHidden: false},
			{Input: "200", ExpectedOutput: "1200.00", IsHidden: true},
			{Input: "300", ExpectedOutput: "3200.00", IsHidden: true},
		},
	},

	// ID 7: Valid Triangle & Type (Debugging)
	{
		ID:          7,
		Type:        TypeDebugging,
		Title:       "Valid Triangle & Type",
		Description: "Check if a triangle is valid and classify it as Equilateral, Isosceles, or Scalene. Fix the logic so that invalid side combinations and types are properly identified.",
		Languages: map[LanguageCode]LanguageSpecificContent{
			LangC:      {InitialCode: q7StarterC},
			LangPython: {InitialCode: q7StarterPy},
		},
		TestCases: []TestCase{
			{Input: "5 5 5", ExpectedOutput: "Equilateral", IsHidden: false},
			{Input: "1 2 10", ExpectedOutput: "Invalid", IsHidden: false},
			{Input: "3 4 5", ExpectedOutput: "Scalene", IsHidden: false},
			{Input: "5 5 8", ExpectedOutput: "Isosceles", IsHidden: false},
			{Input: "5 8 5", ExpectedOutput: "Isosceles", IsHidden: true},
			{Input: "8 5 5", ExpectedOutput: "Isosceles", IsHidden: true},
		},
	},
}

func GetQuestion(id int) *Question {
	for i := range QuestionBank {
		if QuestionBank[i].ID == id {
			return &QuestionBank[i]
		}
	}

	return nil
}