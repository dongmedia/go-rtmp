# go-rtmp

A lightweight Go package for handling RTMP (Real-Time Messaging Protocol) server connections and media streaming.

## Features

- RTMP handshake implementation
- Chunk-based message reading and writing
- Command message handling (connect, createStream, publish)
- Audio streaming support (AAC codec)
- Video streaming support (H.264/AVC codec)
- Media packet parsing and processing
- AMF0 encoding/decoding
- Stream management with channels

## Installation

```bash
go get github.com/dongmedia/go-rtmp
```

## Quick Start

Here's a basic example of setting up an RTMP server:

```go
package main

import (
    "log"
    "net"

    gortmp "github.com/dongmedia/go-rtmp"
)

func main() {
    listener, err := net.Listen("tcp", ":1935")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    log.Println("RTMP server listening on :1935")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Accept error:", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    rtmpConn := gortmp.NewConn(conn)
    rtmpConn.Serve()
}
```

## Architecture

### Core Components

#### Connection Management

The `Conn` type handles RTMP connections including:
- RTMP handshake negotiation
- Chunk reading and writing
- Command message processing
- Media packet routing

#### Stream Processing

The `Stream` type represents an RTMP stream with:
- Separate audio and video channels
- AAC audio configuration
- H.264 SPS/PPS parameter sets

#### Chunk Protocol

The package implements RTMP chunk protocol for:
- Reading chunked messages from the wire
- Parsing chunk headers (Basic Header and Message Header)
- Extracting message payloads

#### Message Types

Supported message types:
- **Type 8**: Audio data (AAC)
- **Type 9**: Video data (H.264)
- **Type 20**: Command messages (AMF0)

### Media Codecs

#### AAC Audio

The package parses AAC audio packets and extracts:
- AAC sequence headers (decoder configuration)
- AAC raw frames

#### H.264 Video

The package parses H.264 video packets and extracts:
- AVC sequence headers (SPS/PPS)
- Video frames (keyframes and inter-frames)
- NAL units

## API Overview

### Main Types

```go
// Create a new RTMP connection
func NewConn(c net.Conn) *Conn

// Start serving the connection
func (c *Conn) Serve()

// Create a new stream
func NewStream(id uint32) *Stream

// Start consuming stream packets
func ConsumeStream(s *Stream)
```

### Message Package

```go
// Parse AAC audio packet
func ParseAAC(data []byte) (*AACPacket, error)

// Parse H.264 video packet
func ParseH264(data []byte) (*H264Packet, error)

// Decode AMF0 command message
func DecodeCommand(data []byte) (*Command, error)
```

### AMF Package

```go
// Encode values using AMF0
func (e *Encoder) Encode(v interface{}) error

// Decode AMF0 values
func (d *Decoder) Decode() (interface{}, error)
```

## Protocol Support

This package implements core RTMP functionality including:

- RTMP handshake (C0, C1, C2, S0, S1, S2)
- Chunk stream protocol
- AMF0 encoding/decoding
- Command messages (connect, createStream, publish)
- Audio/Video data messages

## Use Cases

- Building RTMP streaming servers
- Receiving live video streams from encoders (OBS, FFmpeg, etc.)
- Processing real-time audio/video data
- Creating custom streaming applications

## Limitations

- Currently supports server-side functionality only
- Implements Format 0 chunks (full message headers)
- AAC and H.264 codecs only
- Basic command message support

## License

This project is licensed under the Boost Software License 1.0. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome. Please ensure your code follows Go best practices and includes appropriate tests.

## Links

- [RTMP Specification](https://www.adobe.com/devnet/rtmp.html)
- [Go Documentation](https://pkg.go.dev/github.com/dongmedia/go-rtmp)
