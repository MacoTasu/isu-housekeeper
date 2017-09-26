#!/bin/sh

branch=$1
repository_name=`cat /home/isucon/repository_name`
slack_token=`cat /home/isucon/slack_token`

{
  echo "branch: $1"
  cd /home/isucon/$repository_name
  git fetch --prune
  git checkout $branch
  git pull origin $branch
  echo `git log --oneline -1`
} &> /home/isucon/deploy.log

log=`cat /home/isucon/deploy.log`
curl -XPOST -d "token=$slack_token" -d 'channel=isucon' -d "text=deploy: $log" 'https://slack.com/api/files.upload'
