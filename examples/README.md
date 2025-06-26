# Examples

This directory contains example data and usage examples for the Go Get Images project.

## Contents

- `sample.csv` - Sample CSV file with image URLs for testing
- `README.md` - This file

## Sample CSV File

The `sample.csv` file contains example data with the following structure:

```csv
id,name,image_url,description
1,Image 1,https://example.com/image1.jpg,First image
2,Image 2,https://example.com/image2.png,Second image
3,Image 3,https://example.com/image3.gif,Third image
```

## Usage Examples

### Basic Usage

```bash
# Run with sample data
./go-get-imgs examples/sample.csv 3

# Run with your own CSV file
./go-get-imgs your-data.csv 2
```

### Creating Your Own CSV

Create a CSV file with at least the number of columns specified by the URL column index. The URL column should contain valid image URLs.

Example CSV structure:
```csv
id,title,image_url,description
1,My Image,https://example.com/image.jpg,Description here
2,Another Image,https://example.com/another.png,Another description
```

### Supported URL Formats

- HTTP: `http://example.com/image.jpg`
- HTTPS: `https://example.com/image.png`
- FTP: `ftp://example.com/image.gif`
- File: `file:///path/to/local/image.webp`

### Supported Image Formats

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)
- BMP (.bmp)
- TIFF (.tiff, .tif) 