package revolutiontt

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	//"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const CAT_HDTV = "42" // HD Tv Shows
const CAT_BR = "10"   // BlueRay Movies

type RevolutionTT struct {
	CookieJar *cookiejar.Jar
	Client    http.Client
	username  string
	password  string
}

type TorrentInfo struct {
	Title       string
	CatId       int
	Category    string
	DetailUrl   string
	DownloadUrl string
}

func (r *RevolutionTT) Connect(username string, password string) {
	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}

	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}

	r.CookieJar = jar
	r.Client = http.Client{Jar: r.CookieJar, Timeout: 0}
	r.username = username
	r.password = password

	r.login()
}

func (r *RevolutionTT) login() {

	_, err := r.Client.Get("https://revolutiontt.me/login.php")
	if err != nil {
		log.Fatal(err)
	}

	values := make(url.Values)
	values.Set("username", r.username)
	values.Set("password", r.password)
	values.Set("submit", "login")

	_, err = r.Client.PostForm("https://revolutiontt.me/takelogin.php", values)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RevolutionTT) Search(search string) ([]TorrentInfo, error) {

	results := []TorrentInfo{}

	movies, err := r.SearchByCategory(search, CAT_BR)
	shows, err := r.SearchByCategory(search, CAT_HDTV)

	results = append(movies, shows...)

	return results, err
}

func (r *RevolutionTT) SearchByCategory(search string, category string) ([]TorrentInfo, error) {
	results := []TorrentInfo{}

	v := url.Values{}
	v.Set("search", search)
	v.Add("cat", category)
	v.Add("titleonly", "1")

	resp, err := r.Client.Get("https://revolutiontt.me/browse.php?" + v.Encode())
	if err != nil {
		return results, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		return results, err
	}

	doc.Find("#torrents-table tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		t := TorrentInfo{Title: s.Find("td.br_right b").Text()}
		t.DownloadUrl, _ = s.Find("td:nth-of-type(4) a").First().Attr("href")
		results = append(results, t)
	})

	return results, nil
}
