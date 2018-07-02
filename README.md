# tagzytout

A tool to tag it all !

## how to

In your commits messages, add a line like this:

```
tagzytout: my-tag-name
```

tagzytout finds the line that starts with `tagzytout:` and extracts the following tag name, used to create a tag on that commit.

### install

Run:

```bash
go get github.com/SiegfriedEhret/tagzytout
```

### run

Run the tool like this:

```bash
tagzytout -path=/home/user/gitrepo/
```

### notes

:warning: if you have multiple commits with the same tag, the tag will be applied on the most recent commit.

## it works on my computer

True story.

You'll have to push your tags manually, after checking everything is ok !

## licence

Licenced under the [WTFPL](http://www.wtfpl.net/), see details [here](./LICENSE.md).