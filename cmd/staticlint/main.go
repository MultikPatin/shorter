// This package provides a mechanism for automatically launching multiple tools for analyzing Go programs, ensuring detection of potential issues, improving performance, and enhancing code quality.
//
// Custom Analyzers:
//
//	osexit:
//		Prohibits using a direct os.Exit call.
//
// Core Analyzers:
//
//	appends:
//		Checks correct usage of built-in slice extension functions like `append`.
//	asmdecl:
//		Verifies compatibility between assembly declarations and corresponding Go definitions.
//	assign:
//		Detects incorrect assignments, especially when copying structs or interfaces.
//	atomic:
//		Ensures safe use of atomic operations.
//	atomicalign:
//		Analyzes memory alignment when working with atomic structures.
//	bools:
//		Optimizes boolean expressions and their handling.
//	buildssa:
//		Collects intermediate representation (SSA) information for further analysis.
//	buildtag:
//		Validates proper usage of build tags.
//	composite:
//		Ensures correctness of composite literals (creation of complex objects).
//	copylock:
//		Finds cases where mutexes are accidentally passed by value.
//	ctrlflow:
//		Analyzes control flow within functions.
//	deepequalerrors:
//		Identifies improper use of `reflect.DeepEqual`.
//	defers:
//		Analyzes correct usage of `defer` statements.
//	directive:
//		Looks for preprocessor directive-related errors.
//	errorsas:
//		Analyzes correct usage of methods `errors.As` and `errors.Is`.
//	fieldalignment:
//		Optimizes field layout in structures to reduce memory consumption.
//	findcall:
//		Searches for specific function calls in the codebase.
//	framepointer:
//		Analyzes stack frame safety.
//	httpmux:
//		Analyzes HTTP mux routing for potential problems.
//	ifaceassert:
//		Verifies interface assertions' correctness.
//	inspect:
//		Inspects the program's abstract syntax tree (AST) for deep analysis.
//	loopclosure:
//		Detects closures inside loops that may lead to unexpected behavior.
//	lostcancel:
//		Ensures contexts with cancelation do not lose references to them.
//	nilfunc:
//		Prevents passing null pointers as arguments expected to be non-null.
//	nilness:
//		Analyzes uncertainty about whether values are nil.
//	pkgfact:
//		Collects facts about packages for subsequent analysis.
//	printf:
//		Analyzes argument correctness for printf-like functions.
//	reflectvaluecompare:
//		Verifies correct usage of reflection operations.
//	shadow:
//		Detects variable shadowing and naming conflicts.
//	shift:
//		Ensures proper usage of shift operators.
//	sigchanyzer:
//		Analyzes signal channels to prevent blocking issues.
//	slog:
//		Analyzes logging done using standard loggers.
//	sortslice:
//		Ensures efficient sorting of slices.
//	stdmethods:
//		Tests implementation of standard interface methods.
//	stdversion:
//		Verifies API version compatibility.
//	stringintconv:
//		Checks conversions between strings and integers.
//	structtag:
//		Verifies correct usage of struct metadata tags.
//	testinggoroutine:
//		Analyzes tests involving goroutines.
//	tests:
//		General testing-related analyses.
//	timeformat:
//		Ensures time formatting correctness.
//	unmarshal:
//		Analyzes unmarshaling and marshaling processes.
//	unreachable:
//		Detects unreachable code paths.
//	unsafeptr:
//		Analyzes incorrect usage of `unsafe.Pointer`.
//	unusedresult:
//		Warns about ignored function results.
//	unusedwrite:
//		Analyzes unused writes into memory.
//	usesgenerics:
//		Verifies the application of generics types in your project.
//
// StaticCheck Analyzers (SA Series):
// SA1xxx**:
//   - Formatting and printing related checks (SA1000-SA1030)
//
// SA2xxx**:
//   - Concurrency and synchronization issues (SA2000-SA2003)
//
// SA3xxx**:
//   - OS/File operations and system interactions (SA3000-SA3001)
//
// SA4xxx**:
//   - Logic and correctness issues (SA4000-SA4031)
//
// SA5xxx**:
//   - Security and correctness of HTTP/TLS/exec operations (SA5000-SA5012)
//
// SA6xxx**:
//   - Performance and optimization issues (SA6000-SA6005)
//
// SA9xxx**:
//   - Code style and maintainability (SA9001-SA9008)
//
// Style Checks (ST1000):
//   - Verifies Go code style conventions and documentation standards
//
// Quick Fixes (QF1001):
//   - Identifies patterns that could be simplified with quick fixes
//
// Unused Code Detection (U1000):
//   - Finds unused functions, variables, constants, and types
//
// Build Command:
// To compile the executable binary, run the following command:
//
//	go build -o bin/analyzer ./internal/analyzer
//
// Run Command:
// To analyze the entire project, execute the script:
//
//	./analyze.sh
//
// Note:
// Before compiling, ensure you are in the project's root directory and have dependencies set up correctly for the analyzers.
// The complete set of analyzers provides over 100 different checks covering:
//   - Correctness
//   - Performance
//   - Security
//   - Style
//   - Maintainability
//   - Error handling
//   - Modern Go features
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"log"
	"main/internal/analyzer"
)

func main() {
	cfg, err := analyzer.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	analyzers := analyzer.Setup(cfg)

	multichecker.Main(analyzers...)
}
