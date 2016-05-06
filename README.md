# Statix

[![Build Status](https://travis-ci.org/sarulabs/statix.svg?branch=master)](https://travis-ci.org/sarulabs/statix)
[![GoDoc](https://godoc.org/github.com/sarulabs/statix?status.svg)](http://godoc.org/github.com/sarulabs/statix)
[![Coverage](http://gocover.io/_badge/github.com/sarulabs/statix?0)](https://gocover.io/github.com/sarulabs/statix)
[![codebeat](https://codebeat.co//badges/0f6455a1-8dde-4da6-aed0-de545bf8f458)](https://codebeat.co/projects/github-com-sarulabs-statix)
[![goreport](https://goreportcard.com/badge/github.com/sarulabs/statix)](https://goreportcard.com/report/github.com/sarulabs/statix)


Statix is an asset manager written in go.

Assets are static files that are needed to run your application but are not part of your code. For example images and javascript files are assets. In web development, theses assets are stored in a public directory where they are accessible to the outside world through a web server.

#### The problem

You can work with your assets directly in your public directory, but this comes with some downsides :

*1) assets are served as they are*

The development version of an asset is also the version used in production. But the needs are different in these two environments.

For example, it is better if your javascript files are minified in production. You also want to reduce the number of javascript files downloaded by the browser to reduce loading time. On the contrary in your development environment you want to use uncompressed javascript and have your code split up in different files for readability.

If you want to combine the benefits of both versions, your working directory should not be the same as your public directory.

*2) the URL of an asset does not change*

The filename of an asset does not change when you update its content. After an update of your website, the client browser may still have the old version of an asset in cache. As its URL has not changed, the browser combines the old version of the asset (javascript, css, ...) with the new HTML. This may be catastrophic.


#### What can statix do ?

Statix can handle these two problems. It allows you to :

- export assets from your working directory to a public directory where your static files are accessible through a web server
- combine assets in one file (javascript files for example)
- apply modifications to your assets (optimization, minification, compilation or anything you can imagine)
- avoid browser cache problems (an asset filename will contain the md5 hash of its content and thus it will be unique)


## Todo list

- Add more unit tests
- Add more alterations

## Documentation

+ [Configuration example](#configuration-example)
+ [Dumping assets](#dumping-assets)
+ [Getting URLs](#getting-urls)
+ [Getting paths](#getting-paths)
+ [Manager.Server and Manager.Servers](#managerserver-and-managerservers)
+ [Manager.Assets](#managerassets)
    - [SingleAsset](#singleasset)
        * Output
        * Input
    - [AssetPack](#assetpack)
        * Output
        * Input
        * Pattern
        * Alterations
+ [Manager.Filters](#managerfilters)
+ [Manager.Input and Manager.Output](#managerinput-and-manageroutput)


## Configuration example

Your assets are defined in a Manager. Here is an example of configuration :

```go
import (
    "github.com/sarulabs/statix"
    "github.com/sarulabs/statix/alteration"
    "github.com/sarulabs/statix/resource"
)

manager := statix.Manager{
    Server: statix.Server{
        // link a directory to a base URL
        Directory: "/output/directory",
        URL:       "http://example.com/static",
    },
    Filters: []statix.Filter{
        {
            // applies uglifyjs (a javascript minifier) to every .js asset
            Alteration: alteration.NewUglifyJs("/usr/local/bin/uglifyjs"),
            Pattern:    statix.NewExtensionPattern("js"),
        },
    },
    // asset definitions, the key of the map is the name of the asset.
    // an asset can be a SingleAsset (1 output file)
    // or an AssetPack (export all the files from a directory to another)
    Assets: map[string]statix.Asset{
        // exports .png files from /input/directory/img to /output/directory/img
        // and applies optipng to all files before writing them
        "images": statix.AssetPack{
            Input:   "/input/directory/img",
            Output:  "/output/directory/img",
            Pattern: statix.NewExtensionPattern("png"),
            Alteration: alteration.NewOptiPng("/usr/bin/optipng", 100),
        },
        // creates an app.js file that contains jquery
        // and a typescript file (app.ts) complied in javascript
        "app-js": statix.SingleAsset{
            Output: "/output/directory/app.js",
            Input: resource.NewCollection(
                resource.NewFile("/input/directory/js/jquery.js"),
                resource.NewAlteredResource(
                    resource.NewFile("/input/directory/app.ts"),
                    alteration.NewTypeScript("/usr/local/bin/tsc"),
                ),
            ),
        },
    },
}
```

You can learn more about each attribute of the Manager in the following parts of the documentation.


## Dumping assets

To export your assets, you simply call the Dump method.

```go
manager.Dump()
```

For each asset, two files are created. For example, an asset that you expect to be dumped in `/output/directory/app.js` will be composed of these two files :
- `/output/directory/app.{MD5}.js` : the actual content of the asset (`{MD5}` being the md5 hash of the file).
- `/output/directory/app.js` : a symlink to the other file.

The md5 hash is added to the filename to avoid cache problems in browsers. The symlink may seem useless but is in fact necessary to get the path of the asset without knowing the md5 hash of its content.


## Getting URLs

You can get the URL of an asset thanks to the `URL` method. For a SingleAsset, you call `URL` with the name of the asset (the key of the map Manager.Assets).

```go
manager.URL("app-js")
```

But for an AssetPack, the name of the asset is not enough. That is because an AssetPack may contain more than one asset. Thus you also need to specify the path of the file in the asset directory.

```go
manager.URL("images", "header/logo.png")
```

With the manager given in example, this will return the URL of the file named `/output/directory/img/header/logo.png`.

For this to work, Manager.Server or Manager.Servers needs to be defined properly. Symlinks to asset files without md5 hash also need to exist.


## Getting paths

You can also get the path of an asset. `Symlink` works the same way `URL` does but returns the filename of the symlink.

```go
symlinkAppJs, _ := manager.Symlink("app-js") // for a SingleAsset
symlinkLogo, _ := manager.Symlink("images", "header/logo.png") // for an AssetPack
```

If you want to know the filename of the asset with the md5 hash, you will have to evaluate the symlink.

```go
filename, _ := filepath.EvalSymlinks(symlink)
```


## Manager.Server and Manager.Servers

Be aware that statix is not a web server. It only dumps assets in a directory where static files can be served through a web server.

However there is a link between the path of a static file and its URL. That is why if you set an URL for your asset directory, statix will be able to give you the URL of every asset in this directory.

In case you serve your assets in only one directory, the configuration will look like this :

```go
statix.Manager{
    Server: statix.Server{
        Directory: "/output/directory",
        URL:       "http://example.com/static",
    },
    // ...
}
```

The URL of `/output/directory/img/header/logo.{MD5}.png` is `http://example.com/static/img/header/logo.{MD5}.png`.

You are not restricted to one server. You can define as many servers as you need :

```go
statix.Manager{
    Servers: []statix.Server{
        {
            Directory: "/output/directory-1",
            URL:       "http://static1.example.com",
        },
        {
            Directory: "/output/directory-2",
            URL:       "http://static2.example.com",
        },
    },
    // ...
}
```


## Manager.Assets

Assets are defined in a map. The key is the name of the asset. The value may be a SingleAsset or an AssetPack.


### SingleAsset

A SingleAsset may be composed of a complex input but will be dumped in one file. It is useful if you need to :
  - combine assets (javascript files or stylesheets)
  - load an asset from a string and not from a file

#### Output

It is the name of the file where the content of the asset will be dumped.

#### Input

The input of a SingleAsset is a Resource. There are different kinds of Resource you may want to use as an input. Their definition is in the `resource package`. They all provide a `Dump` method to get their content.

A resource can be a **string** :

```go
r := resource.NewString("my-resource")
c, err := r.Dump() // c is []byte("my-resource")

resource.NewBytes([]byte("my-resource")) // works too
```

But you probably have your resources stored in **files** :

```go
r := resource.NewFile("/path/to/file")
c, err := r.Dump() // c is the content of the file
```

Sometimes you need to combined **multiple resources** to create a new one :

```go
r := resource.NewCollection(
    resource.NewString("resource-1"),
    resource.NewString("resource-2"),
)
c, err := r.Dump() // c is []byte("resource-1resource-2")
```

You may also want to alter the content of a resource using an **alteration** :

```go
r := resource.NewAlteredResource(
    resource.NewString("my-resource"),
    alteration.Reverse{}, // an alteration, Reverse does not exists but is used for comprehension
)
c, err := r.Dump() // c would be []byte("ecruoser-ym")
```

Alterations like `alteration.Reverse{}` are structures that implement the `Alteration` interface from the `resource package`. They are used to modify the content of an asset. They should have an `Alter` method that takes a resource and returns a new altered one.

You can find some alterations in the `alteration package`, but it is also really easy to create your own. In the `alteration package` you will find the alterations for theses programs :
- jpegoptim
- optipng
- stylus
- typescript
- uglifycss
- uglifyjs


### AssetPack

An asset pack should be used if you want to export the content of a directory to an other. If you only have one asset, it is probably better to use a SingleAsset.

#### Output

This is the directory where you want to dump your assets.

#### Input

This is the directory where your raw assets are located.

#### Pattern

`Pattern` is just a wrapper on top of regular expressions. It will allow you to omit some files in the input directory. To only export files ending with `.ext` can for example use this pattern :

```go
statix.NewPattern("\\\.ext$")
```

As filtering files from their extensions is a frequent use case, you may want to use the `NewExtensionPattern` function.

```go
statix.NewExtensionPattern("png", "jpg") // will filter .png and .jpg files
```

#### Alterations

You can apply a list of alterations to all the assets in the AssetPack :

```go
statix.AssetPack{
    Input:   "/input",
    Output:  "/output",
    Pattern: statix.NewExtensionPattern("png"),
    Alterations: []statix.Alteration{
        alteration.OptiPng("/usr/bin/optipng", 100),
    }
}
```

This will take all `.png` files in `/input` and apply optipng to their content before writing the results in the `/output` directory.

## Manager.Filters

Filters allow you to apply some alterations to all your assets. A Filter is composed of an Alteration and a Pattern. The Pattern limits the effect of the Alteration to the files matching its regular expression.

In this example, `uglifyjs` is executed on every `.js` files :

```go
statix.Manager{
    Filters: []statix.Filter{
        {
            Alteration: alteration.NewUglifyJs("/usr/local/bin/uglifyjs"),
            Pattern:    statix.NewExtensionPattern("js"),
        },
    },
    // ...
}
```

Filters are applied after alterations contained in SingleAsset and AssetPack. Minifying or optimizing files (javascript, css, images) depending on their extension is probably the main use case of Filters.


## Manager.Input and Manager.Output

Manager.Input defines the root of all relative input paths. Identically Manager.Output defines the root of all relative output paths. These root paths apply for all assets and server directories. More precisely, if a server directory, an asset input or an asset output is relative, its path will be rewritten to be based in the Manager.Input or Manager.Output directory. But if its path is absolute nothing will happen and the asset or the server will stay the same.

This example is the same as the first one but with relative asset paths. Server directory is also relative because in this case an empty string equals `.`.

```go
statix.Manager{
    Input:  "/input/directory",
    Output: "/output/directory",
    Server: statix.Server{
        // Directory: "/output/directory"
        Url: "http://example.com/static",
    },
    Assets: map[string]statix.Asset{
        "imgs": statix.AssetPack{
            Input:   "img", // /input/directory/img
            Output:  "img", // /output/directory/img
            Pattern: statix.NewExtensionPattern("jpg", "png"),
        },
        "app-js": statix.SingleAsset{
            Output: "app.js", // /output/directory/app.js
            Input: asset.NewCollection(
                asset.NewFile("js/jquery.js"), // /input/directory/js/jquery.js
                asset.NewFiltering(
                    asset.NewFile("app.ts"), // /input/directory/app.ts
                    filter.NewTypeScript("/usr/local/bin/tsc"),
                ),
            ),
        },
    },
}
```

Manager.Input and Manager.output are not mandatory. There are just here to simplify the Manager definition.
