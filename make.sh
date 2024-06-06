#!/bin/bash

# compile up to date
go build -o mcli

# create tex file
./mcli rulebook -s mistfall.json -t mistfall.tex -d mistfall_design.json
./mcli rulebook -s terra_noctis.json -t terra_noctis.tex -d mistfall_design.json

# run twice...
pdflatex --shell-escape -interaction=nonstopmode mistfall.tex > /dev/null 2>&1
pdflatex --shell-escape -interaction=nonstopmode mistfall.tex > /dev/null 2>&1
pdflatex --shell-escape -interaction=nonstopmode terra_noctis.tex > /dev/null 2>&1
pdflatex --shell-escape -interaction=nonstopmode terra_noctis.tex > /dev/null 2>&1

# cleanup latex files
rm mistfall.aux
rm mistfall.log
rm mistfall.out
rm mistfall.toc
mv mistfall.tex latexGen/
rm terra_noctis.aux
rm terra_noctis.log
rm terra_noctis.out
rm terra_noctis.toc
mv terra_noctis.tex latexGen/

# move to target directory
mv mistfall.pdf pdfs/ || exit 1
mv terra_noctis.pdf pdfs/ || exit 1

# view it
firefox pdfs/terra_noctis.pdf
