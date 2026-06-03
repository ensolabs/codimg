Read more on <https://cekrem.github.io/posts/codimg/>, or try the live demo @ <https://codimg.alwaysdata.net>

This project solves a very specific problem, namely being unable to successfully copy/paste code examples into the Webflow WYSIWYG our consultancy blog (sadly?!😅) runs on. Leaving the "why" aside, **it transforms code to a syntax highlighted SVG that you can embed anywhere!**

And it's completely stateless and pure: code in (as part of query param) -> svg out!

Though I cannot promise to always keep this live or to never have a breaking change, in theory this could/would work "forever" without having to ever store any data.

## Usage

Hit `/code.svg` with your code in the `input` query param and (optionally) a `lang` to pick the syntax highlighting:

```
https://codimg.alwaysdata.net/code.svg?lang=go&input=package%20main%0A%0Afunc%20main()%20%7B%7D
```

`input` is just URL-encoded code, so you can embed the result anywhere an image goes:

```html
<img
  src="https://codimg.alwaysdata.net/code.svg?lang=go&input=package%20main%0A%0Afunc%20main()%20%7B%7D"
  alt="code"
/>
```

```markdown
![code](<https://codimg.alwaysdata.net/code.svg?lang=go&input=package%20main%0A%0Afunc%20main()%20%7B%7D>)
```

Renders like this:

![code](<https://codimg.alwaysdata.net/code.svg?lang=go&input=package%20main%0A%0Afunc%20main()%20%7B%7D>)

For larger snippets, `input` also accepts base64-encoded, deflate-compressed code (transparently auto-detected), which keeps URLs short. `lang` accepts any [Chroma](https://github.com/alecthomas/chroma) lexer name (`go`, `python`, `elm`, `js`, ...); leave it off to auto-detect.

You can see it all in practice in the [playground](https://codimg.alwaysdata.net).
