package woocommerce

import (
	"fmt"
)

const (
	filesBasePath = "download"
)

// FileService is an interface for interfacing with the file endpoints of
// the WooCommerce files restful API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#files
type FileService interface {
	Get(file string) (*File, error)
}

// FileServiceOp handles communication with the files related methods of WooCommerce restful api
type FileServiceOp struct {
	Client *Client
}

// File represent a  wooCommerce file's All  properties columns
type File struct {
	Name    string `json:"filename,omitempty"`
	Content []byte `json:"content,omitempty"`
}

// Get implement for retrieve and view a specific file
// https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-file
func (w *FileServiceOp) Get(file string) (*File, error) {
	path := fmt.Sprintf("%s/%s", filesBasePath, file)
	resource := new(File)
	// Use createAndDoGetHeaders to access response headers
	headers, err := w.Client.createAndDoGetHeaders("GET", path, nil, nil, &resource)
	
	if err == nil {
		w.Client.log.Infof("FileServiceOp.Get success: file=%s, size=%d, headers=%v", file, len(resource.Content), headers)
	}

	return resource, err
}
