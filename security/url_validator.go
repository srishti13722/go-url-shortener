package security

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

// list of known spam or phishing domains
var blacklistedDomains = []string{
	"example-spam.com",
	"malware-site.net",
	"phishing-link.org",
}

func IsValidURL(url string) bool{
	//check if url is well formed
	if !govalidator.IsURL(url){
		return false
	}

	//extract domain
	domain := extractDomain(url)

	//chcek if domain is in the blacklist

	for _, blocked := range blacklistedDomains{
		if strings.Contains(domain, blocked){
			return false
		}
	}

	return true
}

//extractDomain extracts the main domain from a URL
func extractDomain(url string) string {
	urlParts := strings.Split(url, "/")
	if len(urlParts) > 2 {
		return urlParts[2] // extracts "domain.com" from "https://domain.com/path"
	}
	return url
}