---  
title: "网站内容创作指南"  
date: 2025-07-12  
author: "Woo Alex"  
description: "如何使用 Go + Hugo 平台进行内容创作与管理"  
tags: ["hugo", "markdown", "go", "cms"]  
categories: ["指南"]  
---
## 1 Custom Music URL {#custom-music-url}
{{< version 0.2.10 >}} The complete usage of [local resource references](../theme-documentation-content#contents-organization) is supported.
The `music` shortcode has the following named parameters by custom music URL:
* **server** *[required]*
    URL of the custom music.
* **name** *[optional]*
    Name of the custom music.
* **artist** *[optional]*
    Artist of the custom music.
* **cover** *[required]*
    URL of the custom music cover.
Example `music` input by custom music URL:
```markdown
{{</* music url="/music/Wavelength.mp3" name=Wavelength artist=oldmanyoung cover="/images/Wavelength.jpg" */>}}
```
The rendered output looks like this:
{{< music url="/music/Wavelength.mp3" name=Wavelength artist=oldmanyoung cover="/images/Wavelength.jpg" >}}
## 2 Music Platform URL Automatic Identification {#automatic-identification}
The `music` shortcode has one named parameter by music platform URL automatic identification:
* **auto** *[required]* (**first** positional parameter)
    URL of the music platform URL for automatic identification,
    which supports `netease`, `tencent` and `xiami` music platform.
Example `music` input by music platform URL automatic identification:
```markdown
{{</* music auto="https://music.163.com/#/playlist?id=60198" */>}}
Or
{{</* music "https://music.163.com/#/playlist?id=60198" */>}}
```
The rendered output looks like this:
{{< music auto="https://music.163.com/#/playlist?id=60198" >}}
## 3 Custom Server, Type and ID {#custom-server}
The `music` shortcode has the following named parameters by custom music platform:
* **server** *[required]* (**first** positional parameter)
    [`netease`, `tencent`, `kugou`, `xiami`, `baidu`]
    Music platform.
* **type** *[required]* (**second** positional parameter)

    [`song`, `playlist`, `album`, `search`, `artist`]

    Type of the music.

* **id** *[required]* (**third** positional parameter)

    Song ID, or playlist ID, or album ID, or search keyword, or artist ID.

Example `music` input by custom music platform:

```markdown
{{</* music server="netease" type="song" id="1868553" */>}}
Or
{{</* music netease song 1868553 */>}}
```

The rendered output looks like this:

{{< music netease song 1868553 >}}

## 4 Other Parameters {#other-parameters}

The `music` shortcode has other named parameters applying to the above three ways:

* **theme** *[optional]*

    {{< version 0.2.0 changed >}} Main color of the music player, default value is `#448aff`.

* **fixed** *[optional]*

    Whether to enable fixed mode, default value is `false`.

* **mini** *[optional]*

    Whether to enable mini mode, default value is `false`.

* **autoplay** *[optional]*

    Whether to autoplay music, default value is `false`.

* **volume** *[optional]*

    Default volume when the player is first opened, which will be remembered in the browser, default value is `0.7`.

* **mutex** *[optional]*

    Whether to pause other players when this player starts playing, default value is `true`.

The `music` shortcode has the following named parameters only applying to the type of music list:

* **loop** *[optional]*

    [`all`, `one`, `none`]

    Loop mode of the music list, default value is `none`.

* **order** *[optional]*

    [`list`, `random`]

    Play order of the music list, default value is `list`.

* **list-folded** *[optional]*

    Whether the music list should be folded at first, default value is `false`.

* **list-max-height** *[optional]*

    Max height of the music list, default value is `340px`.
