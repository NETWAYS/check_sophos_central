package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ResponseBody struct {
	Items []json.RawMessage `json:"items"`
	Pages PageInfo          `json:"pages"`
}

type PageInfo struct {
	// The key of the first item in the returned page.
	FromKey string `json:"fromKey"`
	// (Optional) The total number of items on all the pages, if pageTotal=true was passed into the request.
	Items int `json:"items"`
	// The maximum page size that can be requested.
	MaxSize int `json:"maxSize"`
	// The key to use when fetching the next page.
	NextKey string `json:"nextKey"`
	// The size of the page being returned.
	Size int `json:"size"`
	// (Optional) The total number of pages that exist, if pageTotal=true in the request.
	Total int `json:"total"`
}

func (c *Client) GetResults(request *http.Request) (items []json.RawMessage, err error) {
	var (
		ctx          = request.Context()
		httpResponse *http.Response
		body         []byte
		nextKey      string
		response     *ResponseBody
	)

	// Set default page size if not set
	if !strings.Contains(request.URL.RawQuery, "pageSize=") {
		if request.URL.RawQuery != "" {
			request.URL.RawQuery += "&"
		}

		request.URL.RawQuery += fmt.Sprintf("pageSize=%d", c.PageSize)
	}

	for {
		r := request.Clone(ctx)
		if nextKey != "" {
			if r.URL.RawQuery != "" {
				r.URL.RawQuery += "&"
			}

			r.URL.RawQuery += "pageFromKey=" + url.QueryEscape(nextKey)
		}

		httpResponse, err = c.Do(r)
		if err != nil {
			return
		}

		// read response body
		body, err = ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			err = fmt.Errorf("could not retrieve response body: %w", err)
			return
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode != 200 {
			log.WithFields(log.Fields{
				"status": httpResponse.StatusCode,
				"body":   string(body),
			}).Debug("HTTP returned non-ok result")

			err = fmt.Errorf("HTTP request returned non-ok status: %s", httpResponse.Status)

			return
		}

		// parse response
		response = &ResponseBody{}

		err = json.Unmarshal(body, response)
		if err != nil {
			err = fmt.Errorf("could not decode JSON from body: %w", err)
			return
		}

		// retrieve items from response
		for _, item := range response.Items {
			items = append(items, item)
		}

		// set nextKey or break iteration when done
		if response.Pages.NextKey == "" {
			break
		} else if response.Pages.NextKey == nextKey {
			err = fmt.Errorf("iteration error in pages, nextKey is the same as fromKey: %s", nextKey)
			return
		} else {
			nextKey = response.Pages.NextKey
		}
	}

	return
}
