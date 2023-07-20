package internal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hitachienergy.com/cr-operator/agent/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GenerateOCILayout() string {
	return "{\"imageLayoutVersion\": \"1.0.0\"}"
}

func ComputeSha256(reader io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Gzip(reader io.Reader, target io.Writer) error {
	archiver := gzip.NewWriter(target)
	defer archiver.Close()

	_, err := io.Copy(archiver, reader)
	return err
}

type Config struct {
}

type RootFs struct {
	Type    string   `json:"type"`
	DiffIds []string `json:"diff_ids"`
}

type HistoryEntry struct {
	Created   time.Time `json:"created"`
	CreatedBy string    `json:"created_by"`
}

type ImageConfig struct {
	Created      time.Time      `json:"created"`
	Architecture string         `json:"architecture"`
	Variant      string         `json:"variant"`
	Os           string         `json:"os"`
	Config       Config         `json:"config"`
	RootFs       RootFs         `json:"rootfs"`
	History      []HistoryEntry `json:"history"`
}

func GenerateImageConfig(diffId string) ([]byte, error) {
	conf := ImageConfig{
		Created:      time.Now(),
		Architecture: "amd64",
		Variant:      "v7",
		Os:           "linux",
		Config:       Config{},
		RootFs: RootFs{
			Type: "layers",
			DiffIds: []string{
				diffId,
			},
		},
		History: []HistoryEntry{
			{
				Created:   time.Now(),
				CreatedBy: "cr-agent",
			},
		},
	}

	return json.Marshal(conf)
}

type BlobReference struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

type ImageManifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        BlobReference     `json:"config"`
	Layers        []BlobReference   `json:"layers"`
	Annotations   map[string]string `json:"annotations"`
}

func GenerateImageManifest(imageConfigDigest string, imageConfigSize int, layerDigest string, layerSize int, containerName string) ([]byte, error) {
	manifest := ImageManifest{
		SchemaVersion: 2,
		MediaType:     "applications/vnd.oci.image.manifest.v1+json",
		Config: BlobReference{
			MediaType: "application/vnd.oci.image.config.v1+json",
			Digest:    imageConfigDigest,
			Size:      imageConfigSize,
		},
		Layers: []BlobReference{
			{
				MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
				Digest:    layerDigest,
				Size:      layerSize,
			},
		},
		Annotations: map[string]string{
			"io.kubernetes.cri-o.annotations.checkpoint.name": containerName,
			"org.opencontainers.image.base.digest":            "",
			"org.opencontainers.image.base.name":              "",
		},
	}

	return json.Marshal(manifest)
}

type Index struct {
	SchemaVersion int             `json:"schemaVersion"`
	Manifests     []BlobReference `json:"manifests"`
}

func GenerateIndex(imageManifestDigest string, imageManifestSize int) ([]byte, error) {
	manifest := Index{
		SchemaVersion: 2,
		Manifests: []BlobReference{
			{
				MediaType: "application/vnd.oci.image.manifest.v1+json",
				Digest:    imageManifestDigest,
				Size:      imageManifestSize,
			},
		},
	}

	return json.Marshal(manifest)
}

func CreateOCIImage(fromPath string, containerName string, checkpointName string) error {
	checkpointFile, err := os.Create(filepath.Join(config.TempDir, checkpointName))
	if err != nil {
		return err
	}
	tarWriter := tar.NewWriter(checkpointFile)
	blobPath := "blobs/sha256/"

	err = tarWriter.WriteHeader(&tar.Header{
		Typeflag:   tar.TypeDir,
		Name:       "blobs/",
		Size:       0,
		Mode:       0755,
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
	})
	err = tarWriter.WriteHeader(&tar.Header{
		Typeflag:   tar.TypeDir,
		Name:       blobPath,
		Size:       0,
		Mode:       0755,
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
	})
	if err != nil {
		return err
	}

	OCILayout := []byte(GenerateOCILayout())
	err = writeToTar(tarWriter, "oci-layout", OCILayout)
	if err != nil {
		return err
	}
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	unzippedDigestHash := sha256.New()
	firstInputReader := io.TeeReader(from, unzippedDigestHash)

	zippedBuffer := &bytes.Buffer{}
	err = Gzip(firstInputReader, zippedBuffer)
	if err != nil {
		return err
	}
	unzippedDigest := fmt.Sprintf("%x", unzippedDigestHash.Sum(nil))

	zippedBytes := zippedBuffer.Bytes()
	zippedDigest := ComputeSha256(bytes.NewReader(zippedBytes))
	err = writeToTar(tarWriter, blobPath+zippedDigest, zippedBytes)
	if err != nil {
		return err
	}
	imageConfig, err := GenerateImageConfig("sha256:" + unzippedDigest)
	if err != nil {
		return err
	}
	imageConfigDigest := ComputeSha256(bytes.NewReader(imageConfig))
	err = writeToTar(tarWriter, blobPath+imageConfigDigest, imageConfig)
	if err != nil {
		return err
	}
	manifest, err := GenerateImageManifest("sha256:"+imageConfigDigest, len(imageConfig), "sha256:"+zippedDigest, len(zippedBytes), containerName)
	if err != nil {
		return err
	}
	manifestDigest := ComputeSha256(bytes.NewReader(manifest))
	err = writeToTar(tarWriter, blobPath+manifestDigest, manifest)
	if err != nil {
		return err
	}
	index, err := GenerateIndex("sha256:"+manifestDigest, len(manifest))
	if err != nil {
		return err
	}
	err = writeToTar(tarWriter, "index.json", index)
	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	return os.Remove(fromPath)
}

func writeToTar(tarWriter *tar.Writer, filename string, bytes []byte) error {
	err := tarWriter.WriteHeader(&tar.Header{
		Typeflag:   tar.TypeReg,
		Name:       filename,
		Size:       int64(len(bytes)),
		Mode:       0644,
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
	})
	if err != nil {
		return err
	}
	_, err = tarWriter.Write(bytes)
	return err
}
