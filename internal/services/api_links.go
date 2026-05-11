package services

import "net/url"

func getPageFromAPILink(urlParse string) string {
	u, _ := url.Parse(urlParse)

	page := u.Query().Get("page")

	return page
}
