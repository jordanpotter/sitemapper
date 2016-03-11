# SiteMapper
Produces a site map that contains the publicly reachable links and assets for all pages within the specified domain.

A URL is considered to be in the specified domain if the protocol and host match exactly. For example, `https://foo.com` and `https://foo.com/docs` are considered to be in the same domain, however `http://foo.com` and `https://bar.com` are not.

At each step of the web crawl, we retrieve the HTML content for the page and parse it for all links and assets. These individual page maps are compiled together to create the final site map. Note that while _all_ links and assets are included in a page map, only links that belong to the specified domain are crawled and thus produce their own page map.

Two methods are provided to create a site map for a particular domain, which are detailed below.

## CLI
A CLI is provided to create a site map easily from the command line. Assuming this package has been installed via `go install`

	cli --site INITIAL_URL --workers NUM_WORKERS --file OUTPUT_FILE

As an example, to produce a site map for `foo.com` using `100` workers and save it to `sitemap.json`

	cli --site https://foo.com --workers 100 --file sitemap.json

## API
A REST API has also been provided. Assuming this package has been installed via `go install`

	api --port PORT

As an example, to run the server on port `8000`

	api --port 8000

With the server running, site maps are created via a simple `GET` request

	GET http://localhost:8000/sitemap?site=INITIAL_URL&workers=NUM_WORKERS

As an example, to produce a site map for `foo.com` using `100` workers

	GET http://localhost:8000/sitemap?site=https://foo.com&workers=100

## Prototype - GUI
When running the API server, additionally specify the path to the static `gui` directory of this repository. For example

	api --port 8000 --static sitemapper/gui

With the API server running, visit `http://localhost:8000` in your preferred browser. After filling out the form to specify the initial URL and number of workers, the site map will be displayed as an interactive graph.

Dark blue nodes represent web pages that belong to the domain, while gray nodes represent pages outside the domain. Light blue nodes are assets, such as images and javascript libraries. Hovering over a node displays the URL for that node.

![site map example](http://i.imgur.com/rQbMVyb.png)
