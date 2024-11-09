# Downloads Folder Organizer

A lightweight Go application that automatically organizes files in your Downloads folder by categorizing them based on their file extensions.

## Features

- Automatically sorts files into appropriate categories:

  - Audios (mp3, wav, flac, etc.)
  - Videos (mp4, avi, mkv, etc.)
  - Images (jpg, png, gif, etc.)
  - Documents (pdf, doc, txt, etc.)
  - Data files (zip, csv, xml, etc.)
  - Executables (exe, msi, app, etc.)
  - Others (uncategorized files)

- Handles duplicate filenames by adding numerical suffixes
- Processes files concurrently for better performance
- Maintains original file names
- Provides logging of operations and errors

## Installation

### Option 1: Pre-built binaries

Download the appropriate zip file for your system from the `distributions` folder:

- Windows: `download-manager-windows.zip`
- macOS: `download-manager-macos.zip`
- Linux: `download-manager-linux.zip`

Extract the zip file and you're ready to go!

### Option 2: Build from source

If you prefer to build the application yourself:

1. Make sure you have Go installed on your system
2. Clone this repository
3. Build the application:

```bash
go build -o download-manager
```

## Usage

Simply run the executable:

```bash
./download-manager
```

The program will:

1. Create category folders in your Downloads directory if they don't exist
2. Move files to their respective category folders based on extension
3. Log any errors that occur during the process

## Requirements

- Go 1.16 or higher (only if building from source)
- Write permissions in your Downloads folder

## Note

The program organizes files in your system's Downloads folder by default. Make sure to backup important files before running the organizer for the first time.

## Contributing

Feel free to open issues or submit pull requests for improvements.

## License

MIT License - feel free to use and modify as needed.
