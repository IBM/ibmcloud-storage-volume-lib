#! /usr/bin/env python

import os
import sys
from slackclient import SlackClient

slack_client = SlackClient('SLACK_TOKEN_REPLACE')

def channel_info(channel_id):
    channel_info = slack_client.api_call("channels.info", channel=channel_id)
    if channel_info:
        return channel_info['channel']
    return None

def send_file (FilePath):
    gopath = os.environ["GOPATH"]
    PathToFile = gopath + "/src/github.com/IBM/ibmcloud-storage-volume-lib/" + FilePath
    slack_client.api_call(
         "files.upload",
        channel="storage-test-runs",
        file=(FilePath, open(PathToFile, 'rb'), 'txt'),
        filename=FilePath,
        username='IBM VPC Storage',
        title="VPC storage common library  e2e test full logs"
  )

def send_message (FilePath):
    gopath = os.environ["GOPATH"]
    PathToFile = gopath + "/src/github.com/IBM/ibmcloud-storage-volume-lib/" + FilePath
    with open(PathToFile, 'r') as content_file:
         content = content_file.read()
    slack_client.api_call(
         "chat.postMessage",
        channel="storage-test-runs",
        text=content,
        username='IBM VPC storage common library e2e test runs and results',
  )

if __name__ == '__main__':
    try:
        filePath = os.environ["FILEPATH"]
        isFile = os.environ["ISFILE"]
    except KeyError:
        print "Please set the environment variables FILEPATH and ISFILE"
        sys.exit(1)

    if isFile == "True" or isFile == "true":
        send_file(os.environ["FILEPATH"])
    else:
        send_message(os.environ["FILEPATH"])
