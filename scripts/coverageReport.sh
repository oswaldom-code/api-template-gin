#!/bin/bash

COVERAGE_FILE="coverage.out"
HTML_REPORT="coverage.html"

echo "Executing tests and generating coverage report"
go test -coverprofile="$COVERAGE_FILE"

if [[ $? -eq 0 ]]; then
	echo "Generating HTML file..."

	go tool cover -html=$COVERAGE_FILE -o $HTML_REPORT

	if command -v xdg-open &>/dev/null; then
		xdg-open $HTML_REPORT # Linux
	elif command -v open &>/dev/null; then
		open $HTML_REPORT # macOS
	else
		echo "Cannot open report automatically. Open the file $HTML_REPORT in your browser ."
	fi
else
	echo "An error occurred while trying to generate the coverage report."
fi
