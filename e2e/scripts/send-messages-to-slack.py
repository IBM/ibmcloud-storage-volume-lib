import os
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
        channels="storage-test-runs",
        file=(FilePath, open(PathToFile, 'rb'), 'txt'),
        filename=FilePath,
        username='IBM VPC storage e2e test runs and results',
        title="VPC storage e2e logs"
  )

if __name__ == '__main__':
   send_file("e2e_logs.txt")
