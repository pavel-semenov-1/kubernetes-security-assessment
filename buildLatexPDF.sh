#!/bin/bash

cd latex
pdflatex -shell-escape main.tex
if [ $? -ne 0 ]; then
    exit 1
fi
bibtex main
if [ $? -ne 0 ]; then
    exit 1
fi
pdflatex -shell-escape main.tex
cd -