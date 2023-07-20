package internal

import (
	"context"
	"fmt"
	"github.com/containers/buildah"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
	"hitachienergy.com/cr-operator/agent/config"
	"os"
	"path/filepath"
	"time"
)

func CreateCheckpointImage(ctx context.Context, checkpointPath string, containerName string, checkpointName string) (string, string, error) {
	startTime := time.Now().UnixMilli()
	checkpointPrefix := "cr"
	checkpointImageName := fmt.Sprintf("%s/%s:latest", checkpointPrefix, checkpointName)
	buildStoreOptions, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	fmt.Println("[perf] defaultStoreOptionsTime: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	if err != nil {
		fmt.Println("storage.DefaultStoreOptions")
		return "", "", err
	}
	buildStore, err := storage.GetStore(buildStoreOptions)
	fmt.Println("[perf] GetStoreTime: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	if err != nil {
		fmt.Println("storage.GetStore")
		return "", "", err
	}
	defer buildStore.Shutdown(false)
	builderOpts := buildah.BuilderOptions{
		FromImage: "scratch",
	}
	builder, err := buildah.NewBuilder(ctx, buildStore, builderOpts)
	fmt.Println("[perf] NewBuilder: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	if err != nil {
		fmt.Println("buildah.NewBuilder")
		return "", "", err
	}
	defer builder.Delete()

	fmt.Println("create Image: builder + store setup complete")
	err = builder.Add("/", true, buildah.AddAndCopyOptions{}, checkpointPath)
	fmt.Println("[perf] BuilderAdd: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	if err != nil {
		fmt.Println("builder.Add")
		return "", "", err
	}
	fmt.Println("create Image: added archive")
	builder.ImageAnnotations["io.kubernetes.cri-o.annotations.checkpoint.name"] = containerName
	imageRef, err := is.Transport.ParseStoreReference(buildStore, checkpointImageName)
	fmt.Println("[perf] ParseStoreReference: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	if err != nil {
		fmt.Println("is.Transport.ParseStoreReference")
		return "", "", err
	}
	fmt.Println("create Image: generated store reference")
	imageId, _, _, err := builder.Commit(ctx, imageRef, buildah.CommitOptions{})
	fmt.Println("[perf] Commit: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	fmt.Println("create Image: committed")
	if err != nil {
		fmt.Println("builder.Commit")
		return "", "", err
	}

	sysCtx := &types.SystemContext{}
	policy, err := signature.DefaultPolicy(sysCtx)
	if err != nil {
		return "", "", fmt.Errorf("obtaining default signature policy: %w", err)
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return "", "", fmt.Errorf("creating new signature policy context: %w", err)
	}
	copyOpts := &copy.Options{
		DestinationCtx: sysCtx,
	}

	fmt.Println("[perf] PolicyContext: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	fmt.Println("create Image: parsed policies")
	exportName := fmt.Sprintf("oci-archive://%s/%s", config.TempDir, checkpointName)
	destinationRef, err := alltransports.ParseImageName(exportName)
	fmt.Println("ParseImageName: ", time.Now().UnixMilli()-startTime)
	if err != nil {
		fmt.Println("is.Transport.ParseStoreReference")
		return "", "", err
	}
	fmt.Println("create Image: parsed image name")
	_, err = copy.Image(context.TODO(), policyContext, destinationRef, imageRef, copyOpts)
	fmt.Println("[perf] CopyImage: ", time.Now().UnixMilli()-startTime)
	startTime = time.Now().UnixMilli()
	fmt.Println("create Image: copied image")
	if err != nil {
		fmt.Println("copy.Image")
		return "", "", err
	}
	err = os.Remove(checkpointPath)
	fmt.Println("create Image: removed archive")
	if err != nil {
		fmt.Println("error while deleting checkpoint archive")
		return "", "", err
	}

	return imageId, checkpointImageName, nil
}

func ImportImage(path string) error {
	sysCtx := &types.SystemContext{}
	policy, err := signature.DefaultPolicy(sysCtx)
	if err != nil {
		return fmt.Errorf("obtaining default signature policy: %w", err)
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return fmt.Errorf("creating new signature policy context: %w", err)
	}

	checkpointName := filepath.Base(path)
	destinationImageName := fmt.Sprintf("containers-storage:localhost/%s:latest", checkpointName)
	destinationRef, err := alltransports.ParseImageName(destinationImageName)
	if err != nil {
		return err
	}
	sourceImageName := fmt.Sprintf("oci-archive://%s/%s", config.TempDir, checkpointName)
	sourceRef, err := alltransports.ParseImageName(sourceImageName)
	if err != nil {
		return err
	}

	err = destinationRef.DeleteImage(context.Background(), sysCtx)
	if err != nil {
		fmt.Println("WARNING: could not delete last image. Continuing...")
	}

	copyOpts := &copy.Options{
		DestinationCtx: sysCtx,
	}
	_, err = copy.Image(context.TODO(), policyContext, destinationRef, sourceRef, copyOpts)
	if err != nil {
		return err
	}
	return nil
}
