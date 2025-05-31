// Package main implements a system for static analysis of Go source code.
//
// ### Overview:
// This package provides a mechanism for automatically launching multiple tools for analyzing Go programs, ensuring detection of potential issues, improving performance, and enhancing code quality.
//
// ### Used Analyzers:
//
// Below are descriptions of all the analyzers connected by this package:
//
// #### Core Analyzers:
// - **appends**: Checks correct usage of built-in slice extension functions like `append`.
// - **asmdecl**: Verifies compatibility between assembly declarations and corresponding Go definitions.
// - **assign**: Detects incorrect assignments, especially when copying structs or interfaces.
// - **atomic**: Ensures safe use of atomic operations.
// - **atomicalign**: Analyzes memory alignment when working with atomic structures.
// - **bools**: Optimizes boolean expressions and their handling.
// - **buildssa**: Collects intermediate representation (SSA) information for further analysis.
// - **buildtag**: Validates proper usage of build tags.
// - **composite**: Ensures correctness of composite literals (creation of complex objects).
// - **copylock**: Finds cases where mutexes are accidentally passed by value.
// - **ctrlflow**: Analyzes control flow within functions.
// - **deepequalerrors**: Identifies improper use of `reflect.DeepEqual`.
// - **defers**: Analyzes correct usage of `defer` statements.
// - **directive**: Looks for preprocessor directive-related errors.
// - **errorsas**: Analyzes correct usage of methods `errors.As` and `errors.Is`.
// - **fieldalignment**: Optimizes field layout in structures to reduce memory consumption.
// - **findcall**: Searches for specific function calls in the codebase.
// - **framepointer**: Analyzes stack frame safety.
// - **httpmux**: Analyzes HTTP mux routing for potential problems.
// - **ifaceassert**: Verifies interface assertions' correctness.
// - **inspect**: Inspects the program's abstract syntax tree (AST) for deep analysis.
// - **loopclosure**: Detects closures inside loops that may lead to unexpected behavior.
// - **lostcancel**: Ensures contexts with cancelation do not lose references to them.
// - **nilfunc**: Prevents passing null pointers as arguments expected to be non-null.
// - **nilness**: Analyzes uncertainty about whether values are nil.
// - **pkgfact**: Collects facts about packages for subsequent analysis.
// - **printf**: Analyzes argument correctness for printf-like functions.
// - **reflectvaluecompare**: Verifies correct usage of reflection operations.
// - **shadow**: Detects variable shadowing and naming conflicts.
// - **shift**: Ensures proper usage of shift operators.
// - **sigchanyzer**: Analyzes signal channels to prevent blocking issues.
// - **slog**: Analyzes logging done using standard loggers.
// - **sortslice**: Ensures efficient sorting of slices.
// - **stdmethods**: Tests implementation of standard interface methods.
// - **stdversion**: Verifies API version compatibility.
// - **stringintconv**: Checks conversions between strings and integers.
// - **structtag**: Verifies correct usage of struct metadata tags.
// - **testinggoroutine**: Analyzes tests involving goroutines.
// - **tests**: General testing-related analyses.
// - **timeformat**: Ensures time formatting correctness.
// - **unmarshal**: Analyzes unmarshaling and marshaling processes.
// - **unreachable**: Detects unreachable code paths.
// - **unsafeptr**: Analyzes incorrect usage of `unsafe.Pointer`.
// - **unusedresult**: Warns about ignored function results.
// - **unusedwrite**: Analyzes unused writes into memory.
// - **usesgenerics**: Verifies the application of generics types in your project.
//
// #### StaticCheck Analyzers (SA Series):
// - **SA1xxx**: Formatting and printing related checks (SA1000-SA1030)
// - **SA2xxx**: Concurrency and synchronization issues (SA2000-SA2003)
// - **SA3xxx**: OS/File operations and system interactions (SA3000-SA3001)
// - **SA4xxx**: Logic and correctness issues (SA4000-SA4031)
//   - Includes checks for: infinite loops, self-assignments, redundant operations,
//     incorrect comparisons, ineffective operations, and suspicious constructs
//
// - **SA5xxx**: Security and correctness of HTTP/TLS/exec operations (SA5000-SA5012)
// - **SA6xxx**: Performance and optimization issues (SA6000-SA6005)
// - **SA9xxx**: Code style and maintainability (SA9001-SA9008)
//
// #### Style Checks (ST1000):
// - Verifies Go code style conventions and documentation standards
//
// #### Quick Fixes (QF1001):
// - Identifies patterns that could be simplified with quick fixes
//
// #### Unused Code Detection (U1000):
// - Finds unused functions, variables, constants, and types
//
// #### Build Command:
// To compile the executable binary, run the following command:
//
//	go build -o bin/analyzer ./internal/analyzer
//
// #### Run Command:
// To analyze the entire project, execute the script:
//
//	./analyze.sh
//
// ### Note:
// Before compiling, ensure you are in the project's root directory and have dependencies set up correctly for the analyzers.
// The complete set of analyzers provides over 100 different checks covering:
// - Correctness
// - Performance
// - Security
// - Style
// - Maintainability
// - Error handling
// - Modern Go features
package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
	"os"
	"path/filepath"
)

// Config is the filename for the JSON configuration
const Config = `static-check_conf.json`

// ConfigData defines the structure of the JSON configuration file.
// It contains a list of staticcheck analyzer names to enable.
type ConfigData struct {
	StaticCheck []string `json:"static-check,omitempty"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	analyzers := setupAnalyzers(cfg)

	multichecker.Main(analyzers...)
}

// loadConfig reads and parses the JSON configuration file.
// It returns the configuration data or an error if loading fails.
func loadConfig() (*ConfigData, error) {
	configPath := filepath.Join(".", Config)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg ConfigData
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// setupAnalyzers configures all analyzers based on the provided configuration.
// It combines standard Go analyzers with selected staticcheck analyzers.
func setupAnalyzers(cfg *ConfigData) []*analysis.Analyzer {
	analyzers := []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpmux.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stdversion.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

	enabledChecks := make(map[string]bool)
	for _, v := range cfg.StaticCheck {
		enabledChecks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if enabledChecks[v.Analyzer.Name] {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	return analyzers
}
