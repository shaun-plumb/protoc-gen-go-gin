#!/bin/bash

# Remove local stale tags : git tag -l | xargs git tag -d && git fetch -t

OLDTAG=$(gh release list -L 1 --json tagName -q '.[0].tagName')

echo $OLDTAG

NEWTAG=$(echo $OLDTAG | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}')

echo $NEWTAG

git tag $NEWTAG
git push origin --tags

if [[ $NEWTAG != v* ]]; then NEWTAG=v$NEWTAG; fi

gh release create $NEWTAG -t $NEWTAG -n ""