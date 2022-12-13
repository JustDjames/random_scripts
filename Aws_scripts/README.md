# AWS Scripts

This README contains descriptions of the scripts in this repo

## Get Admins

Go script used to get list of all AWS admins in AWS account. gets config from AWS CLI config

## Cost breakdown

This Python script breaks down the costs in a AWS account for a specified period by Usage Type. it then takes this information and places it in a excel spreadsheet

## Volume Check

goes through all of the EBS volumes currently in the AWS account and checks whether they are in use or not. A volume is in the if it is attached to the running or pending instance. It is not in use if it is no attached to a instance or it is attached to a instance that is stopped, terminated or shutting-down. the results are placed into a spreadsheet
