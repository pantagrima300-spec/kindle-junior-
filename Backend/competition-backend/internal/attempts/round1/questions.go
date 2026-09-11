package round1

import (
	"errors"
	"strings"
)

// Question is the authoritative Round 1 question model.
// CorrectAnswer is kept server-side and must never be returned to participants.
type Question struct {
	ID            int      `json:"id"`
	Domain        string   `json:"domain"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"-"`
}

// QuestionBank contains the 60 authoritative Round 1 questions for C.
// CorrectAnswer uses zero-based indexes matching the option arrays.
var CQuestionBank = []Question{
	{
		ID:       1,
		Domain:   "C",
		Question: "C is known as a __________ language because compiled code requires recompilation to run on a different operating system or hardware architecture.",
		Options: []string{
			"Platform-independent",
			"Platform-dependent",
			"Interpreter-based",
			"None of the above",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       2,
		Domain:   "C",
		Question: "To write and run a C program, the source code must first be translated into machine code by a ___________.",
		Options: []string{
			"Compiler",
			"Interpreter",
			"Linker",
			"Loader",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       3,
		Domain:   "C",
		Question: "The process of combining various object files and library functions into a single executable file in C is performed by the _______.",
		Options: []string{
			"Preprocessor",
			"Compiler",
			"Linker",
			"Loader",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       4,
		Domain:   "C",
		Question: "Before a C program is compiled, lines beginning with the # symbol (like #include) are processed by the __________.",
		Options: []string{
			"Interpreter",
			"Linker",
			"Preprocessor",
			"Assembler",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       5,
		Domain:   "C",
		Question: "A _________ produces an undesired output but without abrupt termination of the execution of the program.",
		Options: []string{
			"Syntax error",
			"Logical error",
			"Runtime error",
			"Linker error",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       6,
		Domain:   "C",
		Question: "A ___________ causes abnormal termination of the program while it is executing (for example, division by zero).",
		Options: []string{
			"Syntax error",
			"Logical error",
			"Runtime error",
			"Semantic error",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       7,
		Domain:   "C",
		Question: "Missing a semicolon (;) at the end of a statement in C will result in a __________.",
		Options: []string{
			"Syntax error",
			"Logical error",
			"Runtime error",
			"Linker error",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       8,
		Domain:   "C",
		Question: "The process of identifying and removing errors from a computer program is called _______.",
		Options: []string{
			"Debugging",
			"Compiling",
			"Linking",
			"Executing",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       9,
		Domain:   "C",
		Question: "Every C program must contain a _________ function, which acts as the starting point of execution.",
		Options: []string{
			"start()",
			"main()",
			"execute()",
			"init()",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       10,
		Domain:   "C",
		Question: "Comments in a C program (written as // or /* */) are __________.",
		Options: []string{
			"Non-executable statements",
			"Executable statements",
			"Processed by the linker",
			"Converted into machine code",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       11,
		Domain:   "C",
		Question: "C uses __________ to define blocks of code (such as the body of a loop, if-statement, or function).",
		Options: []string{
			"Indentation",
			"Parentheses ()",
			"Curly braces {}",
			"Square brackets []",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       12,
		Domain:   "C",
		Question: "__________ is a user-defined name given to a variable, function, or other entity in a program.",
		Options: []string{
			"Keyword",
			"Identifier",
			"Data type",
			"Token",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       13,
		Domain:   "C",
		Question: "Which of the following is NOT a valid fundamental data type in C?",
		Options: []string{
			"int",
			"float",
			"char",
			"string",
		},
		CorrectAnswer: 3,
	},
	{
		ID:       14,
		Domain:   "C",
		Question: "In C, variables whose values cannot be changed after they are initialized must be declared using the __________ keyword.",
		Options: []string{
			"constant",
			"immutable",
			"const",
			"static",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       15,
		Domain:   "C",
		Question: "_________ in C is a data structure used to store a fixed-size sequential collection of elements of the same data type.",
		Options: []string{
			"Structure",
			"Array",
			"Pointer",
			"Enum",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       16,
		Domain:   "C",
		Question: "___________ in C is a user-defined data type that can hold multiple variables of different data types under a single name.",
		Options: []string{
			"Array",
			"Structure",
			"Pointer",
			"Enum",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       17,
		Domain:   "C",
		Question: "Which format specifier is used to print a standard decimal integer value in C?",
		Options: []string{
			"%d",
			"%f",
			"%c",
			"%s",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       18,
		Domain:   "C",
		Question: "What is the standard size of the char data type in C?",
		Options: []string{
			"1 byte",
			"2 bytes",
			"4 bytes",
			"8 bytes",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       19,
		Domain:   "C",
		Question: "Which function is used to read formatted input from the standard input (keyboard) in C?",
		Options: []string{
			"printf()",
			"scanf()",
			"get()",
			"read()",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       20,
		Domain:   "C",
		Question: "Which unformatted console function is used to read a single character from the standard input?",
		Options: []string{
			"getchar()",
			"putchar()",
			"scanf()",
			"printf()",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       21,
		Domain:   "C",
		Question: "What does the modulus operator % yield?",
		Options: []string{
			"The quotient of integer division",
			"The fractional part of float division",
			"The remainder after integer division",
			"The square root of a number",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       22,
		Domain:   "C",
		Question: "Which operator has the highest precedence in C?",
		Options: []string{
			"||",
			"&&",
			"! (Logical NOT)",
			"==",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       23,
		Domain:   "C",
		Question: "What is the value of b after this code executes: int a = 0, b = 5; if (a && ++b) {}?",
		Options: []string{
			"4",
			"5",
			"6",
			"0",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       24,
		Domain:   "C",
		Question: "What is the value of y after this code executes: int x = 1, y = 10; if (x || ++y) {}?",
		Options: []string{
			"9",
			"10",
			"11",
			"1",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       25,
		Domain:   "C",
		Question: "What does the expression 1 || 0 && 0 evaluate to in C?",
		Options: []string{
			"1 (True)",
			"0 (False)",
			"Compile Error",
			"Garbage Value",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       26,
		Domain:   "C",
		Question: "What is the output of the expression 5 > 4 > 3 in C?",
		Options: []string{
			"1 (True)",
			"0 (False)",
			"Compile Error",
			"Runtime Error",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       27,
		Domain:   "C",
		Question: "What is the exact integer output of the expression 10 + 5 * 2 % 3?",
		Options: []string{
			"10",
			"11",
			"15",
			"0",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       28,
		Domain:   "C",
		Question: "What does the expression !!5 evaluate to in C?",
		Options: []string{
			"0",
			"1",
			"5",
			"-5",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       29,
		Domain:   "C",
		Question: "Which operator is used for Logical AND?",
		Options: []string{
			"&",
			"&&",
			"|",
			"||",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       30,
		Domain:   "C",
		Question: "What is the output of the following code: int a = 5; if (a = 0) printf(\"Yes\"); else printf(\"No\");?",
		Options: []string{
			"Yes",
			"No",
			"Compile Error",
			"Garbage Value",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       31,
		Domain:   "C",
		Question: "How do you rewrite if (x > y) max = x; else max = y; using the conditional (ternary) operator?",
		Options: []string{
			"max = (x > y) ? x : y;",
			"max = (x > y) : x ? y;",
			"max = (x > y) ? y : x;",
			"max = (x > y) & x : y;",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       32,
		Domain:   "C",
		Question: "Which operator is used for Bitwise OR?",
		Options: []string{
			"|",
			"||",
			"&",
			"^",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       33,
		Domain:   "C",
		Question: "What is the value of x << 1 (Bitwise Left Shift) if x = 4?",
		Options: []string{
			"2",
			"4",
			"8",
			"16",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       34,
		Domain:   "C",
		Question: "What is the result of 2 & 3 (Bitwise AND) in C?",
		Options: []string{
			"1",
			"2",
			"3",
			"5",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       35,
		Domain:   "C",
		Question: "The comma operator , evaluates from left to right and returns the value of the ________ operand.",
		Options: []string{
			"Leftmost",
			"Rightmost",
			"Largest",
			"Smallest",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       36,
		Domain:   "C",
		Question: "For a leap year calculation, which conditional expression correctly isolates century years?",
		Options: []string{
			"year / 100 == 0",
			"year % 100 == 0",
			"year % 400 == 0",
			"year / 400 == 0",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       37,
		Domain:   "C",
		Question: "If x < 0 and y < 0, which logical operator should combine these to confirm the point is in the Third Quadrant?",
		Options: []string{
			"!",
			"^",
			"||",
			"&&",
		},
		CorrectAnswer: 3,
	},
	{
		ID:       38,
		Domain:   "C",
		Question: "In C, to which if does an else statement attach if there are no curly braces?",
		Options: []string{
			"The very first if",
			"The nearest preceding if",
			"It causes a syntax error",
			"Both ifs",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       39,
		Domain:   "C",
		Question: "Which data type cannot be evaluated inside a C switch statement?",
		Options: []string{
			"int",
			"char",
			"float",
			"short",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       40,
		Domain:   "C",
		Question: "If case 1: lacks a break; statement, what happens when the variable equals 1?",
		Options: []string{
			"The program crashes",
			"It skips the rest of the switch",
			"It executes case 1 and falls through to execute the next case",
			"It loops infinitely",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       41,
		Domain:   "C",
		Question: "Is the default case mandatory in a switch statement?",
		Options: []string{
			"Yes, always",
			"No, it is optional",
			"Yes, but only for integers",
			"No, unless break is omitted",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       42,
		Domain:   "C",
		Question: "What is the output of this code: int x = 10; if (1) { int x = 20; } printf(\"%d\", x);?",
		Options: []string{
			"10",
			"20",
			"30",
			"Compile Error",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       43,
		Domain:   "C",
		Question: "What happens when if (5 > 3); printf(\"Hello\"); is executed?",
		Options: []string{
			"Prints Hello",
			"Syntax Error",
			"No output",
			"Infinite Loop",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       44,
		Domain:   "C",
		Question: "Which conditional correctly guards against division by zero for z = x / y?",
		Options: []string{
			"if (y == 0)",
			"if (y != 0)",
			"if (y > 0)",
			"if (x != 0)",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       45,
		Domain:   "C",
		Question: "If an electricity bill slab applies $10 for units > 200 and $5 for units > 100, why must the > 200 condition be checked first in an if-else if ladder?",
		Options: []string{
			"Syntax rules require descending order",
			"Checking > 100 first would accidentally catch values like 250",
			"The C compiler optimizes larger numbers",
			"It does not matter",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       46,
		Domain:   "C",
		Question: "What are the three components of a for loop declaration separated by?",
		Options: []string{
			"Commas ,",
			"Semicolons ;",
			"Colons :",
			"Spaces",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       47,
		Domain:   "C",
		Question: "A do-while loop is guaranteed to execute its block of code at least how many times?",
		Options: []string{
			"0",
			"1",
			"2",
			"Depends on the condition",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       48,
		Domain:   "C",
		Question: "What does the break statement do inside a loop?",
		Options: []string{
			"Restarts the loop",
			"Skips the current iteration",
			"Immediately terminates the loop",
			"Pauses execution",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       49,
		Domain:   "C",
		Question: "What does the continue statement do inside a loop?",
		Options: []string{
			"Terminates the program",
			"Terminates the loop",
			"Skips the remaining code in the current iteration and moves to the next iteration",
			"Ignores all future errors",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       50,
		Domain:   "C",
		Question: "Which of the following creates an infinite loop in C?",
		Options: []string{
			"for(;;)",
			"while(1)",
			"do {} while(1);",
			"All of the above",
		},
		CorrectAnswer: 3,
	},
	{
		ID:       51,
		Domain:   "C",
		Question: "Given int x = 0; while (x++ < 3) {}, what is the value of x when the loop finally terminates?",
		Options: []string{
			"2",
			"3",
			"4",
			"5",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       52,
		Domain:   "C",
		Question: "What is the exact output of this code: int a = 1; printf(\"%d\", ++a);?",
		Options: []string{
			"1",
			"2",
			"3",
			"Compile Error",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       53,
		Domain:   "C",
		Question: "If int i = 0; is declared inside the parentheses of a for loop in C99, can it be accessed outside the loop block?",
		Options: []string{
			"Yes",
			"No",
			"Only if it's an infinite loop",
			"Yes, but only as a read-only variable",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       54,
		Domain:   "C",
		Question: "What happens if you forget to increment the counter variable inside a while loop that relies on it to terminate?",
		Options: []string{
			"Compile Error",
			"Syntax Error",
			"Infinite Loop",
			"Loop skips entirely",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       55,
		Domain:   "C",
		Question: "In C, any non-zero integer evaluated in a loop condition is treated as what boolean concept?",
		Options: []string{
			"False",
			"True",
			"Null",
			"Void",
		},
		CorrectAnswer: 1,
	},
	{
		ID:       56,
		Domain:   "C",
		Question: "What does this loop print? for(int i=1; i<=4; i++) { if(i%2==0) printf(\"%d\", i); }",
		Options: []string{
			"24",
			"13",
			"1234",
			"024",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       57,
		Domain:   "C",
		Question: "How many times will int i = 3; while(i > 0) { i--; } execute its body?",
		Options: []string{
			"2",
			"4",
			"3",
			"0",
		},
		CorrectAnswer: 2,
	},
	{
		ID:       58,
		Domain:   "C",
		Question: "What is the final value of sum? int sum=0; for(int i=1; i<=3; i++) { sum = sum + i; }",
		Options: []string{
			"6",
			"3",
			"5",
			"10",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       59,
		Domain:   "C",
		Question: "Is for (int i = 0, j = 10; i < j; i++, j--) a valid C statement?",
		Options: []string{
			"Yes",
			"No",
			"Only in C++",
			"Throws a syntax error",
		},
		CorrectAnswer: 0,
	},
	{
		ID:       60,
		Domain:   "C",
		Question: "What happens here: int i; for(i=0; i<3; i++); printf(\"%d\", i);?",
		Options: []string{
			"Prints 012",
			"Prints 0123",
			"Prints 3",
			"Compile Error",
		},
		CorrectAnswer: 2,
	},
}

// PythonQuestionBank contains the 60 authoritative Round 1 questions for Python.
var PythonQuestionBank = []Question{
	{
		ID:            1,
		Domain:        "Python",
		Question:      "An ordered set of instructions to be executed by a computer to carry out a specific task is called a ___________.",
		Options:       []string{"Program", "Instruction", "Code", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            2,
		Domain:        "Python",
		Question:      "Computers understand the language of 0s and 1s which is called __________.",
		Options:       []string{"Machine language", "Low level language", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            3,
		Domain:        "Python",
		Question:      "A program written in a high-level language is called _________.",
		Options:       []string{"Language", "Source code", "Machine code", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            4,
		Domain:        "Python",
		Question:      "An interpreter read the program statements _________.",
		Options:       []string{"All the source at a time", "One by one", "Both a) and b)", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            5,
		Domain:        "Python",
		Question:      "Python is a __________.",
		Options:       []string{"Low level language", "High level language", "Machine level language", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            6,
		Domain:        "Python",
		Question:      "Python is platform independent, meaning ___________.",
		Options:       []string{"It can run various operating systems", "It can run various hardware platforms", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            7,
		Domain:        "Python",
		Question:      "Python uses indentation for __________.",
		Options:       []string{"Blocks", "Nested Blocks", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            8,
		Domain:        "Python",
		Question:      "To write and run (execute) a Python program, we need to have a ___________.",
		Options:       []string{"Python interpreter installed on the computer", "We can use any online python interpreter", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            9,
		Domain:        "Python",
		Question:      "The interpreter is also called python _______.",
		Options:       []string{"Shell", "Cell", "Program", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            10,
		Domain:        "Python",
		Question:      "To work in the interactive mode, we can simply type a Python statement on the ________ prompt directly.",
		Options:       []string{">>>", ">>", ">", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            11,
		Domain:        "Python",
		Question:      "In the script mode, we can write a Python program in a ________, save it and then use the interpreter to execute it.",
		Options:       []string{"Prompt", "File", "Folder", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            12,
		Domain:        "Python",
		Question:      "By default the python extension is __________.",
		Options:       []string{".py", ".ppy", ".pp", ".pyy"},
		CorrectAnswer: 0,
	},
	{
		ID:            13,
		Domain:        "Python",
		Question:      "__________ are reserved words in python.",
		Options:       []string{"Keyword", "Interpreter", "Program", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:       14,
		Domain:   "Python",
		Question: "The rules for naming an identifier in Python are ________.",
		Options: []string{
			"Name should begin with an uppercase, lowercase or underscore",
			"It can be of any length",
			"It should not be a keyword",
			"All of the above",
		},
		CorrectAnswer: 3,
	},
	{
		ID:            15,
		Domain:        "Python",
		Question:      "To define variables in python _______ special symbols is not allowed.",
		Options:       []string{"@", "# and !", "$ and %", "All of the above"},
		CorrectAnswer: 3,
	},
	{
		ID:            16,
		Domain:        "Python",
		Question:      "A variable in a program is uniquely identified by a name __________.",
		Options:       []string{"Identifier", "Keyword", "Code", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            17,
		Domain:        "Python",
		Question:      "Variable in python refers to an ________.",
		Options:       []string{"Keyword", "Object", "Alphabets", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            18,
		Domain:        "Python",
		Question:      "The variable message holds string type value and so its content is assigned within _________.",
		Options:       []string{"Double quotes \"\"", "Single quotes ''", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            19,
		Domain:        "Python",
		Question:      "_________ must always be assigned values before they are used in expressions.",
		Options:       []string{"Keyword", "Variable", "Code", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            20,
		Domain:        "Python",
		Question:      "___________ are used to add a remark or a note in the source code.",
		Options:       []string{"Keyword", "Source", "Comment", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            21,
		Domain:        "Python",
		Question:      "__________ are not executed by a python interpreter.",
		Options:       []string{"Keyword", "Source", "Comment", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            22,
		Domain:        "Python",
		Question:      "In python comments start from _______.",
		Options:       []string{"#", "@", "%", "$"},
		CorrectAnswer: 0,
	},
	{
		ID:            23,
		Domain:        "Python",
		Question:      "Python treats every value or data item whether numeric, string, or other type as an _________.",
		Options:       []string{"Object", "Variable", "Keyword", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            24,
		Domain:        "Python",
		Question:      "________ data type cannot have duplicate entries.",
		Options:       []string{"List", "Set", "String", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            25,
		Domain:        "Python",
		Question:      "None is a special data type with a single value. It is used to signify the absence of value in a situation.",
		Options:       []string{"List", "Set", "String", "None"},
		CorrectAnswer: 3,
	},
	{
		ID:            26,
		Domain:        "Python",
		Question:      "___________ in Python holds data items in key-value pairs.",
		Options:       []string{"Dictionary", "Set", "String", "None"},
		CorrectAnswer: 0,
	},
	{
		ID:            27,
		Domain:        "Python",
		Question:      "Items in a dictionary are enclosed in __________.",
		Options:       []string{"Parenthesis ()", "Brackets []", "Curly brackets {}", "All of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            28,
		Domain:        "Python",
		Question:      "In the dictionary every key is separated from its value using a _________.",
		Options:       []string{"Colon (:)", "Semicolon (;)", "Comma (,)", "All of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            29,
		Domain:        "Python",
		Question:      "Variables whose values can be changed after they are created and assigned are called __________.",
		Options:       []string{"Immutable", "Mutable", "Both a) and b)", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            30,
		Domain:        "Python",
		Question:      "____________ that are used to perform the four basic arithmetic operations as well as modular division, floor division and exponentiation.",
		Options:       []string{"Arithmetic Operator", "Logical Operator", "Relational Operator", "All of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            31,
		Domain:        "Python",
		Question:      "__________ calculation on operands. That is, raise the operand on the left to the power of the operand on the right.",
		Options:       []string{"Floor Division (//)", "Exponent (**)", "Modulus (%)", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            32,
		Domain:        "Python",
		Question:      "_____________ compares the values of the operands on either side and determines the relationship among them.",
		Options:       []string{"Arithmetic Operator", "Logical Operator", "Relational Operator", "All of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            33,
		Domain:        "Python",
		Question:      "_________ assigns or changes the value of the variable on its left.",
		Options:       []string{"Relational Operator", "Assignment Operator", "Logical Operator", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            34,
		Domain:        "Python",
		Question:      "Membership operators are used to check if a value is a member of the given sequence or not.",
		Options:       []string{"Identity Operator", "Membership Operators", "Relational Operators", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            35,
		Domain:        "Python",
		Question:      "An __________ is defined as a combination of constants, variables, and operators.",
		Options:       []string{"Expressions", "Precedence", "Both a) and b)", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            36,
		Domain:        "Python",
		Question:      "Evaluation of the expression is based on __________ of operators.",
		Options:       []string{"Expressions", "Precedence", "Both a) and b)", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            37,
		Domain:        "Python",
		Question:      "Example of membership operators.",
		Options:       []string{"in", "not in", "in and not in", "All of the above"},
		CorrectAnswer: 3,
	},
	{
		ID:            38,
		Domain:        "Python",
		Question:      "Example of identity operators.",
		Options:       []string{"is", "is not", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            39,
		Domain:        "Python",
		Question:      "In Python, we have the ____________ function for taking the user input.",
		Options:       []string{"prompt()", "input()", "in()", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            40,
		Domain:        "Python",
		Question:      "In Python, we have the ___________ function for displaying the output.",
		Options:       []string{"prompt()", "output()", "print()", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            41,
		Domain:        "Python",
		Question:      "____________, also called type casting, happens when data type conversion takes place because the programmer forced it in the program.",
		Options:       []string{"Explicit conversion", "Implicit conversion", "Both a) and b)", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            42,
		Domain:        "Python",
		Question:      "__________, also known as coercion, happens when data type conversion is done automatically by Python and is not instructed by the programmer.",
		Options:       []string{"Explicit conversion", "Implicit conversion", "Both a) and b)", "None of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            43,
		Domain:        "Python",
		Question:      "A programmer can make mistakes while writing a program, and hence, the program may not execute or may generate wrong output. The process of identifying and removing such mistakes is also known as __________.",
		Options:       []string{"Bugs", "Errors", "Both a) and b)", "None of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            44,
		Domain:        "Python",
		Question:      "Identifying and removing bugs or errors from the program is also known as __________.",
		Options:       []string{"Debugging", "Mistakes", "Error", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            45,
		Domain:        "Python",
		Question:      "Which of the following errors occur in python programs.",
		Options:       []string{"Syntax error", "Logical error", "Runtime error", "All of the above"},
		CorrectAnswer: 3,
	},
	{
		ID:            46,
		Domain:        "Python",
		Question:      "A _________ produces an undesired output but without abrupt termination of the execution of the program.",
		Options:       []string{"Syntax error", "Logical error", "Runtime error", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            47,
		Domain:        "Python",
		Question:      "A ___________ causes abnormal termination of the program while it is executing.",
		Options:       []string{"Syntax error", "Logical error", "Runtime error", "All of the above"},
		CorrectAnswer: 2,
	},
	{
		ID:            48,
		Domain:        "Python",
		Question:      "Python is __________ language that can be used for a multitude of scientific and non-scientific computing purposes.",
		Options:       []string{"Open-source", "High level", "Interpreter-based", "All of the above"},
		CorrectAnswer: 3,
	},
	{
		ID:            49,
		Domain:        "Python",
		Question:      "Comments are __________ statements in a program.",
		Options:       []string{"Non-executable", "Executable", "Both a) and b)", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            50,
		Domain:        "Python",
		Question:      "__________ is a user defined name given to a variable or a constant in a program.",
		Options:       []string{"Keyword", "Identifier", "Data type", "All of the above"},
		CorrectAnswer: 1,
	},
	{
		ID:            51,
		Domain:        "Python",
		Question:      "The process of identifying and removing errors from a computer program is called _______.",
		Options:       []string{"Debugging", "Mistakes", "Error", "None of the above"},
		CorrectAnswer: 0,
	},
	{
		ID:            52,
		Domain:        "Python",
		Question:      "Which of the following is an invalid variable naming syntax in Python?",
		Options:       []string{"_my_var = 10", "myVar_2 = 10", "2nd_var = 10", "MYVAR = 10"},
		CorrectAnswer: 2,
	},
	{
		ID:            53,
		Domain:        "Python",
		Question:      "What is the output of the following string concatenation syntax: print(\"5\" + \"5\")?",
		Options:       []string{"10", "\"55\"", "55 (as integer)", "Syntax Error"},
		CorrectAnswer: 1,
	},
	{
		ID:            54,
		Domain:        "Python",
		Question:      "What is the output of the following boolean syntax: print(bool(\"\"), bool(\"False\"))?",
		Options:       []string{"False, False", "True, True", "False, True", "True, False"},
		CorrectAnswer: 2,
	},
	{
		ID:            55,
		Domain:        "Python",
		Question:      "Which of the following is the correct syntax for floor division in Python?",
		Options:       []string{"/", "//", "%", "\\"},
		CorrectAnswer: 1,
	},
	{
		ID:            56,
		Domain:        "Python",
		Question:      "What is the output of the chained comparison syntax: print(1 < 5 < 10)?",
		Options:       []string{"True", "False", "Compile Error", "1"},
		CorrectAnswer: 0,
	},
	{
		ID:            57,
		Domain:        "Python",
		Question:      "Which keyword is used to test if a specific character or substring exists within a string?",
		Options:       []string{"exists", "has", "in", "contains"},
		CorrectAnswer: 2,
	},
	{
		ID:            58,
		Domain:        "Python",
		Question:      "What is the output of the range() function syntax in this loop: for i in range(2, 5): print(i, end=\"\")?",
		Options:       []string{"2345", "234", "345", "24"},
		CorrectAnswer: 1,
	},
	{
		ID:            59,
		Domain:        "Python",
		Question:      "What is the result of the following string slicing syntax: print(\"Python\"[-3:])?",
		Options:       []string{"Pyt", "hon", "tho", "nohtyP"},
		CorrectAnswer: 1,
	},
	{
		ID:            60,
		Domain:        "Python",
		Question:      "What error is produced by the following variable unpacking syntax: x, y = 5, 10, 15?",
		Options:       []string{"TypeError", "SyntaxError", "ValueError", "NameError"},
		CorrectAnswer: 2,
	},
}

// ErrInvalidLanguage indicates an unsupported domain was requested.
var ErrInvalidLanguage = errors.New("invalid language specified")

// ErrQuestionOutOfRange indicates the question number is outside 1..len(bank).
var ErrQuestionOutOfRange = errors.New("invalid question number: out of bounds")

func GetQuestion(language string, questionNumber int) *Question {
	language = strings.ToLower(strings.TrimSpace(language))

	var bank []Question

	switch language {
	case "python":
		bank = PythonQuestionBank
	case "c":
		bank = CQuestionBank
	default:
		return nil
	}

	// Question numbers are 1-based.
	if questionNumber < 1 || questionNumber > len(bank) {
		return nil
	}

	return &bank[questionNumber-1]
}

func GetTotalQuestions(language string) int {
	language = strings.ToLower(strings.TrimSpace(language))

	switch language {
	case "python":
		return len(PythonQuestionBank)
	case "c":
		return len(CQuestionBank)
	default:
		return 0
	}
}
