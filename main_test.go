package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestCreateInfrastructure(t *testing.T) {
	config := map[string]string{
		"bucket:fileName": "index-dev.html",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		infrastructureObjects, err := createInfrastructure(ctx)
		if err != nil {
			return err
		}
		bucket := infrastructureObjects.bucket
		pulumi.All(bucket.Location).ApplyT(func(all []interface{}) error {
			location := all[0].(string)

			const expectedBucketLocation = "US"

			if location != expectedBucketLocation {
				t.Errorf("invalid bucket location got: %s, want: %s", location, expectedBucketLocation)
			}

			return nil
		})

		bucketIAMBinging := infrastructureObjects.IAMBindingPolicy
		pulumi.All(bucketIAMBinging.Role, bucketIAMBinging.Members).ApplyT(func(all []interface{}) error {
			role := all[0].(string)
			members := all[1].([]string)

			const expectedRole = "roles/storage.objectViewer"
			const expectedMembersLength = 1
			const expectedMember = "allUsers"

			if role != expectedRole {
				t.Errorf("invalid role got: %s, want: %s", role, expectedRole)
			}

			if len(members) != expectedMembersLength {
				t.Errorf("invalid members length got: %v, want: %v", len(members), expectedMembersLength)
			}

			if members[0] != expectedMember {
				t.Errorf("invalid member got: %s, want: %s", members[0], expectedMember)
			}

			return nil
		})

		bucketObject := infrastructureObjects.bucketObject
		pulumi.All(bucketObject.ContentType, bucketObject.Source).ApplyT(func(all []interface{}) error {
			contentType := all[0].(string)
			source := all[1].(pulumi.Asset)

			const expectedContentType = "text/html"
			const expectedFilePath = "index-dev.html"

			if contentType != expectedContentType {
				t.Errorf("invalid content type got: %s, want: %s", contentType, expectedContentType)
			}

			if source.Path() != expectedFilePath {
				t.Errorf("invalid file path got: %s, want: %s", source.Path(), expectedFilePath)
			}

			return nil
		})

		return nil
	}, WithMocksAndConfig("unit-tests", "stack", config, &Mocks{}))

	if err != nil {
		t.Error(t, err)
	}
}
