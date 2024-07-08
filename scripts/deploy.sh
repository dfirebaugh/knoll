#!/bin/bash

cd output
git init
git add .
git commit -m "Deploy gh-pages"
if git remote | grep origin; then
    git remote remove origin
fi
git remote add origin git@github.com:dfirebaugh/knoll.git
git checkout -b gh-pages
git push -f origin gh-pages
cd ..

