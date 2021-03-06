// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"time"

	"github.com/miniflux/miniflux/crypto"
	"github.com/miniflux/miniflux/http"
	"github.com/miniflux/miniflux/http/handler"
	"github.com/miniflux/miniflux/logger"
)

// ImageProxy fetch an image from a remote server and sent it back to the browser.
func (c *Controller) ImageProxy(ctx *handler.Context, request *handler.Request, response *handler.Response) {
	// If we receive a "If-None-Match" header we assume the image in stored in browser cache
	if request.Request().Header.Get("If-None-Match") != "" {
		response.NotModified()
		return
	}

	encodedURL := request.StringParam("encodedURL", "")
	if encodedURL == "" {
		response.HTML().BadRequest(errors.New("No URL provided"))
		return
	}

	decodedURL, err := base64.URLEncoding.DecodeString(encodedURL)
	if err != nil {
		response.HTML().BadRequest(errors.New("Unable to decode this URL"))
		return
	}

	client := http.NewClient(string(decodedURL))
	resp, err := client.Get()
	if err != nil {
		logger.Error("[Controller:ImageProxy] %v", err)
		response.HTML().NotFound()
		return
	}

	if resp.HasServerFailure() {
		response.HTML().NotFound()
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	etag := crypto.HashFromBytes(body)

	response.Cache(resp.ContentType, etag, body, 72*time.Hour)
}
