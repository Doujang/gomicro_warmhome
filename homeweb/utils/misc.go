package utils

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_img_addr + "/" + url

	return domain_url
}
