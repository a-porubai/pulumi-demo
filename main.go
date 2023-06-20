package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := createInfrastructure(ctx)

		return err
	})
}

type infrastructure struct {
	bucket           *storage.Bucket
	IAMBindingPolicy *storage.BucketIAMBinding
	bucketObject     *storage.BucketObject
}

func createInfrastructure(ctx *pulumi.Context) (*infrastructure, error) {
	bucket, err := storage.NewBucket(ctx, "my-bucket", &storage.BucketArgs{
		Location: pulumi.String("US"),
	})
	if err != nil {
		return nil, err
	}

	IAMBindingPolicy, err := storage.NewBucketIAMBinding(ctx, "my-bucket-IAMBinding", &storage.BucketIAMBindingArgs{
		Bucket: bucket.Name,
		Role:   pulumi.String("roles/storage.objectViewer"),
		Members: pulumi.StringArray{
			pulumi.String("allUsers"),
		},
	})
	if err != nil {
		return nil, err
	}

	fileName := config.Get(ctx, "bucket:fileName")

	bucketObject, err := storage.NewBucketObject(ctx, fileName, &storage.BucketObjectArgs{
		Bucket:      bucket.Name,
		ContentType: pulumi.String("text/html"),
		Source:      pulumi.NewFileAsset(fileName),
	})
	if err != nil {
		return nil, err
	}

	bucketEndpoint := pulumi.Sprintf("https://storage.googleapis.com/%s/%s", bucket.Name, bucketObject.Name)
	ctx.Export("bucketEndpoint", bucketEndpoint)

	return &infrastructure{
		bucket:           bucket,
		IAMBindingPolicy: IAMBindingPolicy,
		bucketObject:     bucketObject,
	}, nil
}
