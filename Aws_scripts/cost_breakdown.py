import boto3
import argparse
from pandas import DataFrame
# also requires pip package openpyxl

# get aws profile from arguments and configure boto to use it

parser = argparse.ArgumentParser(description="breakdown the usage costs per service and usage types")
parser.add_argument("aws_profile", help="AWS profile from your aws config to use in script")
parser.add_argument("start_period", help="start date for the time period in format yyyy-mm-dd")
parser.add_argument("end_period", help="end date for the time period in format yyyy-mm-dd. add extra day to ensure values match console values")

args = parser.parse_args()

session = boto3.Session(profile_name=args.aws_profile)

cost = session.client('ce')

# get all services from this period

dimensionRes = cost.get_dimension_values(
    TimePeriod = {
        "Start": args.start_period,
        "End": args.end_period
        },
    Dimension = "SERVICE"
)

services = []

for each in dimensionRes["DimensionValues"]:
    services.append(each["Value"])

# print(services)

excel_service = []
excel_usageGroup = []
excel_cost = []


# loop through services to get cost breakdown by Usage Group
for service in services:
    print(f"getting {service} cost breakdown by Usage Group")
    filter = {
        "Dimensions": {
            "Key": "SERVICE",
            "Values": [service]
        }
    }

    costRes = cost.get_cost_and_usage(
        TimePeriod = {
            "Start": args.start_period,
            "End": args.end_period
        },
        Granularity = "MONTHLY",
        Filter = filter,
        Metrics=["UnblendedCost"],
        GroupBy = [
            {
                "Type": "DIMENSION",
                "Key": "USAGE_TYPE"
            }
            ]
    )

    # print(costRes["ResultsByTime"][0]["Groups"]) 

    # Save results to populate excel sheet

    for val in costRes["ResultsByTime"][0]["Groups"]:
        excel_service.append(service)
        excel_usageGroup.append(val["Keys"][0])
        excel_cost.append(val["Metrics"]["UnblendedCost"]["Amount"])


# build and create excel sheet
print("building excel sheet")

df = DataFrame(
    {
        "Service": excel_service,
        "Usage Type": excel_usageGroup,
        "Cost ($)": excel_cost
    }
)    

df["Cost ($)"] = df["Cost ($)"].astype(float)

df.to_excel('{env}_costs_{start}_{end}.xlsx'.format(env=args.aws_profile,start=args.start_period,end=args.end_period), index=False, float_format="%.2f")        
        
        
print("Done!")