package toolbox

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

// 爬虫工具
type crawlerUtils struct {
}

func (u *utils) NewCrawlerUtils() *crawlerUtils {
	return &crawlerUtils{}
}

func (craw *crawlerUtils) Client() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(NewUtils().NewUAUtils().GetRandomUA()),
		//colly.AllowedDomains(conf.URL),
	)
	return c

}

func (craw *crawlerUtils) ProxyClient(proxyUrl string) (*colly.Collector, error) {
	c := colly.NewCollector(
		colly.UserAgent(NewUtils().NewUAUtils().GetRandomUA()),
		//colly.AllowedDomains(conf.URL),
	)
	err := c.SetProxy(proxyUrl)
	return c, err
}

func (craw *crawlerUtils) ProxySwitcherClient(proxyUrls ...string) (*colly.Collector, error) {
	c := colly.NewCollector(
		colly.UserAgent(NewUtils().NewUAUtils().GetRandomUA()),
		//colly.AllowedDomains(conf.URL),
	)
	p, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	if err != nil {
		return c, err
	}
	c.SetProxyFunc(p)
	return c, nil
}
