// from https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/ListAdmins/ListAdminsv2.go with some minor changes
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func GetNumUsersAndAdmins(c context.Context, client *iam.Client) (string, string, error) {
	users := ""
	admins := ""

	filters := make([]types.EntityType, 1)
	filters[0] = types.EntityTypeUser

	input := &iam.GetAccountAuthorizationDetailsInput{
		Filter: filters,
	}

	fmt.Println("getting info on All Users,Groups, Roles and Policies in account")
	resp, err := client.GetAccountAuthorizationDetails(c, input)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Got info! Beginning Admin check")

	adminName := "AdministratorAccess"

	for _, user := range resp.UserDetailList {
		fmt.Println("checking if", *user.UserName, "is a admin")
		isAdmin, err := isUserAdmin(c, client, user, adminName)

		if err != nil {
			return "", "", err
		}
		users += " " + *user.UserName

		if isAdmin {
			admins += " " + *user.UserName
		}
	}

	for resp.IsTruncated {
		input := &iam.GetAccountAuthorizationDetailsInput{
			Filter: filters,
			Marker: resp.Marker,
		}

		resp, err = client.GetAccountAuthorizationDetails(c, input)
		if err != nil {
			return "", "", err
		}

		for _, user := range resp.UserDetailList {
			fmt.Println("checking if", *user.UserName, "is a admin")
			isAdmin, err := isUserAdmin(c, client, user, adminName)

			if err != nil {
				return "", "", err
			}
			users += " " + *user.UserName

			if isAdmin {
				admins += " " + *user.UserName
			}
		}
	}
	return users, admins, nil
}

func isUserAdmin(c context.Context, client *iam.Client, user types.UserDetail, admin string) (bool, error) {
	policyHasAdmin := userPolicyHasAdmin(user, admin)
	if policyHasAdmin {
		return true, nil
	}

	attachedPolicyHasAdmin := attachedUserPolicyHasAdmin(user, admin)
	if attachedPolicyHasAdmin {
		return true, nil
	}

	userGroupsHaveAdmin, err := usersGroupsHaveAdmin(c, client, user, admin)
	if err != nil {
		return false, err
	}
	if userGroupsHaveAdmin {
		return true, nil
	}

	return false, nil
}

func userPolicyHasAdmin(user types.UserDetail, admin string) bool {
	for _, policy := range user.UserPolicyList {
		if *policy.PolicyName == admin {
			return true
		}
	}
	return false
}

func attachedUserPolicyHasAdmin(user types.UserDetail, admin string) bool {
	for _, policy := range user.AttachedManagedPolicies {
		if *policy.PolicyName == admin {
			return true
		}
	}
	return false
}

func usersGroupsHaveAdmin(c context.Context, client *iam.Client, user types.UserDetail, admin string) (bool, error) {
	input := &iam.ListGroupsForUserInput{
		UserName: user.UserName,
	}

	result, err := client.ListGroupsForUser(c, input)
	if err != nil {
		return false, err
	}

	for _, group := range result.Groups {
		groupPolicyHasAdmin, err := groupPolicyHasAdmin(c, client, group, admin)
		if err != nil {
			return false, err
		}

		if groupPolicyHasAdmin {
			return true, nil
		}

		attachedGroupPolicyHasAdmin, err := attachedGroupPolicyHasAdmin(c, client, group, admin)
		if err != nil {
			return false, err
		}

		if attachedGroupPolicyHasAdmin {
			return true, nil
		}
	}

	return false, nil
}

func groupPolicyHasAdmin(c context.Context, client *iam.Client, group types.Group, admin string) (bool, error) {
	input := &iam.ListGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	result, err := client.ListGroupPolicies(c, input)
	if err != nil {
		return false, err
	}

	for _, policyName := range result.PolicyNames {
		if policyName == admin {
			return true, nil
		}
	}

	return false, nil
}

func attachedGroupPolicyHasAdmin(c context.Context, client *iam.Client, group types.Group, admin string) (bool, error) {
	input := &iam.ListAttachedGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	result, err := client.ListAttachedGroupPolicies(c, input)
	if err != nil {
		return false, err
	}

	for _, policy := range result.AttachedPolicies {
		if *policy.PolicyName == admin {
			return true, nil
		}
	}

	return false, nil
}

func main() {

	listAdmins := flag.Bool("A", false, "list names of Admins")
	listUsers := flag.Bool("U", false, "list names of Users")
	profile := flag.String("profile", "default", "AWS CLI profile to use")
	flag.Parse()

	fmt.Println("Using profile", *profile)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(*profile), config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
		options.TokenProvider = stscreds.StdinTokenProvider
	}))

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := iam.NewFromConfig(cfg)

	users, admins, err := GetNumUsersAndAdmins(context.TODO(), client)
	if err != nil {
		fmt.Println("got an error finding users who are admins:")
		fmt.Println(err)
		return
	}
	fmt.Println("checked all users!")
	userList := strings.Split(users, " ")
	adminList := strings.Split(admins, " ")

	fmt.Println("")
	fmt.Println("Found", len(adminList)-1, "admin(s) out of", len(userList)-1, "user(s)")

	if *listAdmins {
		fmt.Println("")
		fmt.Println("Admins:")
		for _, a := range adminList {
			fmt.Println(" " + a)
		}
	}

	if *listUsers {
		fmt.Println("")
		fmt.Println("Users:")
		for _, u := range userList {
			fmt.Println("  " + u)
		}
	}
}
