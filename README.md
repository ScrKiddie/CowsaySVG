# CowsaySVG

<div align="center">
  <img 
       src="https://cowsay-svg.vercel.app/?colors=%23FF6B6B,%23FFD93D,%236BCB77,%234D96FF,%23A66DD4,%23FFB5E8,%23FF9CEE,%23FF6B6B&duration=4&timing=linear&randomCow=true" 
       style="max-height: 500px; height: auto; width: auto;"
     />
</div>
<br/>

CowsaySVG transforms ASCII art from [cowsay](https://github.com/Code-Hex/Neo-cowsay) into colorful, scalable vector graphics with optional animations. Perfect for websites, READMEs, and dashboards.
## Features

- **Dynamic Color Control**: Specify any number of colors for gradients
- **Smooth Animations**: Fully customizable timing and duration
- **40+ Cowsay Characters**: All characters from [cowsay](https://github.com/Code-Hex/Neo-cowsay) supported
- **API Integration**: Fallback text source support
- **Precise Layout Control**: Adjust spacing and dimensions pixel-perfectly

## Options

| Parameter       | Description                                                            | Default                        | Examples                               |
|-----------------|------------------------------------------------------------------------|--------------------------------|----------------------------------------|
| `text`          | Custom message                                                         | _(empty, uses API_URL if set)_ | `Hello%20World`                          |
| `cow`           | Cowsay character from [cowsay](https://github.com/Code-Hex/Neo-cowsay) | `default`                      | `tux`, `dragon`                        |
| `colors`        | Comma-separated colors for gradient (hex/css names)                    | `%23000000` (black)            | `red,blue`, `%2300ff00,%2300cc00`      |
| `duration`      | Animation duration in seconds (0 = static)                             | `0` (static)                   | `0`, `2.5`                             |
| `timing`        | CSS animation timing function                                          | `steps(1,end)`                 | `linear`, `ease-in-out`                |
| `cascadeDirection` | Sets the cascade direction for character color animation (options: orthogonal, diagonal, center-based). | `rtl`                          | `ltr`, `rtl`, `ttb`, `btt`, `diag-tlbr`, `diag-trbl`, `diag-bltr`, `diag-brtl`, `center-out`, `edges-in` |
| `ballonWidth`   | Max characters per line in bubble                                      | `40`                           | `30`, `80`                             |
| `charWidth`     | Horizontal character spacing (higher = wider)                          | `10`                           | `8`, `15`                              |
| `lineHeight`    | Vertical line spacing (higher = taller output)                         | `20`                           | `16`, `24`                             |
| `eyes`          | Customizes the cow's eyes (max 2 characters)                           | Cow's default                  | `oo`, `^^`, `xx`                       |
| `tongue`        | Customizes the cow's tongue (max 2 characters)                         | Cow's default                  | `U`, `_`, `vv`                         |
| `think`         | Displays message as a thought bubble                                   | `false` (speech bubble)        | `true`, `1`                            |
| `thoughtsChar`  | Custom character for thought bubble trail                              | `\`                            | `*`, `.`                               |
| `noWrap`        | Disables automatic word wrapping in bubble                             | `false` (wrapping enabled)     | `true`, `1`                            |
| `randomCow`     | Selects a random cow character and will overrides `cow` parameter      | `false`                        | `true`, `1`                            |

## Deploy to Vercel

Click the button below to instantly deploy this project to Vercel:

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/ScrKiddie/CowsaySVG)

## Environment
| Variable  | Description                                                                        | Example Value                     |
|-----------| ---------------------------------------------------------------------------------- | --------------------------------- |
| `API_URL` | Fallback URL for text (must return plain text) | `https://api.quotable.io/random` |
| `MAX_TEXT_LENGTH` | Max characters for input text (query/API). Truncated if exceeded. Unlimited if not set | `150` |


