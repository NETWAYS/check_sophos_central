package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
		response     ResponseBody
	)

	for {
		r := request.Clone(ctx)
		// TODO: cleaner way to do this?
		r.URL.RawQuery += "&pageFromKey=" + url.QueryEscape(nextKey)

		httpResponse, err = c.HttpClient.Do(r)
		if err != nil {
			err = fmt.Errorf("HTTP request failed: %w", err)
			return
		}

		// read response body
		body, err = ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode != 200 {
			err = fmt.Errorf("HTTP request returned non-ok status: %s", httpResponse.Status)
			return
		}

		// parse response
		err = json.Unmarshal(body, &response)
		if err != nil {
			return
		}

		// retrieve items from response
		for _, item := range response.Items {
			items = append(items, item)
		}

		// set nextKey or break iteration when done
		nextKey = response.Pages.NextKey
		if nextKey == "" {
			break
		}
	}

	return
}
