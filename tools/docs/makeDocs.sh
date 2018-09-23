#!/usr/bin/env bash

pandoc md/index.md    -o ../../docs/index.html
pandoc md/tutorial.md -o ../../docs/tutorial.html

../../build/golem makeDocs.glm

#pandoc md/index.md    --table-of-contents -s -o ../../docs/index.html