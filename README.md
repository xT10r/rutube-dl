
# RutubeDL - Video Downloader

`RutubeDL` is a Go-based video downloader that allows you to download videos from Rutube using either a list of video IDs or a direct URL link. This tool can be used as both a compiled application and as an external package in your own Go projects.

## Features
- Supports downloading videos using a list of video IDs or direct video links.
- **Multi-language support**: Automatically detects OS language (English/Russian) with cross-platform compatibility.
- **Smart filename handling**: Automatically sanitizes filenames for different operating systems.
- **Transliteration support**: Option to convert Cyrillic characters to Latin equivalents.
- **Automatic FFmpeg management**: Downloads and manages FFmpeg automatically when needed.
- **Organized output**: Creates video-specific directories and generates final MP4 files.
- Uses concurrent downloads for faster performance.
- Implements retry logic to handle transient network issues.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
  - [As a Compiled Application](#as-a-compiled-application)
  - [As an External Package](#as-an-external-package)
- [Flags](#flags)
- [Example Commands](#example-commands)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use this project, you need to have [Go](https://golang.org/dl/) installed on your machine.

1. Clone the repository:
   ```bash
   git clone https://github.com/StanislavKH/rutube-dl.git
   cd rutube-dl
   ```

2. Build the application:
   ```bash
   go build -o rutubedl ./cmd/rutube-dl/main.go
   ```

This will create a binary named `rutubedl` (or `rutubedl.exe` on Windows) in your project directory.

## Usage

### As a Compiled Application

Once the binary is built, you can use it from the command line to download videos by list ID or direct link. You can also specify the directory where files and temporary chunks will be stored.

#### Command-line Options
- `-list_id` : Specify the list ID to download multiple videos.
- `-file_link` : Specify the direct URL to download a single video.
- `-dir` : Specify the directory to store downloaded files and temporary chunks (defaults to `downloads` if not specified).
- `-workers` : Number of concurrent workers for downloading (default: 1, recommended: 3-5).
- `-with_ffmpeg` : Use FFmpeg for more reliable video concatenation (automatically downloads if needed).
- `-from_episode` : Start downloading from a specific episode number (only with `-list_id`).
- `-transliterate` : Convert Cyrillic characters in video names to Latin equivalents.

### Example Commands

- Download videos using a list ID and store them in the `videos` directory:
  ```bash
  ./rutubedl -list_id=123456 -dir=videos
  ```

- Download a video using a direct URL link with transliteration and FFmpeg:
  ```bash
  ./rutubedl -file_link=https://example.com/video-url -transliterate -with_ffmpeg
  ```

- Download from a playlist starting from episode 5 with 3 workers:
  ```bash
  ./rutubedl -list_id=123456 -from_episode=5 -workers=3 -dir=season2
  ```

- Download a video with custom directory and transliteration:
  ```bash
  ./rutubedl -file_link=https://example.com/video-url -dir=custom-directory -transliterate
  ```

### As an External Package

You can also use this project as an external package in your own Go code. To do so, follow these steps:

1. Import the package into your project:
   ```go
   import "github.com/StanislavKH/rutube-dl/pkg/rutubedl"
   ```

2. Use the functions provided by the `rutubedl` package:
   ```go
   package main

   import (
       "log"
       "github.com/StanislavKH/rutube-dl/pkg/rutubedl"
   )

   func main() {
       fileLink := "https://example.com/video-url"
       dir := "videos"
       workers := 5
       withFFmpeg := true
       transliterate := false
       err := rutubedl.DownloadFile(fileLink, &dir, workers, withFFmpeg, transliterate)
       if err != nil {
           log.Fatalf("Failed to download file: %v", err)
       }
   }
   ```

3. Run your Go application as usual:
   ```bash
   go run yourapp.go
   ```

## Flags

| Flag            | Description                                                                            | Example Usage                                      |
|-----------------|----------------------------------------------------------------------------------------|----------------------------------------------------|
| `-list_id`      | The ID of the list to download videos from                                             | `./rutubedl -list_id=123456`                       |
| `-file_link`    | The direct URL of the file to download                                                 | `./rutubedl -file_link=https://...`                |
| `-dir`          | The directory to store files and temporary chunks (optional)                           | `./rutubedl -file_link=https://... -dir=videos`    |
| `-workers`      | The number of workers for chunk download (optional, default: 1)                       | `./rutubedl -file_link=https://... -workers=4`     |
| `-with_ffmpeg`  | Use external ffmpeg for more reliable chunk concatenation (optional)                   | `./rutubedl -file_link=https://... -with_ffmpeg`   |
| `-from_episode` | Download everything starting from this episode number (optional, only with `-list_id`) | `./rutubedl -list_id=123456 -from_episode=5`       |
| `-transliterate`| Convert Cyrillic characters in video names to Latin equivalents (optional)            | `./rutubedl -file_link=https://... -transliterate` |

## Error Handling
If an error occurs during the download process, the program will log the error and proceed to the next item in the list. For individual file downloads, an error will stop the process with an appropriate error message.

## New Features (v2.0)

### Internationalization
- **Automatic language detection**: The application automatically detects your OS language and displays messages accordingly.
- **Supported languages**: English (default) and Russian.
- **Cross-platform**: Works on Windows, Linux, macOS, and other Unix-like systems.

### Smart Filename Handling
- **OS-specific sanitization**: Automatically removes or replaces characters that are not allowed in filenames on different operating systems.
- **Unicode support**: Properly handles Unicode characters in video titles.
- **Length limiting**: Ensures filenames don't exceed filesystem limits.

### Transliteration Support
- **Cyrillic to Latin**: Use the `-transliterate` flag to convert Russian text to Latin characters.
- **Safe filenames**: Ensures compatibility across different systems and reduces encoding issues.

### Automatic FFmpeg Management
- **Multi-source support**: Downloads from multiple reliable sources (BtbN-Builds, FFmpeg-Static)
- **Intelligent build selection**: Prefers static builds to avoid DLL dependency issues
- **DLL handling**: Automatically extracts required DLL files for shared builds on Windows
- **Intelligent fallback**: If download fails, automatically tries system FFmpeg
- **Auto-download**: Automatically downloads the latest FFmpeg version when using `-with_ffmpeg` flag
- **Version-specific storage**: Downloads are stored in `bin/ffmpeg/[version]` directories to avoid re-downloading
- **Version management**: Checks for updates and downloads newer versions automatically
- **Platform detection**: Downloads the appropriate FFmpeg build for your operating system and architecture
- **Robust error handling**: Multiple extraction patterns and detailed debugging information

### Improved Output Structure
- **Direct output**: MP4 files are saved directly to the specified `-dir` parameter instead of creating subdirectories
- **Clean filenames**: Uses `[video_name].mp4` format for final output files
- **Temporary file management**: Creates temporary directories for download segments that are cleaned up automatically

## Testing

This project includes comprehensive tests for all components:

### Unit Tests
Run unit tests for individual packages:
```bash
# Test FFmpeg management
go test -v ./pkg/ffmpeg

# Test filename utilities
go test -v ./pkg/utils

# Test internationalization
go test -v ./pkg/i18n

# Run all unit tests
go test -v ./pkg/...
```

### Integration Tests
Run the comprehensive integration test that includes FFmpeg download testing:
```bash
go run ./test/integration.go
```

This integration test will:
- Test FFmpeg manager initialization and platform detection
- **Download FFmpeg automatically** if not present (saved permanently)
- Test internationalization with language switching
- Test filename sanitization and transliteration
- Verify all components work together correctly

**Note**: The integration test will actually download FFmpeg (about 100MB) if it's not already installed. The downloaded FFmpeg is saved permanently and reused for future tests and application runs.

### Short Tests (Skip Network Operations)
For quick testing without network dependencies:
```bash
go test -short -v ./pkg/...
```

## Contributing
Feel free to submit issues or pull requests for improvements or bug fixes.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
