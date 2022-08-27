# BeReal TimeLapse

Compiles a timelapse from your BeReal memories.

![](https://i.imgur.com/fltmvVE.gif)

## Usage:

```bash
./br-timelapse render --phone_number "+34XXXXXXXXX" -f 5 -o "bereal.mp4"
```

- `--phone_number`: Phone number linked to your BeReal account
- `--fps`: Frame rate, default is 5
- `--output`: Output filename

```bash
Render timelapse

Usage:
  bereal-timelapse render [flags]

Flags:
  -f, --fps int               Frames per second (default 5)
  -h, --help                  help for render
  -o, --output string         Output filename (default "render.mp4")
  -p, --phone_number string   Phone Number: +XXYYYYYYYYY
```

## Dependencies:

FFmpeg in your path

## Credits:

@NotMarek's [BeFake](https://github.com/notmarek/BeFake) for API endpoints
