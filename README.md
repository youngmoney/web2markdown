# web2markdown

Basic web fetching and conversion to markdown. For simple archiving.

``` bash
web2markdown https://example.com
```

## Advanced

### User Agent

To specify a user agent:

``` bash
web2markdown --user-agent 'Custom Agent' https://example.com
```

### Min Content

Sometimes a website returns an error page or shortened page. To return
an error when that happens, specificy a minimum number of characters in
the final markdown content. If less are present, the command will fail.

``` bash
web2markdown --min-content 200 https://example.com/does-not-exist
```
