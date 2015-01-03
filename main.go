package main

import (
    "fmt"
    "os"
    "io"
    "path/filepath"
    "github.com/kr/fs"
    "crypto/sha1"
    "crypto/md5"
    "flag"
    "math"
    "hash"
)

const filechunk = 8192

func hasher(filename string, method string) string {
    file, err := os.Open(filename)
    if err != nil {
        panic(err.Error())
    }
    defer file.Close()
    info, _ := file.Stat()
    filesize := info.Size()
    blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
    var h hash.Hash
    switch method {
    case "sha1":
        h = sha1.New()
    case "md5":
        h = md5.New()
    default:
        h = sha1.New()
    }
    for i := uint64(0); i < blocks; i++ {
        blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
        buf := make([] byte, blocksize)
        file.Read(buf)
        io.WriteString(h, string(buf)) //append to the hash
    }
    return fmt.Sprintf("%X", h.Sum(nil))
}

func marcher(dir string, method string, dest string) map[string]string {
    hmap := make(map[string]string) // {hash:path}
    walker := fs.Walk(dir)
    for walker.Step() {
        // Start walking
        if err := walker.Err(); err != nil {
                fmt.Fprintln(os.Stderr, err)
                continue
                }
        // Check if it is a file
        finfo, err := os.Stat(walker.Path())
        if err != nil {
            fmt.Println(err)
            continue
            }
        if finfo.IsDir() {
            continue // it's a dir so pass and continue
        } else {
            // it's a file so process
            path := walker.Path()
            hash := hasher(walker.Path(), method) 
            search, ok := hmap[hash] 
            if ok {
                 _, filename := filepath.Split(path)
                if err := os.Rename(path, filepath.Join(dest, filename)); err != nil {
                fmt.Println(err)
                continue
                }
                fmt.Println("Duplicates moved =>", search)
            } else {
                hmap[hash] = path
                fmt.Println(hash, "=>", path)
            }
        }
    }
    return hmap
}

func main() {
    source := flag.String("s", "test", "Directory to scan")
    method := flag.String("m", "sha1", "Choose hashing method : md5 or sha1")
    destination := flag.String("d", "doublon", "Choose duplicates destination")
    flag.Parse()
    marcher(*source, *method, *destination)

}