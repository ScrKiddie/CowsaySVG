# CowsaySVG

<div align="center">
  <img 
       src="https://cowsay-svg.vercel.app/?colors=red,orange,yellow,green,blue,indigo,violet&duration=3" 
       style="max-height: 500px; height: auto; width: auto;"
     />
</div>
<br/>

CowsaySVG transforms ASCII art from [neo-cowsay](https://github.com/Code-Hex/Neo-cowsay) into colorful, scalable vector graphics with optional animations. Perfect for websites, READMEs, and dashboards.

## Features

- **Dynamic Color Control**: Specify any number of colors for gradients
- **Smooth Animations**: Fully customizable timing and duration
- **40+ Cowsay Characters**: All characters from [neo-cowsay](https://github.com/Code-Hex/Neo-cowsay) supported
- **API Integration**: Fallback text source support
- **Precise Layout Control**: Adjust spacing and dimensions pixel-perfectly

## Options

| Parameter       | Description                                                                   | Default               | Examples                          |
|-----------------|-------------------------------------------------------------------------------|-----------------------|-----------------------------------|
| `text`          | Custom message                                                  |                       | `Hello%20World`                     |
| `cow`           | Cowsay character from [neo-cowsay](https://github.com/Code-Hex/Neo-cowsay)   | `default`             | `tux`, `dragon`                   |
| `colors`        | Comma-separated colors for gradient (hex/css names)                   | `%23000000` | `red,blue`, `%2300ff00,%2300cc00` |
| `duration`      | Animation duration in seconds (0 = static)                                   | `1.0`                 | `0` (static), `2.5`               |
| `timing`        | CSS animation timing function                                                | `steps(1,end)`        | `linear`, `ease-in-out`           |
| `ballonWidth`   | Max characters per line                                                    | `40`                  | `30`, `80`                        |
| `charWidth`     | Horizontal character spacing (higher = wider)                                                      | `10`                  | `8`, `15`                         |
| `lineHeight`    | Vertical line spacing (higher = taller output)                                                           | `20`                  | `16`, `24`                        |

## Deploy to Vercel

Click the button below to instantly deploy this project to Vercel:

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/ScrKiddie/CowsaySVG)

## Environment
| Variable  | Description                                                                        | Example Value                     |
|-----------| ---------------------------------------------------------------------------------- | --------------------------------- |
| `API_URL` | Fallback URL for text (must return plain text) | `https://api.quotable.io/random` |




