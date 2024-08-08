# Simple novel crawler
## Usage
### Build
```
go build
```
### Run
```
./crawler [novel-slug] [file-name] [max-chapter] [batch-size]
```
## Limitation
This only crawls data as a text file. Need to use an external tool to convert it into other formats like PDF, ...
## Plan
- Separate command per website, e.g. `crawler truyenfull truyenA truyenA 10 5`
- Convert to PDF
## Supported websites
- https://truyenfull.vn
