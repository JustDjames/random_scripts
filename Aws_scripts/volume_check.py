# also requires pip package openpyxl
import boto3
import argparse
from pandas import DataFrame

parser = argparse.ArgumentParser(description="list volumes and what they are attached to")
parser.add_argument("aws_profile", help="AWS profile from your aws config to use in script")

args = parser.parse_args()

session = boto3.Session(profile_name=args.aws_profile)

ec2 = session.client('ec2')

print("getting volumes")

volumeRes = ec2.describe_volumes()

volumes = volumeRes['Volumes']

volume_id = []
in_use = []
attached_to = []

print("checking whether volumes are attached to instances. this will take awhile")

for volume in volumeRes['Volumes']:
    volume_id.append(volume["VolumeId"])
    
    # print(len(volume['Attachments']))
    if len(volume['Attachments']) == 0:
        attached_to.append("Volume isn't attached")
        in_use.append("N")
    else:
        attached_instances = []
        
        for each in volume['Attachments']:
            attached_instances.append(each['InstanceId'])

        instanceRes = ec2.describe_instances(InstanceIds=attached_instances)

        instance_status = []
        
        for each in instanceRes['Reservations']:
            for instance in each["Instances"]:
                instance_status.append(instance['State']['Name'])
        
        if 'pending' in instance_status or 'running' in instance_status:
            attached_to.append("Volume is attached to running or pending instances: " + ', '.join(attached_instances))
            in_use.append("Y")
        else:
            attached_to.append("Volume is attached to stopped, terminated or shutting-down instances: " + ', '.join(attached_instances))
            in_use.append("N")

print("building excel sheet")

df = DataFrame(
    {
        "Volume": volume_id,
        "In Use?": in_use,
        "Attached To": attached_to
    }
)

df.to_excel('{env}_volume_check.xlsx'.format(env=args.aws_profile), index=False)