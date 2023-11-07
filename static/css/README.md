# Setting Up TailwindCSS

Using Node, install tailwind and initialize it:

```bash
npm install -D tailwindcss
npx tailwindcss init
```

Add the path to the template files to the `content` array in `tailwind.config.js`.

Create an `input.css` file and fill it with whatever tailwind classes required, e.g.:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Finally, start a watcher:

```bash
npx tailwindcss -i input.css -o output.css --watch
```

To create a minified css file, use:

```bash
npx tailwindcss -i input.css -o output.css --minify
```

Include in your templates the `output.css` file:

```html
<link rel='stylesheet' type='text/css' href='../css/output.css'>
```

# Adding DaisyUI

To add DaisyUI, simply run

```bash
npm i -D daisyui@latest
```

Then add DaisyUI to the `module.exports` plugins field:

```js
module.exports = {
  //...
  plugins: [require("daisyui")],
}
```