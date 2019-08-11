#!/usr/bin/env bash

golint
go fmt
git add .
git commit -m "Implement and test feature: provide automated link to wrecking ball slack"
git push
