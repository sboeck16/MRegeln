#!/bin/bash

# compile up to date
go build -o mcli

# create tex file
./mcli rulebook -s mistfall.json -t mistfall.tex -d mistfall_design.json

# run twice...
pdflatex --shell-escape -interaction=nonstopmode mistfall.tex > /dev/null 2>&1
pdflatex --shell-escape -interaction=nonstopmode mistfall.tex > /dev/null 2>&1

# cleanup latex files
rm mistfall.aux
rm mistfall.log
rm mistfall.out
rm mistfall.toc
mv mistfall.tex latexGen/

# move to target directory
mv mistfall.pdf pdfs/ || exit 1

# view it
firefox pdfs/mistfall.pdf
