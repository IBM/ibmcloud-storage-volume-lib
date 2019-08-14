#! /usr/bin/env python

import os
import sys
from slackclient import SlackClient

slack_client = SlackClient(token=os.environ['SLACK_API_TOKEN'])
channelName = "ibm-vpc-block-csi-e2e"
def channel_info(channel_id):
    channels_info = slack_client.api_call("channels.info", channel=channel_id)
    if channels_info:
        return channels_info
    return None

def send_file (FilePath):
    with open(FilePath, 'r') as content_file:
        content = content_file.read()

    response = slack_client.api_call(
        "files.upload",
        channel=channelName,
        file=content,
        filename="e2e_test_job_logs.txt",
        username='IBM VPC storage common library e2e test results',
        title="VPC storage common library e2e test full logs")
    print response["ok"]

def send_message (FilePath):
    with open(FilePath, 'r+') as file_content:
        content = file_content.read()

    content = '```' + content + '```'
    slack_client.api_call(
        "chat.postMessage",
        channel=channelName,
        text=content,
        username='IBM VPC storage common library e2e test results')

if __name__ == '__main__':
    gopath = os.environ["GOPATH"]
    filePath = gopath + "/src/github.com/IBM/ibmcloud-storage-volume-lib/" + sys.argv[1]
    isAttachment = sys.argv[2]
    if filePath is None or isAttachment is None:
        print "Please send the variables filePath and isAttachment"
        sys.exit(1)

    if isAttachment == "True" or isAttachment == "true":
        send_file(filePath)
    else:
        send_message(filePath)
