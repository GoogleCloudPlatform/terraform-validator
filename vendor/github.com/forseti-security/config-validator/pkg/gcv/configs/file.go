// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package configs helps with loading and parsing configuration files
package configs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var (
	globals struct {
		// once for only running GCS client setup once
		once   sync.Once
		client *storage.Client
	}
)

// configGCSClient sets up the GCS client when needed.
func configGCSClient() {
	ctx := context.Background()

	var err error
	globals.client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewPath returns a new Path to a local or gcs file.
func NewPath(path string) (Path, error) {
	fileURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if fileURL.Scheme == "gs" {
		globals.once.Do(configGCSClient)
		return &gcsPath{
			bucket: fileURL.Host,
			path:   strings.TrimLeft(fileURL.Path, "/"),
		}, nil
	}

	// local fileIface could be dirIface or fileIface
	return &localPath{path: path}, nil
}

// File represents the contents of a file
type File struct {
	// Path is the path to the file.
	Path string
	// Content is the full contents for the file.
	Content []byte
}

// readPredicate is a predicate function for ReadAll to determine whether to read a file
type readPredicate func(path string) bool

// SuffixPredicate returns read predicate that returns true if the file name has the specified suffix.
func SuffixPredicate(suffix string) readPredicate {
	return func(path string) bool {
		return strings.HasSuffix(path, suffix)
	}
}

func matchesPredicates(path string, predicates []readPredicate) bool {
	for _, predicate := range predicates {
		if !predicate(path) {
			return false
		}
	}
	return true
}

// Path represents a path to a file or directory.
type Path interface {
	// ReadAll will read the given file, or recursively read all files under the specified directory.
	ReadAll(ctx context.Context, predicates ...readPredicate) ([]File, error)
}

// localPath handles local file paths.
type localPath struct {
	path string
}

// ReadAll implements Path
func (p *localPath) ReadAll(ctx context.Context, predicates ...readPredicate) ([]File, error) {
	var files []File
	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "error visiting path %s", path)
		}
		if f.IsDir() {
			return nil
		}
		if !matchesPredicates(path, predicates) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s", path)
		}
		files = append(files, File{Path: path, Content: content})
		return nil
	}
	err := filepath.Walk(p.path, visit)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read files in %s", p.path)
	}
	return files, nil
}

// gcsPath represents an object or prefix on GCS.
type gcsPath struct {
	bucket string
	path   string
}

// read reads an object from GCS
func (p *gcsPath) read(ctx context.Context, bucket *storage.BucketHandle, name string) (File, error) {
	fileName := fmt.Sprintf("gs://%s/%s", p.bucket, name)
	glog.V(2).Infof("Listing GCS Object %s", fileName)

	reader, err := bucket.Object(name).NewReader(ctx)
	if err != nil {
		return File{}, errors.Wrapf(err, "failed to read object %s", fileName)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			glog.Warningf("failed to close %s: %s", fileName, err)
		}
	}()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return File{}, errors.Wrapf(err, "failed to read %s", fileName)
	}
	return File{
		Content: data,
		Path:    fileName,
	}, nil
}

// ReadAll implements Path
func (p *gcsPath) ReadAll(ctx context.Context, predicates ...readPredicate) ([]File, error) {
	var files []File

	bucket := globals.client.Bucket(p.bucket)
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: p.path,
	})
	glog.V(2).Infof("Listing files in GCS at host %s and path %s", p.bucket, p.path)
	for {
		attrs, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		if !matchesPredicates(attrs.Name, predicates) {
			continue
		}

		file, err := p.read(ctx, bucket, attrs.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "")
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, errors.Errorf("no objects found at gs://%s/%s", p.bucket, p.path)
	}
	return files, nil
}
