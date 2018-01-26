# scrape
Scrape HTML using template strings

Yes, another package for scraping HTML (or any kind of Reader). What sets this one apart from the others is that you can use HTML fragments as a template for targeting the data you want to scrape.

Template strings can be nested or not, and `scrape` will do it's best job at finding the approriate data.

```go
    template := scrape.NewTemplate(`
        <title>{{title}}</title>
        <img src="{{cats|hasCats()}}/>
        <div class="all-items">
            <ul>
                <li>{{item}}</li>
            </ul>
        </div>
    `)

    template.Validator("hasCats", func(v string) bool {
        return strings.Contains(v, "cat")
    })

    resp, _ := http.Get(theURL)
    defer resp.Body.Close()

    data, _ := template.Scrape(resp.Body)

    if title, ok := data["title"]; ok {
        // got the scraped title as title[0]
    }

    if images, ok := data["cats"]; ok {
        for _, img := range images {
            // got pictures of cats!
        }
    }

    // so on and so forth
```

Within the template strings you use `{{...}}` to provide directives for the scraper. By default it's the key/name assigned to the scraped data. 

You can also specify a validator if it's followed by parantheses, such as `{{hasCats()}}`, this validator must also be provided with the same exact name via `template.Validator(string, func(v string) bool)`. 

If a validator fails then that HTML node and all of its nested nodes/tags are skipped. If you want to scrape the data _after_ it's been validated as true then you pipe the two like `{{cats|hasCats()}}`. 

This package works and it's easy setting up by using HTML templates. But I'm not satisfied with how you retrieve the scraped data. Using a map of string arrays is the obvious choice, but it gets pretty verbose when testing for and using the scraped data.

---
I think eventually I want this package to use a `struct` that specifies the template by using tags. Then you can provide the struct to the scraper and it will be populated with the correct data. 

An idea would be like:
```go
type ScrapedData struct {
    Title   string   `scrape:"<title>{{*}}</title>"`
    Cats    []string `scrape:"<img src=\"{{*|hasCats()}}\"/>"`
    // etc
}

data := &ScrapedData{}
scrape.Marshal(reader, data)
```

But, it's still kind of ugly. Chances are good that Go version 2 will address struct tags and may provide a better way of decorating structs, so I'll wait and see what gets rolled out first.