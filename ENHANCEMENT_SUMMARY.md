# Enhancement Summary: Playlist Numbering and Improved FFmpeg Error Handling

## 🎯 **Implemented Solutions**

### 1. **Playlist Numbering System**
- **Problem**: Playlists were downloaded without numbering, making sorting difficult
- **Solution**: Automatic file numbering in playlists

#### Features:
- **Smart Index Formatting**: Format adapts to playlist size:
  - 1-9 files: `1 - Video.mp4`
  - 10-99 files: `01 - Video.mp4`
  - 100-999 files: `001 - Video.mp4`
  - 1000+ files: `0001 - Video.mp4`
- **Improved Filtering**: Proper episode filtering by `from_episode` parameter
- **Better Logging**: Detailed information about playlist download process

#### New Functions Added:
```go
// Formats index based on total count
func FormatPlaylistIndex(index, total int) string

// Generates filename with playlist number
func GetPlaylistFilename(index, total int, title, extension string, transliterate bool) string

// Downloads file from playlist with proper numbering
func DownloadPlaylistFile(fileLink string, customOutputDir *string, numWorkers int, withFfmpeg bool, transliterate bool, index, total int, customTitle string) error
```

### 2. **Enhanced FFmpeg Error Handling**
- **Problem**: Error `exit status 0xfffffffe` without diagnostics
- **Solution**: Improved diagnostics and more reliable merging method

#### Improvements:
- **Concat Demuxer**: Uses `concat demuxer` instead of `concat protocol`
- **File List Method**: Creates temporary file with segment list
- **Error Diagnostics**: Detailed stdout/stderr FFmpeg logs
- **Segment Validation**: Check existence and size of segments
- **Path Escaping**: Proper path escaping for FFmpeg

#### Enhanced Error Reporting:
```go
// New improved mergeSegmentsWithFfmpeg function:
- Check existence of all segments
- Create temporary file list
- Capture stdout/stderr for diagnostics
- Check empty/corrupted segments
- Detailed error messages
```

## 🧪 **Testing Results**

### Playlist Numbering Test:
```
Small playlist (5 items):
  1 - Video.mp4
  
Medium playlist (50 items):
  01 - Video.mp4
  
Large playlist (500 items):
  001 - Video.mp4
  
Huge playlist (5000 items):
  0001 - Video.mp4
```

### Unit Tests:
- ✅ **FormatPlaylistIndex**: All test cases pass
- ✅ **GetPlaylistFilename**: Russian/transliteration tests pass
- ✅ **All existing tests**: No regressions introduced

## 📁 **Updated File Structure**

### Before:
```
downloads/
├── Video1.mp4
├── Video2.mp4
└── Video3.mp4
```

### After (Playlist):
```
downloads/
├── 01 - Video1.mp4
├── 02 - Video2.mp4
└── 03 - Video3.mp4
```

## 🔧 **Updated Logic Flow**

### Playlist Download:
1. **Fetch playlist items** → `GetItemsListFromFeedURI()`
2. **Filter by episode** → Improved filtering logic
3. **Calculate total count** → For proper numbering
4. **Download with numbering** → `DownloadPlaylistFile()` with index
5. **Generate numbered filenames** → `GetPlaylistFilename()`

### FFmpeg Merge:
1. **Validate segments** → Check existence and size
2. **Create file list** → Temporary concat list
3. **Execute FFmpeg** → With concat demuxer
4. **Capture diagnostics** → Full stdout/stderr logging
5. **Report detailed errors** → Enhanced error messages

## 🎯 **Benefits**

### Playlist Numbering:
1. **Organized Downloads**: Files automatically sorted in correct order
2. **Easy Navigation**: Clear episode numbering for large playlists
3. **Flexible Formatting**: Index format adapts to playlist size
4. **Better Filtering**: Proper handling of `from_episode` parameter

### FFmpeg Error Handling:
1. **Reliable Merging**: Concat demuxer more stable than concat protocol
2. **Better Diagnostics**: Detailed error reporting for troubleshooting
3. **Segment Validation**: Early detection of corrupted/missing segments
4. **Path Safety**: Proper escaping handles special characters in paths

## 🚀 **Ready for Production**

Both enhancements are fully implemented and tested:
- ✅ Code changes complete
- ✅ Unit tests passing (100% success rate)
- ✅ Integration tests successful
- ✅ Documentation updated
- ✅ No breaking changes introduced
- ✅ Backward compatibility maintained

## 💡 **Usage Examples**

### Download Playlist with Numbering:
```bash
./rutubedl -list_id=817135 -dir=videos -with_ffmpeg -transliterate
```

### Result:
```
01 - Little Cars Turned Into Powerful EXCAVATORS!.mp4
02 - Helper Trucks Became Colorful! Clean Up Trash and Build Roads!.mp4
03 - A Stone Flew Into the Taxi Window! What Happened to the Kitten?.mp4
...
```

### Enhanced Error Diagnostics:
```
FFmpeg command failed: ffmpeg -f concat -safe 0 -i filelist.txt -c copy output.mp4
FFmpeg stderr: [detailed error message]
Number of segments: 150
Warning: Segment 45 is empty: segment_045.ts
```