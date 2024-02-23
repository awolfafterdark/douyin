package utils

import (
"bytes"
"compress/gzip"
)

func GzipEncode(input []byte) ([]byte, error) {
//Create a new byte output stream
var buf bytes.Buffer
//Create a new gzip output stream

//NoCompression = flate.NoCompression // No compression
//BestSpeed = flate.BestSpeed // Fastest speed
//BestCompression = flate.BestCompression // Best compression ratio
//DefaultCompression = flate.DefaultCompression //Default compression ratio
//gzip.NewWriterLevel()
gzipWriter := gzip.NewWriter(&buf)
//Write the input byte array to this output stream
_, err := gzipWriter.Write(input)
if err != nil {
_ = gzipWriter.Close()
return nil, err
}
if err := gzipWriter.Close(); err != nil {
return nil, err
}
// Return the compressed bytes array
return buf.Bytes(), nil
}

func GzipDecode(input []byte) ([]byte, error) {
//Create a new gzip.Reader
bytesReader := bytes.NewReader(input)
gzipReader, err := gzip.NewReader(bytesReader)
if err != nil {
return nil, err
}
defer func() {
// Close gzipReader in defer
_ = gzipReader.Close()
}()
buf := new(bytes.Buffer)
//Read data from Reader
if _, err := buf.ReadFrom(gzipReader); err != nil {
return nil, err
}
return buf.Bytes(), nil
}
