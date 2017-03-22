# Grit

## What?

Grit keeps track of your local clones of Git repositories, allowing you to
quickly switch to clone directories in your terminal based on the repository
name or a portion thereof.

## Why?

I spend most of my day working on projects that live in Git repositories. Many
of them are hosted on GitHub.com, but many more are in my employer's private
GitHub Enterprise and BitBucket installations.

Keeping track of hundreds of clones can be a little tedious, so some time back
I adopted a basic directory naming convention and wrote some shell scripts to
handle cloning in a consistent way.

This worked well for a while, until the list of places I needed to clone from
increased, and I started working more heavily in [Go](http://golang.org),
which places it's [own requirements](https://github.com/golang/go/wiki/GOPATH)
on the location of your Git clones.

Grit is the logical evolution of those original scripts into a standalone
project that handles multiple Git sources and Go's peculiarities. It's hacked
together, there are no tests, and there are other more general solutions for
navigating your filesystem; but this works for me. I've published it to GitHub
in case it works for you too.

## How?

If you want to get a sense of what Grit does without having to read or
understand anything, try this:

1. Download `grit` from the [releases page](https://github.com/jmalloc/grit/releases)
1. Clone a repository with `grit clone <repo>` maybe try `grit clone jmalloc/grit`
1. Find a repository with `grit index find <repo>`

Joke's on you! Grit only knows about that one repository. Now you have to read
a bunch of stuff anyway ...

### Configuration

Grit is useful because it knows where you want to clone repositories from. But
it's not a magician, it only knows after you tell it, which is less impressive.

By default, Grit looks for a configuration file at `~/.grit/config.toml`, which
at its most basic is a list of places to clone from:

```toml
[clone.sources]
my-company = "git@git.example.com:*.git"
```

If you only use GitHub you don't need to define any additional sources, and so
the configuration file can be omitted entirely. For a complete list of the
available configuration directives, see the [example](etc/example.toml)
configuration file.

From now on, I'll assume you're using the configuration shown above, which is a
bit silly, but here we are. Rest assured, all paths and filenames used by Grit
are probably 100% configurable, more or less.

### Cloning a repository

There are two parts to Grit. One of them is the `clone` command which clones
repositories. It's like `git clone` but without any of the flexibility or power.
That's minimalism.

The clone command accepts a single argument, the repository "slug". The slug
is the part of the repository URL that a standard-issue human would use to
identify a repository. For GitHub, this is the familiar `username/repository`
syntax. For sources configured in `config.toml`, it's the part represented
by the asterisk in the URL.

Try it out:

    grit clone jmalloc/grit

Grit will print the absolute path of the local clone, which should look
something like:

    /Users/james/grit/github.com/jmalloc/grit

I'm assuming a lot about your system configuration and your given name, just
roll with it.

Grit does a lot of this, printing paths. It does this on `STDOUT` whenever
it thinks you might want to know where something is.

### Querying the index

Whenever you clone a repository with Grit, the repository is added to the index.
The index is a database mapping repository slugs and names to directories. The
index is also that second part of Grit that I mentioned back when we were
talking about cloning. I haven't forgotten.

The index can be queried to find a repository by slug:

    grit index find jmalloc/grit

Or by name alone:

    grit index find grit

Just like the `clone` command, `index find` prints the paths on `STDOUT`. If
there is more than one matching path, it just prints them all. It's relentless.

If nothing is found, Grit exits with a non-zero exit code. The universal signal
that you should re-think your actions.

### What about Go?

We've already cloned `jmalloc/grit` into `~/grit`, but Grit is written in Go,
so it needs to live somewhere special, and that somewhere is `$GOPATH`.

Sure, you could just `mv` the directory; but `mv` isn't written in Go, and
therefore it isn't web-scale!

Try this instead:

    grit clone --go jmalloc/grit

*Again*, Grit ruthlessly prints the clone path to `STDOUT`. This time however,
you'll notice that the clone is in a subfolder of `$GOPATH` instead of `~/grit`.
What a relief!

If you're an avid Gopher, you might be wondering "Why not use `go get`?". Well,
that doesn't update the Grit index, of course!

So now you've got two clones of `jmalloc/grit` on your system. This increases
the counter on my traffic graph, making me feel quite special. Beyond that, it
also brings us to our next example ...

### Selecting between multiple clones

Try this:

    grit cd jmalloc/grit

If you've followed all the steps until now, and I haven't messed up the examples
too badly, you should be presented with a list of matching directories:

      1) [go] src/github.com/jmalloc/grit
      2) [grit] github.com/jmalloc/grit
    >

This time, Grit hasn't printed them to `STDOUT`. Oh, no! This time they've been
sent to `STDERR` and Grit patiently awaits your decision. Enter your selection
(numerically) and press enter.

Grit diligently prints the absolute path to your selection to `STDOUT`.
Interminable!

## Ok, so Grit prints paths to STDOUT, I get it.

Grit is designed to be consumed by shell scripts, take this Bash function,
for example:

```bash
grit-cd() {
    local dir
    dir=$(grit cd $1) && cd $dir
}
```

If you load this function in your `.bash_profile` (or just paste it into your
terminal) and run:

    grit-cd grit

You will be prompted for a choice just as before, but after you make your
selection you will be *whisked* away to the directory that you selected.

This is only possible by running the `cd` command in the interactive login
shell, which is why Grit is all about paths on `STDOUT`. Grit is UNIXY!

## What next?

You can expect auto-completion of command names, parameters and indexed
repository slugs. These are all features present my the shell scripts I
mentioned way back at the beginning of our time together.

I'm in the process of porting the scripts in my [dotfiles repo](https://github.com/jmalloc/dotfiles)
to use Grit instead of the cobbled together, in-memory cache they use currently.

By the way, have you ever tried to write an in-memory cache in Bash? In a
version of Bash available on OS X out of the box? You have? Oh.
