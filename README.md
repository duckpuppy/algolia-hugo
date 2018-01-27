## Algolia search uploader for Hugo

[Algolia](https://www.algolia.com) is a hosted search platform geared towards
business users, but provides a free limited Community tier of their product.
It integrates with Hugo quite well if you instruct Hugo to build a JSON site
index file.

First, sign up for an Algolia account.  Create an index for your website. Take
note of your Admin API key and App ID.  You will need both of those, as well
as the name of the index you created, to use this tool.

The configuration file for this tool lives at
`$XDG_CONFIG_HOME/algolia-hugo/algolia-hugo.yaml`, which by default points to
`$HOME/.config/algolia-hugo/algolia-hugo.yaml`.  Plase the following contents
into that file.

```yaml
algolia_app_id: <your app id>
algolia_api_key: <your api key>
algolia_index_name: <your index name>
```

## Usage 

This tool has very little functionality right now.  There are three commands.

### update

This will update your search objects in the Algolia index.  It does this by
deleting the existing objects and uploading the entire index again. In this way
it alleviates the problem of old search objects that may have been edited or
deleted remaining in the index. Future versions of this tool may support
updating existing objects when I find both time and need.

By default this tool will look for a file named `public/index.json` relative
to the current directory. This can be overridden with the `-f` or `--file`
arguments to the `update` command.

### clear

This command simply clears your search index on Algolia, leaving you with an
empty index.

### config

This command is used to show the config file found by this tool, as well as
showing the configuration fields read in.


### help

Full help is available by running the `help` command, or by executing the tool
without any subcommands.
